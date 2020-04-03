
// func getStateByRange(stub shim.ChaincodeStubInterface, startKey string, endKey string) (result map[string][]byte, iStateErr Error) {
// 	result = make(map[string][]byte)
// 	iterator, err := stub.GetStateByRange(startKey, endKey)
// 	if err != nil {
// 		iStateErr = NewError(err, 3006)
// 		return
// 	}
// 	defer iterator.Close()
// 	for i := 0; iterator.HasNext(); i++ {
// 		iteratorResult, err := iterator.Next()
// 		if err != nil {
// 			iStateErr = NewError(err, 3007)
// 			return
// 		}
// 		keyRef := string(iteratorResult.GetValue())
// 		// If key already present in result, then can avoid re-getstate()
// 		if _, ok := result[keyRef]; ok {
// 			continue
// 		}

// 		valBytes, err := stub.GetState(keyRef)
// 		if err != nil {
// 			iStateErr = NewError(err, 3008)
// 			return
// 		}
// 		result[keyRef] = valBytes
// 	}
// 	return
// }