// Copyright 2020 <>. All rights reserved.

package istate

import (
	"fmt"
	"reflect"
	"strings"
)

// "." is restricted eg: .docType
func (iState *iState) generateRelationalTables(obj map[string]interface{}, keyref string, isQuery bool) (relationalTables [][]map[string]interface{}, iStateErr Error) {
	iStateLogger.Debugf("Inside generateRelationalTables")
	defer iStateLogger.Debugf("Exiting generateRelationalTables")

	for index, val := range obj {
		var newTable []map[string]interface{}
		field := reflect.TypeOf(iState.structRef).Field(iState.fieldJSONIndexMap[index])
		tableName := field.Tag.Get(iStateTag)
		if tableName == "" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			iStateErr = NewError(nil, 2001, field.Name, reflect.TypeOf(iState.structRef))
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
	iStateLogger.Debugf("Inside traverseAndGenerateRelationalTable")
	defer iStateLogger.Debugf("Exiting traverseAndGenerateRelationalTable")
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
			tableNameString := joinStringInterfaceSlice(tableName, seperator)
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
			// fmt.Printf("1 ADVANCED COUNTER: CURRENT DEPTH KEY: %v, For Row: %v\n", currentDepth, newRow)
			newTable = append(newTable, newRow)
		}
	case reflect.Map:
		mapKeys := reflect.ValueOf(val).MapKeys()
		currentDepth := joinStringInterfaceSlice(append([]interface{}{jsonTag}, genericTableName...), seperator) + seperator
		// fmt.Println("SEARCHED CURRENT DEPTH: ", currentDepth)

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
				case meta[0].(int) > 1 && !isQuery:
					newRow := make(map[string]interface{})
					newGenericTableName := []interface{}{iStateTag}
					newRow[docTypeField] = append(newGenericTableName, genericTableName...)
					newRow[valueField] = convertedMapKey
					newRow[keyRefField] = keyref
					currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, genericTableName...))
					newRow[fieldNameField] = currentDepth
					// fmt.Printf("2.1 ADVANCED COUNTER: CURRENT DEPTH KEY: %v, For Row: %v\n", currentDepth, newRow)
					newTable = append(newTable, newRow)
				default:
					newRow := make(map[string]interface{})
					currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, genericTableName...))
					newRow[fieldNameField] = currentDepth
					// fmt.Printf("2.2 ADVANCED COUNTER: CURRENT DEPTH KEY: %v, For Row: %v\n", currentDepth, newRow)
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
			tableNameString := joinStringInterfaceSlice(tableName, seperator)
			tableNameString = removeLastSeparators(tableNameString)
			if tableNameString == iStateTag {
				break
			}
			newRow := make(map[string]interface{})
			newRow[docTypeField] = tableName
			newRow[valueField] = ""
			newRow[keyRefField] = keyref
			// currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, tableName[1:]...))
			//newRow[fieldNameField] = currentDepth
			// fmt.Printf("3 ADVANCED COUNTER: CURRENT DEPTH KEY: %v, For Row: %v\n", currentDepth, newRow)
			newTable = append(newTable, newRow)
		}
	default:
		if val == nil {
			break
		}

		newGenericTableName := []interface{}{iStateTag}
		newGenericTableName = append(newGenericTableName, meta[1].([]interface{})...)

		newGenericTableNameString := joinStringInterfaceSlice(newGenericTableName, seperator) // Here genericPrefix is not taken, as "_" will be added when encoding.
		// Making inner elements searchable
		tableNameString := joinStringInterfaceSlice(tableName, seperator)
		addedGenericRow := false
		if meta[0].(int) > 0 && newGenericTableNameString != tableNameString && !isQuery {
			newRow := make(map[string]interface{})
			newRow[docTypeField] = newGenericTableName
			newRow[valueField] = reflect.ValueOf(val).Interface()
			newRow[keyRefField] = keyref
			currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, genericTableName...))
			newRow[fieldNameField] = currentDepth
			// fmt.Printf("4 ADVANCED COUNTER: CURRENT DEPTH KEY: %v, For Row: %v\n", currentDepth, newRow)
			newTable = append(newTable, newRow)
			addedGenericRow = true
		}

		newRow := make(map[string]interface{})
		newRow[docTypeField] = tableName
		newRow[valueField] = reflect.ValueOf(val).Interface()
		newRow[keyRefField] = keyref
		switch !isQuery {
		case true && !addedGenericRow:
			currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, tableName[1:]...))
			newRow[fieldNameField] = currentDepth
			// fmt.Printf("5 ADVANCED COUNTER: CURRENT DEPTH KEY: %v, For Row: %v\n", currentDepth, newRow)
		default:
			currentDepth := joinStringInterfaceSliceWithDotStar(append([]interface{}{jsonTag}, genericTableName...))
			newRow[fieldNameField] = currentDepth
		}
		newTable = append(newTable, newRow)
	}
	return
}

