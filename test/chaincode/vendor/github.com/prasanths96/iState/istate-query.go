// Copyright 2020 <>. All rights reserved.

package istate

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"strconv"
	"strings"
)

// Syntax and Keywords
const (
	eq    = "eq"
	neq   = "neq"
	gt    = "gt"
	lt    = "lt"
	gte   = "gte"
	lte   = "lte"
	cmplx = "cmplx"
	or    = "or"
	and   = "and"

	seq    = "*eq"
	sneq   = "*neq"
	sgt    = "*gt"
	slt    = "*lt"
	sgte   = "*gte"
	slte   = "*lte"
	scmplx = "*cmplx"

	splitDot = "."

	asciiLast = "~"
)

// Query function is used to
// "," separated objects are considered "or" always
// queryString = [{"docType":"eq USERPROFILE_DOCTYPE", "doctor.whatever": "cmplx or(and(gt bla, lt bla),or(eq a, eq b))", "groups":["cmplx and(neq doctor, neq patient)"]}, {"docType":"eq USERPROFILE_DOCTYPE", "groups":["eq patient"]}]
func (iState *iState) Query(stub shim.ChaincodeStubInterface, queryString string) (finalResult interface{}, iStateErr Error) {
	iStateLogger.Infof("Inside Query")
	defer iStateLogger.Infof("Exiting Query")

	var uQuery []map[string]interface{}
	err := json.Unmarshal([]byte(queryString), &uQuery)
	if err != nil {
		iStateErr = NewError(err, 3002)
		return
	}

	fmt.Println(uQuery)
	results := make([]map[string][]byte, len(uQuery), len(uQuery))
	for i := 0; i < len(uQuery); i++ {
		results[i], iStateErr = iState.parseAndEvalSingle(stub, uQuery[i])
		if iStateErr != nil {
			return
		}
	}

	// Or operation over results
	combinedResults := orOperation(results...)
	fmt.Println("FINAL COMBINE RESULTS MAP: ", combinedResults)
	finalResult = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(iState.structRef)), len(combinedResults), len(combinedResults))
	i := 0
	for _, val := range combinedResults {
		singleElem := reflect.New(reflect.TypeOf(iState.structRef)).Interface()
		err := json.Unmarshal(val, &singleElem)
		if err != nil {
			iStateErr = NewError(err, 3002)
			return
		}
		finalResult.(reflect.Value).Index(i).Set(reflect.ValueOf(singleElem).Elem())
		i++
	}

	finalResult = finalResult.(reflect.Value).Interface()
	return
}

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

