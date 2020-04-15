//

package istate

import (
	"reflect"
	"strconv"
	"strings"
)

func dotsToActualDepth(splitFieldName []string, val interface{}, curIndex ...int) (actualMap map[string]interface{}) {
	actualMap = make(map[string]interface{})
	if len(curIndex) == 0 {
		curIndex = []int{0}
	}

	if len(splitFieldName)-1 > curIndex[0] {
		actualMap[splitFieldName[curIndex[0]]] = dotsToActualDepth(splitFieldName, val, curIndex[0]+1)
	} else {
		actualMap[splitFieldName[curIndex[0]]] = val
	}

	return
}

func getRightPrimitiveType(fieldName string, jsonFieldKindMap map[string]reflect.Kind, mapKeyKindMap map[string]reflect.Kind) (primKind reflect.Kind, iStateErr Error) {

	splitFieldName := strings.Split(fieldName, splitDot)
	if len(splitFieldName) == 0 {
		iStateErr = NewError(nil, 3010)
		return
	}

	curField := splitFieldName[0]
	nextIndex := 1
	for {
		kind, ok := jsonFieldKindMap[curField]
		if !ok {
			iStateErr = NewError(nil, 3016, curField)
			return
		}
	SpecialFlow:
		switch kind {
		case reflect.Array, reflect.Slice:
			if len(splitFieldName) <= nextIndex {
				iStateErr = NewError(nil, 3015, fieldName)
				return
			}
			curField = curField + splitDot + star
			nextIndex++
			continue
		case reflect.Struct:
			if len(splitFieldName) <= nextIndex {
				iStateErr = NewError(nil, 3015, fieldName)
				return
			}
			curField = curField + splitDot + splitFieldName[nextIndex]
			nextIndex++
			continue
		case reflect.Map:
			// Is field only to be searched, or value too
			// If this is the last index
			switch len(splitFieldName) == nextIndex {
			case true:
				// Notice kind and ok are changed
				kind, ok = mapKeyKindMap[curField]
				if !ok {
					iStateErr = NewError(nil, 3016, curField)
					return
				}
				goto SpecialFlow
			default:
				if len(splitFieldName) <= nextIndex {
					iStateErr = NewError(nil, 3015, fieldName)
					return
				}
				//prevField := curField
				curField = curField + splitDot + star
				// kind, ok = jsonFieldKindMap[curField]
				//curField = prevField + splitDot + splitFieldName[nextIndex]
				nextIndex++
			}
			continue
		default: // Primitive type
			primKind = kind
		}
		break
	}
	return
}

func getKeyWordAndVal(queryToEvaluate string) (keyword, val string, iStateErr Error) {
	// Trim Space
	queryToEvaluate = strings.TrimSpace(queryToEvaluate)
	firstSpaceIndex := strings.Index(queryToEvaluate, " ")
	if firstSpaceIndex == -1 {
		iStateErr = NewError(nil, 3004, queryToEvaluate)
		return
	}
	keyword = queryToEvaluate[:firstSpaceIndex]
	val = queryToEvaluate[firstSpaceIndex+1:]
	return
}

func convertToRightType(fieldName string, toConvert string, jsonFieldKindMap map[string]reflect.Kind, mapKeyKindMap map[string]reflect.Kind) (convertedVal interface{}, iStateErr Error) {

	splitFieldName := strings.Split(fieldName, splitDot)
	if len(splitFieldName) == 0 {
		iStateErr = NewError(nil, 3010)
		return
	}

	curField := splitFieldName[0]
	nextIndex := 1
	for {
		kind, ok := jsonFieldKindMap[curField]
		if !ok {
			iStateErr = NewError(nil, 3016, curField)
			return
		}
	SpecialFlow:
		switch kind {
		case reflect.Array, reflect.Slice:
			if len(splitFieldName) <= nextIndex {
				iStateErr = NewError(nil, 3015, fieldName)
				return
			}
			curField = curField + splitDot + star
			nextIndex++
			continue
		case reflect.Struct:
			if len(splitFieldName) <= nextIndex {
				iStateErr = NewError(nil, 3015, fieldName)
				return
			}
			curField = curField + splitDot + splitFieldName[nextIndex]
			nextIndex++
			continue
		case reflect.Map:
			// Is field only to be searched, or value too
			// If this is the last index
			switch len(splitFieldName) == nextIndex {
			case true:
				// Notice kind and ok are changed
				kind, ok = mapKeyKindMap[curField]
				if !ok {
					iStateErr = NewError(nil, 3016, curField)
					return
				}
				goto SpecialFlow
			default:
				if len(splitFieldName) <= nextIndex {
					iStateErr = NewError(nil, 3015, fieldName)
					return
				}
				//prevField := curField
				curField = curField + splitDot + star
				// kind, ok = jsonFieldKindMap[curField]
				//curField = prevField + splitDot + splitFieldName[nextIndex]
				nextIndex++
			}
			continue
		default: // Primitive type
			convertedVal, iStateErr = convertToPrimitiveType(toConvert, kind)
			if iStateErr != nil {
				return
			}
		}
		break
	}
	return
}

func convertToPrimitiveType(toConvert string, kind reflect.Kind) (convertedVal interface{}, iStateErr Error) {
	// When generating table for query, * need not be converted
	if toConvert == star {
		convertedVal = star
		return
	}
	var err error
	switch kind {
	case reflect.Bool:
		convertedVal, err = strconv.ParseBool(toConvert)
		if err != nil {
			iStateErr = NewError(err, 3009)
			return
		}
	case reflect.Int:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.Atoi(toConvert)
		if err != nil {
			iStateErr = NewError(err, 3011)
			return
		}
	case reflect.Int8:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int8(convertedVal.(int64))
	case reflect.Int16:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int16(convertedVal.(int64))
	case reflect.Int32:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int32(convertedVal.(int64))
	case reflect.Int64:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
	case reflect.Uint:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint(convertedVal.(uint64))
	case reflect.Uint8:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint8(convertedVal.(uint64))
	case reflect.Uint16:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint16(convertedVal.(uint64))
	case reflect.Uint32:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint32(convertedVal.(uint64))
	case reflect.Uint64:
		toConvert, iStateErr = decodeScientificNotation(toConvert)
		if iStateErr != nil {
			return
		}
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
	case reflect.Float32:
		convertedVal, err = strconv.ParseFloat(toConvert, 64)
		if err != nil {
			iStateErr = NewError(err, 3014)
			return
		}
		convertedVal = float32(convertedVal.(float64))
	case reflect.Float64:
		convertedVal, err = strconv.ParseFloat(toConvert, 64)
		if err != nil {
			iStateErr = NewError(err, 3014)
			return
		}
	case reflect.String:
		convertedVal = toConvert
	default:
		iStateErr = NewError(nil, 3017, kind)
		return
	}
	return
}
