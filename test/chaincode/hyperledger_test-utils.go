package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	// "strings"
	"testing"
)

func InitChaincode(test *testing.T) *shim.MockStub {
	stub := shim.NewMockStub("testingStub", new(TestSmartContract))
	result := stub.MockInit("000", nil)

	if result.Status != shim.OK {
		test.FailNow()
	}
	return stub
}

func Invoke(test *testing.T, stub *shim.MockStub, txID int, function string, args ...string) pb.Response {

	cc_args := make([][]byte, 1+len(args))
	cc_args[0] = []byte(function)
	for i, arg := range args {
		cc_args[i+1] = []byte(arg)
	}
	result := stub.MockInvoke(fmt.Sprintf("%d", txID), cc_args)
	// fmt.Println("Call:    ", function, "(", strings.Join(args, ","), ")")
	// fmt.Println("RetCode: ", result.Status)
	// fmt.Println("RetMsg:  ", result.Message)
	// fmt.Println("Payload: ", string(result.Payload))
	return result
}

func MarshalAndStringify(input interface{}) (string, error) {
	marshalledInput, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(marshalledInput), nil
}
