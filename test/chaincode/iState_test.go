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
