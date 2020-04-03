//

package istate

import (
	"strings"
)

func filter(keyEncKVMap map[string]map[string][]byte, encKey string, evalFunc func(string, map[string][]byte) bool) {
	for key, encKV := range keyEncKVMap {
		if !evalFunc(encKey, encKV) {
			delete(keyEncKVMap, key)
		}
	}
}

func evalEq(encKey string, encKV map[string][]byte) (found bool) {
	for index := range encKV {
		if strings.HasPrefix(index, encKey) {
			found = true
			break
		}
	}
	return
}
