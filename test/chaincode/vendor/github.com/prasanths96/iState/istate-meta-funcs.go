//

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
