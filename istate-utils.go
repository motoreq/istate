/*
	Copyright 2020 Motoreq Infotech Pvt Ltd

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package istate

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/trees/btree"
	"reflect"
	"strconv"
	"strings"
)

// "." is restricted eg: .docType
func (iState *iState) generateRelationalTables(obj map[string]interface{}, keyref string, isQuery bool) (relationalTables [][]map[string]interface{}, iStateErr Error) {
	// iStateLogger.Debugf("Inside generateRelationalTables")
	// defer iStateLogger.Debugf("Exiting generateRelationalTables")

	for index, val := range obj {
		var newTable []map[string]interface{}
		field := reflect.TypeOf(iState.structRef).Field(iState.fieldJSONIndexMap[index])
		tableName := field.Tag.Get(iStateTag)
		if tableName == "" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			iStateErr = newError(nil, 2001, field.Name, reflect.TypeOf(iState.structRef))
			return
		}
		newTable, iStateErr = iState.traverseAndGenerateRelationalTable(val, []interface{}{tableName}, jsonTag, tableName, keyref, isQuery)
		if iStateErr != nil {
			return
		}
		if len(newTable) != 0 {
			relationalTables = append(relationalTables, newTable)
		}
	}
	return
}

//
func (iState *iState) traverseAndGenerateRelationalTable(val interface{}, tableName []interface{}, jsonTag string, iStateTag string, keyref string, isQuery bool, meta ...interface{}) (newTable []map[string]interface{}, iStateErr Error) {
	// iStateLogger.Debugf("Inside traverseAndGenerateRelationalTable")
	// defer iStateLogger.Debugf("Exiting traverseAndGenerateRelationalTable")
	// meta[0] Universal depth
	// meta[1] appended with "_" (Previous key)
	genericTableName := []interface{}{}
	switch len(meta) > 1 {
	case true:
		genericTableName = meta[1].([]interface{})
	default:
		meta = []interface{}{0, tableName}

	}
	switch kind := reflect.ValueOf(val).Kind(); kind {
	case reflect.Slice, reflect.Array:
		sliceLen := reflect.ValueOf(val).Len()
		newTableName := append(tableName, "")
		for i := 0; i < sliceLen; i++ {
			innerVal := reflect.ValueOf(val).Index(i).Interface()
			var innerTables []map[string]interface{}
			innerTables, iStateErr = iState.traverseAndGenerateRelationalTable(innerVal, newTableName, jsonTag, iStateTag, keyref, isQuery, meta[0].(int)+1, append(genericTableName, ""))
			if iStateErr != nil {
				return
			}
			if len(innerTables) != 0 {
				newTable = append(newTable, innerTables...)
			}
		}
		// For empty slice as leaf value
		if sliceLen == 0 && !isQuery {
			tableNameString := joinStringInterfaceSlice(tableName, separator)
			tableNameString = removeLastSeparators(tableNameString)
			if tableNameString == iStateTag {
				break
			}
			newRow := make(map[string]interface{})
			newRow[docTypeField] = tableName
			newRow[valueField] = ""
			newRow[keyRefField] = keyref
			currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, tableName[1:]...))
			newRow[fieldNameField] = currentDepth
			newTable = append(newTable, newRow)
		}
	case reflect.Map:
		mapKeys := reflect.ValueOf(val).MapKeys()
		currentDepth := joinStringInterfaceSlice(append([]interface{}{jsonTag}, genericTableName...), separator) + separator

		for i := 0; i < len(mapKeys); i++ {
			var nextInitialTableName []interface{}
			var convertedMapKey interface{} = mapKeys[i].String()
			// This is enough to see if it is realmap or not
			switch kind, realMap := iState.depthKindMap[currentDepth]; realMap {
			case true:
				convertedMapKey, iStateErr = convertToPrimitiveType(mapKeys[i].String(), kind)
				if iStateErr != nil {
					return
				}
				nextInitialTableName = append(genericTableName, "")
				switch {
				default:
					newRow := make(map[string]interface{})
					currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, genericTableName...))
					newRow[fieldNameField] = currentDepth
					newTable = append(newTable, newRow)
				}
			default:
				nextInitialTableName = append(genericTableName, convertedMapKey)
			}

			innerVal := reflect.ValueOf(val).MapIndex(mapKeys[i]).Interface()
			var innerTables []map[string]interface{}
			innerTables, iStateErr = iState.traverseAndGenerateRelationalTable(innerVal, append(tableName, convertedMapKey), jsonTag, iStateTag, keyref, isQuery, meta[0].(int)+1, nextInitialTableName)
			if iStateErr != nil {
				return
			}
			if len(innerTables) != 0 {
				newTable = append(newTable, innerTables...)
			}

		}
		// For empty structs as leaf value
		if len(mapKeys) == 0 && !isQuery {
			tableNameString := joinStringInterfaceSlice(tableName, separator)
			tableNameString = removeLastSeparators(tableNameString)
			if tableNameString == iStateTag {
				break
			}
			newRow := make(map[string]interface{})
			newRow[docTypeField] = tableName
			newRow[valueField] = ""
			newRow[keyRefField] = keyref
			newTable = append(newTable, newRow)
		}
	default:
		if val == nil {
			break
		}
		var currentDepthStar string
		var currentDepth string
		// newGenericTableName := []interface{}{iStateTag}
		// newGenericTableName = append(newGenericTableName, meta[1].([]interface{})...)
		newRow := make(map[string]interface{})
		switch !isQuery {
		case true: //&& !addedGenericRow:

			currentDepthStar = joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, tableName[1:]...))
			currentDepth = joinStringInterfaceSlice(append([]interface{}{jsonTag}, tableName[1:]...), separator)
			newRow[fieldNameField] = currentDepthStar

			kind, ok := iState.depthKindMap[currentDepth]
			// Newly added
			if !ok {
				currentDepth = joinStringInterfaceSlice(append([]interface{}{jsonTag}, genericTableName...), separator)
				kind, ok = iState.depthKindMap[currentDepth]
			}
			// Newly added end
			if !ok {
				iStateErr = newError(nil, 2016)
				return
			}
			valAsString := fmt.Sprintf("%v", reflect.ValueOf(val).Interface())
			var convertedVal interface{}
			convertedVal, iStateErr = convertToPrimitiveType(valAsString, kind)
			if iStateErr != nil {
				return
			}
			newRow[docTypeField] = tableName
			newRow[valueField] = convertedVal
			newRow[keyRefField] = keyref
		default:

			currentDepthStar = joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, genericTableName...))
			// currentDepth = joinStringInterfaceSlice(append([]interface{}{jsonTag}, genericTableName...), separator)
			newRow[fieldNameField] = currentDepthStar
			newRow[docTypeField] = tableName
			newRow[valueField] = reflect.ValueOf(val).Interface()
			newRow[keyRefField] = keyref
		}

		newTable = append(newTable, newRow)
	}
	return
}

//
func (iState *iState) getPrimaryKey(object interface{}) (key string) {
	// iStateLogger.Debugf("Inside getPrimaryKey")
	// defer iStateLogger.Debugf("Exiting getPrimaryKey")

	reflectVal := reflect.ValueOf(object)
	key = fmt.Sprintf("%v", reflectVal.Field(iState.primaryIndex).Interface())
	return
}

//
func (iState *iState) encodeState(oMap map[string]interface{}, keyref string, hashString string, flags ...interface{}) (encodedKeyValPairs map[string][]byte, docsCounter map[string]int, encKeyDocNameMap map[string]string, iStateErr Error) {
	// flags[0]: int -> 0: no separation, 1: keyRefSeparated index, 2: keyRef&ValueSeparated index (default: 0)
	// flags[1]: isQuery (default: false)
	// iStateLogger.Debugf("Inside encodeState")
	// defer iStateLogger.Debugf("Exiting encodeState")
	switch len(flags) {
	case 0:
		flags = []interface{}{0, false}
	case 1:
		flags = []interface{}{flags[0], false}
	}
	separation := flags[0].(int)
	isQuery := flags[1].(bool)
	docsCounter = make(map[string]int)
	encodedKeyValPairs = make(map[string][]byte)
	encKeyDocNameMap = make(map[string]string)
	// RQ-Index key - value
	relationalTables, iStateErr := iState.generateRelationalTables(oMap, keyref, isQuery)
	if iStateErr != nil {
		return
	}

	for i := 0; i < len(relationalTables); i++ {
		for j := 0; j < len(relationalTables[i]); j++ {
			repeatedRow := false
			encodedKey := ""
			if docType, ok := relationalTables[i][j][docTypeField]; ok {
				encodedKeyParts := docType.([]interface{})
				for k := 0; k < len(encodedKeyParts); k++ {
					encodedKeyParts[k], iStateErr, _ = encode(encodedKeyParts[k])
					if iStateErr != nil {
						return
					}
				}
				encodedVal := ""
				// isNum := false
				valInterface := relationalTables[i][j][valueField]
				encodedVal, iStateErr, _ = encode(valInterface)
				if iStateErr != nil {
					return
				}

				// if isNum {
				// 	encodedVal = encodedVal
				// }

				switch separation {
				case 0:
					encodedKey = joinStringInterfaceSlice(encodedKeyParts, separator) + separator + encodedVal + separator + keyref
				case 1:
					encodedKey = joinStringInterfaceSlice(encodedKeyParts, separator) + separator + encodedVal + separator
				case 2:
					encodedKey = joinStringInterfaceSlice(encodedKeyParts, separator) + separator
				}

				switch _, ok := encodedKeyValPairs[encodedKey]; ok {
				case true:
					repeatedRow = true
				default:
					// encodedKeyValPairs[encodedKey] = []byte(keyref)
					// encodedKeyValPairs[encodedKey] = nullByte
					derivedKeys := deriveIndexKeys(encodedKey, isQuery)
					switch isQuery {
					case true:
						// Selecting the one with less stars
						bestSoFar := encodedKey
						for i := 0; i < len(derivedKeys); i++ {
							if hasLessOrEqStars(bestSoFar, derivedKeys[i]) {
								bestSoFar = derivedKeys[i]
							}
							//encodedKeyValPairs[derivedKeys[i]] = []byte(hashString)
						}
						encodedKeyValPairs[bestSoFar] = []byte(hashString)
					default:
						encodedKeyValPairs[encodedKey] = []byte(hashString)
						for i := 0; i < len(derivedKeys); i++ {
							encodedKeyValPairs[derivedKeys[i]] = []byte(hashString)
						}
					}
				}
			}

			if fieldNameRow, ok := relationalTables[i][j][fieldNameField]; ok && !repeatedRow {
				if encodedKey != "" {
					encKeyDocNameMap[encodedKey] = fieldNameRow.(string)
				}
				docsCounter[fieldNameRow.(string)]++
			}
		}
	}

	return
}

func encode(value interface{}) (encodedVal string, iStateErr Error, isNum bool) {
	switch kind := reflect.ValueOf(value).Kind(); kind {
	case reflect.Bool:
		switch value.(bool) {
		case true:
			encodedVal = boolTrue
		default:
			encodedVal = boolFalse
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		isNum = true
		numString := fmt.Sprintf("%v", value)
		numEncodePrefix := ""
		switch strings.HasPrefix(numString, "-") {
		case true:
			numEncodePrefix += negativeNum
			numString = numString[1:]
		default:
			numEncodePrefix += positiveNum
		}
		if _, ok := numDigits[len(numString)]; !ok {
			iStateErr = newError(nil, 2009, len(numString))
			return
		}
		numEncodePrefix += numDigits[len(numString)]
		encodedVal = numSym + numEncodePrefix + numSeparator + numString
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		isNum = true
		numString := fmt.Sprintf("%v", value)
		if _, ok := numDigits[len(numString)]; !ok {
			iStateErr = newError(nil, 2009, len(numString))
			return
		}
		numEncodePrefix := positiveNum + numDigits[len(numString)]
		encodedVal = numSym + numEncodePrefix + numSeparator + numString
	case reflect.Float32, reflect.Float64:
		isNum = true
		numString := fmt.Sprintf("%v", value)
		numEncodePrefix := ""
		switch strings.HasPrefix(numString, "-") {
		case true:
			numEncodePrefix += negativeNum
			numString = numString[1:]
		default:
			numEncodePrefix += positiveNum
		}
		wholeNum := strings.Split(numString, ".")[0]
		if _, ok := numDigits[len(wholeNum)]; !ok {
			iStateErr = newError(nil, 2009, len(wholeNum))
			return
		}
		numEncodePrefix += numDigits[len(wholeNum)]
		encodedVal = numSym + numEncodePrefix + numSeparator + numString
	case reflect.String:
		encodedVal = value.(string)
	default:
		iStateErr = newError(nil, 2005, kind)
		return
	}
	return
}

func isNum(indexVal string) (isNum bool) {
	if len(indexVal) > 0 {
		if string(indexVal[0]) == numSym {
			isNum = true
		}
	}
	return
}

func isPositive(numVal string) (positive bool, iStateErr Error) {
	if len(numVal) < 2 {
		iStateErr = newError(nil, 2017)
		return
	}
	// numVal[1] because numVal[0] has numSym / num marker
	if numVal[1] == positiveNum[0] {
		positive = true
	}
	return
}

//
func (iState *iState) findDifference(sourceObjMap map[string]interface{}, targetObjMap map[string]interface{}) (appendOrModifyMap map[string]interface{}, deleteMap map[string]interface{}, iStateErr Error) {
	count := 0
	deleteMap = make(map[string]interface{})
	appendOrModifyMap = make(map[string]interface{})
	// Helps for recursive call with real map
	combinedKeysMap := make(map[string]struct{})
	for fieldName := range sourceObjMap {
		combinedKeysMap[fieldName] = struct{}{}
	}
	for fieldName := range targetObjMap {
		combinedKeysMap[fieldName] = struct{}{}
	}

	for fieldName := range combinedKeysMap {

		count++
		switch {
		case reflect.ValueOf(sourceObjMap[fieldName]).Kind() == reflect.Slice || reflect.ValueOf(targetObjMap[fieldName]).Kind() == reflect.Slice:
			var appendS, deleteS interface{}
			appendS, deleteS, iStateErr = findSliceDifference(sourceObjMap[fieldName], targetObjMap[fieldName])
			if iStateErr != nil {
				return
			}
			if appendS != nil && reflect.ValueOf(appendS).Len() != 0 {
				appendOrModifyMap[fieldName] = appendS
			}
			if deleteS != nil && reflect.ValueOf(deleteS).Len() != 0 {
				deleteMap[fieldName] = deleteS
			}

		case fmt.Sprintf("%T", sourceObjMap[fieldName]) == "map[string]interface {}" || fmt.Sprintf("%T", targetObjMap[fieldName]) == "map[string]interface {}":

			var appendM, deleteM interface{}

			appendM, deleteM, iStateErr = findMapDifference(sourceObjMap[fieldName], targetObjMap[fieldName])
			if iStateErr != nil {
				return
			}
			if appendM != nil && reflect.ValueOf(appendM).Len() != 0 {
				appendOrModifyMap[fieldName] = appendM
			}
			if deleteM != nil && reflect.ValueOf(deleteM).Len() != 0 {
				deleteMap[fieldName] = deleteM
			}

		default:
			switch k := reflect.ValueOf(sourceObjMap[fieldName]).Kind(); k {
			case reflect.Bool:
				fallthrough
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fallthrough
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fallthrough
			case reflect.Float32, reflect.Float64:
				fallthrough
			case reflect.String:
				if !reflect.DeepEqual(sourceObjMap[fieldName], targetObjMap[fieldName]) {
					appendOrModifyMap[fieldName] = targetObjMap[fieldName]
					deleteMap[fieldName] = sourceObjMap[fieldName]
				}
			case reflect.Invalid:
				switch kt := reflect.ValueOf(targetObjMap[fieldName]).Kind(); kt {
				// If both are invalid, continue (empty complex type objects returns invalid)
				case reflect.Invalid:
					continue
				default:
					appendOrModifyMap[fieldName] = targetObjMap[fieldName]
				}
			default:
				iStateErr = newError(nil, 2010, k)
				return
			}
		}

	}
	return
}

func findSliceDifference(sourceS interface{}, targetS interface{}) (appendS interface{}, deleteS interface{}, iStateErr Error) {
	sourceType := fmt.Sprintf("%v", reflect.TypeOf(sourceS))
	targetType := fmt.Sprintf("%v", reflect.TypeOf(targetS))

	switch sourceType == "<nil>" || targetType == "<nil>" {
	case true:
		switch sourceType {
		// If both are nil
		case targetType:
			return
		case "<nil>":
			sourceS = reflect.New(reflect.TypeOf(targetS)).Elem().Interface()
		default:
			targetS = reflect.New(reflect.TypeOf(sourceS)).Elem().Interface()
		}

	}

	if reflect.TypeOf(sourceS) != reflect.TypeOf(targetS) {
		iStateErr = newError(nil, 2011, reflect.TypeOf(sourceS), reflect.TypeOf(targetS))
		return
	}
	if reflect.ValueOf(sourceS).Kind() != reflect.Slice {
		iStateErr = newError(nil, 2012, reflect.TypeOf(sourceS), reflect.TypeOf(targetS))
		return
	}

	if reflect.DeepEqual(sourceS, targetS) {
		return
	}

	indexSMap := make(map[int]struct{})
	indexTMap := make(map[int]struct{})

	lenS := reflect.ValueOf(sourceS).Len()
	lenT := reflect.ValueOf(targetS).Len()

	for i := 0; i < lenS; i++ {
		indexSMap[i] = struct{}{}
	}
	for i := 0; i < lenT; i++ {
		indexTMap[i] = struct{}{}
	}

	for indexS := range indexSMap {
		for indexT := range indexTMap {
			if reflect.DeepEqual(reflect.ValueOf(sourceS).Index(indexS).Interface(), reflect.ValueOf(targetS).Index(indexT).Interface()) {
				delete(indexSMap, indexS)
				delete(indexTMap, indexT)
			}
		}
	}

	appendS = reflect.MakeSlice(reflect.ValueOf(sourceS).Type(), len(indexTMap), len(indexTMap))
	deleteS = reflect.MakeSlice(reflect.ValueOf(targetS).Type(), len(indexSMap), len(indexSMap))
	// Whatever remaining in source is for deleting
	i := 0
	for index := range indexSMap {
		deleteS.(reflect.Value).Index(i).Set(reflect.ValueOf(sourceS).Index(index))
		i++
	}
	// Whatever remaining in target is for appending
	i = 0
	for index := range indexTMap {
		appendS.(reflect.Value).Index(i).Set(reflect.ValueOf(targetS).Index(index))
		i++
	}

	appendS = appendS.(reflect.Value).Interface()
	deleteS = deleteS.(reflect.Value).Interface()
	return
}

func findMapDifference(sourceM interface{}, targetM interface{}) (appendM interface{}, deleteM interface{}, iStateErr Error) {

	if reflect.DeepEqual(sourceM, targetM) {
		return
	}

	isNilSource := reflect.ValueOf(reflect.ValueOf(sourceM)).Field(0).IsNil()
	isNilTarget := reflect.ValueOf(reflect.ValueOf(targetM)).Field(0).IsNil()
	if !(isNilSource || isNilTarget) {
		mapKeysS := reflect.ValueOf(sourceM).MapKeys()
		for i := 0; i < len(mapKeysS); i++ {
			sourceVal := reflect.ValueOf(sourceM).MapIndex(mapKeysS[i])
			targetVal := reflect.ValueOf(targetM).MapIndex(mapKeysS[i])
			//notEmptyMap := !(fmt.Sprintf("%v", sourceVal) == "map[]" || fmt.Sprintf("%v", targetVal) == "map[]")
			//notInvalid := !(sourceVal.Kind() == reflect.Invalid || targetVal.Kind() == reflect.Invalid)
			isNilSource := reflect.ValueOf(sourceVal).Field(0).IsNil()
			isNilTarget := reflect.ValueOf(targetVal).Field(0).IsNil()
			switch isNilSource || isNilTarget {
			case true:
				switch isNilSource {
				case isNilTarget:
					// Delete field
					reflect.ValueOf(targetM).SetMapIndex(mapKeysS[i], reflect.Zero(targetVal.Type()).Elem())
					reflect.ValueOf(sourceM).SetMapIndex(mapKeysS[i], reflect.Zero(sourceVal.Type()).Elem())
				}
			default:
				if reflect.DeepEqual(sourceVal.Interface(), targetVal.Interface()) {
					// Delete field
					reflect.ValueOf(targetM).SetMapIndex(mapKeysS[i], reflect.Zero(targetVal.Type()).Elem())
					reflect.ValueOf(sourceM).SetMapIndex(mapKeysS[i], reflect.Zero(sourceVal.Type()).Elem())
				}
			}
		}
	}

	appendM = targetM
	deleteM = sourceM

	return
}

func isRealMap(source map[string]interface{}, target map[string]interface{}) (realMap bool) {
	if len(source) != len(target) {
		realMap = true
		return
	}
	combined := make(map[string]struct{})
	for index := range source {
		combined[index] = struct{}{}
	}
	for index := range target {
		combined[index] = struct{}{}
	}
	if len(combined) != len(source) {
		realMap = true
		return
	}
	return
}

//
func fillZeroValue(structRef interface{}) (filledStruct interface{}) {
	// Panics for empty interface
	filledStruct = reflect.New(reflect.TypeOf(structRef)).Elem()
	refVal := reflect.ValueOf(structRef)
	switch refVal.Kind() {
	case reflect.Struct:
		for i := 0; i < refVal.NumField(); i++ {
			curField := refVal.Field(i)
			// If empty interface
			// if fmt.Sprintf("%v", reflect.TypeOf(curField.Interface())) == "<nil>" {
			// 	continue
			// }
			filledField := fillZeroValue(curField.Interface())
			filledStruct.(reflect.Value).Field(i).Set(reflect.ValueOf(filledField))
		}
	case reflect.Slice:
		elemType := refVal.Type().Elem()
		tempSlice := reflect.New(elemType).Elem()
		filledSlice := fillZeroValue(tempSlice.Interface())

		filledStruct.(reflect.Value).Set(reflect.Append(filledStruct.(reflect.Value), reflect.ValueOf(filledSlice)))

	case reflect.Map:
		newMap := reflect.MakeMap(refVal.Type())
		keyType := refVal.Type().Key()
		elemType := refVal.Type().Elem()
		newKey := reflect.New(keyType).Elem()
		newElem := reflect.New(elemType).Elem()
		filledElem := fillZeroValue(newElem.Interface())
		newMap.SetMapIndex(newKey, reflect.ValueOf(filledElem))
		filledStruct = newMap
	default:
	}
	filledStruct = filledStruct.(reflect.Value).Interface()
	return
}

//
func generateFieldJSONIndexMap(structRef interface{}, fieldJSONIndexMap map[string]int, parentName ...interface{}) (iStateErr Error) {
	prefix := ""
	// Generate map of json tag :FieldIndex
	if len(parentName) != 0 {
		prefix = parentName[0].(string) + "_"
	}
	switch kind := reflect.ValueOf(structRef).Kind(); kind {
	case reflect.Slice, reflect.Array:
		sliceLen := reflect.ValueOf(structRef).Len()
		switch {
		case len(parentName) > 1:
			parentName[1] = parentName[1].(int) + 1
		default:
			parentName = append(parentName, 0)
		}
		for i := 0; i < sliceLen; i++ {
			innerVal := reflect.ValueOf(structRef).Index(i).Interface()
			//parentName[0] = prefix + strconv.Itoa(parentName[1].(int))
			iStateErr = generateFieldJSONIndexMap(innerVal, fieldJSONIndexMap, parentName...)
			if iStateErr != nil {
				return iStateErr
			}
		}
	case reflect.Map:
		mapKeys := reflect.ValueOf(structRef).MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			//suffix := fmt.Sprintf("%s", mapKeys[i].Interface())
			suffix := ""
			innerVal := reflect.ValueOf(structRef).MapIndex(mapKeys[i]).Interface()
			parentName[0] = prefix + suffix
			iStateErr = generateFieldJSONIndexMap(innerVal, fieldJSONIndexMap, parentName...)
			if iStateErr != nil {
				return iStateErr
			}
		}
	case reflect.Struct:
		for i := 0; i < reflect.ValueOf(structRef).NumField(); i++ {
			field := reflect.TypeOf(structRef).Field(i)
			jsonTag := prefix + field.Tag.Get("json")
			if field.Tag.Get("json") == "" {
				iStateErr = newError(nil, 2001, field.Name, reflect.TypeOf(structRef))
				return
			}
			fieldJSONIndexMap[jsonTag] = i
			iStateErr = generateFieldJSONIndexMap(reflect.ValueOf(structRef).Field(i).Interface(), fieldJSONIndexMap, jsonTag)
			if iStateErr != nil {
				return iStateErr
			}
		}
	default:
	}
	return
}

//
func getPrimaryFieldIndex(object interface{}) (outi int, iStateErr Error) {
	// iStateLogger.Debugf("Inside getPrimaryFieldIndex")
	// defer iStateLogger.Debugf("Exiting getPrimaryFieldIndex")
	outi = -1
	reflectVal := reflect.ValueOf(object)
	for i := 0; i < reflectVal.NumField(); i++ {
		field := reflect.TypeOf(object).Field(i)
		if field.Tag.Get(iStatePrimaryTag) == iStatePrimaryTrueVal {
			outi = i
			break
		}
	}
	if outi == -1 {
		iStateErr = newError(nil, 2002, reflect.TypeOf(object))
		return
	}
	return
}

//
func generatejsonFieldKindMap(structRef interface{}, jsonFieldKindMap map[string]reflect.Kind, mapKeyKindMap map[string]reflect.Kind, prev ...string) (iStateErr Error) {
	prefix := ""
	if len(prev) > 0 {
		prefix = prev[0] + dot
	}
	refVal := reflect.ValueOf(structRef)
	switch kind := refVal.Kind(); kind {
	case reflect.Slice, reflect.Array:
		sliceLen := refVal.Len()
		jsonTag := prefix + star
		for i := 0; i < sliceLen; i++ {
			innerVal := refVal.Index(i).Interface()
			jsonFieldKindMap[jsonTag] = refVal.Index(i).Kind()
			iStateErr = generatejsonFieldKindMap(innerVal, jsonFieldKindMap, mapKeyKindMap, jsonTag)
			if iStateErr != nil {
				return
			}
		}
	case reflect.Map:
		mapKeys := refVal.MapKeys()
		jsonTag := prefix + star
		for i := 0; i < len(mapKeys); i++ {
			innerVal := refVal.MapIndex(mapKeys[i]).Interface()
			jsonFieldKindMap[jsonTag] = refVal.MapIndex(mapKeys[i]).Kind()
			mapKeyKindMap[prev[0]] = mapKeys[i].Kind()
			iStateErr = generatejsonFieldKindMap(innerVal, jsonFieldKindMap, mapKeyKindMap, jsonTag)
			if iStateErr != nil {
				return
			}
		}
	case reflect.Struct:
		for i := 0; i < refVal.NumField(); i++ {
			field := reflect.TypeOf(structRef).Field(i)
			jsonTag := prefix + field.Tag.Get("json")
			if field.Tag.Get("json") == "" {
				iStateErr = newError(nil, 2001, field.Name, reflect.TypeOf(structRef))
				return
			}
			jsonFieldKindMap[jsonTag] = refVal.Field(i).Kind()
			iStateErr = generatejsonFieldKindMap(refVal.Field(i).Interface(), jsonFieldKindMap, mapKeyKindMap, jsonTag)
			if iStateErr != nil {
				return
			}

		}
	default:

	}
	return
}

//
func generateDepthKindMap(structRef interface{}, depthKindMap map[string]reflect.Kind, meta ...interface{}) (iStateErr Error) {
	// meta[0] = previousKey
	// meta[1] = depth
	prefix := ""
	switch len(meta) > 0 {
	case true:
		prefix = meta[0].(string) + separator
	default:
		meta = []interface{}{"", 0}
	}
	refVal := reflect.ValueOf(structRef)
	switch kind := refVal.Kind(); kind {
	case reflect.Slice, reflect.Array:
		sliceLen := refVal.Len()
		fieldName := prefix
		switch fmt.Sprintf("%s", refVal.Type()) == "[]uint8" {
		case true:
			iStateErr = newError(nil, 2020)
			return
		default:
			depthKindMap[fieldName] = refVal.Kind()
			break
		}
		for i := 0; i < sliceLen; i++ {
			innerVal := refVal.Index(i).Interface()
			// At this depth, the value will be slice's index like 0,1,2, its type Int
			depthKindMap[fieldName] = reflect.Int
			iStateErr = generateDepthKindMap(innerVal, depthKindMap, fieldName, meta[1].(int)+1)
			if iStateErr != nil {
				return
			}
		}
	case reflect.Map:
		mapKeys := refVal.MapKeys()
		fieldName := prefix
		for i := 0; i < len(mapKeys); i++ {
			innerVal := refVal.MapIndex(mapKeys[i]).Interface()
			//depthKindMap[fieldName] = refVal.MapIndex(mapKeys[i]).Kind()
			depthKindMap[fieldName] = mapKeys[i].Kind()
			iStateErr = generateDepthKindMap(innerVal, depthKindMap, fieldName, meta[1].(int)+1)
			if iStateErr != nil {
				return
			}
		}
	case reflect.Struct:
		for i := 0; i < refVal.NumField(); i++ {
			field := reflect.TypeOf(structRef).Field(i)
			jsonTag := prefix + field.Tag.Get("json")
			if field.Tag.Get("json") == "" {
				iStateErr = newError(nil, 2001, field.Name, reflect.TypeOf(structRef))
				return
			}
			depthKindMap[jsonTag] = refVal.Field(i).Kind()
			iStateErr = generateDepthKindMap(refVal.Field(i).Interface(), depthKindMap, jsonTag, meta[1].(int)+1)
			if iStateErr != nil {
				return
			}

		}
	default:
		fieldName := meta[0].(string) // Not using prefix here, to remove "_" at end
		depthKindMap[fieldName] = refVal.Kind()
	}
	return
}

// Plural
func removeLastSeparators(input string) (val string) {
	val = input
	endIndex := -1
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] != separator[0] { // separator[0] is '_' (not "_")
			break
		}
		endIndex = i
	}
	if endIndex != -1 {
		val = input[:endIndex]
	}
	return
}

func joinStringInterfaceSlice(slice []interface{}, separatorString string) (joinedString string) {
	if len(slice) > 0 {
		joinedString = fmt.Sprintf("%v", slice[0])
	}
	for i := 1; i < len(slice); i++ {
		joinedString += separatorString + fmt.Sprintf("%v", slice[i])
	}
	return
}

func joinStringInterfaceSliceWithDotStar(slice []interface{}) (joinedString string) {
	if len(slice) > 0 {
		joinedString = fmt.Sprintf("%v", slice[0])
	}
	for i := 1; i < len(slice); i++ {
		stringified := fmt.Sprintf("%v", slice[i])
		switch stringified == "" {
		case true:
			joinedString += dot + star
		default:
			joinedString += dot + fmt.Sprintf("%v", slice[i])
		}
	}
	return
}

func (iState *iState) convertIndexToQueryFieldName(index string) (joinedString string) {
	splitFields := strings.Split(index, separator)
	if len(splitFields) > 0 {
		// If permuted index, it will begin with ~<binary><separator><actual-istateTag>
		var istateTagIndex int
		if splitFields[0][0] == asciiLast[0] {
			istateTagIndex = 1
		}

		var ok bool
		joinedString, ok = iState.istateJSONMap[splitFields[istateTagIndex]]
		if !ok {
			// Not supposed to happen
			panic(fmt.Sprintf("istate tag not found in istateJSONMap: %v", splitFields[0]))
		}
	}
	for i := 1; i < len(splitFields); i++ {
		switch splitFields[i] == "" {
		case true:
			joinedString += dot + star
		default:
			joinedString += dot + splitFields[i]
		}
	}
	return
}

func initQueryEnv(qEnv *queryEnv) {
	qEnv.ufetchedKVMap = make(map[string][]byte)
	qEnv.uindecesMap = make(map[string]map[string][]byte)
}

// Singular
func removeLastSeparator(key string) (outKey string) {
	outKey = key
	if len(key) != 0 {
		if key[len(key)-1] == separator[0] { // separator[0] is '_' (not "_")
			outKey = key[:len(key)-1]
		}
	}
	return
}

func (iState *iState) generateActualStructure(field string, convertedVal interface{}) (newField string, newVal interface{}, iStateErr Error) {

	// Dot to Map
	splitfield := strings.Split(field, splitDot)
	switch len(splitfield) > 1 {
	case true:
		newField = splitfield[0]
		newVal = dotsToActualDepth(splitfield[1:], convertedVal)
	default:
		newField = field
		newVal = convertedVal
	}
	return
}

func decodeScientificNotation(sciNot string) (decoded string, iStateErr Error) {
	decoded = sciNot
	// If trying to convert scientific notation to int
	if strings.Contains(sciNot, "e") {
		floatVal, err := strconv.ParseFloat(sciNot, 64)
		if err != nil {
			iStateErr = newError(err, 3014)
			return
		}
		decoded = fmt.Sprintf("%.0f", floatVal)
	}
	return
}

func (iState *iState) getBestEncodedKeyFunc(querySet querys) (bestKey string, fetchFunc fetchFuncType, encodedKVSet encodedKVs, iStateErr Error) {
	encodedKVSet = encodedKVs{
		eq:    make(map[string][]byte),
		neq:   make(map[string][]byte),
		gt:    make(map[string][]byte),
		lt:    make(map[string][]byte),
		gte:   make(map[string][]byte),
		lte:   make(map[string][]byte),
		cmplx: make(map[string][]byte),

		seq:  make(map[string][]byte),
		sneq: make(map[string][]byte),
		sgt:  make(map[string][]byte),
		slt:  make(map[string][]byte),
		sgte: make(map[string][]byte),
		slte: make(map[string][]byte),
		// scmplx: make(map[string][]byte),
	}

	tree := btree.NewWithIntComparator(3)
	safe := &safeKeyFunc{}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.eq, encodedKVSet.eq, iState.fetchEq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.neq, encodedKVSet.neq, iState.fetchNeq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.gt, encodedKVSet.gt, iState.fetchGt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.lt, encodedKVSet.lt, iState.fetchLt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.gte, encodedKVSet.gte, iState.fetchGte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.lte, encodedKVSet.lte, iState.fetchLte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.cmplx, encodedKVSet.cmplx, iState.fetchCmplx, tree, safe)
	if iStateErr != nil {
		return
	}

	iStateErr = iState.generateEncKeysAndAddToTree(querySet.seq, encodedKVSet.seq, iState.fetchSeq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.sneq, encodedKVSet.sneq, iState.fetchSneq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.sgt, encodedKVSet.sgt, iState.fetchSgt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.slt, encodedKVSet.slt, iState.fetchSlt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.sgte, encodedKVSet.sgte, iState.fetchSgte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.slte, encodedKVSet.slte, iState.fetchSlte, tree, safe)
	if iStateErr != nil {
		return
	}
	// iStateErr = iState.generateEncKeysAndAddToTree(querySet.scmplx, encodedKVSet.scmplx, iState.fetchSlte, tree, safe)
	// if iStateErr != nil {
	// 	return
	// }

	leftNodeVal := tree.LeftValue()
	switch leftNodeVal == nil {
	case true:
		bestKey = safe.encKey
		fetchFunc = safe.fetchFunc
		beforei := (*safe.relatedQueryP)[:safe.i]
		afteri := (*safe.relatedQueryP)[safe.i+1:]
		*safe.relatedQueryP = append(beforei, afteri...)
		delete(safe.relatedEncKV, bestKey)
	default:

		bestKey = leftNodeVal.(efficientKeyType).enckey
		fetchFunc = leftNodeVal.(efficientKeyType).fetchFunc
		queryP := leftNodeVal.(efficientKeyType).relatedQueryP
		beforei := (*queryP)[:safe.i]
		afteri := (*queryP)[safe.i+1:]
		*queryP = append(beforei, afteri...)
		delete(leftNodeVal.(efficientKeyType).relatedEncKV, bestKey)
	}
	return
}

func hasLessStars(key1 string, key2 string) (ok bool) {
	count1 := strings.Count(key1, "*")
	count2 := strings.Count(key2, "*")
	if count2 < count1 {
		ok = true
	}
	return
}

func hasLessOrEqStars(key1 string, key2 string) (ok bool) {
	count1 := strings.Count(key1, "*")
	count2 := strings.Count(key2, "*")
	if count2 <= count1 {
		ok = true
	}
	return
}

func (iState *iState) generateEncKeysAndAddToTree(query []map[string]interface{}, encKVMap map[string][]byte, fetchFunc fetchFuncType, tree *btree.Tree, safe *safeKeyFunc) (iStateErr Error) {
	keyref := ""
	for i := 0; i < len(query); i++ {
		var encKeyDocNameMap map[string]string
		var encodedKeyVal map[string][]byte
		encodedKeyVal, _, encKeyDocNameMap, iStateErr = iState.encodeState(query[i], keyref, "", 1, true) // separation: 1, isQuery: true
		if iStateErr != nil {
			return
		}
		for index, val := range encodedKeyVal {

			genericField := encKeyDocNameMap[index]
			if numDocs, ok := iState.readDocsCounter(genericField); ok {
				addToTree(tree, index, genericField, numDocs, fetchFunc, encKVMap, &query, i)
			}

			encKVMap[index] = val
			if safe.fetchFunc == nil {
				safe.encKey = index
				safe.fetchFunc = fetchFunc
				safe.relatedQueryP = &query
				safe.relatedEncKV = encKVMap
				safe.i = i
			}
		}
	}
	return
}

func addToTree(tree *btree.Tree, encKey string, genericField string, numDocs int, fetchFunc fetchFuncType, relatedEncKV map[string][]byte, relatedQueryP *[]map[string]interface{}, i int) {
	efficientKey := efficientKeyType{
		enckey:        encKey,
		genericField:  genericField,
		fetchFunc:     fetchFunc,
		relatedEncKV:  relatedEncKV,
		relatedQueryP: relatedQueryP,
		i:             i,
	}
	switch val, ok := tree.Get(numDocs); ok {
	case true:
		if ok := hasLessStars(val.(efficientKeyType).enckey, efficientKey.enckey); ok {
			tree.Put(numDocs, efficientKey)
		}
	default:
		tree.Put(numDocs, efficientKey)
	}
}

// Clean this...
func generateistateJSONMap(obj interface{}, fieldJSONIndexMap map[string]int) (istateJSONMap map[string]string, iStateErr Error) {
	// iStateLogger.Debugf("Inside generateRelationalTables")
	// defer iStateLogger.Debugf("Exiting generateRelationalTables")
	istateJSONMap = make(map[string]string)
	var objMap map[string]interface{}
	mo, err := json.Marshal(obj)
	if err != nil {
		iStateErr = newError(err, 2018)
		return
	}
	err = json.Unmarshal(mo, &objMap)
	if err != nil {
		iStateErr = newError(err, 2019)
		return
	}

	for index := range objMap {
		field := reflect.TypeOf(obj).Field(fieldJSONIndexMap[index])
		istateTag := field.Tag.Get(iStateTag)
		if istateTag == "" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			iStateErr = newError(nil, 2001, field.Name, reflect.TypeOf(obj))
			return
		}
		istateJSONMap[istateTag] = jsonTag
	}
	return
}
