/*
	Copyright 2020 Prasanth Sundaravelu

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
)

// Cleanup -> stub is in meta data
func (iState *iState) fetchEq(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	start := encodedKey
	end := encodedKey + asciiLast
	iStateErr = iState.getStateByRange(stub, start, end, kindecesMap)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchNeq(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	start1 := partIndex
	end1 := partIndex + removedVals[1]
	iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
	if iStateErr != nil {
		return
	}

	start2 := partIndex + incLastChar(removedVals[1])
	end2 := partIndex + asciiLast
	iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
	if iStateErr != nil {
		return
	}

	return
}

func (iState *iState) fetchGt(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		var positive bool
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestPNum
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + removedVals[1]
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1] + incChar
		end2 := partIndex + asciiLast
		iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchLt(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		var positive bool
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + pNumPrefix
			end1 := partIndex + removedVals[1]
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestNNum
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1]
		iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return

}

func (iState *iState) fetchGte(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		var positive bool
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestPNum
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1]
		end2 := partIndex + asciiLast
		iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchLte(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		var positive bool
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + pNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestNNum
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1] + incChar
		iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchCmplx(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)

	partIndexRemovedSeparator := removeLastSeparator(partIndex)

	fieldName = iState.convertIndexToQueryFieldName(partIndexRemovedSeparator)
	// Get Type of val
	kind, iStateErr := getRightPrimitiveType(fieldName, iState.jsonFieldKindMap, iState.mapKeyKindMap)
	if iStateErr != nil {
		return
	}

	kindecesMap, iStateErr = iState.parseCmplxAndFetch(partIndex, removedVals[1], kind)
	if iStateErr != nil {
		return
	}

	return
}

// TODO Cache startKey, endKey too?
func (iState *iState) getStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string, kindecesMap map[string]map[string][]byte) (iStateErr Error) {

	// // Compact Index
	// cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))

	// cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	// if iStateErr != nil {
	// 	return
	// }
	// for keyRef, hashString := range cIndexV {

	// 	iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, hashString)
	// 	if iStateErr != nil {
	// 		return
	// 	}
	// }

	// Normal Index
	iterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		iStateErr = newError(err, 3006)
		return
	}

	defer iterator.Close()
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = newError(err, 3007)
			return
		}
		indexkey := iteratorResult.GetKey()
		hashBytes := iteratorResult.GetValue()
		keyRef := getKeyFromIndex(indexkey)

		iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, string(hashBytes))
		if iStateErr != nil {
			return
		}
	}

	return
}
func (iState *iState) getStateByRangeWithPagination(stub shim.ChaincodeStubInterface, startKey string, endKey string, qEnv *queryEnv, pageSize int32, bookmark string) (fetchedKVMap map[string][]byte, iStateErr Error) {

	fetchedKVMap = make(map[string][]byte)
	// // Compact Index
	// cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))
	// cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	// if iStateErr != nil {
	// 	return
	// }

	// for keyRef := range cIndexV {
	// 	iStateErr = loadFetchedKV(stub, fetchedKVMap, keyRef, qEnv)
	// 	if iStateErr != nil {
	// 		return
	// 	}
	// }

	// Normal Index
	iterator, _, err := stub.GetStateByRangeWithPagination(startKey, endKey, pageSize, bookmark)
	if err != nil {
		iStateErr = newError(err, 3006)
		return
	}

	defer iterator.Close()
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = newError(err, 3007)
			return
		}
		indexkey := iteratorResult.GetKey()
		keyRef := getKeyFromIndex(indexkey)

		iStateErr = loadFetchedKV(stub, fetchedKVMap, keyRef, qEnv)
		if iStateErr != nil {
			return
		}
	}

	return
}

func loadFetchedKV(stub shim.ChaincodeStubInterface, fetchedKVMap map[string][]byte, keyRef string, qEnv *queryEnv) (iStateErr Error) {
	if _, ok := qEnv.ufetchedKVMap[keyRef]; ok {
		if uValBytes, ok := fetchedKVMap[keyRef]; !ok {
			fetchedKVMap[keyRef] = uValBytes
		}
		return
	}
	// Doesn't fetch if already fetched before
	valBytes, err := stub.GetState(keyRef)
	if err != nil {
		iStateErr = newError(err, 3008)
		return
	}
	if valBytes != nil {
		fetchedKVMap[keyRef] = valBytes
		qEnv.ufetchedKVMap[keyRef] = valBytes
	}
	return
}
func (iState *iState) loadkindecesMap(stub shim.ChaincodeStubInterface, kindecesMap map[string]map[string][]byte, keyRef string, newHash string) (iStateErr Error) {
	// Hash validation for cache consistency
	hashString, iStateErr := iState.getkvHash(keyRef)
	if iStateErr != nil {
		return
	}
	if hashString != newHash {
		iStateErr = iState.removeCache(keyRef)
		if iStateErr != nil {
			return
		}
	}
	indeces, iStateErr := iState.getIndeces(keyRef)
	if iStateErr != nil {
		return
	}
	kindecesMap[keyRef] = indeces

	return
}
