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
