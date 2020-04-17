//

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

func evalAndFilterEq(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalEq)
	}
}

func evalAndFilterNeq(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalNeq)
	}
}

func evalAndFilterGt(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalGt)
	}
}

func evalAndFilterLt(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalLt)
	}
}

func evalAndFilterGte(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalGte)
	}
}

func evalAndFilterLte(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalLte)
	}
}

func (iState *iState) evalAndFilterCmplx(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) (iStateErr Error) {
	for encKeyWithStar := range encQKeyVal {
		iStateErr = filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), iState.evalCmplx)
		if iStateErr != nil {
			return
		}
	}
	return
}

func evalAndFilterSeq(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalSeq)
	}
}

func evalAndFilterSneq(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalSneq)
	}
}

func evalAndFilterSgt(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalSgt)
	}
}

func evalAndFilterSlt(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalSlt)
	}
}

func evalAndFilterSgte(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalSgte)
	}
}

func evalAndFilterSlte(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithStar := range encQKeyVal {
		filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalSlte)
	}
}

// keyRefSeparated index only expected
func filter(keyEncKVMap map[string]map[string][]byte, encQKey string, evalFunc func(string, map[string][]byte) (bool, Error)) (iStateErr Error) {
	for key, encKV := range keyEncKVMap {
		var ok bool
		ok, iStateErr = evalFunc(encQKey, encKV)
		if iStateErr != nil {
			return
		}
		if !ok {
			delete(keyEncKVMap, key)
		}
	}
	return
}

// Evaluators:

func evalEq(encQKey string, encKV map[string][]byte) (found bool, iState Error) {
	_, found = encKV[encQKey]
	return
}

func evalNeq(encQKey string, encKV map[string][]byte) (ok bool, iState Error) {
	qpartIndex, _ := removeNValsFromIndex(encQKey, 2)
	for index := range encKV {
		// If this is the field
		if strings.Contains(index, qpartIndex) {
			// If neq matches
			if index != encQKey {
				ok = true
			}
		}
	}
	return
}

func evalGt(encQKey string, encKV map[string][]byte) (ok bool, iState Error) {
	qpartIndex, qremovedVals := removeNValsFromIndex(encQKey, 2)
	switch isNum(qremovedVals[1]) {
	case true:
		postive1, iStateErr := isPositive(qremovedVals[1])
		if iStateErr != nil {
			panic(iStateErr)
		}
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				_, removedVals := removeNValsFromIndex(index, 2)
				postive2, iStateErr := isPositive(removedVals[1])
				if iStateErr != nil {
					panic(iStateErr)
				}
				switch !postive1 && !postive2 {
				// If both negative
				case true:
					if index < encQKey {
						ok = true
						return
					}
				//
				default:
					if index > encQKey {
						ok = true
						return
					}
				}
			}
		}

	default:
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				if index > encQKey {
					ok = true
					return
				}
			}
		}

	}
	return
}

func evalLt(encQKey string, encKV map[string][]byte) (ok bool, iState Error) {
	qpartIndex, qremovedVals := removeNValsFromIndex(encQKey, 2)
	switch isNum(qremovedVals[1]) {
	case true:
		postive1, iStateErr := isPositive(qremovedVals[1])
		if iStateErr != nil {
			panic(iStateErr)
		}
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				_, removedVals := removeNValsFromIndex(index, 2)
				postive2, iStateErr := isPositive(removedVals[1])
				if iStateErr != nil {
					panic(iStateErr)
				}
				switch !postive1 && !postive2 {
				// If both negative
				case true:
					if index > encQKey {
						ok = true
						return
					}
				//
				default:
					if index < encQKey {
						ok = true
						return
					}
				}
			}
		}

	default:
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				if index < encQKey {
					ok = true
					return
				}
			}
		}

	}
	return
}

