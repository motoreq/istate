package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/prasanths96/istate"
)

type TestSmartContract struct {
	TestStructiState istate.Interface
	SomeStructiState istate.Interface
}

// Init initializes chaincode.
func (sc *TestSmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	TestStructiState, err := istate.NewiState(TestStruct{})
	if err != nil {
		return shim.Error(err.Error())
	}
	sc.TestStructiState = TestStructiState
	return shim.Success(nil)
}

//Invoke is the entry point for any transaction
func (sc *TestSmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return sc.handleFunctions(stub)
}

//Query is the entry point for any data retrival
func (sc *TestSmartContract) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return sc.handleFunctions(stub)
}

func (sc *TestSmartContract) handleFunctions(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()
	fmt.Println("Function:")
	fmt.Println(function)

	switch function {
	case "CreateState":
		return sc.CreateState(stub)
	case "UpdateState":
		return sc.UpdateState(stub)
	case "ReadState":
		return sc.ReadState(stub)
	case "DeleteState":
		return sc.DeleteState(stub)
	case "QueryState":
		return sc.QueryState(stub)
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

	stateBytes, err := sc.TestStructiState.ReadState(stub, input.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(stateBytes)

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

	result, err := sc.TestStructiState.Query(stub, input.QueryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	mr, err := json.Marshal(result.([]TestStruct))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(mr)

}
