## Under Development

### iState [![GoDoc](https://godoc.org/github.com/prasanths96/iState?status.svg)](https://godoc.org/github.com/prasanths96/iState) [![Go Report Card](https://goreportcard.com/badge/github.com/prasanths96/iState)](https://goreportcard.com/report/github.com/prasanths96/iState)


iState is a state management package for Hyperledger fabric chaincode. It can be used to enable encryption of states and high performance rich queries on leveldb.

### Features

* Rich Query in levelDB.

* On an average case, query is **~7 times** faster than CouchDB's Rich Query with Index enabled.

* In-memory caching using ARC (Adaptive Replacement Cache) algorithm.

* Cache consistency is maintained. Data returned by query is always consistent.

* Easy to use.

### Installation

#### Using govendor

* Initialize vendor folder in the chaincode directory using ```govendor init```

* Get the dependent packages using the following commands:

  1. ```govendor fetch github.com/prasanths96/iState```
  2. ```govendor fetch github.com/bluele/gcache```
  3. ```govendor fetch github.com/emirpasic/gods```
  4. ```govendor fetch github.com/prasanths96/gostack```

#### Mannual method (No tools needed)

* Clone this repository in a preferred location using ```git clone https://github.com/prasanths96/iState.git```.

* Copy the ```.go``` files in this repo and paste inside ```chaincode/vendor/github.com/prasanths96/istate/``` 
*(Note: No need to copy files inside folders.)* 

* Copy the vendor folder in this repo and merge it with ```chaincode/vendor```

Thats all, iState is ready to be imported in the chaincode.

### Example

#### Adding tags to struct

The following tags must be added only to the struct types which is getting stored in state db.

- ```primary``` tag is used to denote the primary key / id in the struct. This field **must** contain universal unique value and is handled externally.
- ```istate``` tag is used to denote the fields that is allowed to be queried. It is recommended to add this tag to all fields.
- Value of ```istate``` tag must be universally unique with other structs in the chaincode. Recommended format: ```<structname>_<fieldname>```

```go
type TestStruct struct {
	ID      string  `json:"id" istate:"TestStruct_id" primary:"true"`
	AString string  `json:"docType" istate:"TestStruct_aString"`
	AnInt   int64   `json:"anInt" istate:"TestStruct_anInt"`
}
```

#### Init

```go
package main 

import (
  "github.com/prasanths96/istate"
)

type TestSmartContract struct {
	TestStructiState istate.Interface
}

// Init initializes chaincode.
func (sc *TestSmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := sc.init()
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (sc *TestSmartContract) init() error {
	iStateOpt := istate.Options{
		CacheSize:             1000000,
		DefaultCompactionSize: 10000,
	}
	TestStructiState, err := istate.NewiState(TestStruct{}, iStateOpt)
	if err != nil {
		return err
	}
	sc.TestStructiState = TestStructiState
	return nil
}
```

### Reference

`godoc` or https://godoc.org/github.com/prasanths96/iState

## License <a name="license"></a>

iState Project source code files are made available under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file. 
