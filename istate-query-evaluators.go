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

func evalAndFilterCmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

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

func evaluateAndFilterScmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

// keyRefSeparated index only expected
func filter(keyEncKVMap map[string]map[string][]byte, encQKey string, evalFunc func(string, map[string][]byte) bool) {
	for key, encKV := range keyEncKVMap {
		if !evalFunc(encQKey, encKV) {
			delete(keyEncKVMap, key)
		}
	}
}

// Evaluators:

func evalEq(encQKey string, encKV map[string][]byte) (found bool) {
	_, found = encKV[encQKey]
	return
}

func evalNeq(encQKey string, encKV map[string][]byte) (ok bool) {
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

func evalGt(encQKey string, encKV map[string][]byte) (ok bool) {
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

func evalLt(encQKey string, encKV map[string][]byte) (ok bool) {
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

func evalGte(encQKey string, encKV map[string][]byte) (ok bool) {
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

func evalLte(encQKey string, encKV map[string][]byte) (ok bool) {
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
func evalSeq(encQKey string, encKV map[string][]byte) (match bool) {
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

func evalSneq(encQKey string, encKV map[string][]byte) (match bool) {
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

func evalSgt(encQKey string, encKV map[string][]byte) (match bool) {
	blackListFound := evalLte(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match = evalGt(encQKey, encKV)
	return
}
func evalSlt(encQKey string, encKV map[string][]byte) (match bool) {
	blackListFound := evalGte(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match = evalLt(encQKey, encKV)
	return
}
func evalSgte(encQKey string, encKV map[string][]byte) (match bool) {
	blackListFound := evalLt(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match = evalGte(encQKey, encKV)
	return
}
func evalSlte(encQKey string, encKV map[string][]byte) (match bool) {
	blackListFound := evalGt(encQKey, encKV)
	if blackListFound {
		match = false
		return
	}

	match = evalLte(encQKey, encKV)
	return
}
