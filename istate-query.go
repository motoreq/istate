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
)

// Syntax and Keywords
const (
	eq    = "eq"
	neq   = "neq"
	gt    = "gt"
	lt    = "lt"
	gte   = "gte"
	lte   = "lte"
	cmplx = "cmplx"
	or    = "or"
	and   = "and"

	seq  = "^eq"
	sneq = "^neq"
	sgt  = "^gt"
	slt  = "^lt"
	sgte = "^gte"
	slte = "^lte"
	// scmplx = "^cmplx"

	splitDot = "."
)

type encodedKVs struct {
	eq    map[string][]byte
	neq   map[string][]byte
	gt    map[string][]byte
	lt    map[string][]byte
	gte   map[string][]byte
	lte   map[string][]byte
	cmplx map[string][]byte

	seq  map[string][]byte
	sneq map[string][]byte
	sgt  map[string][]byte
	slt  map[string][]byte
	sgte map[string][]byte
	slte map[string][]byte
	// scmplx map[string][]byte
}

type querys struct {
	eq    []map[string]interface{}
	neq   []map[string]interface{}
	gt    []map[string]interface{}
	lt    []map[string]interface{}
	gte   []map[string]interface{}
	lte   []map[string]interface{}
	cmplx []map[string]interface{}

	seq  []map[string]interface{}
	sneq []map[string]interface{}
	sgt  []map[string]interface{}
	slt  []map[string]interface{}
	sgte []map[string]interface{}
	slte []map[string]interface{}
	// scmplx []map[string]interface{}
}

type efficientKeyType struct {
	enckey        string
	genericField  string
	fetchFunc     func(shim.ChaincodeStubInterface, string, string) (map[string]map[string][]byte, Error)
	relatedEncKV  map[string][]byte
	relatedQueryP *[]map[string]interface{}
	i             int
}

type safeKeyFunc struct {
	encKey        string
	fetchFunc     func(shim.ChaincodeStubInterface, string, string) (map[string]map[string][]byte, Error)
	relatedEncKV  map[string][]byte
	relatedQueryP *[]map[string]interface{}
	i             int
}

type queryEnv struct {
	ufetchedKVMap map[string][]byte
	uindecesMap   map[string]map[string][]byte
}

// Query function is used to
// "," separated objects are considered "or" always
// queryString = [{"docType":"eq USERPROFILE_DOCTYPE", "doctor.whatever": "cmplx or(and(gt bla, lt bla),or(eq a, eq b))", "groups":["cmplx and(neq doctor, neq patient)"]}, {"docType":"eq USERPROFILE_DOCTYPE", "groups":["eq patient"]}]
func (iState *iState) Query(stub shim.ChaincodeStubInterface, queryString string) (finalResult interface{}, iStateErr Error) {
	iStateLogger.Infof("Inside Query")
	defer iStateLogger.Infof("Exiting Query")
	iState.setStub(&stub)

	qEnv := &queryEnv{}
	initQueryEnv(qEnv)

	var uQuery []map[string]interface{}
	err := json.Unmarshal([]byte(queryString), &uQuery)
	if err != nil {
		iStateErr = NewError(err, 3002)
		return
	}

	filteredKeys := make([]map[string]map[string][]byte, len(uQuery), len(uQuery))
	for i := 0; i < len(uQuery); i++ {
		filteredKeys[i], iStateErr = iState.parseAndEvalSingle(stub, uQuery[i])
		if iStateErr != nil {
			return
		}
	}
	combinedResults := orOperation(filteredKeys...)

	finalResult = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(iState.structRef)), len(combinedResults), len(combinedResults))
	i := 0
	for key := range combinedResults {
		var uObj reflect.Value
		uObj, iStateErr = iState.getuObj(key)
		if iStateErr != nil {
			return
		}
		finalResult.(reflect.Value).Index(i).Set(uObj)
		i++
	}

	finalResult = finalResult.(reflect.Value).Interface()
	return
}

