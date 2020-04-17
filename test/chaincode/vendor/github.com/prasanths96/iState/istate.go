//

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var iStateLogger = shim.NewLogger("SmartContractMain")

// Init function is called during initialization stage of
// the program that uses this package.
// Init is used to set several initialization parameters
// to istate package eg: Logging level.
func Init(level shim.LoggingLevel) {
	iStateLogger.SetLevel(level)
	iStateLogger.Infof("iStateLogger logging level set to: %v", level)
}
