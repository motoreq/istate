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

* Initialize vendor folder in your chaincode directory using ```govendor init```

* Get the dependent packages using the following commands:

  1. ```govendor fetch github.com/prasanths96/iState```
  2. ```govendor fetch github.com/bluele/gcache```
  3. ```govendor fetch github.com/emirpasic/gods```
  4. ```govendor fetch github.com/prasanths96/gostack```

#### Mannual method (No tools needed)

* Clone this repository in a preferred location using ```git clone https://github.com/prasanths96/iState.git```.

* Copy the ```.go``` files in this repo and paste inside ```yourchaincode/vendor/github.com/prasanths96/istate/``` 
*(Note: No need to copy files inside folders.)* 

* Copy the vendor folder in this repo and merge it with ```yourchaincode/vendor```

Thats all, you're read to import this package in your chaincode.

### Reference

`godoc` or https://godoc.org/github.com/prasanths96/iState

## License <a name="license"></a>

iState Project source code files are made available under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file. 
