// Copyright 2020 <>. All rights reserved.

package istate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	ez "github.com/prasanths96/hyperledger/easycompositestate"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// ===================================================================================
// ADVANCED QUERY
// ===================================================================================
// Syntax and Keywords
const EQ = "eq"
const NEQ = "neq"
const GT = "gt"
const LT = "lt"
const GTE = "gte"
const LTE = "lte"
const CMPLX = "cmplx"
const OR = "or"
const AND = "and"

//
func (iState *iState) AdvancedQuery(stub shim.ChaincodeStubInterface, nativeStruct interface{}, queryString string) (interface{}, error) {
	// queryString = [{"docType":"eq USERPROFILE_DOCTYPE", "doctor.whatever": "cmplx or(and(gt bla, lt bla),or(eq a, eq b))", "groups":["cmplx and(neq doctor, neq patient)"]}, "or", {"docType":"eq USERPROFILE_DOCTYPE", "groups":["eq patient"]}]
	// Assert that structure is not a pointer
	structRefVal := reflect.ValueOf(nativeStruct)
	if kind := structRefVal.Kind(); kind == reflect.Ptr {
		return nil, errors.New("Value of structure expected, not pointer.")
	}
	//	// Form positionedSearchAttributes to call queryCompositeState
	var parsedQuery []interface{}
	err := json.Unmarshal([]byte(queryString), &parsedQuery)
	if err != nil {
		return nil, err
	}
	querySize := len(parsedQuery)
	// Check if linkers are of same type eg: only or's / only and's
	linker := ""
	if querySize > 2 {
		linker = parsedQuery[1].(string)
	}
	for i := 1; i < querySize; i = i + 2 {
		if parsedQuery[i].(string) != linker {
			return nil, errors.New("Only the logical operators of same kind are allowed between object queries at the moment. Eg: [X,or,Y,or,Z] or [X,and,Y,and,Z]")
		}
	}
	// Initialize with len and cap aswell
	positionedSearchAttributes := make([][]string, querySize, querySize)
	for i, numField := 0, structRefVal.NumField(); i < querySize; i = i + 2 { // Odd indexes are link keywords like - or, and
		positionedSearchAttributes[i] = make([]string, numField)
	}
	for i := 0; i < structRefVal.NumField(); i++ {
		fieldTag := structRefVal.Type().Field(i).Tag.Get("json")
		// fmt.Println(fieldTag)
		for j := 0; j < querySize; j = j + 2 {
			//	if reflect.TypeOf(parsedQuery[j]).String() == "string" {
			//		fmt.Println("This should be ignored since it is a link statement like - or", j)
			//	}
			if v, ok := parsedQuery[j].(map[string]interface{})[fieldTag]; ok {
				if reflect.TypeOf(v).String() != "string" {
					break
				}
				parts := strings.Split(strings.TrimSpace(v.(string)), " ")
				if parts[0] != EQ {
					break
				}
				//fmt.Println("Found ", fieldTag, " at ", j, " value is ", parts[1])
				positionedSearchAttributes[j][i] = parts[1]
				// Delete this field from parsedQuery
				delete(parsedQuery[j].(map[string]interface{}), fieldTag)
			}

		}

	}
	// fmt.Println(positionedSearchAttributes)
	resultInterfaceSlice := make([]interface{}, querySize, querySize)
	for i := 0; i < querySize; i = i + 2 {
		_, resultInterface, err := ez.QueryCompositeState(stub, nativeStruct, positionedSearchAttributes[i])
		if err != nil {
			return nil, err
		}
		resultInterfaceSlice[i] = resultInterface
	}

	//	// Step 2: Rich Query
	results := make([]map[string]interface{}, querySize, querySize) // 0 is for query0, 2 is for query2
	for i := 0; i < querySize; i = i + 2 {
		// Init individual slice
		results[i] = make(map[string]interface{})

		// For each result obtained
		resultInterfaceSliceRefVal := reflect.ValueOf(resultInterfaceSlice[i])
		// thisresultInterface := resultInterfaceSlice[i].([]interface{})
		queryGotField := len(parsedQuery[i].(map[string]interface{})) > 0
		for j, size := 0, resultInterfaceSliceRefVal.Len(); j < size; j++ {
			// for j, size := 0, len(thisresultInterface); j < size; j ++ {
			match := true
			singleResult := resultInterfaceSliceRefVal.Index(j).Interface()
			// singleResult := thisresultInterface[j]
			if queryGotField {
				singleResultMap := make(map[string]interface{})
				mresult, err := json.Marshal(singleResult)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(mresult, &singleResultMap)
				if err != nil {
					return nil, err
				}

				// For each field in query
				for key := range parsedQuery[i].(map[string]interface{}) {
					// parsedQuery[i].(map[string]interface {})[key] = valueToMatch
					fieldMatch := false
					obtainedValue, err := fetchKeyValueFromResultMap(key, singleResultMap)
					if err != nil {
						return nil, err
					}
					var arrayFlag bool
					var valueToMatchString string
					// fmt.Println("DEBUG: reflect.TypeOf(parsedQuery[i])", reflect.TypeOf(parsedQuery[i].(map[string]interface{})[key]).String())
					// fmt.Println("DEBUG: ", parsedQuery[i])
					if reflect.TypeOf(parsedQuery[i].(map[string]interface{})[key]).String() == "[]interface {}" { // its originally []string after unmarshal, it became []interface {}
						arrayFlag = true
						valueToMatchString = parsedQuery[i].(map[string]interface{})[key].([]interface{})[0].(string)
					} else {
						valueToMatchString = parsedQuery[i].(map[string]interface{})[key].(string)
					}
					fieldMatch, err = parseAndEvaluate(obtainedValue, valueToMatchString, arrayFlag)
					if err != nil {
						return nil, err
					}
					// Fields inside query object are default "AND" operation
					// match = match && fieldMatch
					if !fieldMatch {
						match = false
						break
					}
				}
			}

			if match {
				// Appending to map
				results[i][fmt.Sprintf("%v", singleResult)] = singleResult
			}
		}
	}

	// Link - OR = Union, AND = Intersection
	finalResultMap := results[0]
	if querySize > 2 {
		for i := 2; i < querySize; i = i + 2 {
			switch linker {
			case OR:
				unionToFirstMap(finalResultMap, results[i])
			case AND:
				intersectToFirstMap(finalResultMap, results[i])
			}
		}
	}

	finalResultLen := len(finalResultMap)
	finalResultSliceRefVal := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(nativeStruct)), finalResultLen, finalResultLen)
	i := 0
	for key := range finalResultMap {
		finalResultSliceRefVal.Index(i).Set(reflect.ValueOf(finalResultMap[key]))
		i++
	}

	return finalResultSliceRefVal.Interface(), nil

}

