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
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/prasanths96/gostack"
	"reflect"
	"strings"
)

// Clean this up...
// Can make it smarter too...
func (iState *iState) parseCmplxAndFetch(stub shim.ChaincodeStubInterface, partIndex, queryToParse string, valKind reflect.Kind, forceFetchDB bool) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	operatorStack := gostack.NewStack()
	resultStack := gostack.NewStack()
	firstArgumentFlag := false
	vocabStartIndex := 0
	queryToParseLen := len(queryToParse)
	for i := 0; i < len(queryToParse); i++ {
		c := queryToParse[i]
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
			if len(queryToEvaluate) > 0 {
				// Get keyword and val
				var keyword, secondPart string
				keyword, secondPart, iStateErr = getKeyWordAndVal(queryToEvaluate)
				if iStateErr != nil {
					return
				}
				// Convert and Encode
				var convertedVal interface{}
				convertedVal, iStateErr = convertToPrimitiveType(secondPart, valKind)
				if iStateErr != nil {
					return
				}
				encodedVal := ""
				encodedVal, iStateErr, _ = encode(convertedVal)
				if iStateErr != nil {
					return
				}
				// Full Index
				indexKey := partIndex + encodedVal + separator

				// Fetch only if it is the first arg / or operator
				curOperator, err := operatorStack.ReadString()
				if err != nil {
					iStateErr = newError(err, 3020)
					return
				}
				fetch := firstArgumentFlag || (curOperator == or)
				switch fetch {
				case true:
					var result map[string]map[string][]byte
					result, iStateErr = iState.selectAndFetch(stub, keyword, indexKey, forceFetchDB)
					if iStateErr != nil {
						return
					}

					// Init with first eval
					if firstArgumentFlag {
						resultStack.Push(result)
					} else {
						// It is not first argument + it is not and operator
						popped, err := resultStack.Pop()
						if err != nil {
							iStateErr = newError(err, 3021)
							return
						}
						curOperator, err := operatorStack.ReadString()
						if err != nil {
							iStateErr = newError(err, 3020)
							return
						}
						// and operator will already be evaluated when fetching (by filtering)
						// hence do nothing
						// if curOperator == or {
						// control enters this block only for or operator
						if curOperator != or {
							panic("Unexpected operator")
						}
						result = orOperation(popped.(map[string]map[string][]byte), result)
						resultStack.Push(result)
						// }
					}

				default:
					// And operator evaluation happens here
					// to avoid unnecessary fetching from db, when
					// it can already be done with existing fetched result.
					popped, err := resultStack.Pop()
					if err != nil {
						iStateErr = newError(err, 3021)
						return
					}
					iStateErr = iState.selectAndFilter(stub, keyword, indexKey, popped.(map[string]map[string][]byte))
					if iStateErr != nil {
						return
					}
					// popped = orOperation(popped.(map[string]map[string][]byte), result)
					resultStack.Push(popped)
				}
			}

			if c == ')' {
				_, err := operatorStack.Pop()
				if err != nil {
					iStateErr = newError(err, 3022)
					return
				}
				// In Multiple level cmplx queries, when exiting a level, there might be result
				// from other cmplx query from same level in stack. Evaluating it.
				if resultStack.Size() > 1 {
					poppedfirst, err := resultStack.Pop()
					if err != nil {
						iStateErr = newError(err, 3021)
						return
					}
					poppedsecond, err := resultStack.Pop()
					if err != nil {
						iStateErr = newError(err, 3021)
						return
					}

					curOperator, err := operatorStack.ReadString()
					if err != nil {
						iStateErr = newError(err, 3020)
						return
					}
					// here, and operation may also be needed
					var result map[string]map[string][]byte
					switch curOperator {
					case or:
						result = orOperation(poppedfirst.(map[string]map[string][]byte), poppedsecond.(map[string]map[string][]byte))
					case and:
						result = andOperation(poppedfirst.(map[string]map[string][]byte), poppedsecond.(map[string]map[string][]byte))
					}
					resultStack.Push(result)
				}
			}
			firstArgumentFlag = false
		}
	}

	fullEvalResult, err := resultStack.Pop()
	if err != nil {
		iStateErr = newError(err, 3021)
		return
	}
	kindecesMap = fullEvalResult.(map[string]map[string][]byte)
	return
}