func (iState *iState) parseAndEvalSingle(stub shim.ChaincodeStubInterface, uQuery map[string]interface{}) (filteredKeys map[string]map[string][]byte, iStateErr Error) {

	// Fields will be declared automatically and make not needed
	querySet := querys{}

	for index, val := range uQuery {
		if val, ok := val.(string); !ok {
			iStateErr = NewError(nil, 3003, reflect.TypeOf(val))
			return
		}

		// Get keyword and val
		var keyword, secondPart string
		keyword, secondPart, iStateErr = getKeyWordAndVal(val.(string))
		if iStateErr != nil {
			return
		}
		// For cmplx, cannot convert to right type
		var convertedVal interface{} = secondPart
		if keyword != cmplx { //&& keyword != scmplx {
			// Convert the string value to appropriate type
			convertedVal, iStateErr = convertToRightType(index, secondPart, iState.jsonFieldKindMap, iState.mapKeyKindMap)
			if iStateErr != nil {
				return
			}
		}
		newIndex := ""
		var newVal interface{}
		newIndex, newVal, iStateErr = iState.generateActualStructure(index, convertedVal)
		if iStateErr != nil {
			return
		}

		switch keyword {
		case eq:
			querySet.eq = addKeyWithoutOverLap(querySet.eq, newIndex, newVal)
		case neq:
			querySet.neq = addKeyWithoutOverLap(querySet.neq, newIndex, newVal)
		case gt:
			querySet.gt = addKeyWithoutOverLap(querySet.gt, newIndex, newVal)
		case lt:
			querySet.lt = addKeyWithoutOverLap(querySet.lt, newIndex, newVal)
		case gte:
			querySet.gte = addKeyWithoutOverLap(querySet.gte, newIndex, newVal)
		case lte:
			querySet.lte = addKeyWithoutOverLap(querySet.lte, newIndex, newVal)
		case cmplx:
			querySet.cmplx = addKeyWithoutOverLap(querySet.cmplx, newIndex, newVal)
		case seq:
			querySet.seq = addKeyWithoutOverLap(querySet.seq, newIndex, newVal)
		case sneq:
			querySet.sneq = addKeyWithoutOverLap(querySet.sneq, newIndex, newVal)
		case sgt:
			querySet.sgt = addKeyWithoutOverLap(querySet.sgt, newIndex, newVal)
		case slt:
			querySet.slt = addKeyWithoutOverLap(querySet.slt, newIndex, newVal)
		case sgte:
			querySet.sgte = addKeyWithoutOverLap(querySet.sgte, newIndex, newVal)
		case slte:
			querySet.slte = addKeyWithoutOverLap(querySet.slte, newIndex, newVal)
		// case scmplx:
		// 	querySet.scmplx = addKeyWithoutOverLap(querySet.scmplx, newIndex, newVal)
		default:
			iStateErr = NewError(nil, 3005, keyword)
			return
		}

	}

	var bestKey string
	var fetchFunc func(shim.ChaincodeStubInterface, string, string) (map[string]map[string][]byte, Error)
	var queryEncodedKVset encodedKVs

	bestKey, fetchFunc, queryEncodedKVset, iStateErr = iState.getBestEncodedKeyFunc(querySet)
	if iStateErr != nil {
		return
	}

	kindecesMap, iStateErr := fetchFunc(stub, removeStarFromKey(bestKey), "index")
	evalAndFilterEq(stub, queryEncodedKVset.eq, kindecesMap)
	evalAndFilterNeq(stub, queryEncodedKVset.neq, kindecesMap)
	evalAndFilterGt(stub, queryEncodedKVset.gt, kindecesMap)
	evalAndFilterLt(stub, queryEncodedKVset.lt, kindecesMap)
	evalAndFilterGte(stub, queryEncodedKVset.gte, kindecesMap)
	evalAndFilterLte(stub, queryEncodedKVset.lte, kindecesMap)
	evalAndFilterSeq(stub, queryEncodedKVset.seq, kindecesMap)
	evalAndFilterSneq(stub, queryEncodedKVset.sneq, kindecesMap)
	evalAndFilterSgt(stub, queryEncodedKVset.sgt, kindecesMap)
	evalAndFilterSlt(stub, queryEncodedKVset.slt, kindecesMap)
	evalAndFilterSgte(stub, queryEncodedKVset.sgte, kindecesMap)
	evalAndFilterSlte(stub, queryEncodedKVset.slte, kindecesMap)

	iStateErr = iState.evalAndFilterCmplx(stub, queryEncodedKVset.cmplx, kindecesMap)
	if iStateErr != nil {
		return
	}

	// and operation between fields

	filteredKeys = kindecesMap
	return
}