//
func (iState *iState) getPrimaryKey(object interface{}) (key string) {
	iStateLogger.Debugf("Inside getPrimaryKey")
	defer iStateLogger.Debugf("Exiting getPrimaryKey")

	reflectVal := reflect.ValueOf(object)
	key = fmt.Sprintf("%v", reflectVal.Field(iState.primaryIndex).Interface())
	return
}

//
func (iState *iState) encodeState(oMap map[string]interface{}, keyref string, isQuery ...bool) (encodedKeyValPairs map[string][]byte, docsCounter map[string]int, encKeyDocNameMap map[string]string, iStateErr Error) {
	iStateLogger.Debugf("Inside encodeState")
	defer iStateLogger.Debugf("Exiting encodeState")
	if len(isQuery) == 0 {
		isQuery = []bool{false}
	}
	docsCounter = make(map[string]int)
	encodedKeyValPairs = make(map[string][]byte)
	encKeyDocNameMap = make(map[string]string)
	// RQ-Index key - value
	relationalTables, iStateErr := iState.generateRelationalTables(oMap, keyref, isQuery[0])
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
					encodedKeyParts[k], iStateErr = encode(encodedKeyParts[k])
					if iStateErr != nil {
						return
					}
				}
				encodedVal := ""
				valInterface := relationalTables[i][j][valueField]
				encodedVal, iStateErr = encode(valInterface)
				if iStateErr != nil {
					return
				}

				encodedKey = joinStringInterfaceSlice(encodedKeyParts, seperator) + seperator + encodedVal + seperator + keyref

				switch _, ok := encodedKeyValPairs[encodedKey]; ok {
				case true:
					repeatedRow = true
				default:
					encodedKeyValPairs[encodedKey] = []byte(keyref)
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

func encode(value interface{}) (encodedVal string, iStateErr Error) {
	switch kind := reflect.ValueOf(value).Kind(); kind {
	case reflect.Bool:
		switch value.(bool) {
		case true:
			encodedVal = boolTrue
		default:
			encodedVal = boolFalse
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
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
			iStateErr = NewError(nil, 2009, len(numString))
			return
		}
		numEncodePrefix += numDigits[len(numString)]
		encodedVal = numEncodePrefix + seperator + numString
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		numString := fmt.Sprintf("%v", value)
		if _, ok := numDigits[len(numString)]; !ok {
			iStateErr = NewError(nil, 2009, len(numString))
			return
		}
		numEncodePrefix := positiveNum + numDigits[len(numString)]
		encodedVal = numEncodePrefix + seperator + numString
	case reflect.Float32, reflect.Float64:
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
			iStateErr = NewError(nil, 2009, len(wholeNum))
			return
		}
		numEncodePrefix += numDigits[len(wholeNum)]
		encodedVal = numEncodePrefix + seperator + numString
	case reflect.String:
		encodedVal = value.(string)
	default:
		iStateErr = NewError(nil, 2005, kind)
		return
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
			// if sourceObjMap[fieldName] == nil || targetObjMap[fieldName] == nil {
			// 	switch sourceObjMap[fieldName] {
			// 	case targetObjMap[fieldName]:
			// 		continue
			// 	case nil:
			// 		appendOrModifyMap[fieldName] = targetObjMap[fieldName]
			// 		continue
			// 	default:
			// 		deleteMap[fieldName] = sourceObjMap[fieldName]
			// 		continue
			// 	}
			// }
			// _, ok := iState.fieldJSONIndexMap[fieldName]

			switch {
			//case isRealMap(sourceObjMap[fieldName].(map[string]interface{}), targetObjMap[fieldName].(map[string]interface{})):
			// Its a real map, if the field name does not exist in structref
			// case !ok:
			default:
				var appendM, deleteM interface{}
				// isNilSource := reflect.ValueOf(reflect.ValueOf(sourceObjMap[fieldName])).Field(0).IsNil()
				// isNilTarget := reflect.ValueOf(reflect.ValueOf(targetObjMap[fieldName])).Field(0).IsNil()
				// if isNilSource || isNilTarget {
				// 	switch {
				// 	case isNilSource == isNilTarget:
				// 		continue
				// 	case isNilSource:
				// 		appendOrModifyMap[fieldName] = targetObjMap[fieldName]
				// 		continue
				// 	default:
				// 		deleteMap[fieldName] = sourceObjMap[fieldName]
				// 		continue
				// 	}
				// }
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
				// default:
				// 	var tempAppendOrModifyMap, tempDeleteMap map[string]interface{}
				//
				// 	tempAppendOrModifyMap, tempDeleteMap, iStateErr = iState.findDifference(sourceObjMap[fieldName].(map[string]interface{}), targetObjMap[fieldName].(map[string]interface{}))
				// 	if iStateErr != nil {
				// 		return
				// 	}
				// 	if reflect.ValueOf(tempAppendOrModifyMap).Len() != 0 {
				// 		appendOrModifyMap[fieldName] = tempAppendOrModifyMap
				// 	}
				// 	if reflect.ValueOf(tempDeleteMap).Len() != 0 {
				// 		deleteMap[fieldName] = tempDeleteMap
				// 	}
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
			//case reflect.Invalid:
			//	fallthrough
			case reflect.String:
				if !reflect.DeepEqual(sourceObjMap[fieldName], targetObjMap[fieldName]) {
					appendOrModifyMap[fieldName] = targetObjMap[fieldName]
					deleteMap[fieldName] = sourceObjMap[fieldName]
				}
			// case reflect.Invalid:
			// 	appendOrModifyMap[fieldName] = targetObjMap[fieldName]
			default:
				iStateErr = NewError(nil, 2010, k)
			}
		}

	}
	return
}