func evalGte(encQKey string, encKV map[string][]byte) (ok bool, iState Error) {
	qpartIndex, qremovedVals := removeNValsFromIndex(encQKey, 2)
	switch isNum(qremovedVals[1]) {
	case true:
		postive1, iStateErr := isPositive(qremovedVals[1])
		if iStateErr != nil {
			panic(iStateErr)
		}
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				_, removedVals := removeNValsFromIndex(index, 2)
				postive2, iStateErr := isPositive(removedVals[1])
				if iStateErr != nil {
					panic(iStateErr)
				}
				switch !postive1 && !postive2 {
				// If both negative
				case true:
					if index <= encQKey {
						ok = true
						return
					}
				//
				default:
					if index >= encQKey {
						ok = true
						return
					}
				}
			}
		}

	default:
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				if index >= encQKey {
					ok = true
					return
				}
			}
		}

	}
	return
}

func evalLte(encQKey string, encKV map[string][]byte) (ok bool, iState Error) {
	qpartIndex, qremovedVals := removeNValsFromIndex(encQKey, 2)
	switch isNum(qremovedVals[1]) {
	case true:
		postive1, iStateErr := isPositive(qremovedVals[1])
		if iStateErr != nil {
			panic(iStateErr)
		}
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				_, removedVals := removeNValsFromIndex(index, 2)
				postive2, iStateErr := isPositive(removedVals[1])
				if iStateErr != nil {
					panic(iStateErr)
				}
				switch !postive1 && !postive2 {
				// If both negative
				case true:
					if index >= encQKey {
						ok = true
						return
					}
				//
				default:
					if index <= encQKey {
						ok = true
						return
					}
				}
			}
		}

	default:
		for index := range encKV {
			// If this is the field
			if strings.Contains(index, qpartIndex) {
				if index <= encQKey {
					ok = true
					return
				}
			}
		}

	}
	return
}

//
func evalSeq(encQKey string, encKV map[string][]byte) (match bool, iState Error) {
	// Remove those which has neq 1
	_, ok := encKV[encQKey]
	if !ok {
		return
	}

	// Field name
	partIndex, _ := removeNValsFromIndex(encQKey, 2)

	for index := range encKV {
		// If this is the field
		if strings.Contains(index, partIndex) {
			// But it has different value
			if !strings.Contains(index, encQKey) {
				return
			}
		}
	}

	match = true
	return
}

func evalSneq(encQKey string, encKV map[string][]byte) (match bool, iState Error) {
	// Remove those which has eq 1
	_, ok := encKV[encQKey]
	if ok {
		return
	}

	// Field name
	partIndex, _ := removeNValsFromIndex(encQKey, 2)

	for index := range encKV {
		// If this is the field
		if strings.Contains(index, partIndex) {
			// But it has the same match value
			if strings.Contains(index, encQKey) {
				return
			}
		}
	}

	match = true
	return
}

func evalSgt(encQKey string, encKV map[string][]byte) (match bool, iState Error) {
	blackListFound, _ := evalLte(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match, _ = evalGt(encQKey, encKV)
	return
}
func evalSlt(encQKey string, encKV map[string][]byte) (match bool, iState Error) {
	blackListFound, _ := evalGte(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match, _ = evalLt(encQKey, encKV)
	return
}
func evalSgte(encQKey string, encKV map[string][]byte) (match bool, iState Error) {
	blackListFound, _ := evalLt(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match, _ = evalGte(encQKey, encKV)
	return
}
func evalSlte(encQKey string, encKV map[string][]byte) (match bool, iState Error) {
	blackListFound, _ := evalGt(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match, _ = evalLte(encQKey, encKV)
	return
}
func (iState *iState) evalCmplx(encQKey string, encKV map[string][]byte) (match bool, iStateErr Error) {

	partIndex, removedVals := removeNValsFromIndex(encQKey, 2)

	partIndexRemovedSeparator := removeLastSeparator(partIndex)

	fieldName := iState.convertIndexToQueryFieldName(partIndexRemovedSeparator)
	// Get Type of val
	kind, iStateErr := getRightPrimitiveType(fieldName, iState.jsonFieldKindMap, iState.mapKeyKindMap)
	if iStateErr != nil {
		return
	}

	match, iStateErr = iState.parseCmplxAndEval(partIndex, removedVals[1], kind, encKV)
	if iStateErr != nil {
		return
	}
	return
}
