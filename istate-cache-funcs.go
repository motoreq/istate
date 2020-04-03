//

package istate

import (
	"sync"
)

//
type kvCache struct {
	mux     sync.Mutex
	kvCache map[string][]byte
}

//
type keyEncKVCache struct {
	mux           sync.Mutex
	keyEncKVCache map[string]map[string][]byte
}

//
func (iState *iState) readkvCache(key string) (valBytes []byte, ok bool) {
	iState.kvCache.mux.Lock()
	defer iState.kvCache.mux.Unlock()
	valBytes, ok = iState.kvCache.kvCache[key]
	return
}

//
func (iState *iState) addkvCache(key string, valBytes []byte) {
	iState.kvCache.mux.Lock()
	defer iState.kvCache.mux.Unlock()
	iState.kvCache.kvCache[key] = valBytes
	return
}

//
func (iState *iState) readkeyEncKVCache(key string) (encKV map[string][]byte, ok bool) {
	iState.keyEncKVCache.mux.Lock()
	defer iState.keyEncKVCache.mux.Unlock()
	encKV, ok = iState.keyEncKVCache.keyEncKVCache[key]
	return
}

//
func (iState *iState) addkeyEncKVCache(key string, encKV map[string][]byte) {
	iState.keyEncKVCache.mux.Lock()
	defer iState.keyEncKVCache.mux.Unlock()
	iState.keyEncKVCache.keyEncKVCache[key] = encKV
	return
}

//
func (iState *iState) incDocsCounter(key string, count int) {
	iState.mux.Lock()
	defer iState.mux.Unlock()

	iState.docsCounter[key] += count
}

//
func (iState *iState) decDocsCounter(key string, count int) {
	iState.mux.Lock()
	defer iState.mux.Unlock()

	iState.docsCounter[key] -= count
}

//
func (iState *iState) readDocsCounter(key string) (count int, ok bool) {
	iState.mux.Lock()
	defer iState.mux.Unlock()

	count, ok = iState.docsCounter[key]
	return
}