func unionToFirstMap(map1 map[string]interface{}, map2 map[string]interface{}) {
	for key := range map2 {
		map1[key] = map2[key]
	}
}

func intersectToFirstMap(map1 map[string]interface{}, map2 map[string]interface{}) {
	for key := range map1 {
		if _, ok := map2[key]; !ok {
			delete(map1, key)
		}
	}
}

func fetchKeyValueFromResultMap(key string, result map[string]interface{}) (interface{}, error) {
	// Split key into parts (for dotted field names)
	keyParts := strings.Split(key, ".")
	var finalValue interface{} = result
	for _, v := range keyParts {
		if val, ok := finalValue.(map[string]interface{})[v]; ok {
			finalValue = val
		} else {
			return nil, errors.New("Invalid field name in query string: " + key)
		}
	}
	// fmt.Println("DEBUG: This is the finalValue: ", finalValue, reflect.TypeOf(finalValue))
	return finalValue, nil
}

func parseAndEvaluate(obtainedValue interface{}, queryToEvaluate string, arrayFlag bool) (bool, error) {
	var err error
	var match bool
	var matchAllFlag bool
	// Trim Space
	queryToEvaluate = strings.TrimSpace(queryToEvaluate)
	firstSpaceIndex := strings.Index(queryToEvaluate, " ")
	if firstSpaceIndex == -1 {
		return false, errors.New(fmt.Sprintf("Syntax error: <Space> not found in %v ", queryToEvaluate))
	}
	keyword := queryToEvaluate[:firstSpaceIndex]
	secondPart := queryToEvaluate[firstSpaceIndex+1:]
	if keyword[0] == '*' { // Equivalent of '*' in rune is 42
		keyword = keyword[1:]
		matchAllFlag = true
	}

	if len(secondPart) == 0 {
		return false, errors.New(fmt.Sprintf("Query syntax error: %v", queryToEvaluate))
	}

	if !arrayFlag {
		switch keyword {
		case EQ:
			convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue).String())
			if err != nil {
				return false, err
			}
			match = convertedValue.Eq(obtainedValue)
		case NEQ:
			convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue).String())
			if err != nil {
				return false, err
			}
			match = convertedValue.Neq(obtainedValue)
		case GT:
			convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue).String())
			if err != nil {
				return false, err
			}
			match = convertedValue.Gt(obtainedValue)
		case LT:
			convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue).String())
			if err != nil {
				return false, err
			}
			match = convertedValue.Lt(obtainedValue)
		case GTE:
			convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue).String())
			if err != nil {
				return false, err
			}
			match = convertedValue.Gte(obtainedValue)
		case LTE:
			convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue).String())
			if err != nil {
				return false, err
			}
			match = convertedValue.Lte(obtainedValue)
		case CMPLX:
			match, err = parseCmplxAndEvaluate(obtainedValue, secondPart, arrayFlag)
			if err != nil {
				return false, err
			}
		}
	} else {
		switch keyword {
		case EQ:
			if reflect.ValueOf(obtainedValue).Len() > 0 {
				convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue.([]interface{})[0]).String())
				if err != nil {
					return false, err
				}
				// fmt.Println("DEBUG: This is before calling Eq: ", obtainedValue, reflect.TypeOf(obtainedValue))
				match, err = arrayMatch(convertedValue.Eq, obtainedValue, matchAllFlag) // obtainedValue is an array, matchAll = false for "contains" operation
			}
		case NEQ:
			if reflect.ValueOf(obtainedValue).Len() > 0 {
				convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue.([]interface{})[0]).String())
				if err != nil {
					return false, err
				}
				match, err = arrayMatch(convertedValue.Neq, obtainedValue, matchAllFlag) // matchAll = true
			} else {
				match = true // If empty array, then NEQ <value> is always true (does not contain)
			}
		case GT:
			if reflect.ValueOf(obtainedValue).Len() > 0 {
				convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue.([]interface{})[0]).String())
				if err != nil {
					return false, err
				}
				match, err = arrayMatch(convertedValue.Gt, obtainedValue, matchAllFlag)
			}
		case LT:
			if reflect.ValueOf(obtainedValue).Len() > 0 {
				convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue.([]interface{})[0]).String())
				if err != nil {
					return false, err
				}
				match, err = arrayMatch(convertedValue.Lt, obtainedValue, matchAllFlag)
			}
		case GTE:
			if reflect.ValueOf(obtainedValue).Len() > 0 {
				convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue.([]interface{})[0]).String())
				if err != nil {
					return false, err
				}
				match, err = arrayMatch(convertedValue.Gte, obtainedValue, matchAllFlag)
			}
		case LTE:
			if reflect.ValueOf(obtainedValue).Len() > 0 {
				convertedValue, err := convertToRightType(secondPart, reflect.TypeOf(obtainedValue.([]interface{})[0]).String())
				if err != nil {
					return false, err
				}
				match, err = arrayMatch(convertedValue.Lte, obtainedValue, matchAllFlag)
			}
		case CMPLX:
			match, err = parseCmplxAndEvaluate(obtainedValue, secondPart, arrayFlag)
		}
		if err != nil {
			return false, err
		}

	}

	return match, nil

}

