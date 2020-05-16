package main

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

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/motoreq/istate"
)

type TestSmartContract struct {
	TestStructiState istate.Interface
	SomeStructiState istate.Interface
}

// Init initializes chaincode.
func (sc *TestSmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := sc.init()
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (sc *TestSmartContract) init() error {
	iStateOpt := istate.Options{
		CacheSize:             1000000,
		DefaultCompactionSize: 10000,
	}
	TestStructiState, err := istate.NewiState(TestStruct{}, iStateOpt)
	if err != nil {
		return err
	}
	sc.TestStructiState = TestStructiState
	return nil
}

//Invoke is the entry point for any transaction
func (sc *TestSmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	if sc.TestStructiState == nil {
		err := sc.init()
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	return sc.handleFunctions(stub)
}

//Query is the entry point for any data retrival
func (sc *TestSmartContract) Query(stub shim.ChaincodeStubInterface) pb.Response {
	if sc.TestStructiState == nil {
		err := sc.init()
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	return sc.handleFunctions(stub)
}

func (sc *TestSmartContract) handleFunctions(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()
	// fmt.Println("Function:")
	// fmt.Println(function)

	switch function {
	case "CreateState":
		return sc.CreateState(stub)
	case "UpdateState":
		return sc.UpdateState(stub)
	case "PartialUpdateState":
		return sc.PartialUpdateState(stub)
	case "ReadState":
		return sc.ReadState(stub)
	case "DeleteState":
		return sc.DeleteState(stub)
	case "QueryState":
		return sc.QueryState(stub)
		// case "CompactIndex":
		// 	return sc.CompactIndex(stub)
	}

	return shim.Error(fmt.Sprintf("Invalid function provided: %v", function))
}

func main() {
	sc := &TestSmartContract{}
	err := shim.Start(sc)
	if err != nil {
		fmt.Printf("Error starting  chaincode: %v\n", err)
	}

}

// ====================================================================================
// Test
// ====================================================================================

//
func (sc *TestSmartContract) CreateState(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Invalid no. of arguments - Expecting 1.")
	}

	var input TestStruct
	inputJSON := args[0]

	err = json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = sc.TestStructiState.CreateState(stub, input)
	if err != nil {
		return shim.Error(err.Error())
	}
	output := fmt.Sprintf("Successfully saved: %v", input)
	return shim.Success([]byte(output))

}

//
func (sc *TestSmartContract) UpdateState(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Invalid no. of arguments - Expecting 1.")
	}

	var input TestStruct
	inputJSON := args[0]

	err = json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = sc.TestStructiState.UpdateState(stub, input)
	if err != nil {
		return shim.Error(fmt.Sprintf("err %v, err==nil %v, type %T", err, err == nil, err))
	}
	output := fmt.Sprintf("Successfully updated: %v", input)
	return shim.Success([]byte(output))

}

//
func (sc *TestSmartContract) ReadState(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Invalid no. of arguments - Expecting 1.")
	}

	var input ReadStateInput
	inputJSON := args[0]

	err = json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	stateInterface, err := sc.TestStructiState.ReadState(stub, input.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	ms, err := json.Marshal(stateInterface)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(ms)

}

//
func (sc *TestSmartContract) DeleteState(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Invalid no. of arguments - Expecting 1.")
	}

	var input DeleteStateInput
	inputJSON := args[0]

	err = json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = sc.TestStructiState.DeleteState(stub, input.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	output := fmt.Sprintf("Deleted state successfully: %v", input.ID)
	return shim.Success([]byte(output))

}

//
func (sc *TestSmartContract) QueryState(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Invalid no. of arguments - Expecting 1.")
	}

	var input QueryInput
	inputJSON := args[0]

	err = json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	result, err := sc.TestStructiState.Query(stub, input.QueryString, input.IsInvoke)
	if err != nil {
		return shim.Error(err.Error())
	}

	out := QueryOut{
		Result: result.([]TestStruct),
		Count:  len(result.([]TestStruct)),
	}
	mr, err := json.Marshal(out)
	if err != nil {
		return shim.Error(err.Error())
	}
	// _ = mr
	return shim.Success(mr)

}

//
func (sc *TestSmartContract) PartialUpdateState(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Invalid no. of arguments - Expecting 1.")
	}

	var input PartialUpdateInput
	inputJSON := args[0]

	err = json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = sc.TestStructiState.PartialUpdateState(stub, input.PrimaryKey, input.PartialObject)
	if err != nil {
		return shim.Error(fmt.Sprintf("err %v, err==nil %v, type %T", err, err == nil, err))
	}
	output := fmt.Sprintf("Successfully updated: %v", input)
	return shim.Success([]byte(output))

}

// //
// func (sc *TestSmartContract) CompactIndex(stub shim.ChaincodeStubInterface) pb.Response {
// 	var err error
// 	err = sc.TestStructiState.CompactIndex(stub)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success([]byte("Success"))
// }
