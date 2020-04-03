//

package istate

import (
	"strings"
)

func filter(keyEncKVMap map[string]map[string]string, encKey string, evalFunc func(string, map[string]string) bool) {
	for key, encKV := range keyEncKVMap {
		if !evalFunc(encKey, encKV) {
			delete(keyEncKVMap, key)
		}
	}
}

func evalEq(encKey string, encKV map[string]string) (found bool) {
	for index := range encKV {
		if strings.HasPrefix(index, encKey) {
			found = true
			break
		}
	}
	return
}
