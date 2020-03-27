// Copyright 2020 <>. All rights reserved.

package istate

import (
	"fmt"
)

const errorPrefix = "iState Error: "

// errorCodes stores a map of Error Code to Error message.
var errorCodes = map[int]string{
	0: "Error Test",

	1001: "CreateState: PutStateError",
	1002: "Marshal error at SaveState",
	1003: "Unmarshal error at SaveState",
	1004: "UpdateState: Found no change with update",
	1005: "ReadState: GetStateError",
	1006: "Marshal error at UpdateState",
	1007: "Unmarshal error at UpdateState",
	1008: "UpdateState: DelStateError",
	1009: "UpdateState: PutState",
	1010: "Marshal error at DeleteState",
	1011: "Unmarshal error at DeleteState",
	1012: "DeleteState: DelStateError",
	1013: "DeleteState: Key does not exist: %s",

	// Util errors
	2001: "json tag is not set for struct field: %s of type: %v",
	2002: "primary tag is not set for any field in struct type: %v (or) primary field value is empty",
	2003: "Marshal error at encodeState",
	2004: "Unmarshal error at encodeState",
	2005: "encodeState: Unsupported kind: %v",
	2006: "Marshal error at generateRelationalTables",
	2007: "Unmarshal error at generateRelationalTables",
	2008: "generateRelationalTables: Pointer of structure expected. Received: %v",
	2009: "encodeState: Integer overflow. Number with digits: %d",
	2010: "findDifference: Unsupported kind received: %v",
	2011: "findSliceDifference: Source and Target type are not same, Received: %v %v",
	2012: "findSliceDifference: Only slice kind is expected, Received: %v %v",
	2013: "findMapDifference: Source and Target type are not same, Received: %v %v",
	2014: "findMapDifference: Only map kind is expected, Received: %v %v",
}

// Error is the interface of this package.
type Error interface {
	Error() string
	GetCode() int
}

// iStateError is the error type of this package.
type iStateError struct {
	Err  string
	Code int
}

// Error function is defined to let Error implement error interface.
func (err iStateError) Error() string {
	return err.Err
}

// Error function is defined to let Error implement error interface.
func (err iStateError) GetCode() int {
	return err.Code
}

// NewError function is to create Errors in a more readable way.
func NewError(err error, code int, params ...interface{}) Error {
	iStateLogger.Debugf("Inside NewError")
	defer iStateLogger.Debugf("Exiting NewError")
	msg := fmt.Sprintf("%d: ", code) + fmt.Sprintf(errorCodes[code], params...)
	if err != nil {
		msg = errorPrefix + msg + ": " + err.Error()
	}
	return iStateError{
		Err:  msg,
		Code: code,
	}
}
