//

package istate

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"sync"
	"time"
)

func (iState *iState) fetchEq(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (fetchedKVMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	fetchedKVMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchNeq(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchGt(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchLt(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchGte(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchLte(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchCmplx(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSeq(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSneq(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSgt(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSlt(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSgte(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchSlte(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

func (iState *iState) fetchScmplx(stub shim.ChaincodeStubInterface, encodedKey string, qEnv *queryEnv) (keyValMap map[string][]byte, iStateErr Error) {
	start := encodedKey
	end := encodedKey + asciiLast
	keyValMap, iStateErr = iState.getStateByRange(stub, start, end, qEnv)
	if iStateErr != nil {
		return
	}
	return
}

// TODO Cache startKey, endKey too?
func (iState *iState) getStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string, qEnv *queryEnv) (fetchedKVMap map[string][]byte, iStateErr Error) {
	start := time.Now()
	fmt.Println("Before 4:", start)
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
	fmt.Println("After 4: ", time.Now().Sub(start))

	start = time.Now()
	fmt.Println("Before 5:", start)
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
		// keyRef := string(iteratorResult.GetValue())
		indexkey := iteratorResult.GetKey()
		keyRef := getKeyFromIndex(indexkey)

		iStateErr = loadFetchedKV(stub, fetchedKVMap, keyRef, qEnv)
		if iStateErr != nil {
			return
		}
	}
	fmt.Println("After 5: ", time.Now().Sub(start))
	return
}

func loadFetchedKV(stub shim.ChaincodeStubInterface, fetchedKVMap map[string][]byte, keyRef string, qEnv *queryEnv) (iStateErr Error) {
	var wg sync.WaitGroup
	if _, ok := qEnv.ufetchedKVMap[keyRef]; ok {
		if uValBytes, ok := fetchedKVMap[keyRef]; !ok {
			fetchedKVMap[keyRef] = uValBytes
		}
		return
	}
	wg.Add(1)
	go func() {
		// Doesn't fetch if already fetched before
		valBytes, err := stub.GetState(keyRef)
		if err != nil {
			iStateErr = NewError(err, 3008)
			return
		}
		fetchedKVMap[keyRef] = valBytes
		qEnv.ufetchedKVMap[keyRef] = valBytes
		wg.Done()
	}()
	wg.Wait()
	return
}
