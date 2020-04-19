/*
	Copyright 2020 Prasanth Sundaravelu

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

//
type Interface interface {
	// CreateState function is used to create a new state in the state db
	// It must be called via the Interface returned by NewiState function
	// It takes chaincode stub and the actual structure value as input params
	// Note: This function does not do state validations such as, checking whether state exists or not
	// before performing the operation
	CreateState(shim.ChaincodeStubInterface, interface{}) Error
	// ReadState function is used to read a state from state db
	// It must be called via the Interface returned by NewiState function
	// It takes chaincode stub and value of primary key as input params
	// It returns the actual structure value as an interface{}. The returned value can
	// be type asserted to the actual struct type before using
	// Note: This function does not do state validations such as, checking whether state exists or not
	// before performing the operation
	ReadState(shim.ChaincodeStubInterface, interface{}) (interface{}, Error)
	// UpdateState function is used to update a state from statedb
	// It must be called via the Interface returned by NewiState function
	// It takes chaincode stub and the actual structure value as input params
	// Note: This function does not do state validations such as, checking whether state exists or not
	// before performing the operation
	UpdateState(shim.ChaincodeStubInterface, interface{}) Error
	// DeleteState function is used to delete a state from state db
	// It must be called via the Interface returned by NewiState function
	// It takes chaincode stub and value of primary key as input params
	// Note: This function does not do state validations such as, checking whether state exists or not
	// before performing the operation
	DeleteState(shim.ChaincodeStubInterface, interface{}) Error
	// Query function is used to perform Rich Queries over a state type in state db
	// It must be called via the Interface returned by NewiState function
	// It takes stub, query string and isInvoke (bool) flag as input params
	// isInvoke is set true, if this func is called in an invoke transaction
	// It returns slice of actual structure values as an interface{}. The returned value can
	// be type asserted to the actual slice of struct type before using
	// Learn more about the Rich query language formats in README.md
	Query(shim.ChaincodeStubInterface, string, bool) (interface{}, Error)
	// CompactIndex(shim.ChaincodeStubInterface) Error
}
