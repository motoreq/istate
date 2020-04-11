//

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func evalAndFilterEq(stub shim.ChaincodeStubInterface, encQKeyVal map[string][]byte, keyEncKVMap map[string]map[string][]byte) {
	for encKeyWithoutStar := range encQKeyVal {
		filter(keyEncKVMap, encKeyWithoutStar, evalEq)
	}
}

func evalAndFilterNeq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evalAndFilterGt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evalAndFilterLt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evalAndFilterGte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evalAndFilterLte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evalAndFilterCmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evalAndFilterSeq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
}

func evalAndFilterSneq(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evaluateAndFilterSgt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {
}

func evaluateAndFilterSlt(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evaluateAndFilterSgte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evaluateAndFilterSlte(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}

func evaluateAndFilterScmplx(stub shim.ChaincodeStubInterface, query map[string]interface{}) {

}