func arrayMatch(matchFunc func(interface{}) bool, obtainedValue interface{}, matchAll bool) (bool, error) {
	// Checking if obtainedValue is actually a slice / array
	if reflect.TypeOf(obtainedValue).String()[0] != '[' { // Since all slice types start with '[' eg: []string, []int
		return false, errors.New("Array query on a non-array element.")
	}

	// Update this looping and matching algorithm for better performance
	refVal := reflect.ValueOf(obtainedValue)
	var defaultResult bool
	if matchAll {
		defaultResult = true
		for i := 0; i < refVal.Len(); i++ {
			if !matchFunc(refVal.Index(i).Interface()) { // Not a match
				return false, nil
			}
		}
	} else {
		// defaultResult = false // Zero val of bool is already false
		for i := 0; i < refVal.Len(); i++ {
			if matchFunc(refVal.Index(i).Interface()) { // Match
				return true, nil
			}
		}
	}
	return defaultResult, nil
}

func parseCmplxAndEvaluate(obtainedValue interface{}, queryToParse string, arrayFlag bool) (bool, error) {
	var err error
	// traversalMemory := queryToParse
	operatorStack := NewStack()
	resultStack := NewStack()
	firstArgumentFlag := false
	vocabStartIndex := 0
	queryToParseLen := len(queryToParse)

	for i, c := range queryToParse {
		switch c { // c is a rune / int32 here
		case '(': // ' ' for rune representation
			// Logical operator always comes before ( eg: and(or(lt 10, gt 5), eq 7)
			operatorStack.Push(strings.TrimSpace(queryToParse[vocabStartIndex:i]))
			//traversalMemory = queryToParse[i + 1 :]
			firstArgumentFlag = true
			vocabStartIndex = i + 1
		case ',', ')':
			queryToEvaluate := strings.TrimSpace(queryToParse[vocabStartIndex:i])
			for {
				if (i + 1) < queryToParseLen {
					if queryToParse[i+1] == ',' || queryToParse[i+1] == ' ' {
						i++
					} else {
						break
					}
				} else {
					break
				}
			}
			vocabStartIndex = i + 1
			//traversalMemory = queryToParse[i + 1 :]
			evalResult, err := parseAndEvaluate(obtainedValue, queryToEvaluate, arrayFlag)
			if err != nil {
				return false, err
			}
			// Init with first eval
			if firstArgumentFlag {
				resultStack.Push(evalResult)
			} else {
				popped, err := resultStack.PopBool()
				if err != nil {
					return false, err
				}
				curOperator, err := operatorStack.ReadString()
				if err != nil {
					return false, err
				}
				popped, err = logicalEval(popped, evalResult, curOperator)
				if err != nil {
					return false, err
				}
				resultStack.Push(popped)
			}
			if c == ')' {
				_, err = operatorStack.Pop()
				if err != nil {
					return false, err
				}
				if resultStack.Size() > 1 {
					poppedfirst, err := resultStack.PopBool()
					poppedsecond, err := resultStack.PopBool()
					if err != nil {
						return false, err
					}
					curOperator, err := operatorStack.ReadString()
					if err != nil {
						return false, err
					}
					finalResult, err := logicalEval(poppedfirst, poppedsecond, curOperator)
					if err != nil {
						return false, err
					}
					resultStack.Push(finalResult)
				}
			}
			firstArgumentFlag = false
			// default: continue
		}
	}

	fullEvalResult, err := resultStack.ReadBool()
	if err != nil {
		return false, err
	}

	return fullEvalResult, nil
}

