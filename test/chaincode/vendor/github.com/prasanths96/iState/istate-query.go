//

package istate

import (
	"encoding/json"
	// "fmt"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"strconv"
	"strings"
)

// Syntax and Keywords
const (
	eq    = "eq"
	neq   = "neq"
	gt    = "gt"
	lt    = "lt"
	gte   = "gte"
	lte   = "lte"
	cmplx = "cmplx"
	or    = "or"
	and   = "and"

	seq    = "*eq"
	sneq   = "*neq"
	sgt    = "*gt"
	slt    = "*lt"
	sgte   = "*gte"
	slte   = "*lte"
	scmplx = "*cmplx"

	splitDot = "."

	asciiLast = "~"
)

type encodedKVs struct {
	eq    map[string]string
	neq   map[string]string
	gt    map[string]string
	lt    map[string]string
	gte   map[string]string
	lte   map[string]string
	cmplx map[string]string

	seq    map[string]string
	sneq   map[string]string
	sgt    map[string]string
	slt    map[string]string
	sgte   map[string]string
	slte   map[string]string
	scmplx map[string]string
}

type querys struct {
	eq    []map[string]interface{}
	neq   []map[string]interface{}
	gt    []map[string]interface{}
	lt    []map[string]interface{}
	gte   []map[string]interface{}
	lte   []map[string]interface{}
	cmplx []map[string]interface{}

	seq    []map[string]interface{}
	sneq   []map[string]interface{}
	sgt    []map[string]interface{}
	slt    []map[string]interface{}
	sgte   []map[string]interface{}
	slte   []map[string]interface{}
	scmplx []map[string]interface{}
}

type efficientKeyType struct {
	enckey       string
	genericField string
	fetchFunc    func(shim.ChaincodeStubInterface, string, *queryEnv) (map[string][]byte, Error)
	relatedEncKV map[string]string
}

type safeKeyFunc struct {
	encKey    string
	fetchFunc func(shim.ChaincodeStubInterface, string, *queryEnv) (map[string][]byte, Error)
}

type queryEnv struct {
	// For Query - u for universal
	ufetchedKVMap map[string][]byte
	ukeyEncKVMap  map[string]map[string]string
}

// Query function is used to
// "," separated objects are considered "or" always
// queryString = [{"docType":"eq USERPROFILE_DOCTYPE", "doctor.whatever": "cmplx or(and(gt bla, lt bla),or(eq a, eq b))", "groups":["cmplx and(neq doctor, neq patient)"]}, {"docType":"eq USERPROFILE_DOCTYPE", "groups":["eq patient"]}]
func (iState *iState) Query(stub shim.ChaincodeStubInterface, queryString string) (finalResult interface{}, iStateErr Error) {
	iStateLogger.Infof("Inside Query")
	defer iStateLogger.Infof("Exiting Query")

	qEnv := &queryEnv{}
	initQueryEnv(qEnv)

	var uQuery []map[string]interface{}
	err := json.Unmarshal([]byte(queryString), &uQuery)
	if err != nil {
		iStateErr = NewError(err, 3002)
		return
	}

	// fmt.Println(uQuery)
	filteredKeys := make([]map[string]map[string]string, len(uQuery), len(uQuery))
	for i := 0; i < len(uQuery); i++ {
		filteredKeys[i], iStateErr = iState.parseAndEvalSingle(stub, uQuery[i], qEnv)
		if iStateErr != nil {
			return
		}
	}

	// Or operation over results
	combinedResults := orOperation(filteredKeys...)

	finalResult = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(iState.structRef)), len(combinedResults), len(combinedResults))
	i := 0
	for key := range combinedResults {
		singleElem := reflect.New(reflect.TypeOf(iState.structRef)).Interface()
		err := json.Unmarshal(qEnv.ufetchedKVMap[key], &singleElem)
		if err != nil {
			iStateErr = NewError(err, 3002)
			return
		}
		finalResult.(reflect.Value).Index(i).Set(reflect.ValueOf(singleElem).Elem())
		i++
	}

	finalResult = finalResult.(reflect.Value).Interface()
	return
}

