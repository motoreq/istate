// Copyright 2020 <>. All rights reserved.

/*
	Package iState is used to easily manage perform CRUD operations
	on states/assets in Hyperledger Fabric chaincode.
	It also can be used to easily enable encryption when storing
	states and auto decryption when reading from state db.
	The main purpose of this package is to enable high performance
	Rich Queries when using levelDB as state db.
	Note: To enable high performance queries, it has an indexing mechanism
	that may take extra storage space.

	Known Limitations and Issues:
	Indexing: A map with integer / number as key type will still be
	considered as string when indexing.

*/

package istate