func (iState *iState) parseAndEvalSingle(stub shim.ChaincodeStubInterface, uQuery map[string]interface{}) (result map[string][]byte, iStateErr Error) {
	// I dont have to make [], it will be made at addKeyWithoutOverLap
	var eqQuery []map[string]interface{}
	neqQuery := make(map[string]interface{})
	gtQuery := make(map[string]interface{})
	ltQuery := make(map[string]interface{})
	gteQuery := make(map[string]interface{})
	lteQuery := make(map[string]interface{})
	cmplxQuery := make(map[string]interface{})

	seqQuery := make(map[string]interface{})
	sneqQuery := make(map[string]interface{})
	sgtQuery := make(map[string]interface{})
	sltQuery := make(map[string]interface{})
	sgteQuery := make(map[string]interface{})
	slteQuery := make(map[string]interface{})
	scmplxQuery := make(map[string]interface{})

	for index, val := range uQuery {
		if val, ok := val.(string); !ok {
			iStateErr = NewError(nil, 3003, reflect.TypeOf(val))
			return
		}

		// Trim Space
		queryToEvaluate := strings.TrimSpace(val.(string))
		firstSpaceIndex := strings.Index(queryToEvaluate, " ")
		if firstSpaceIndex == -1 {
			iStateErr = NewError(nil, 3004, queryToEvaluate)
			return
		}
		keyword := queryToEvaluate[:firstSpaceIndex]
		secondPart := queryToEvaluate[firstSpaceIndex+1:]
		fmt.Println("SECOND PART: ", secondPart)
		// Convert the string value to appropriate type
		var convertedVal interface{}
		convertedVal, iStateErr = convertToRightType(index, secondPart, iState.jsonFieldKindMap, iState.mapKeyKindMap)
		if iStateErr != nil {
			return
		}

		// Dot to Map
		newIndex := ""
		var newVal interface{}
		splitFieldName := strings.Split(index, splitDot)
		switch len(splitFieldName) > 1 {
		case true:
			newIndex = splitFieldName[0]
			newVal = dotsToActualDepth(splitFieldName[1:], convertedVal)
		default:
			newIndex = index
			newVal = convertedVal
		}

		switch keyword {
		case eq:
			// eqQuery[newIndex] = newVal
			eqQuery = addKeyWithoutOverLap(eqQuery, newIndex, newVal)
		case neq:
			neqQuery[newIndex] = newVal
		case gt:
			gtQuery[newIndex] = newVal
		case lt:
			ltQuery[newIndex] = newVal
		case gte:
			gteQuery[newIndex] = newVal
		case lte:
			lteQuery[newIndex] = newVal
		case cmplx:
			cmplxQuery[newIndex] = newVal
		case seq:
			seqQuery[newIndex] = newVal
		case sneq:
			sneqQuery[newIndex] = newVal
		case sgt:
			sgtQuery[newIndex] = newVal
		case slt:
			sltQuery[newIndex] = newVal
		case sgte:
			sgteQuery[newIndex] = newVal
		case slte:
			slteQuery[newIndex] = newVal
		case scmplx:
			scmplxQuery[newIndex] = newVal
		default:
			iStateErr = NewError(nil, 3005, keyword)
			return
		}

	}

	var resultEq map[string][]byte

	resultEq, iStateErr = iState.evaluateEq(stub, eqQuery)
	iState.evaluateNeq(stub, neqQuery)
	iState.evaluateGt(stub, gtQuery)
	iState.evaluateLt(stub, ltQuery)
	iState.evaluateGte(stub, gteQuery)
	iState.evaluateLte(stub, lteQuery)
	iState.evaluateCmplx(stub, cmplxQuery)

	iState.evaluateSeq(stub, seqQuery)
	iState.evaluateSneq(stub, sneqQuery)
	iState.evaluateSgt(stub, sgtQuery)
	iState.evaluateSlt(stub, sltQuery)
	iState.evaluateSgte(stub, sgteQuery)
	iState.evaluateSlte(stub, slteQuery)
	iState.evaluateScmplx(stub, scmplxQuery)

	// and operation between fields

	result = resultEq
	fmt.Println("RESULT EQ: ", resultEq)
	return
}

func (iState *iState) evaluateEq(stub shim.ChaincodeStubInterface, query []map[string]interface{}) (result map[string][]byte, iStateErr Error) {
	fmt.Println("EQ QUERY:", query)
	keyref := ""
	encodedKeyValue := make(map[string][]byte)
	for i := 0; i < len(query); i++ {
		var tempEncodedKeyValue map[string][]byte
		tempEncodedKeyValue, iStateErr = iState.encodeState(query[i], keyref, true)
		if iStateErr != nil {
			return
		}
		for index, val := range tempEncodedKeyValue {
			encodedKeyValue[index] = val
		}
	}

	fmt.Println("ENCODED BEFORE * REMOVE:", encodedKeyValue)

	// Remove * results from encodedKeys
	removeStarFromKeys(encodedKeyValue)

	keyResults := make([]map[string][]byte, len(encodedKeyValue), len(encodedKeyValue))
	fmt.Println("ENCODED:", encodedKeyValue)
	i := 0
	for index := range encodedKeyValue {
		start := index
		end := index + asciiLast
		// Each key
		keyResults[i], iStateErr = getStateByRange(stub, start, end)
		if iStateErr != nil {
			return
		}
		i++
	}

	fmt.Println("EVALUATE EQ BEFORE AND OPERATION: ", keyResults)
	// and operation between results of each keys
	result = andOperation(keyResults...)
	fmt.Println("EVALUATE EQ: ", result)
	return
}

func (iState *iState) evaluateNeq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("NEQ QUERY:", query)

}

func (iState *iState) evaluateGt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("GT QUERY:", query)

}

func (iState *iState) evaluateLt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("LT QUERY:", query)

}

func (iState *iState) evaluateGte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("GTE QUERY:", query)

}

func (iState *iState) evaluateLte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("LTE QUERY:", query)

}

func (iState *iState) evaluateCmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("CMPLX QUERY:", query)

}

