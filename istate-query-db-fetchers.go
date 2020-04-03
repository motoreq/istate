//

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
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

// TODO Cache startKey, endKey too
func (iState *iState) getStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string, qEnv *queryEnv) (fetchedKVMap map[string][]byte, iStateErr Error) {
	fetchedKVMap = make(map[string][]byte)
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
		keyRef := indexkey[strings.LastIndex(indexkey, null)+1:]
		if _, ok := qEnv.ufetchedKVMap[keyRef]; ok {
			if uValBytes, ok := fetchedKVMap[keyRef]; !ok {
				fetchedKVMap[keyRef] = uValBytes
			}
			continue
		}
		// Doesn't fetch if already fetched before
		// Do Multi thread?
		valBytes, err := stub.GetState(keyRef)
		if err != nil {
			iStateErr = NewError(err, 3008)
		}
		fetchedKVMap[keyRef] = nil
		qEnv.ufetchedKVMap[keyRef] = valBytes
	}
	return
}
