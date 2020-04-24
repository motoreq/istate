/*
	Copyright 2020 Motoreq Infotech Pvt Ltd

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
