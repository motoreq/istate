//

package istate

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

func (iState *iState) fetchEq(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	start := encodedKey
	end := encodedKey + asciiLast
	fmt.Println("Start:", []byte(start))
	fmt.Println("End  :", []byte(end))
	fmt.Println("Start:", start)
	fmt.Println("End  :", end)
	iStateErr = iState.getStateByRange(stub, start, end, kindecesMap)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchNeq(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
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

func (iState *iState) fetchGt(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
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
			fmt.Println("num Start: ", start1)
			fmt.Println("Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + removedVals[1]
			fmt.Println("neg num Start: ", start1)
			fmt.Println("neg Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			fmt.Println("pos num Start: ", start2)
			fmt.Println("pos Num End: ", end2)
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1] + incChar
		end2 := partIndex + asciiLast
		fmt.Println("Not num Start: ", start2)
		fmt.Println("Not Num End: ", end2)
		iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchLt(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
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
			fmt.Println("num Start: ", start1)
			fmt.Println("Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			fmt.Println("neg num Start: ", start2)
			fmt.Println("neg Num End: ", end2)
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + incLastChar(removedVals[1])
			end1 := partIndex + biggestNNum
			fmt.Println("neg num Start: ", start1)
			fmt.Println("neg Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

		}

	default:
		start1 := partIndex
		end1 := partIndex + removedVals[1]
		fmt.Println("Start: ", []byte(start1))
		fmt.Println("End  : ", []byte(end1))
		fmt.Println("Start: ", start1)
		fmt.Println("End  : ", end1)
		iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return

}

func (iState *iState) fetchGte(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
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
			fmt.Println("num Start: ", start1)
			fmt.Println("Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + nNumPrefix
			end1 := partIndex + incLastChar(removedVals[1])
			fmt.Println("neg num Start: ", start1)
			fmt.Println("neg Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + pNumPrefix
			end2 := partIndex + pNumPrefix + asciiLast
			fmt.Println("pos num Start: ", start2)
			fmt.Println("pos Num End: ", end2)
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		}

	default:
		start2 := partIndex + removedVals[1]
		end2 := partIndex + asciiLast
		fmt.Println("Not num Start: ", start2)
		fmt.Println("Not Num End: ", end2)
		iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
		if iStateErr != nil {
			return
		}
	}

	return
}

func (iState *iState) fetchLte(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
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
			fmt.Println("num Start: ", start1)
			fmt.Println("Num End: ", end1)
			iStateErr = iState.getStateByRange(stub, start1, end1, kindecesMap)
			if iStateErr != nil {
				return
			}

			start2 := partIndex + nNumPrefix
			end2 := partIndex + nNumPrefix + asciiLast
			fmt.Println("neg num Start: ", start2)
			fmt.Println("neg Num End: ", end2)
			iStateErr = iState.getStateByRange(stub, start2, end2, kindecesMap)
			if iStateErr != nil {
				return
			}
		default:
			start1 := partIndex + removedVals[1]
			end1 := partIndex + biggestNNum
			fmt.Println("neg num Start: ", start1)
			fmt.Println("neg Num End: ", end1)
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

func (iState *iState) fetchCmplx(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	iStateErr = iState.getStateByRange(stub, start, end, kindecesMap)
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
		keyRef := getKeyFromIndex(indexkey)

		iStateErr = loadFetchedKV(stub, fetchedKVMap, keyRef, qEnv)
		if iStateErr != nil {
			return
		}
	}

	return
}
func (iState *iState) getStateByRangeSeq(stub shim.ChaincodeStubInterface, startKey string, endKey string, matchIndex string, kindecesMap map[string]map[string][]byte) (iStateErr Error) {

	// // Check this!!
	// // Compact Index
	// cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))
	// fmt.Println("Generated C Index Key:", cIndexKey)
	// cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	// if iStateErr != nil {
	// 	return
	// }
	// for keyRef, hashString := range cIndexV {
	// 	fmt.Println("Compact Index: ", keyRef)
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

	matchedKeyRefs := make(map[string][]byte)
	matchedKeyAndIndex := make(map[string]string)
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = NewError(err, 3007)
			return
		}
		indexkey := iteratorResult.GetKey()

		hashBytes := iteratorResult.GetValue()
		keyRef := getKeyFromIndex(indexkey)

		if _, ok := matchedKeyRefs[keyRef]; ok {
			// If key already present in matched list, it means, different values are present for that key, so match already failed
			return
		}

		matchedKeyRefs[keyRef] = hashBytes
		partIndex, _ := removeNValsFromIndex(indexkey, 1)
		matchedKeyAndIndex[keyRef] = partIndex

	}

	for keyRef, hashBytes := range matchedKeyRefs {
		// If index does not match the index with val, it means it already failed, so dont load that key
		fmt.Println("This is what we're checking", matchedKeyAndIndex[keyRef], matchIndex)
		if matchedKeyAndIndex[keyRef] != matchIndex {
			continue
		}
		iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, string(hashBytes))
		if iStateErr != nil {
			return
		}
	}

	return
}
func (iState *iState) getStateByRangeSneq(stub shim.ChaincodeStubInterface, startKey string, endKey string, matchIndex string, kindecesMap map[string]map[string][]byte) (iStateErr Error) {

	// // Check this!!
	// // Compact Index
	// cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))
	// fmt.Println("Generated C Index Key:", cIndexKey)
	// cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	// if iStateErr != nil {
	// 	return
	// }
	// for keyRef, hashString := range cIndexV {
	// 	fmt.Println("Compact Index: ", keyRef)
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

	matchedKeyRefs := make(map[string][]byte)
	matchedKeyAndIndex := make(map[string]string)
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = NewError(err, 3007)
			return
		}
		indexkey := iteratorResult.GetKey()

		hashBytes := iteratorResult.GetValue()
		keyRef := getKeyFromIndex(indexkey)

		matchedKeyRefs[keyRef] = hashBytes
		matchedKeyAndIndex[keyRef] = indexkey

	}

	for keyRef, hashBytes := range matchedKeyRefs {
		// If index match the index with val, it means it already failed, so dont load that key
		fmt.Println("This is what we're checking", matchedKeyAndIndex[keyRef], matchIndex)
		if strings.Contains(matchedKeyAndIndex[keyRef], matchIndex) {
			continue
		}
		iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, string(hashBytes))
		if iStateErr != nil {
			return
		}
	}

	return
}
func (iState *iState) getStateByRangeSgt(stub shim.ChaincodeStubInterface, startKey string, endKey string, matchIndex string, kindecesMap map[string]map[string][]byte) (iStateErr Error) {

	// // Check this!!
	// // Compact Index
	// cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))
	// fmt.Println("Generated C Index Key:", cIndexKey)
	// cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	// if iStateErr != nil {
	// 	return
	// }
	// for keyRef, hashString := range cIndexV {
	// 	fmt.Println("Compact Index: ", keyRef)
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

	matchedKeyRefs := make(map[string][]byte)
	matchedKeyAndIndex := make(map[string]string)
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = NewError(err, 3007)
			return
		}
		indexkey := iteratorResult.GetKey()

		hashBytes := iteratorResult.GetValue()
		keyRef := getKeyFromIndex(indexkey)

		matchedKeyRefs[keyRef] = hashBytes
		matchedKeyAndIndex[keyRef] = indexkey

	}

	for keyRef, hashBytes := range matchedKeyRefs {
		fmt.Println("This is what we're checking", matchedKeyAndIndex[keyRef], matchIndex)
		// If index match the index with val, it means it already failed, so dont load that key
		if strings.Contains(matchedKeyAndIndex[keyRef], matchIndex) {
			continue
		}
		iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, string(hashBytes))
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
		iStateErr = NewError(err, 3008)
		return
	}
	if valBytes != nil {
		fetchedKVMap[keyRef] = valBytes
		qEnv.ufetchedKVMap[keyRef] = valBytes
	}
	return
}
func (iState *iState) loadkindecesMap(stub shim.ChaincodeStubInterface, kindecesMap map[string]map[string][]byte, keyRef string, newHash string) (iStateErr Error) {
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
