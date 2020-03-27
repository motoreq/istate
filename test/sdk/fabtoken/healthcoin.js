'use strict';

const path = require('path');
const Fabric_Client = require(path.resolve(__dirname, 'node_modules/fabric-client'));
const util = require('util');
const os = require('os');
const fs = require('fs-extra');

const { makeAbsolutePath } = require(path.resolve(__dirname, '../routes/makeAbsolutePath'));
const ccpPath = path.resolve(__dirname, '../routes/connection.json');
const walletPath = path.resolve(__dirname, '../routes/wallet/');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
var ccp = JSON.parse(ccpJSON);

const channel_name = "tokenchannel"
ccp = makeAbsolutePath(ccp);

// Important Constants
const SUCCESS_RESULT = 'Token transaction has been submitted.'
const FAIL_RESULT = 'Token transaction failed.'

// console.log(ccp);


// Create fabric client, channel, orderer, and peer instances.
// These are needed for SDK to invoke token operations.
// createFabricClient();
async function createFabricClient() {
	// fabric client instance
    // starting point for all interactions with the fabric network
    const fabric_client = new Fabric_Client();
    // const fabric_client = Fabric_Client.loadFromConfig(ccp);
    fabric_client.loadFromConfig(ccp);
    const channel = fabric_client.getChannel(channel_name);
    // console.log(channel);
    // const fabric_client = new Fabric_Client();

	// -- channel instance to represent the ledger
	// const channel = fabric_client.newChannel(channel_name);
	// console.log(' Created client side object to represent the channel');
    console.log(channel);
	// // // -- peer instance to represent a peer on the channel
	// const peer = fabric_client.newPeer('grpcs://localhost:8051');
	// console.log(' Created client side object to represent the peer');

	// // // -- orderer instance to reprsent the channel's orderer
	// const orderer = fabric_client.newOrderer('grpcs://localhost:7050')
	// console.log(' Created client side object to represent the orderer');

	// // // add peer and orderer to the channel
	// channel.addPeer(peer);
	// channel.addOrderer(orderer);
    // console.log("Here:", fabric_client._network_config._network_config.channels);
	return {client: fabric_client, channel: channel};
}

// Issue token to the user with args [type, quantity]
// It uses "admin" to issue tokens, but other users can also issue tokens as long as they have the permission.
const issue = async function issue(userName, args) {
	try {
		console.log('Start token issue with args ' + args);
		
		// create fabric client and related instances
		// starting point for all interactions with the fabric network
		const {client, channel} = await createFabricClient();

		let admin = await createUser('admin');
		let user = await createUser(userName);

		await client.setUserContext(admin, true);

		const tokenClient = client.newTokenClient(channel, 'peer0.patient.com');

		// build the request to issue tokens to the user
		const txId = client.newTransactionID();
		const param = {
			owner: user.getIdentity().serialize(),
			type: args[0],
			quantity: args[1]
		};
		const request = {
			params: [param],
			txId: txId,
		};

		var result = await tokenClient.issue(request);
		result.txId = txId._transaction_id;
		return {
			success: true,
			result: SUCCESS_RESULT,
			msg: result,
		}
	}
	catch (e) {
		console.log(e);
		return {
            success: false,
            result: FAIL_RESULT,
            msg: e.toString(),
        }
	}
}

// List tokens for the user
const list = async function list(userName) {
	try {
		// create fabric client and related instances
		// starting point for all interactions with the fabric network
		const {client, channel} = await createFabricClient();
		let user = await createUser(userName);		
		await client.setUserContext(user, true);
		const tokenClient = client.newTokenClient(channel, 'peer0.patient.com');
		var result = await tokenClient.list()
		console.log(result);

		return {
			success: true,
			result: SUCCESS_RESULT,
			msg: JSON.stringify(result),
		}
	}
	catch(e) {
		console.log(e);
		return {
            success: false,
            result: FAIL_RESULT,
            msg: e.toString(),
        }
	}
}

// Create user by loading crypto files
// createUser("prasanths96");
async function createUser(user) {
	// This sample application will read user idenitity information from
	// pre-generated crypto files and create users. It will use a client object as
	// an easy way to create the user objects from known cyrpto material.

	const client = new Fabric_Client();

    // load user
    let userWalletPath = path.join(walletPath, user);
    let file_path = path.join(userWalletPath, user);
    console.log(file_path);
    let idfile = fs.readFileSync(file_path);
    // console.log(JSON.parse(idfile.toString()).enrollment.identity.certificate);
    // console.log(JSON.parse(idfile.toString()).enrollment.signingIdentity);
    let identityObj = JSON.parse(idfile.toString());
    let fileName = identityObj.enrollment.signingIdentity;
    let privKeyFileName = fileName + '-priv';
    let keyPEM = fs.readFileSync(path.join(userWalletPath, privKeyFileName)).toString();
    console.log(keyPEM);
    
	// certPath = path.join(__dirname, '../../basic-network/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts');
    // certPEM = readAllFiles(certPath)[0];
    let certPEM = identityObj.enrollment.identity.certificate;
    console.log(certPEM);

	let user_opts = {
		username: user,
		mspid: 'patientMSP',
		skipPersistence: true,
		cryptoContent: {
			privateKeyPEM: keyPEM,
			signedCertPEM: certPEM
		}
	};
	const user1 = await client.createUser(user_opts);

	return user1;
}

function readAllFiles(dir) {
	const files = fs.readdirSync(dir);
	const certs = [];
	files.forEach((file_name) => {
		const file_path = path.join(dir, file_name);
		const data = fs.readFileSync(file_path);
		certs.push(data);
	});
	return certs;
}

let username = "prasanths96";
let args = ["MED", "50"];
// start(username, args);
// list("prasanths96");

async function start(func, username, args) {
	try {
        switch(func) {
			case "issue": await issue(username, args);
			break;
			case "list": await list(username);
			break;
			case "default": console.log("Invalid");
		}

	} catch(error) {
		console.log('Problem with fabric token ::'+ error.toString());
		process.exit(1);
	}
};
start(process.argv[2], process.argv[3], process.argv[4]);


module.exports = {
	issue,
	list,
}