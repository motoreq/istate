//

package istate

import (
// "strings"
)

// keyRefSeparated index only expected
func filter(keyEncKVMap map[string]map[string][]byte, encQKey string, evalFunc func(string, map[string][]byte) bool) {
	for key, encKV := range keyEncKVMap {
		if !evalFunc(encQKey, encKV) {
			delete(keyEncKVMap, key)
		}
	}
}

// func evalEq(encQKey string, encKV map[string][]byte) (found bool) {
// 	for index := range encKV {
// 		if strings.HasPrefix(index, encQKey) {
// 			found = true
// 			break
// 		}
// 	}
// 	return
// }

func evalEq(encQKey string, encKV map[string][]byte) (found bool) {
	_, found = encKV[encQKey]
	return
}