func (iState *iState) parseCmplxAndEval(partIndex, queryToParse string, valKind reflect.Kind, encKV map[string][]byte) (finalResult bool, iStateErr Error) {
	operatorStack := gostack.NewStack()
	resultStack := gostack.NewStack()
	firstArgumentFlag := false
	vocabStartIndex := 0
	queryToParseLen := len(queryToParse)
	for i := 0; i < len(queryToParse); i++ {
		c := queryToParse[i]
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
			if len(queryToEvaluate) > 0 {
				// Get keyword and val
				var keyword, secondPart string
				keyword, secondPart, iStateErr = getKeyWordAndVal(queryToEvaluate)
				if iStateErr != nil {
					return
				}
				// Convert and Encode
				var convertedVal interface{}
				convertedVal, iStateErr = convertToPrimitiveType(secondPart, valKind)
				if iStateErr != nil {
					return
				}
				encodedVal := ""
				encodedVal, iStateErr, _ = encode(convertedVal)
				if iStateErr != nil {
					return
				}
				// Full Index
				indexKey := partIndex + encodedVal + separator

				var match bool
				match, iStateErr = selectAndEval(keyword, indexKey, encKV)
				if iStateErr != nil {
					return
				}

				// Init with first eval
				if firstArgumentFlag {
					resultStack.Push(match)
				} else {
					popped, err := resultStack.PopBool()
					if err != nil {
						iStateErr = newError(err, 3021)
						return
					}
					curOperator, err := operatorStack.ReadString()
					if err != nil {
						iStateErr = newError(err, 3020)
						return
					}

					match, iStateErr = logicalEval(popped, match, curOperator)
					if iStateErr != nil {
						return
					}
					resultStack.Push(match)
				}
			}

			if c == ')' {
				_, err := operatorStack.Pop()
				if err != nil {
					iStateErr = newError(err, 3022)
					return
				}
				// In Multiple level cmplx queries, when exiting a level, there might be result
				// from other cmplx query from same level in stack. Evaluating it.
				if resultStack.Size() > 1 {
					poppedfirst, err := resultStack.PopBool()
					if err != nil {
						iStateErr = newError(err, 3021)
						return
					}
					poppedsecond, err := resultStack.PopBool()
					if err != nil {
						iStateErr = newError(err, 3021)
						return
					}

					curOperator, err := operatorStack.ReadString()
					if err != nil {
						iStateErr = newError(err, 3020)
						return
					}
					poppedfirst, iStateErr = logicalEval(poppedfirst, poppedsecond, curOperator)
					if iStateErr != nil {
						return
					}
					resultStack.Push(poppedfirst)
				}
			}
			firstArgumentFlag = false
		}
	}

	finalResult, err := resultStack.PopBool()
	if err != nil {
		iStateErr = newError(err, 3021)
		return
	}
	return
}

