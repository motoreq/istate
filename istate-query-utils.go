//

package istate

import (
	"strings"
)

func removeKeysWithStar(keyValue map[string][]byte) {
	for index := range keyValue {
		if strings.Contains(index, star) {
			delete(keyValue, index)
		}
	}
}

func removeStarFromKeys(keyValue map[string][]byte) {
	for index := range keyValue {
		// Replace is used as ReplaceAll isnt available in go version used in fabric image
		newIndex := strings.Replace(index, star, "", len(index))
		if newIndex != index {
			keyValue[newIndex] = keyValue[index]
			delete(keyValue, index)
		}

	}
}

func removeStarFromKey(key string) (newKey string) {
	// Replace is used as ReplaceAll isnt available in go version used in fabric image
	newKey = strings.Replace(key, star, "", len(key))
	return
}

func addKeyWithoutOverLap(query []map[string]interface{}, index string, value interface{}) (newQuery []map[string]interface{}) {
	newQuery = query
	successFlag := false
	for i := 0; i < len(newQuery); i++ {
		if _, ok := newQuery[i][index]; !ok {
			newQuery[i][index] = value
			successFlag = true
			break
		}
	}
	if !successFlag {
		tempMap := make(map[string]interface{})
		tempMap[index] = value
		newQuery = append(newQuery, tempMap)
	}
	return
}

func dotsToActualDepth(splitFieldName []string, val interface{}, curIndex ...int) (actualMap map[string]interface{}) {
	actualMap = make(map[string]interface{})
	if len(curIndex) == 0 {
		curIndex = []int{0}
	}

	if len(splitFieldName)-1 > curIndex[0] {
		actualMap[splitFieldName[curIndex[0]]] = dotsToActualDepth(splitFieldName, val, curIndex[0]+1)
	} else {
		actualMap[splitFieldName[curIndex[0]]] = val
	}

	return
}
