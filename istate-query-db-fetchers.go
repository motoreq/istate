//

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (iState *iState) fetchEq(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchNeq(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchGt(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchLt(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchGte(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchLte(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchCmplx(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSeq(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSneq(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSgt(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSlt(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSgte(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSlte(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchScmplx(stub shim.ChaincodeStubInterface, encodedKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	kindecesMap, iStateErr = iState.getStateByRange(stub, start, end)
	if iStateErr != nil {
		return
	}
	return
}

// TODO Cache startKey, endKey too?
func (iState *iState) getStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string) (kindecesMap map[string]map[string][]byte, iStateErr Error) {
	kindecesMap = make(map[string]map[string][]byte)
	// Compact Index
	cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))

	cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	if iStateErr != nil {
		return
	}
	for keyRef, hashString := range cIndexV {

		iStateErr = iState.loadkindecesMap(stub, kindecesMap, keyRef, hashString)
		if iStateErr != nil {
			return
		}
	}

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
	// Compact Index
	cIndexKey, _ := generateCIndexKey(removeLastSeparator(startKey))
	cIndexV, iStateErr := fetchCompactIndex(stub, cIndexKey)
	if iStateErr != nil {
		return
	}

	for keyRef := range cIndexV {
		iStateErr = loadFetchedKV(stub, fetchedKVMap, keyRef, qEnv)
		if iStateErr != nil {
			return
		}
	}

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
