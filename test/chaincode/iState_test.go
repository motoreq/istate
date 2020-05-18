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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var stub *shim.MockStub
var txCounter int

func TestInit(test *testing.T) {
	stub = InitChaincode(test)
}

func TestCreateState(test *testing.T) {
	asrt := assert.New(test)
	for i := 0; i < createStateCount; i++ {
		txCounter++
		input := TestStruct{
			DocType: "bleh",
			ID:      "bleh" + strconv.Itoa(i),
			AnInt:   createStateCount,
		}
		resp := CreateState(test, stub, input, txCounter)
		asrt.Equal(200, int(resp.Status), resp.Message)
	}
}

func TestReadStateById(test *testing.T) {
	var id string
	asrt := assert.New(test)
	for i := 0; i < createStateCount; i++ {
		id = "bleh" + strconv.Itoa(i)
		var input = ReadStateInput{id}
		resp := StateOperation(test, stub, input, txCounter, "ReadState")
		asrt.Equal(200, int(resp.Status), resp.Message)
	}
}

func TestQueryState(test *testing.T) {
	asrt := assert.New(test)
	for i := 0; i < queryNum; i++ {
		txCounter++
		input := QueryInput{
			QueryString: `[{"anInt":"eq ` + strconv.Itoa(createStateCount) + `"}]`,
		}
		resp := StateOperation(test, stub, input, txCounter, "QueryState")
		asrt.Equal(200, int(resp.Status), resp.Message)
		var output QueryOut
		json.Unmarshal(resp.Payload, &output)
		asrt.Equal(createStateCount, len(output.Result), fmt.Sprintf("Should be %d", createStateCount))
	}
}

func TestUpdateState(test *testing.T) {
	asrt := assert.New(test)
	for i := 0; i < createStateCount; i++ {
		txCounter++
		input := TestStruct{
			DocType: "bleh",
			ID:      "bleh" + strconv.Itoa(i),
			AnInt:   -2, // New value for AnInt
		}
		resp := StateOperation(test, stub, input, txCounter, "UpdateState")
		asrt.Equal(200, int(resp.Status), resp.Message)
	}
}

func TestPartialUpdateState(test *testing.T) {
	asrt := assert.New(test)
	var partialObject = make(map[string]interface{})
	partialObject["anInt"] = -1
	for i := 0; i < createStateCount; i++ {
		txCounter++
		input := PartialUpdateInput{
			PrimaryKey:    "bleh" + strconv.Itoa(i),
			PartialObject: partialObject, // New value for AnInt
		}
		resp := StateOperation(test, stub, input, txCounter, "PartialUpdateState")
		asrt.Equal(200, int(resp.Status), resp.Message)
	}
}

func TestDeleteState(test *testing.T) {
	asrt := assert.New(test)
	for i := 0; i < createStateCount; i++ {
		txCounter++
		input := DeleteStateInput{
			ID: "bleh" + strconv.Itoa(i),
		}
		resp := StateOperation(test, stub, input, txCounter, "DeleteState")
		asrt.Equal(200, int(resp.Status), resp.Message)
	}
}
