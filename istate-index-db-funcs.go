//
package istate

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//
func getIndex(stub shim.ChaincodeStubInterface, key string) (indexv indexValType, iStateErr Error) {
	valBytes, err := stub.GetState(key)
	if err != nil {
		iStateErr = NewError(err, 2022)
		return
	}
	switch valBytes != nil {
	case true:
		err := json.Unmarshal(valBytes, &indexv)
		if err != nil {
			iStateErr = NewError(err, 2016)
			return
		}
	default:
		// indexv = make(indexValType)
	}
	return
}

//
func putIndex(stub shim.ChaincodeStubInterface, key string, keyRef string) (iStateErr Error) {
	var indexv indexValType
	valBytes, err := stub.GetState(key)
	switch valBytes != nil {
	case true:
		err := json.Unmarshal(valBytes, &indexv)
		if err != nil {
			iStateErr = NewError(err, 2018)
			return
		}
		if _, ok := indexv[keyRef]; ok {
			return
		}
	default:
		indexv = make(indexValType)
	}

	indexv[keyRef] = struct{}{}

	mi, err := json.Marshal(indexv)
	if err != nil {
		iStateErr = NewError(err, 2017)
		return
	}

	err = stub.PutState(key, mi)
	if err != nil {
		iStateErr = NewError(err, 2021)
	}
	return
}

//
func delIndex(stub shim.ChaincodeStubInterface, key string, keyRef string) (iStateErr Error) {
	var indexv indexValType
	valBytes, err := stub.GetState(key)
	if err != nil {
		iStateErr = NewError(err, 2023)
		return
	}
	switch valBytes != nil {
	case true:
		err := json.Unmarshal(valBytes, &indexv)
		if err != nil {
			iStateErr = NewError(err, 2017)
			return
		}
		delete(indexv, keyRef)
	default:
		return
	}

	switch len(indexv) == 0 {
	case true:
		err := stub.DelState(key)
		if err != nil {
			iStateErr = NewError(err, 2019)
			return
		}
	default:
		mi, err := json.Marshal(indexv)
		if err != nil {
			iStateErr = NewError(err, 2020)
			return
		}
		err = stub.PutState(key, mi)
		if err != nil {
			iStateErr = NewError(err, 2022)
		}
	}

	return
}