func findSliceDifference(sourceS interface{}, targetS interface{}) (appendS interface{}, deleteS interface{}, iStateErr Error) {
	if reflect.TypeOf(sourceS) != reflect.TypeOf(targetS) {
		iStateErr = NewError(nil, 2011, reflect.TypeOf(sourceS), reflect.TypeOf(targetS))
		return
	}
	if reflect.ValueOf(sourceS).Kind() != reflect.Slice {
		iStateErr = NewError(nil, 2012, reflect.TypeOf(sourceS), reflect.TypeOf(targetS))
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

	// if reflect.TypeOf(sourceM) != reflect.TypeOf(targetM) {
	// 	iStateErr = NewError(nil, 2013, reflect.TypeOf(sourceM), reflect.TypeOf(targetM))
	// 	return
	// }
	// if sKind, tKind := reflect.ValueOf(sourceM).Kind(), reflect.ValueOf(targetM).Kind(); (sKind != reflect.Map || tKind != reflect.Map) {
	// 	if
	// 	iStateErr = NewError(nil, 2014, reflect.TypeOf(sourceM), reflect.TypeOf(targetM))
	// 	return
	// }

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
	filledStruct = reflect.New(reflect.TypeOf(structRef)).Elem()
	refVal := reflect.ValueOf(structRef)
	switch refVal.Kind() {
	case reflect.Struct:
		for i := 0; i < refVal.NumField(); i++ {
			curField := refVal.Field(i)
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
				iStateErr = NewError(nil, 2001, field.Name, reflect.TypeOf(structRef))
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
	iStateLogger.Debugf("Inside getPrimaryFieldIndex")
	defer iStateLogger.Debugf("Exiting getPrimaryFieldIndex")
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
		iStateErr = NewError(nil, 2002, reflect.TypeOf(object))
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
				iStateErr = NewError(nil, 2001, field.Name, reflect.TypeOf(structRef))
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
		prefix = meta[0].(string) + seperator
	default:
		meta = []interface{}{"", 0}
	}
	refVal := reflect.ValueOf(structRef)
	switch kind := refVal.Kind(); kind {
	case reflect.Slice, reflect.Array:
		sliceLen := refVal.Len()
		fieldName := prefix
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
				iStateErr = NewError(nil, 2001, field.Name, reflect.TypeOf(structRef))
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

func removeLastSeparators(input string) (val string) {
	val = input
	endIndex := -1
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] != seperator[0] { // separator[0] is '_' (not "_")
			break
		}
		endIndex = i
	}
	if endIndex != -1 {
		val = input[:endIndex]
	}
	return
}

func joinStringInterfaceSlice(slice []interface{}, seperatorString string) (joinedString string) {
	if len(slice) > 0 {
		joinedString = fmt.Sprintf("%v", slice[0])
	}
	for i := 1; i < len(slice); i++ {
		joinedString += seperatorString + fmt.Sprintf("%v", slice[i])
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

func initQueryEnv(qEnv *queryEnv) {
	qEnv.ufetchedKVMap = make(map[string][]byte)
	qEnv.ukeyEncKVMap = make(map[string]map[string][]byte)
}
