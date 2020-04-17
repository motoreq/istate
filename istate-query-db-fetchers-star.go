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

func (iState *iState) fetchSeq(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {

	// Get Blacklisted Keys:
	blackList, iStateErr := fetchSeqBlackList(stub, encodedKey)
	if iStateErr != nil {
		return
	}

	// Fetch Non-Blacklisted
	kindecesMap = make(map[string]map[string][]byte)
	start := encodedKey
	end := encodedKey + asciiLast
	iStateErr = iState.loadNonBlackListStateByRange(stub, start, end, kindecesMap, blackList)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSneq(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys:
	blackList, iStateErr := fetchSneqBlackList(stub, encodedKey)
	if iStateErr != nil {
		return
	}

	// Fetch Non-Blacklisted
	kindecesMap = make(map[string]map[string][]byte)

	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	start1 := partIndex
	end1 := partIndex + removedVals[1]
	iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
	if iStateErr != nil {
		return
	}

	start2 := partIndex + incLastChar(removedVals[1])
	end2 := partIndex + asciiLast
	iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSgt(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {

	// Get Blacklisted Keys:
	blackList, iStateErr := fetchSgtBlackList(stub, encodedKey)
	if iStateErr != nil {
		return
	}

	// Fetch Non-Blacklisted
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestPNum
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + removedVals[1]
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1] + incChar
		end2 := partIndex + asciiLast
		iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchSlt(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys:
	blackList, iStateErr := fetchSltBlackList(stub, encodedKey)
	if iStateErr != nil {
		return
	}

	// Fetch Non-Blacklisted
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + pNumPrefix
			end1 := partIndex + removedVals[1]
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestNNum
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1]
		iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchSgte(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys:
	blackList, iStateErr := fetchSgteBlackList(stub, encodedKey)
	if iStateErr != nil {
		return
	}

	// Fetch Non-Blacklisted
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestPNum
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1]
		end2 := partIndex + asciiLast
		iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
		if iStateErr != nil {
			return
		}
	}
	return
}

func (iState *iState) fetchSlte(stub shim.ChaincodeStubInterface, encodedKey string, fieldName string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys:
	blackList, iStateErr := fetchSlteBlackList(stub, encodedKey)
	if iStateErr != nil {
		return
	}

	// Fetch Non-Blacklisted
	kindecesMap = make(map[string]map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + pNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			iStateErr = iState.loadNonBlackListStateByRange(stub, start2, end2, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestNNum
			iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1] + incChar
		iStateErr = iState.loadNonBlackListStateByRange(stub, start1, end1, kindecesMap, blackList)
		if iStateErr != nil {
			return
		}
	}
	return
}

func fetchSeqBlackList(stub shim.ChaincodeStubInterface, encodedKey string) (blackList map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys: (neq)
	blackList = make(map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	start1 := partIndex
	end1 := partIndex + removedVals[1]
	iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
	if iStateErr != nil {
		return
	}

	start2 := partIndex + incLastChar(removedVals[1])
	end2 := partIndex + asciiLast
	iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
	if iStateErr != nil {
		return
	}
	return
}

func fetchSneqBlackList(stub shim.ChaincodeStubInterface, encodedKey string) (blackList map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys: (eq)
	blackList = make(map[string][]byte)
	start := encodedKey
	end := encodedKey + asciiLast
	iStateErr = loadKeyRefByRange(stub, start, end, blackList)
	if iStateErr != nil {
		return
	}
	return
}

func fetchSgtBlackList(stub shim.ChaincodeStubInterface, encodedKey string) (blackList map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys: (lte)
	blackList = make(map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + pNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestNNum
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1] + incChar
		iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
		if iStateErr != nil {
			return
		}
	}
	return
}
func fetchSltBlackList(stub shim.ChaincodeStubInterface, encodedKey string) (blackList map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys: (gte)
	blackList = make(map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestPNum
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1]
		end2 := partIndex + asciiLast
		iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
		if iStateErr != nil {
			return
		}
	}

	return
}
func fetchSgteBlackList(stub shim.ChaincodeStubInterface, encodedKey string) (blackList map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys: (lt)
	blackList = make(map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + pNumPrefix
			end1 := partIndex + removedVals[1]
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestNNum
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1]
		iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
		if iStateErr != nil {
			return
		}
	}
	return
}
func fetchSlteBlackList(stub shim.ChaincodeStubInterface, encodedKey string) (blackList map[string][]byte, iStateErr Error) {
	// Get Blacklisted Keys: (gt)
	blackList = make(map[string][]byte)
	partIndex, removedVals := removeNValsFromIndex(encodedKey, 2)
	switch isNum(removedVals[1]) {
	case true:
		//+1   //-1   //0
		positive := false
		positive, iStateErr = isPositive(removedVals[1])
		if iStateErr != nil {
			return
		}
		switch positive {
		case true:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestPNum
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + removedVals[1]
			iStateErr = loadKeyRefByRange(stub, start1, end1, blackList)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1] + incChar
		end2 := partIndex + asciiLast
		iStateErr = loadKeyRefByRange(stub, start2, end2, blackList)
		if iStateErr != nil {
			return
		}
	}

	return
}

//
func loadKeyRefByRange(stub shim.ChaincodeStubInterface, startKey, endKey string, fetchedKVMap map[string][]byte) (iStateErr Error) {
	iterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		iStateErr = NewError(err, 3006)
		return
	}
	defer iterator.Close()
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = NewError(err, 4006)
			return
		}
		key := iteratorResult.GetKey()
		keyRef := getKeyFromIndex(key)
		fetchedKVMap[keyRef] = nil
	}
	return
}

func (iState *iState) loadNonBlackListStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string, kindecesMap map[string]map[string][]byte, blackList map[string][]byte) (iStateErr Error) {

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
		indexkey := iteratorResult.GetKey()
		hashBytes := iteratorResult.GetValue()
		keyRef := getKeyFromIndex(indexkey)

		if _, ok := blackList[keyRef]; ok {
			continue
		}
		iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, string(hashBytes))
		if iStateErr != nil {
			return
		}
	}

	return
}
