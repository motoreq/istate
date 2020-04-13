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

	Restrictions:
		- Cannot use the following ascii characters in the struct names or field values:
			- "\000"
			- "\001"
			- "\002"
			- "~" (or) "\176"
			- "\177"

	Known Limitations and Issues:


	Fixed:
		- Indexing: A map with integer / number as key type will still be
    considered as string when indexing.

*/

package istate

// Debts:
// Encryption support?
// Clean errors / error.go
// 1. Include data in index as optional
// 2. Options to activate / deactivate / load cache
// 3. Protobuf
// 4. Trying out GetMultipleStates()
// 5. Fix fieldJSONIndexMap and other meta data
// 6. Load Docs Counter from db
