//

package istate

import (
	// "fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func evalAndFilterEq(stub shim.ChaincodeStubInterface, encodedKeyVal map[string]string, keyEncKVMap map[string]map[string]string) {
	for encKeyWithoutStar := range encodedKeyVal {
		filter(keyEncKVMap, encKeyWithoutStar, evalEq)
	}
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
