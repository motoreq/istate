// Copyright 2020 <>. All rights reserved.

package istate

import (
	"fmt"
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"strconv"
	"strings"
)

// "." is restricted eg: .docType
func (iState *iState) generateRelationalTables(obj map[string]interface{}, keyref string) (relationalTables [][]map[string]interface{}, iStateErr Error) {
	iStateLogger.Debugf("Inside generateRelationalTables")
	defer iStateLogger.Debugf("Exiting generateRelationalTables")

	for index, val := range obj {
		var newTable []map[string]interface{}
		field := reflect.TypeOf(iState.structRef).Field(iState.fieldJSONIndexMap[index])
		tableName := field.Tag.Get(iStateTag)
		if tableName == "" {
			continue
		}
		newTable = iState.traverseAndGenerateRelationalTable(val, tableName, keyref, index)
		if len(newTable) != 0 {
			relationalTables = append(relationalTables, newTable)
		}
	}
	return
}

//
func (iState *iState) traverseAndGenerateRelationalTable(val interface{}, tableName string, keyref string, previousKey string, depth ...int) (newTable []map[string]interface{}) {
	iStateLogger.Debugf("Inside traverseAndGenerateRelationalTable")
	defer iStateLogger.Debugf("Exiting traverseAndGenerateRelationalTable")
	switch kind := reflect.ValueOf(val).Kind(); kind {
	case reflect.Slice, reflect.Array:
		sliceLen := reflect.ValueOf(val).Len()
		newTableName := tableName
		switch {
		case len(depth) != 0:
			newTableName += separator + strconv.Itoa(depth[0])
			depth[0] = depth[0] + 1
		default:
			depth = []int{0}
		}
		for i := 0; i < sliceLen; i++ {
			innerVal := reflect.ValueOf(val).Index(i).Interface()
			innerTables := iState.traverseAndGenerateRelationalTable(innerVal, newTableName, keyref, previousKey, depth...)
			if len(innerTables) != 0 {
				newTable = append(newTable, innerTables...)
			}
		}
	case reflect.Map:
		prefix := previousKey + separator
		mapKeys := reflect.ValueOf(val).MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			fmt.Println("INITIAL PREFIX", prefix)
			curKey := prefix + mapKeys[i].String()
			keyToSendNext := curKey
			// If currentKey doesn't exist in fieldjsonindexmap, then its a real map, not a struct
			_, ok := iState.fieldJSONIndexMap[curKey]
			fmt.Printf("IS IT A MAP? %v, %v\n", !ok, curKey)
			if _, ok := iState.fieldJSONIndexMap[curKey]; !ok {
				newRow := make(map[string]interface{})
				newRow[docTypeField] = tableName
				newRow[valueField] = mapKeys[i].String()
				newRow[keyRefField] = keyref
				newTable = append(newTable, newRow)
				keyToSendNext = prefix
			}

			innerVal := reflect.ValueOf(val).MapIndex(mapKeys[i]).Interface()
			suffix := fmt.Sprintf("%s", mapKeys[i].Interface())
			fmt.Println("KEY TO SEND NEXT: ", keyToSendNext)
			innerTables := iState.traverseAndGenerateRelationalTable(innerVal, tableName+separator+suffix, keyref, keyToSendNext)
			if len(innerTables) != 0 {
				newTable = append(newTable, innerTables...)
			}

		}
	default:
		if val == nil {
			break
		}
		newRow := make(map[string]interface{})
		newRow[docTypeField] = tableName
		newRow[valueField] = reflect.ValueOf(val).Interface()
		newRow[keyRefField] = keyref
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
func (iState *iState) encodeState(oMap map[string]interface{}, keyref string) (encodedKeyValPairs map[string][]byte, iStateErr Error) {
	iStateLogger.Debugf("Inside encodeState")
	defer iStateLogger.Debugf("Exiting encodeState")

	encodedKeyValPairs = make(map[string][]byte)

	// RQ-Index key - value
	relationalTables, iStateErr := iState.generateRelationalTables(oMap, keyref)
	if iStateErr != nil {
		return
	}

	for i := 0; i < len(relationalTables); i++ {
		for j := 0; j < len(relationalTables[i]); j++ {
			encodedVal := ""
			valInterface := relationalTables[i][j][valueField]
			switch kind := reflect.ValueOf(valInterface).Kind(); kind {
			case reflect.Bool:
				switch valInterface.(bool) {
				case true:
					encodedVal = boolTrue
				default:
					encodedVal = boolFalse
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				numString := fmt.Sprintf("%v", valInterface)
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
				encodedVal = numEncodePrefix + separator + numString
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				numString := fmt.Sprintf("%v", valInterface)
				if _, ok := numDigits[len(numString)]; !ok {
					iStateErr = NewError(nil, 2009, len(numString))
					return
				}
				numEncodePrefix := positiveNum + numDigits[len(numString)]
				encodedVal = numEncodePrefix + separator + numString
			case reflect.Float32, reflect.Float64:
				numString := fmt.Sprintf("%v", valInterface)
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
				encodedVal = numEncodePrefix + separator + numString
			case reflect.String:
				encodedVal = valInterface.(string)
			default:
				iStateErr = NewError(nil, 2005, kind)
				return
			}

			encodedKey := relationalTables[i][j][docTypeField].(string) + separator + encodedVal + separator + keyref
			encodedKeyValPairs[encodedKey] = []byte(keyref)
		}
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
			_, ok := iState.fieldJSONIndexMap[fieldName]
			fmt.Printf("THIS IS IN A REAL MAP?: %v, %v\n", !ok, fieldName)
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
				// 	fmt.Println("This map is sent to recurse again:", sourceObjMap[fieldName], targetObjMap[fieldName])
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
			fmt.Printf("Default: %v, %v, Kinds: %v, %v\n", sourceObjMap[fieldName], targetObjMap[fieldName], reflect.ValueOf(sourceObjMap[fieldName]).Kind(), reflect.ValueOf(targetObjMap[fieldName]).Kind())
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
	fmt.Println("findMapDifference SOURCEM: ", sourceM)
	fmt.Println("findMapDifference TARGETM: ", targetM)
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
		fmt.Printf("findMapDifference: Deepequal: %v, %v\n", sourceM, targetM)
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
	fmt.Println("APPENDM: ", appendM)
	fmt.Println("DELETEM: ", deleteM)
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
func getPrimaryFieldIndex(object interface{}) (outi int, iStateError Error) {
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
		iStateError = NewError(nil, 2002, reflect.TypeOf(object))
		return
	}
	return
}
