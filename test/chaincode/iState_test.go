package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"testing"
)

var stub *shim.MockStub

func TestInit(test *testing.T) {
	stub = InitChaincode(test)
}

func TestFlows(test *testing.T) {
	var txCounter = 0

	// Create Lists
	var createdList = make([]string, createStateCount, createStateCount)

	// Success cases
	for i := 0; i < createStateCount; i++ {
		txCounter++
		input := TestStruct{
			DocType: "bleh",
			ID:      "bleh" + strconv.Itoa(i),
			AnInt:   createStateCount,
		}
		resp := CreateState(test, stub, input, txCounter)
		if resp.Status != shim.OK {
			test.Fail()
		}
		resultString := string(resp.Payload)
		fmt.Println(resultString)
		createdList[i] = "bleh" + strconv.Itoa(i)
	}

	// Query
	for i := 0; i < queryNum; i++ {
		txCounter++
		input := QueryInput{
			QueryString: `[{"anInt":"eq ` + strconv.Itoa(createStateCount) + `"}]`,
		}
		resp := QueryState(test, stub, input, txCounter)
		if resp.Status != shim.OK {
			test.Fail()
		}
		fmt.Println("Query Resp: ", string(resp.GetPayload()))
	}

}