func (iState *iState) selectAndFetch(stub shim.ChaincodeStubInterface, keyword string, indexKey string, forceFetchDB bool) (kindecesMap map[string]map[string][]byte, iStateErr Error) {

	switch keyword {
	case eq:
		kindecesMap, iStateErr = iState.fetchEq(stub, indexKey, "", forceFetchDB)
		// if iStateErr != nil {return}
	case neq:
		kindecesMap, iStateErr = iState.fetchNeq(stub, indexKey, "", forceFetchDB)
	case gt:
		kindecesMap, iStateErr = iState.fetchGt(stub, indexKey, "", forceFetchDB)
	case lt:
		kindecesMap, iStateErr = iState.fetchLt(stub, indexKey, "", forceFetchDB)
	case gte:
		kindecesMap, iStateErr = iState.fetchGte(stub, indexKey, "", forceFetchDB)
	case lte:
		kindecesMap, iStateErr = iState.fetchLte(stub, indexKey, "", forceFetchDB)
	case seq:
		kindecesMap, iStateErr = iState.fetchSeq(stub, indexKey, "", forceFetchDB)
	case sneq:
		kindecesMap, iStateErr = iState.fetchSneq(stub, indexKey, "", forceFetchDB)
	case sgt:
		kindecesMap, iStateErr = iState.fetchSgt(stub, indexKey, "", forceFetchDB)
	case slt:
		kindecesMap, iStateErr = iState.fetchSlt(stub, indexKey, "", forceFetchDB)
	case sgte:
		kindecesMap, iStateErr = iState.fetchSgte(stub, indexKey, "", forceFetchDB)
	case slte:
		kindecesMap, iStateErr = iState.fetchSlte(stub, indexKey, "", forceFetchDB)
	default:
		iStateErr = newError(nil, 3005, keyword)
		return
	}
	return
}

func (iState *iState) selectAndFilter(stub shim.ChaincodeStubInterface, keyword string, indexKey string, kindecesMap map[string]map[string][]byte) (iStateErr Error) {

	encQKeyVal := make(map[string][]byte)
	encQKeyVal[indexKey] = nil
	switch keyword {
	case eq:
		evalAndFilterEq(stub, encQKeyVal, kindecesMap)
	case neq:
		evalAndFilterNeq(stub, encQKeyVal, kindecesMap)
	case gt:
		evalAndFilterGt(stub, encQKeyVal, kindecesMap)
	case lt:
		evalAndFilterLt(stub, encQKeyVal, kindecesMap)
	case gte:
		evalAndFilterGte(stub, encQKeyVal, kindecesMap)
	case lte:
		evalAndFilterLte(stub, encQKeyVal, kindecesMap)
	case seq:
		evalAndFilterSeq(stub, encQKeyVal, kindecesMap)
	case sneq:
		evalAndFilterSneq(stub, encQKeyVal, kindecesMap)
	case sgt:
		evalAndFilterSgt(stub, encQKeyVal, kindecesMap)
	case slt:
		evalAndFilterSlt(stub, encQKeyVal, kindecesMap)
	case sgte:
		evalAndFilterSgte(stub, encQKeyVal, kindecesMap)
	case slte:
		evalAndFilterSlte(stub, encQKeyVal, kindecesMap)
	default:
		iStateErr = newError(nil, 3005, keyword)
		return
	}
	return
}

func selectAndEval(keyword string, indexKey string, encKV map[string][]byte) (match bool, iStateErr Error) {
	encQKeyVal := make(map[string][]byte)
	encQKeyVal[indexKey] = nil
	switch keyword {
	case eq:
		match, _ = evalEq(indexKey, encKV)
	case neq:
		match, _ = evalNeq(indexKey, encKV)
	case gt:
		match, _ = evalGt(indexKey, encKV)
	case lt:
		match, _ = evalLt(indexKey, encKV)
	case gte:
		match, _ = evalGte(indexKey, encKV)
	case lte:
		match, _ = evalLte(indexKey, encKV)
	case seq:
		match, _ = evalSeq(indexKey, encKV)
	case sneq:
		match, _ = evalSneq(indexKey, encKV)
	case sgt:
		match, _ = evalSgt(indexKey, encKV)
	case slt:
		match, _ = evalSlt(indexKey, encKV)
	case sgte:
		match, _ = evalSgte(indexKey, encKV)
	case slte:
		match, _ = evalSlte(indexKey, encKV)
	default:
		iStateErr = newError(nil, 3005, keyword)
		return
	}
	return
}

func logicalEval(bool1 bool, bool2 bool, operator string) (result bool, iStateErr Error) {
	switch operator {
	case or:
		return bool1 || bool2, nil
	case and:
		return bool1 && bool2, nil
	default:
		iStateErr = newError(nil, 7001, operator)
		return
	}
}
