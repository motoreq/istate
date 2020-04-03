// Copyright 2020 <>. All rights reserved.

package istate

import (
	// "fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (iState *iState) evalAndFilterEq(stub shim.ChaincodeStubInterface, query []map[string]interface{}, keyEncKVMap map[string]map[string][]byte) (iStateErr Error) {
	keyref := ""
	for i := 0; i < len(query); i++ {
		var encodedKeyVal map[string][]byte
		encodedKeyVal, _, _, iStateErr = iState.encodeState(query[i], keyref, true)
		if iStateErr != nil {
			return
		}

		for encKeyWithStar := range encodedKeyVal {
			filter(keyEncKVMap, removeStarFromKey(encKeyWithStar), evalEq)
		}
	}
	return
}

func evalAndFilterNeq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("NEQ QUERY:", query)

}

func evalAndFilterGt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("GT QUERY:", query)

}

func evalAndFilterLt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("LT QUERY:", query)

}

func evalAndFilterGte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("GTE QUERY:", query)

}

func evalAndFilterLte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("LTE QUERY:", query)

}

func evalAndFilterCmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("CMPLX QUERY:", query)

}

func evalAndFilterSeq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*EQ QUERY:", query)

}

func evalAndFilterSneq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*NEQ QUERY:", query)

}

func evaluateAndFilterSgt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*GT QUERY:", query)

}

func evaluateAndFilterSlt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*LT QUERY:", query)

}

func evaluateAndFilterSgte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*GTE QUERY:", query)

}

func evaluateAndFilterSlte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*LTE QUERY:", query)

}

func evaluateAndFilterScmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
	// fmt.Println("*CMPLX QUERY:", query)

}
