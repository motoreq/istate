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
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"strings"
)

//
func (iState *iState) setStub(stub *shim.ChaincodeStubInterface) {
	iState.currentStub = stub
	return
}

// Note: It returns reflect.Value , not interface{}
func (iState *iState) unmarshalToStruct(valBytes []byte) (uObj reflect.Value, iStateErr Error) {
	singleElem := reflect.New(reflect.TypeOf(iState.structRef)).Interface()
	err := json.Unmarshal(valBytes, &singleElem)
	if err != nil {
		iStateErr = newError(err, 4004)
		return
	}
	uObj = reflect.ValueOf(singleElem).Elem()
	return
}

func (iState *iState) getQIndexMap(key string, valBytes []byte) (encodedKV map[string][]byte, iStateErr Error) {
	var tempVar map[string]interface{}
	err := json.Unmarshal(valBytes, &tempVar)
	if err != nil {
		iStateErr = newError(err, 4005)
		return
	}
	encodedKV, _, _, iStateErr = iState.encodeState(tempVar, key, "", 1) // keyRefSeperatedIndex = 1, query = false
	if iStateErr != nil {
		return
	}
	return
}

//
func convertObjToMap(obj interface{}) (uObj map[string]interface{}, iStateErr Error) {
	mo, err := json.Marshal(obj)
	if err != nil {
		iStateErr = newError(err, 4001)
		return
	}
	err = json.Unmarshal(mo, &uObj)
	if err != nil {
		iStateErr = newError(err, 4002)
		return
	}
	return
}

//
func getKeyByRange(stub shim.ChaincodeStubInterface, startKey, endKey string, limit ...int) (fetchedKVMap map[string][]byte, iStateErr Error) {
	if len(limit) == 0 {
		limit = []int{int32Biggest}
	}
	fetchedKVMap = make(map[string][]byte)
	iterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		iStateErr = newError(err, 3006)
		return
	}
	defer iterator.Close()
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = newError(err, 4003)
			return
		}
		key := iteratorResult.GetKey()
		val := iteratorResult.GetValue()
		fetchedKVMap[key] = val

		if i >= limit[0] {
			break
		}
	}
	return
}

func getKeyByRangeWithPagination(stub shim.ChaincodeStubInterface, startKey, endKey string, pagesize int32, bookmark string) (fetchedKVMap map[string][]byte, iStateErr Error) {
	fetchedKVMap = make(map[string][]byte)
	iterator, _, err := stub.GetStateByRangeWithPagination(startKey, endKey, pagesize, bookmark)
	if err != nil {
		iStateErr = newError(err, 3006)
		return
	}
	defer iterator.Close()
	for i := 0; iterator.HasNext(); i++ {
		iteratorResult, err := iterator.Next()
		if err != nil {
			iStateErr = newError(err, 4003)
			return
		}
		key := iteratorResult.GetKey()
		val := iteratorResult.GetValue()
		fetchedKVMap[key] = val
	}
	return
}
func getKeyFromIndex(indexkey string) (keyRef string) {
	splitPosition := strings.LastIndex(indexkey, null)
	if splitPosition != -1 {
		keyRef = indexkey[splitPosition+1:]
	}
	return
}

func splitIndexAndKey(index string) (partindex, keyRef string) {
	partindex = index
	splitPosition := strings.LastIndex(index, null)
	if splitPosition != -1 {
		partindex = index[:splitPosition]
		keyRef = index[splitPosition+1:]
	}
	return
}

func incLastChar(val string) (incVal string) {
	if len(val) == 0 {
		return
	}
	lastChar := val[len(val)-1]
	val = val[:len(val)-1]
	incVal = val + string(lastChar+1)
	return
}

func decLastChar(val string) (decVal string) {
	if len(val) == 0 {
		return
	}
	lastChar := val[len(val)-1]
	val = val[:len(val)-1]
	decVal = val + string(lastChar-1)
	return
}