func logicalEval(bool1 bool, bool2 bool, operator string) (bool, error) {
	switch operator {
	case OR:
		return bool1 || bool2, nil
	case AND:
		return bool1 && bool2, nil
	default:
		return false, errors.New(fmt.Sprintf("Unsupported logical operator: %v", operator))
	}
}

func convertToRightType(valueToConvert string, typeToConvert string) (Comparable, error) {
	var err error
	var convertedValue Comparable

	switch typeToConvert {
	case "int", "[]int":
		var value int
		value, err = strconv.Atoi(valueToConvert)
		convertedValue = Int(value)
	case "int64", "[]int64":
		var value int64
		value, err = strconv.ParseInt(valueToConvert, 10, 64)
		convertedValue = Int64(value)
	case "string", "[]string":
		convertedValue = String(valueToConvert)
	case "float32", "[]float32":
		var value float64
		value, err = strconv.ParseFloat(valueToConvert, 64)
		convertedValue = Float32(value)
	case "bool", "[]bool":
		var value bool
		value, err = strconv.ParseBool(valueToConvert)
		convertedValue = Bool(value)

	default:
		return nil, errors.New(fmt.Sprintf("Unknown / uncomparable type found in query: %v", typeToConvert))
	}
	if err != nil {
		return nil, err
	}

	return convertedValue, nil
}

// ===================================================================================
// Custom Types for implementing comparison functions
// ===================================================================================

type Comparable interface {
	Gt(interface{}) bool
	Lt(interface{}) bool
	Gte(interface{}) bool
	Lte(interface{}) bool
	Eq(interface{}) bool
	Neq(interface{}) bool
}

type Int int
type Int64 int64
type Float32 float32
type String string
type Bool bool

// Value in parameter must be primitive type, not custom type
func (val Int) Gt(val2 interface{}) bool {
	return val > Int(val2.(int))
}
func (val Int64) Gt(val2 interface{}) bool {
	return val > Int64(val2.(int64))
}
func (val Float32) Gt(val2 interface{}) bool {
	return val > Float32(val2.(float32))
}
func (val String) Gt(val2 interface{}) bool {
	return val > String(val2.(string))
}

// Boolean does not support Gt / Lt / Gte / LTe
func (val Bool) Gt(val2 interface{}) bool {
	return false
}

