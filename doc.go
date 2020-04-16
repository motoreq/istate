//

/*
	Package iState is used to easily manage perform CRUD operations
	on states/assets in Hyperledger Fabric chaincode.
	It also can be used to easily enable encryption when storing
	states and auto decryption when reading from state db.
	The main purpose of this package is to enable high performance
	Rich Queries when using levelDB as state db.
	Note: To enable high performance queries, it has an indexing mechanism
	that may take extra storage space.

	Requirement:
		- primary key tag ("primary") must be present in structure
		- primary key should be universally unique and it is not handled by this package.
		  It should be handled by the application that imports this package.
		- "istate" tag must be present for the fields in struct that needs to be available for query support.
		- Original marshalled structures will be stored with primary key as the key in db.
		  External application may fetch the structure based on key directly using GetState() API.

	Restrictions:
		- Cannot use the following ascii characters in the struct names or field values:
			- "\000"
			- "\001"
			- "\002"
			- "~" (or) "\176"
			- "\177"
		- Cannot use these in struct field names:
			- "."
			- "*"
			- ".docType"
			- ".keyref"
			- ".value"
			- ".fieldName"
			- (For future) It is good to avoid having field names starting with "." in the structs

	To be noted:
		- CreateState, ReadState, UpdateState and DeleteState functions does not validate if state exists or not.
			Validation must be handled by the external program.
		- Query:
			- If an array/slice/map of elemets needs to be queried, the following applies:
				- eq -> atleast one element in array/slice/map is equal to the value given.
				- neq -> atleast one element in array/slice/map is not equal to the value given.
				- *eq -> all the elements in array/slice/map must be equal to the value given.
				- *neq -> all the elements in array/slice/map must be not equal to the value given.
				Note: Here, map implies, value part of map needs to be queries without knowing the key part of map.
				Eg: "aMap.*":"eq somevalue" as opposed to "aMap.key1": "eq somevalue"

	Known Limitations and Issues:


	Fixed:
		- Indexing: A map with integer / number as key type will still be
    considered as string when indexing.

*/

package istate

import (
// "github.com/hyperledger/fabric/core/mocks/txvalidator"
)

// Debts:
// Cleanup *stub - remove stub from all internal functions
// Cleanup Errors when have time
// fetchCmplx and evalCmplx function is different from others, and needs extra info, try clean it up for symmetry
// Encryption support?
// Clean errors / error.go

// Enable Compaction support
// 1. Include data in index as optional
// 2. Options to activate / deactivate / load cache
// 3. Protobuf
// 4. Trying out GetMultipleStates()
// 5. Fix fieldJSONIndexMap and other meta data
// 6. Load Docs Counter from db
