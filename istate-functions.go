// Copyright 2020 <>. All rights reserved.

package istate

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type iState struct {
	structRef         interface{}
	fieldJSONIndexMap map[string]int
	primaryIndex      int
}

// NewiState function is used to
func NewiState(object interface{}) (iStateInterface IStateInterface, iStateErr Error) {
	iStateLogger.Infof("Inside NewiState")
	defer iStateLogger.Infof("Exiting NewiState")

	filledRef := fillZeroValue(object)
	fmt.Printf("FILLEDREF: %v\n", filledRef)
	// A map of JSON fieldname & it's position in the struct
	fieldJSONIndexMap := make(map[string]int)
	generateFieldJSONIndexMap(filledRef, fieldJSONIndexMap)
	i, tiStateErr := getPrimaryFieldIndex(filledRef)
	if tiStateErr != nil {
		iStateErr = tiStateErr
		return
	}
	iStateInterface = &iState{
		structRef:         filledRef,
		fieldJSONIndexMap: fieldJSONIndexMap,
		primaryIndex:      i,
	}
	fmt.Printf("Generated fieldJSONIndexMap: %v\n", fieldJSONIndexMap)
	fmt.Println("MEW PRIMARY INDEX:", i)
	return
}

// CreateState function is used to
func (iState *iState) CreateState(stub shim.ChaincodeStubInterface, object interface{}) (iStateErr Error) {
	iStateLogger.Infof("Inside CreateState")
	defer iStateLogger.Infof("Exiting CreateState")

	// Find primary key
	keyref := iState.getPrimaryKey(object)

	mo, err := json.Marshal(object)
	if err != nil {
		iStateErr = NewError(err, 1002)
		return
	}

	var oMap map[string]interface{}
	err = json.Unmarshal(mo, &oMap)
	if err != nil {
		iStateErr = NewError(err, 1003)
		return
	}

	encodedKeyValPairs, iStateErr := iState.encodeState(oMap, keyref)
	if iStateErr != nil {
		return
	}

	// Main key - value
	encodedKeyValPairs[keyref] = mo

	for key, val := range encodedKeyValPairs {
		err = stub.PutState(key, val)
		if err != nil {
			iStateErr = NewError(err, 1001)
			return
		}
	}

	return nil
}

// ReadState function is used to
func (iState *iState) ReadState(stub shim.ChaincodeStubInterface, primaryKey interface{}) (stateBytes []byte, iStateErr Error) {
	iStateLogger.Infof("Inside ReadState")
	defer iStateLogger.Infof("Exiting ReadState")

	primaryKeyString := fmt.Sprintf("%v", primaryKey)

	stateBytes, err := stub.GetState(primaryKeyString)
	if err != nil {
		iStateErr = NewError(err, 1005)
	}

	return
}

// UpdateState function is used to
func (iState *iState) UpdateState(stub shim.ChaincodeStubInterface, object interface{}) (iStateErr Error) {
	iStateLogger.Infof("Inside UpdateState")
	defer iStateLogger.Infof("Exiting UpdateState")

	// Find primary key
	keyref := iState.getPrimaryKey(object)

	stateBytes, iStateErr := iState.ReadState(stub, keyref)
	if iStateErr != nil {
		// outiStateErr = tiStateErr
		return
	}

	if stateBytes == nil {
		return iState.CreateState(stub, object)
	}

	var source map[string]interface{}
	err := json.Unmarshal(stateBytes, &source)
	if err != nil {
		iStateErr = NewError(err, 1007)
		return
	}

	mo, err := json.Marshal(object)
	if err != nil {
		iStateErr = NewError(err, 1006)
		return
	}

	var target map[string]interface{}
	err = json.Unmarshal(mo, &target)
	if err != nil {
		iStateErr = NewError(err, 1007)
		return
	}

	appendOrModifyMap, deleteMap, iStateErr := iState.findDifference(source, target)
	if iStateErr != nil {
		//outiStateErr = iStateErr
		return
	}

	if len(appendOrModifyMap) == 0 && len(deleteMap) == 0 {
		//outiStateErr = NewError(nil, 1004)
		iStateErr = NewError(nil, 1004)
		return
	}
	fmt.Println("Changes: ", appendOrModifyMap, deleteMap)

	// Delete first, so that TestStruct_aMap_user1 (map key) does not get deleted.
	// When updating same key of a map, the above key gets over-writted at first,
	// then when deleting, the over-written key gets deleted.
	deleteEncodedKeyValPairs, iStateErr := iState.encodeState(deleteMap, keyref)
	if iStateErr != nil {
		// iStateErr = iStateErr
		return
	}
	appendEncodedKeyValPairs, iStateErr := iState.encodeState(appendOrModifyMap, keyref)
	if iStateErr != nil {
		// iStateErr = iStateErr
		return
	}
	// Main key - value
	appendEncodedKeyValPairs[keyref] = mo

	for key := range deleteEncodedKeyValPairs {
		err = stub.DelState(key)
		if err != nil {
			iStateErr = NewError(err, 1008)
			return
		}
	}
	for key, val := range appendEncodedKeyValPairs {
		err = stub.PutState(key, val)
		if err != nil {
			iStateErr = NewError(err, 1009)
			return
		}
	}

	return nil
}

// DeleteState function is used to
func (iState *iState) DeleteState(stub shim.ChaincodeStubInterface, primaryKey interface{}) (outiStateErr Error) {
	iStateLogger.Infof("Inside DeleteState")
	defer iStateLogger.Infof("Exiting DeleteState")

	keyref := fmt.Sprintf("%v", primaryKey)

	stateBytes, tiStateErr := iState.ReadState(stub, keyref)
	if tiStateErr != nil {
		outiStateErr = tiStateErr
		return
	}

	if stateBytes == nil {
		outiStateErr = NewError(nil, 1013, keyref)
		return
	}

	var source map[string]interface{}
	err := json.Unmarshal(stateBytes, &source)
	if err != nil {
		outiStateErr = NewError(err, 1011)
		return
	}

	// Delete first, so that TestStruct_aMap_user1 (map key) does not get deleted.
	// When updating same key of a map, the above key gets over-writted at first,
	// then when deleting, the over-written key gets deleted.
	encodedKeyValPairs, iStateErr := iState.encodeState(source, keyref)
	if iStateErr != nil {
		outiStateErr = iStateErr
		return
	}
	// Main key - value
	encodedKeyValPairs[keyref] = stateBytes

	for key := range encodedKeyValPairs {
		err = stub.DelState(key)
		if err != nil {
			outiStateErr = NewError(err, 1012)
			return
		}
	}

	return nil
}
