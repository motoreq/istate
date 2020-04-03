//

package istate

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"sync"
)

type iState struct {
	mux               sync.Mutex
	structRef         interface{}
	fieldJSONIndexMap map[string]int
	jsonFieldKindMap  map[string]reflect.Kind
	mapKeyKindMap     map[string]reflect.Kind
	depthKindMap      map[string]reflect.Kind
	primaryIndex      int
	docsCounter       map[string]int

	//
	kvCache       *kvCache
	keyEncKVCache *keyEncKVCache
}

type kvCache struct {
	mux     sync.Mutex
	kvCache map[string][]byte
}

type keyEncKVCache struct {
	mux           sync.Mutex
	keyEncKVCache map[string]map[string][]byte
}

func (iState *iState) readkvCache(key string) (valBytes []byte, ok bool) {
	iState.kvCache.mux.Lock()
	defer iState.kvCache.mux.Unlock()
	valBytes, ok = iState.kvCache.kvCache[key]
	return
}

func (iState *iState) addkvCache(key string, valBytes []byte) {
	iState.kvCache.mux.Lock()
	defer iState.kvCache.mux.Unlock()
	iState.kvCache.kvCache[key] = valBytes
	return
}

func (iState *iState) readkeyEncKVCache(key string) (encKV map[string][]byte, ok bool) {
	iState.keyEncKVCache.mux.Lock()
	defer iState.keyEncKVCache.mux.Unlock()
	encKV, ok = iState.keyEncKVCache.keyEncKVCache[key]
	return
}

func (iState *iState) addkeyEncKVCache(key string, encKV map[string][]byte) {
	iState.keyEncKVCache.mux.Lock()
	defer iState.keyEncKVCache.mux.Unlock()
	iState.keyEncKVCache.keyEncKVCache[key] = encKV
	return
}

func (iState *iState) incDocsCounter(key string, count int) {
	iState.mux.Lock()
	defer iState.mux.Unlock()

	iState.docsCounter[key] += count
}

func (iState *iState) decDocsCounter(key string, count int) {
	iState.mux.Lock()
	defer iState.mux.Unlock()

	iState.docsCounter[key] -= count
}

func (iState *iState) readDocsCounter(key string) (count int, ok bool) {
	iState.mux.Lock()
	defer iState.mux.Unlock()

	count, ok = iState.docsCounter[key]
	return
}

// Need to Copy for every transaction
func (is *iState) CopyiState() (iStateInterface IStateInterface) {
	is.mux.Lock()
	defer is.mux.Unlock()
	iStateInterface = &iState{
		structRef:         is.structRef,
		fieldJSONIndexMap: is.fieldJSONIndexMap,
		jsonFieldKindMap:  is.jsonFieldKindMap,
		mapKeyKindMap:     is.mapKeyKindMap,
		depthKindMap:      is.depthKindMap,
		primaryIndex:      is.primaryIndex,
		docsCounter:       is.docsCounter,
	}
	return
}

// NewiState function is used to
func NewiState(object interface{}) (iStateInterface IStateInterface, iStateErr Error) {
	iStateLogger.Infof("Inside NewiState")
	defer iStateLogger.Infof("Exiting NewiState")

	filledRef := fillZeroValue(object)
	fmt.Printf("FILLEDREF: %v\n", filledRef)
	// A map of JSON fieldname & it's position in the struct
	fieldJSONIndexMap := make(map[string]int)
	jsonFieldKindMap := make(map[string]reflect.Kind)
	mapKeyKindMap := make(map[string]reflect.Kind)
	depthKindMap := make(map[string]reflect.Kind)
	iStateErr = generateFieldJSONIndexMap(filledRef, fieldJSONIndexMap)
	if iStateErr != nil {
		return
	}
	iStateErr = generatejsonFieldKindMap(filledRef, jsonFieldKindMap, mapKeyKindMap)
	if iStateErr != nil {
		return
	}
	iStateErr = generateDepthKindMap(filledRef, depthKindMap)
	if iStateErr != nil {
		return
	}
	var i int
	i, iStateErr = getPrimaryFieldIndex(filledRef)
	if iStateErr != nil {
		return
	}

	docsCounter := make(map[string]int)
	iStateInterface = &iState{
		structRef:         filledRef,
		fieldJSONIndexMap: fieldJSONIndexMap,
		jsonFieldKindMap:  jsonFieldKindMap,
		mapKeyKindMap:     mapKeyKindMap,
		depthKindMap:      depthKindMap,
		primaryIndex:      i,
		docsCounter:       docsCounter,
		kvCache:           &kvCache{kvCache: make(map[string][]byte)},
		keyEncKVCache:     &keyEncKVCache{keyEncKVCache: make(map[string]map[string][]byte)},
	}
	fmt.Println("JSON FIELD KIND MAP", jsonFieldKindMap)
	fmt.Println("MAP KEY KIND MAP", mapKeyKindMap)
	fmt.Println("DEPTH KIND MAP", depthKindMap)
	fmt.Println("=========================================================================================")
	fmt.Println("DOCSCOUNT MAP:", docsCounter)

	return
}