func (iState *iState) parseAndEvalSingle(stub shim.ChaincodeStubInterface, uQuery map[string]interface{}, qEnv *queryEnv) (filteredKeys map[string]map[string]string, iStateErr Error) {

	// Fields will be declared automatically and make not needed
	querySet := querys{
		// eq:   []map[string]interface{},
		// neq:  []map[string]interface{},
		// gt:   []map[string]interface{},
		// lt:   []map[string]interface{},
		// gte:  []map[string]interface{},
		// lte:  []map[string]interface{},
		// cmplx:[]map[string]interface{},

		// seq:   []map[string]interface{},
		// sneq:  []map[string]interface{},
		// sgt:   []map[string]interface{},
		// slt:   []map[string]interface{},
		// sgte:  []map[string]interface{},
		// slte:  []map[string]interface{},
		// scmplx:[]map[string]interface{},
	}

	for index, val := range uQuery {
		if val, ok := val.(string); !ok {
			iStateErr = NewError(nil, 3003, reflect.TypeOf(val))
			return
		}

		// Trim Space
		queryToEvaluate := strings.TrimSpace(val.(string))
		firstSpaceIndex := strings.Index(queryToEvaluate, " ")
		if firstSpaceIndex == -1 {
			iStateErr = NewError(nil, 3004, queryToEvaluate)
			return
		}
		keyword := queryToEvaluate[:firstSpaceIndex]
		secondPart := queryToEvaluate[firstSpaceIndex+1:]
		// fmt.Println("SECOND PART: ", secondPart)
		// Convert the string value to appropriate type
		var convertedVal interface{}
		convertedVal, iStateErr = convertToRightType(index, secondPart, iState.jsonFieldKindMap, iState.mapKeyKindMap)
		if iStateErr != nil {
			return
		}

		// Dot to Map
		newIndex := ""
		var newVal interface{}
		splitFieldName := strings.Split(index, splitDot)
		switch len(splitFieldName) > 1 {
		case true:
			newIndex = splitFieldName[0]
			newVal = dotsToActualDepth(splitFieldName[1:], convertedVal)
		default:
			newIndex = index
			newVal = convertedVal
		}

		switch keyword {
		case eq:
			// eqQuery[newIndex] = newVal
			querySet.eq = addKeyWithoutOverLap(querySet.eq, newIndex, newVal)
		case neq:
			querySet.neq = addKeyWithoutOverLap(querySet.neq, newIndex, newVal)
		case gt:
			querySet.gt = addKeyWithoutOverLap(querySet.gt, newIndex, newVal)
		case lt:
			querySet.lt = addKeyWithoutOverLap(querySet.lt, newIndex, newVal)
		case gte:
			querySet.gte = addKeyWithoutOverLap(querySet.gte, newIndex, newVal)
		case lte:
			querySet.lte = addKeyWithoutOverLap(querySet.lte, newIndex, newVal)
		case cmplx:
			querySet.cmplx = addKeyWithoutOverLap(querySet.cmplx, newIndex, newVal)
		case seq:
			querySet.seq = addKeyWithoutOverLap(querySet.seq, newIndex, newVal)
		case sneq:
			querySet.sneq = addKeyWithoutOverLap(querySet.sneq, newIndex, newVal)
		case sgt:
			querySet.sgt = addKeyWithoutOverLap(querySet.sgt, newIndex, newVal)
		case slt:
			querySet.slt = addKeyWithoutOverLap(querySet.slt, newIndex, newVal)
		case sgte:
			querySet.sgte = addKeyWithoutOverLap(querySet.sgte, newIndex, newVal)
		case slte:
			querySet.slte = addKeyWithoutOverLap(querySet.slte, newIndex, newVal)
		case scmplx:
			querySet.scmplx = addKeyWithoutOverLap(querySet.scmplx, newIndex, newVal)
		default:
			iStateErr = NewError(nil, 3005, keyword)
			return
		}

	}

	var bestKey string
	var fetchFunc func(shim.ChaincodeStubInterface, string, *queryEnv) (map[string][]byte, Error)
	var queryEncodedKVset encodedKVs

	bestKey, fetchFunc, queryEncodedKVset, iStateErr = iState.getBestEncodedKeyFunc(querySet)
	if iStateErr != nil {
		return
	}

	var fetchedKVMap map[string][]byte
	fetchedKVMap, iStateErr = fetchFunc(stub, bestKey, qEnv)
	keyEncKVMap := make(map[string]map[string]string)
	for key := range fetchedKVMap {
		switch encodedKV, ok := qEnv.ukeyEncKVMap[key]; ok {
		case true:
			keyEncKVMap[key] = encodedKV
		default:
			var tempVar map[string]interface{}
			err := json.Unmarshal(qEnv.ufetchedKVMap[key], &tempVar)
			if err != nil {
				iStateErr = NewError(err, 3019)
				return
			}
			encodedKV, _, _, iStateErr = iState.encodeState(tempVar, key) // query = false
			if iStateErr != nil {
				return
			}
			keyEncKVMap[key] = encodedKV
		}

	}

	evalAndFilterEq(stub, queryEncodedKVset.eq, keyEncKVMap)

	//resultEq, iStateErr = iState.evaluateEq(stub, eqQuery)
	// iState.evaluateNeq(stub, neqQuery)
	// iState.evaluateGt(stub, gtQuery)
	// iState.evaluateLt(stub, ltQuery)
	// iState.evaluateGte(stub, gteQuery)
	// iState.evaluateLte(stub, lteQuery)
	// iState.evaluateCmplx(stub, cmplxQuery)

	// iState.evaluateSeq(stub, seqQuery)
	// iState.evaluateSneq(stub, sneqQuery)
	// iState.evaluateSgt(stub, sgtQuery)
	// iState.evaluateSlt(stub, sltQuery)
	// iState.evaluateSgte(stub, sgteQuery)
	// iState.evaluateSlte(stub, slteQuery)
	// iState.evaluateScmplx(stub, scmplxQuery)

	// and operation between fields

	filteredKeys = keyEncKVMap
	return
}

