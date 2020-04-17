//

package istate

func andOperation(keyValuePairs ...map[string]map[string][]byte) (result map[string]map[string][]byte) {
	switch len(keyValuePairs) {
	case 0:
		return
	case 1:
		result = keyValuePairs[0]
		return
	}

	result = keyValuePairs[0]
	for index := range result {
		neFlag := false
		for i := 1; i < len(keyValuePairs); i++ {
			if _, ok := keyValuePairs[i][index]; !ok {
				neFlag = true
				break
			}
		}
		if neFlag {
			delete(result, index)
		}
	}
	return
}

func orOperation(keyValuePairs ...map[string]map[string][]byte) (result map[string]map[string][]byte) {
	switch len(keyValuePairs) {
	case 0:
		return
	case 1:
		result = keyValuePairs[0]
		return
	}

	result = keyValuePairs[0]
	for i := 1; i < len(keyValuePairs); i++ {
		for index, val := range keyValuePairs[i] {
			if _, ok := result[index]; !ok {
				result[index] = val
			}
		}
	}

	return
}
