/*
	Copyright 2020 Motoreq Infotech Pvt Ltd

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package istate

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var iStateLogger = shim.NewLogger("SmartContractMain")

// Init function is called during initialization stage of
// the program that uses this package.
// Init is used to set several initialization parameters
// to istate package eg: Logging level.
// func Init(level shim.LoggingLevel) {
// 	iStateLogger.SetLevel(level)
// 	iStateLogger.Infof("iStateLogger logging level set to: %v", level)
// }
