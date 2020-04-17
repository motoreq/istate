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
	"fmt"
	"hash/crc64"
	"reflect"
)

type cache struct {
	uObj    reflect.Value
	objHash string
	indeces map[string][]byte
}

// iState - kuVCache - Key, UnmarshalledVal cache: kvCache
//  - keyIndexKVCache - Key, Index key-val : indecesCache
// This loader is for kvCache
func (iState *iState) loader(key interface{}) (val interface{}, iStateErr error) {
	keyString, ok := key.(string)
	if !ok {
		iStateErr = newError(nil, 6001, reflect.TypeOf(key))
		return
	}
	// See what's wrong ..
	if keyString == "" {
		return
	}
	stubP := iState.currentStub
	valBytes, err := (*stubP).GetState(keyString)
	if err != nil {
		iStateErr = newError(err, 6002)
		return
	}
	if valBytes == nil {
		iStateErr = newError(nil, 6003)
	}

	uObj, iStateErr := iState.unmarshalToStruct(valBytes)
	if iStateErr != nil {
		return
	}

	hash := iState.hash(valBytes)

	indexMap, iStateErr := iState.getQIndexMap(keyString, valBytes)
	if iStateErr != nil {
		return
	}

	val = cache{
		uObj:    uObj,
		objHash: hash,
		indeces: indexMap,
	}
	return
}

func (iState *iState) setCache(key string, obj interface{}, valBytes []byte, hashString string) (iStateErr Error) {
	indexMap, iStateErr := iState.getQIndexMap(key, valBytes)
	if iStateErr != nil {
		return
	}
	cache := cache{
		uObj:    reflect.ValueOf(obj),
		objHash: hashString,
		indeces: indexMap,
	}
	iState.kvCache.Set(key, cache)
	return
}

func (iState *iState) removeCache(key string) (iStateErr Error) {
	ok := iState.kvCache.Remove(key)
	if !ok {
		iStateErr = newError(nil, 6007)
		return
	}
	return
}

func (iState *iState) getkvHash(key string) (hashString string, iStateErr Error) {
	val, err := iState.kvCache.Get(key)
	if err != nil {
		iStateErr = newError(err, 6004)
		return
	}
	hashString = val.(cache).objHash
	return
}

func (iState *iState) getuObj(key string) (uObj reflect.Value, iStateErr Error) {
	val, err := iState.kvCache.Get(key)
	if err != nil {
		iStateErr = newError(err, 6005)
		return
	}
	uObj = val.(cache).uObj
	return
}

func (iState *iState) getIndeces(key string) (indeces map[string][]byte, iStateErr Error) {
	val, err := iState.kvCache.Get(key)
	if err != nil {
		iStateErr = newError(err, 6006)
		return
	}
	// See what's wrong
	if val == nil {
		indeces = make(map[string][]byte)
		return
	}
	indeces = val.(cache).indeces
	return
}

func (iState *iState) hash(valBytes []byte) (checkSum string) {
	return fmt.Sprintf("%0x", crc64.Checksum(valBytes, iState.hashTable))
}
