//

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//
type Interface interface {
	CreateState(shim.ChaincodeStubInterface, interface{}) Error
	ReadState(shim.ChaincodeStubInterface, interface{}) ([]byte, Error)
	UpdateState(shim.ChaincodeStubInterface, interface{}) Error
	DeleteState(shim.ChaincodeStubInterface, interface{}) Error
	Query(shim.ChaincodeStubInterface, string) (interface{}, Error)
	CompactIndex(shim.ChaincodeStubInterface) Error
}
