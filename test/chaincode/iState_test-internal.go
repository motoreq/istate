package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"testing"
)

func CreateState(test *testing.T, stub *shim.MockStub, input TestStruct, txID int) pb.Response {
	inputString, err := MarshalAndStringify(input)
	if err != nil {
		test.FailNow()
	}
	return Invoke(test, stub, txID, "CreateState", inputString)
}

func QueryState(test *testing.T, stub *shim.MockStub, input interface{}, txID int) pb.Response {
	inputString, err := MarshalAndStringify(input)
	if err != nil {
		test.FailNow()
	}
	return Invoke(test, stub, txID, "QueryState", inputString)
}
