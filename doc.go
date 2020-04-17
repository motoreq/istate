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
		  External application may fetch the structure based on key directly using GetState() API. // This will not be true once original key optimization is done.

	Restrictions:
		- Cannot use type "interface{}" for fields
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
		- Useful ENV for peer container:
			- CORE_LEDGER_STATE_TOTALQUERYLIMIT=1000000   // Query limit
			- CORE_CHAINCODE_EXECUTETIMEOUT=300s		  // To avoid timeout during compaction
    		- CORE_VM_DOCKER_HOSTCONFIG_MEMORY=5368709120‬ // To raise RAM limit of container

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
// Cleanup Errors when have time (error.go)
// fetchCmplx and evalCmplx function is different from others, and needs extra info, try clean it up for symmetry

// Encryption support?

// Enable Load Cache !! Important !!
// Adding prefix to orig key - saves 500ms for 6000 record fetch
// Enable Compaction support

// 1. Include data in index as optional
// 2. Options to activate / deactivate / load cache
// 3. Protobuf
// 4. Trying out GetMultipleStates()
// 5. Fix fieldJSONIndexMap and other meta data
// 6. Load Docs Counter from db