// CreateState function is used to
func (iState *iState) CreateState(stub shim.ChaincodeStubInterface, object interface{}) (iStateErr Error) {
	iStateLogger.Infof("Inside CreateState")
	defer iStateLogger.Infof("Exiting CreateState")

	if reflect.TypeOf(object) != reflect.TypeOf(iState.structRef) {
		iStateErr = NewError(nil, 1014, reflect.TypeOf(iState.structRef), reflect.TypeOf(object))
		return
	}

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

	encodedKeyValPairs, docsCounter, _, iStateErr := iState.encodeState(oMap, keyref)
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

	for index, count := range docsCounter {
		iState.incDocsCounter(index, count)
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

	if reflect.TypeOf(object) != reflect.TypeOf(iState.structRef) {
		iStateErr = NewError(nil, 1015, reflect.TypeOf(iState.structRef), reflect.TypeOf(object))
		return
	}

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
	// fmt.Println("Changes: ", appendOrModifyMap, deleteMap)

	// Delete first, so that TestStruct_aMap_user1 (map key) does not get deleted.
	// When updating same key of a map, the above key gets over-writted at first,
	// then when deleting, the over-written key gets deleted.
	deleteEncodedKeyValPairs, delDocsCounter, _, iStateErr := iState.encodeState(deleteMap, keyref)
	if iStateErr != nil {
		// iStateErr = iStateErr
		return
	}
	appendEncodedKeyValPairs, addDocsCounter, _, iStateErr := iState.encodeState(appendOrModifyMap, keyref)
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

	for index, count := range delDocsCounter {
		iState.decDocsCounter(index, count)
	}
	for index, count := range addDocsCounter {
		iState.incDocsCounter(index, count)
	}

	return nil
}

// DeleteState function is used to
func (iState *iState) DeleteState(stub shim.ChaincodeStubInterface, primaryKey interface{}) (iStateErr Error) {
	iStateLogger.Infof("Inside DeleteState")
	defer iStateLogger.Infof("Exiting DeleteState")

	keyref := fmt.Sprintf("%v", primaryKey)

	stateBytes, iStateErr := iState.ReadState(stub, keyref)
	if iStateErr != nil {
		return
	}

	if stateBytes == nil {
		iStateErr = NewError(nil, 1013, keyref)
		return
	}

	var source map[string]interface{}
	err := json.Unmarshal(stateBytes, &source)
	if err != nil {
		iStateErr = NewError(err, 1011)
		return
	}

	// Delete first, so that TestStruct_aMap_user1 (map key) does not get deleted.
	// When updating same key of a map, the above key gets over-writted at first,
	// then when deleting, the over-written key gets deleted.
	encodedKeyValPairs, delDocsCounter, _, iStateErr := iState.encodeState(source, keyref)
	if iStateErr != nil {
		return
	}
	// Main key - value
	encodedKeyValPairs[keyref] = stateBytes

	for key := range encodedKeyValPairs {
		err = stub.DelState(key)
		if err != nil {
			iStateErr = NewError(err, 1012)
			return
		}
	}

	for index, count := range delDocsCounter {
		iState.decDocsCounter(index, count)
	}

	return nil
}