func convertToRightType(fieldName string, toConvert string, jsonFieldKindMap map[string]reflect.Kind, mapKeyKindMap map[string]reflect.Kind) (convertedVal interface{}, iStateErr Error) {

	splitFieldName := strings.Split(fieldName, splitDot)
	if len(splitFieldName) == 0 {
		iStateErr = NewError(nil, 3010)
		return
	}
	// curField := fieldName
	curField := splitFieldName[0]
	nextIndex := 1
	for {
		kind, ok := jsonFieldKindMap[curField]
		if !ok {
			iStateErr = NewError(nil, 3016, curField)
			return
		}
	SpecialFlow:
		switch kind {
		case reflect.Array, reflect.Slice:
			if len(splitFieldName) <= nextIndex {
				iStateErr = NewError(nil, 3015, fieldName)
				return
			}
			curField = curField + splitDot + star
			nextIndex++
			continue
		case reflect.Struct:
			if len(splitFieldName) <= nextIndex {
				iStateErr = NewError(nil, 3015, fieldName)
				return
			}
			curField = curField + splitDot + splitFieldName[nextIndex]
			nextIndex++
			continue
		case reflect.Map:
			// Is field only to be searched, or value too
			// If this is the last index
			switch len(splitFieldName) == nextIndex {
			case true:
				// Notice kind and ok are changed
				kind, ok = mapKeyKindMap[curField]
				if !ok {
					iStateErr = NewError(nil, 3016, curField)
					return
				}
				goto SpecialFlow
			default:
				if len(splitFieldName) <= nextIndex {
					iStateErr = NewError(nil, 3015, fieldName)
					return
				}
				//prevField := curField
				curField = curField + splitDot + star
				// kind, ok = jsonFieldKindMap[curField]
				//curField = prevField + splitDot + splitFieldName[nextIndex]
				nextIndex++
			}
			continue
		default: // Primitive type
			convertedVal, iStateErr = convertToPrimitiveType(toConvert, kind)
			if iStateErr != nil {
				return
			}
		}
		break
	}
	return
}

func convertToPrimitiveType(toConvert string, kind reflect.Kind) (convertedVal interface{}, iStateErr Error) {
	// When generating table for query, * need not be converted
	if toConvert == star {
		convertedVal = star
		return
	}
	var err error
	switch kind {
	case reflect.Bool:
		convertedVal, err = strconv.ParseBool(toConvert)
		if err != nil {
			iStateErr = NewError(err, 3009)
			return
		}
	case reflect.Int:
		// fmt.Println("Atoi: Trying to Convert:", toConvert)
		convertedVal, err = strconv.Atoi(toConvert)
		if err != nil {
			iStateErr = NewError(err, 3011)
			return
		}
	case reflect.Int8:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int8(convertedVal.(int64))
	case reflect.Int16:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int16(convertedVal.(int64))
	case reflect.Int32:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
		convertedVal = int32(convertedVal.(int64))
	case reflect.Int64:
		convertedVal, err = strconv.ParseInt(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3012)
			return
		}
	case reflect.Uint:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint(convertedVal.(uint64))
	case reflect.Uint8:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint8(convertedVal.(uint64))
	case reflect.Uint16:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint16(convertedVal.(uint64))
	case reflect.Uint32:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
		convertedVal = uint32(convertedVal.(uint64))
	case reflect.Uint64:
		convertedVal, err = strconv.ParseUint(toConvert, 10, 64)
		if err != nil {
			iStateErr = NewError(err, 3013)
			return
		}
	case reflect.Float32:
		convertedVal, err = strconv.ParseFloat(toConvert, 64)
		if err != nil {
			iStateErr = NewError(err, 3014)
			return
		}
		convertedVal = float32(convertedVal.(float64))
	case reflect.Float64:
		convertedVal, err = strconv.ParseFloat(toConvert, 64)
		if err != nil {
			iStateErr = NewError(err, 3014)
			return
		}
	case reflect.String:
		convertedVal = toConvert
	default:
		iStateErr = NewError(nil, 3017, kind)
		return
	}
	return
}