func (val Int) Lt(val2 interface{}) bool {
	return val < Int(val2.(int))
}
func (val Int64) Lt(val2 interface{}) bool {
	return val < Int64(val2.(int64))
}
func (val Float32) Lt(val2 interface{}) bool {
	return val < Float32(val2.(float32))
}
func (val String) Lt(val2 interface{}) bool {
	return val < String(val2.(string))
}
func (val Bool) Lt(val2 interface{}) bool {
	return false
}

func (val Int) Gte(val2 interface{}) bool {
	return val >= Int(val2.(int))
}
func (val Int64) Gte(val2 interface{}) bool {
	return val >= Int64(val2.(int64))
}
func (val Float32) Gte(val2 interface{}) bool {
	return val >= Float32(val2.(float32))
}
func (val String) Gte(val2 interface{}) bool {
	return val >= String(val2.(string))
}
func (val Bool) Gte(val2 interface{}) bool {
	return false
}

func (val Int) Lte(val2 interface{}) bool {
	return val <= Int(val2.(int))
}
func (val Int64) Lte(val2 interface{}) bool {
	return val <= Int64(val2.(int64))
}
func (val Float32) Lte(val2 interface{}) bool {
	return val <= Float32(val2.(float32))
}
func (val String) Lte(val2 interface{}) bool {
	return val <= String(val2.(string))
}
func (val Bool) Lte(val2 interface{}) bool {
	return false
}

func (val Int) Eq(val2 interface{}) bool {
	return val == Int(val2.(int))
}
func (val Int64) Eq(val2 interface{}) bool {
	return val == Int64(val2.(int64))
}
func (val Float32) Eq(val2 interface{}) bool {
	return val == Float32(val2.(float32))
}
func (val String) Eq(val2 interface{}) bool {
	return val == String(val2.(string))
}
func (val Bool) Eq(val2 interface{}) bool {
	return val == Bool(val2.(bool))
}

func (val Int) Neq(val2 interface{}) bool {
	return val != Int(val2.(int))
}
func (val Int64) Neq(val2 interface{}) bool {
	return val != Int64(val2.(int64))
}
func (val Float32) Neq(val2 interface{}) bool {
	return val != Float32(val2.(float32))
}
func (val String) Neq(val2 interface{}) bool {
	return val != String(val2.(string))
}
func (val Bool) Neq(val2 interface{}) bool {
	return val != Bool(val2.(bool))
}

// ===================================================================================
// Synchronized Stack accepting dynamic type
// ===================================================================================

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []interface{}
}

func NewStack() *stack {
	return &stack{sync.Mutex{}, nil}
}

func (s *stack) Push(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) Read() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	return res, nil
}

func (s *stack) Size() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.s)
}

func (s *stack) PopBool() (bool, error) {
	result, err := s.Pop()
	if err != nil {
		return false, err
	}
	if thisType := reflect.TypeOf(result).String(); thisType != "bool" {
		return false, errors.New(fmt.Sprintf("Cannot pop %v value as bool", thisType))
	}
	return result.(bool), nil
}

func (s *stack) PopString() (string, error) {
	result, err := s.Pop()
	if err != nil {
		return "", err
	}
	if thisType := reflect.TypeOf(result).String(); thisType != "string" {
		return "", errors.New(fmt.Sprintf("Cannot pop %v value as string", thisType))
	}
	return result.(string), nil
}

func (s *stack) PopInt() (int, error) {
	result, err := s.Pop()
	if err != nil {
		return 0, err
	}
	if thisType := reflect.TypeOf(result).String(); thisType != "int" {
		return 0, errors.New(fmt.Sprintf("Cannot pop %v value as int", thisType))
	}
	return result.(int), nil
}

func (s *stack) ReadBool() (bool, error) {
	result, err := s.Read()
	if err != nil {
		return false, err
	}
	if thisType := reflect.TypeOf(result).String(); thisType != "bool" {
		return false, errors.New(fmt.Sprintf("Cannot Read %v value as bool", thisType))
	}
	return result.(bool), nil
}

func (s *stack) ReadString() (string, error) {
	result, err := s.Read()
	if err != nil {
		return "", err
	}
	if thisType := reflect.TypeOf(result).String(); thisType != "string" {
		return "", errors.New(fmt.Sprintf("Cannot Read %v value as string", thisType))
	}
	return result.(string), nil
}

func (s *stack) ReadInt() (int, error) {
	result, err := s.Read()
	if err != nil {
		return 0, err
	}
	if thisType := reflect.TypeOf(result).String(); thisType != "int" {
		return 0, errors.New(fmt.Sprintf("Cannot Read %v value as int", thisType))
	}
	return result.(int), nil
}
