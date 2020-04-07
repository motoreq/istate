package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"testing"
)

// var stub *shim.MockStub
// var createdRuleIds []string

// func TestInit(test *testing.T) {
// 	stub = InitChaincode(test)
// }

// func TestCreateRuleType0(test *testing.T) {
// 	input := "{\"Type\":0, \"OperationType\":0, \"JSONStruct\":\"{\\\"a\\\":\\\"\\\", \\\"b\\\":\\\"\\\"}\", \"Description\":\"\"}"
// 	result := Invoke(test, stub, 1, "CreateRuleType0", input)
// 	if result.Status != shim.OK {
// 		test.Fail()
// 	}
// 	resultString := string(result.Payload)
// 	index := strings.Index(resultString, "rule_")
// 	ruleID := resultString[index:]
// 	createdRuleIds = append(createdRuleIds, ruleID)

// 	input = "{\"Type\":0, \"OperationType\":0, \"JSONStruct\":\"{\\\"b\\\":\\\"\\\", \\\"a\\\":\\\"\\\"}\", \"Description\":\"\"}"
// 	result = Invoke(test, stub, 2, "CreateRuleType0", input)
// 	if result.Status == shim.OK {
// 		test.Fail()
// 	}

// 	input = "{\"Type\":0, \"OperationType\":0, \"JSONStruct\":\"{\\\"a\\\":\\\"\\\"}\", \"Description\":\"\"}"
// 	result = Invoke(test, stub, 3, "CreateRuleType0", input)
// 	if result.Status != shim.OK {
// 		test.Fail()
// 	}
// 	resultString = string(result.Payload)
// 	index = strings.Index(resultString, "rule_")
// 	ruleID = resultString[index:]
// 	createdRuleIds = append(createdRuleIds, ruleID)
// }

// func TestReadRule(test *testing.T) {
// 	fmt.Println("Inside TestReadRule")
// 	fmt.Println(createdRuleIds)
// 	for i := 0; i < len(createdRuleIds); i++ {
// 		input := struct {
// 			ID string `json:"id"`
// 		}{
// 			ID: createdRuleIds[i],
// 		}
// 		marshalledInput, err := json.Marshal(input)
// 		if err != nil {
// 			test.Fail()
// 		}

// 		result := Invoke(test, stub, 1, "ReadRule", string(marshalledInput))
// 		if result.Status != shim.OK {
// 			test.Fail()
// 		}
// 	}
// }

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