func (iState *iState) getBestEncodedKeyFunc(querySet querys) (encodedKey string, fetchFunc func(shim.ChaincodeStubInterface, string, *queryEnv) (map[string][]byte, Error), encodedKVSet encodedKVs, iStateErr Error) {
	encodedKVSet = encodedKVs{
		eq:    make(map[string]string),
		neq:   make(map[string]string),
		gt:    make(map[string]string),
		lt:    make(map[string]string),
		gte:   make(map[string]string),
		lte:   make(map[string]string),
		cmplx: make(map[string]string),

		seq:    make(map[string]string),
		sneq:   make(map[string]string),
		sgt:    make(map[string]string),
		slt:    make(map[string]string),
		sgte:   make(map[string]string),
		slte:   make(map[string]string),
		scmplx: make(map[string]string),
	}

	tree := btree.NewWithIntComparator(3)
	safe := &safeKeyFunc{}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.eq, encodedKVSet.eq, iState.fetchEq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.neq, encodedKVSet.neq, iState.fetchNeq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.gt, encodedKVSet.gt, iState.fetchGt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.lt, encodedKVSet.lt, iState.fetchLt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.gte, encodedKVSet.gte, iState.fetchGte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.lte, encodedKVSet.lte, iState.fetchLte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.cmplx, encodedKVSet.cmplx, iState.fetchCmplx, tree, safe)
	if iStateErr != nil {
		return
	}

	iStateErr = iState.generateEncKeysAndAddToTree(querySet.seq, encodedKVSet.seq, iState.fetchSeq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.sneq, encodedKVSet.sneq, iState.fetchSneq, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.sgt, encodedKVSet.sgt, iState.fetchSgt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.slt, encodedKVSet.slt, iState.fetchSlt, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.sgte, encodedKVSet.sgte, iState.fetchSgte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.slte, encodedKVSet.slte, iState.fetchSlte, tree, safe)
	if iStateErr != nil {
		return
	}
	iStateErr = iState.generateEncKeysAndAddToTree(querySet.scmplx, encodedKVSet.scmplx, iState.fetchSlte, tree, safe)
	if iStateErr != nil {
		return
	}

	// fmt.Println("BEST KEY: ", tree.LeftValue(), tree.LeftKey())
	// Remove * from encodedKeys and replaces with ""
	//removeStarFromKeys(encodedKeyVal)
	// fmt.Println("Updated encodedKeyVal:", encodedKVSet.eq)
	leftNodeVal := tree.LeftValue()
	switch leftNodeVal == nil {
	case true:
		encodedKey = removeStarFromKey(safe.encKey)
		fetchFunc = safe.fetchFunc
	default:
		encodedKey = removeStarFromKey(leftNodeVal.(efficientKeyType).enckey)
		fetchFunc = leftNodeVal.(efficientKeyType).fetchFunc
	}
	return
}

func hasLessStars(key1 string, key2 string) (ok bool) {
	count1 := strings.Count(key1, "*")
	count2 := strings.Count(key2, "*")
	if count2 < count1 {
		ok = true
	}
	return
}

func (iState *iState) generateEncKeysAndAddToTree(query []map[string]interface{}, encKVMap map[string]string, fetchFunc func(shim.ChaincodeStubInterface, string, *queryEnv) (map[string][]byte, Error), tree *btree.Tree, safe *safeKeyFunc) (iStateErr Error) {
	keyref := ""
	for i := 0; i < len(query); i++ {
		var encKeyDocNameMap map[string]string
		var encodedKeyVal map[string]string
		encodedKeyVal, _, encKeyDocNameMap, iStateErr = iState.encodeState(query[i], keyref, true)
		if iStateErr != nil {
			return
		}
		// fmt.Println("Current: encodedKeyVal", encodedKeyVal)
		for index, val := range encodedKeyVal {

			genericField := encKeyDocNameMap[index]
			if numDocs, ok := iState.readDocsCounter(genericField); ok {
				addToTree(tree, index, genericField, numDocs, fetchFunc, encKVMap)
			}

			encKVMap[index] = val
			if safe.fetchFunc == nil {
				safe.encKey = index
				safe.fetchFunc = fetchFunc
			}
		}
	}
	return
}

func addToTree(tree *btree.Tree, encKey string, genericField string, numDocs int, fetchFunc func(shim.ChaincodeStubInterface, string, *queryEnv) (map[string][]byte, Error), relatedEncKV map[string]string) {
	efficientKey := efficientKeyType{
		enckey:       encKey,
		genericField: genericField,
		fetchFunc:    fetchFunc,
		relatedEncKV: relatedEncKV,
	}
	switch val, ok := tree.Get(numDocs); ok {
	case true:
		if ok := hasLessStars(val.(efficientKeyType).enckey, efficientKey.enckey); ok {
			tree.Put(numDocs, efficientKey)
		}
	default:
		tree.Put(numDocs, efficientKey)
	}
}