func (iState *iState) evaluateSeq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*EQ QUERY:", query)

}

func (iState *iState) evaluateSneq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*NEQ QUERY:", query)

}

func (iState *iState) evaluateSgt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*GT QUERY:", query)

}

func (iState *iState) evaluateSlt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*LT QUERY:", query)

}

func (iState *iState) evaluateSgte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*GTE QUERY:", query)

}

func (iState *iState) evaluateSlte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*LTE QUERY:", query)

}

func (iState *iState) evaluateScmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	fmt.Println("*CMPLX QUERY:", query)

}

func getStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string) (result map[string][]byte, iStateErr Error) {
	result = make(map[string][]byte)
	iterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		iStateErr = NewError(err, 3006)
		return
	}
	defer iterator.Close()
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = NewError(err, 3007)
			return
		}
		keyRef := string(iteratorResult.GetValue())
		// If key already present in result, then can avoid re-getstate()
		if _, ok := result[keyRef]; ok {
			continue
		}

		valBytes, err := stub.GetState(keyRef)
		if err != nil {
			iStateErr = NewError(err, 3008)
			return
		}
		result[keyRef] = valBytes
	}
	return
}

func andOperation(keyValuePairs ...map[string][]byte) (result map[string][]byte) {
	switch len(keyValuePairs) {
	case 0:
		return
	case 1:
		result = keyValuePairs[0]
		return
	}

	result = keyValuePairs[0]
	for index := range result {
		neFlag := false
		for i := 1; i < len(keyValuePairs); i++ {
			if _, ok := keyValuePairs[i][index]; !ok {
				neFlag = true
				break
			}
		}
		if neFlag {
			delete(result, index)
		}
	}
	return
}

func orOperation(keyValuePairs ...map[string][]byte) (result map[string][]byte) {
	switch len(keyValuePairs) {
	case 0:
		return
	case 1:
		result = keyValuePairs[0]
		return
	}

	result = keyValuePairs[0]
	for i := 1; i < len(keyValuePairs); i++ {
		for index, val := range keyValuePairs[i] {
			if _, ok := result[index]; !ok {
				result[index] = val
			}
		}
	}

	return
}

func convertToRightType(fieldName string, toConvert string, jsonFieldKindMap map[string]reflect.Kind, mapKeyKindMap map[string]reflect.Kind) (convertedVal interface{}, iStateErr Error) {

	splitFieldName := strings.Split(fieldName, splitDot)
	if len(splitFieldName) == 0 {
		iStateErr = NewError(nil, 3010)
		return
	}
	// curField := fieldName
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
		fmt.Println("Atoi: Trying to Convert:", toConvert)
		convertedVal, err = strconv.Atoi(toConvert)
		if err != nil {
			iStateErr = NewError(err, 3011)
			return
		}
	case reflect.Int8:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int8(convertedVal.(int64))
	case reflect.Int16:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int16(convertedVal.(int64))
	case reflect.Int32:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int32(convertedVal.(int64))
	case reflect.Int64:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
	case reflect.Uint:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint(convertedVal.(uint64))
	case reflect.Uint8:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint8(convertedVal.(uint64))
	case reflect.Uint16:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint16(convertedVal.(uint64))
	case reflect.Uint32:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint32(convertedVal.(uint64))
	case reflect.Uint64:
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

func removeKeysWithStar(keyValue map[string][]byte) {
	for index := range keyValue {
		if strings.Contains(index, star) {
			delete(keyValue, index)
		}
	}
}

func removeStarFromKeys(keyValue map[string][]byte) {
	for index := range keyValue {
		// Replace is used as ReplaceAll isnt available in go version used in fabric image
		newIndex := strings.Replace(index, star, "", len(index))
		if newIndex != index {
			keyValue[newIndex] = keyValue[index]
			delete(keyValue, index)
		}

	}
}

func addKeyWithoutOverLap(query []map[string]interface{}, index string, value interface{}) (newQuery []map[string]interface{}) {
	newQuery = query
	successFlag := false
	for i := 0; i < len(newQuery); i++ {
		if _, ok := newQuery[i][index]; !ok {
			newQuery[i][index] = value
			successFlag = true
			break
		}
	}
	if !successFlag {
		tempMap := make(map[string]interface{})
		tempMap[index] = value
		newQuery = append(newQuery, tempMap)
	}
	return
}
