/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';
const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const S3Wallet = require(path.resolve(__dirname, './modified-sdk/Wallet'))
const { makeAbsolutePath } = require('./makeAbsolutePath');

const ccpPath = path.resolve(__dirname, 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
var ccp = JSON.parse(ccpJSON);

const mainPath = path.resolve(__dirname);
// ccp = makeAbsolutePath(ccp);

// Important Constants
const SUCCESS_RESULT = 'Transaction has been submitted.'
const FAIL_RESULT = 'Transaction failed.'

const registerUser = async function main(user, pass) {
    try {
        console.log("within in register user")
        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.users.medisot.com'].url;
        const ca = new FabricCAServices(caURL);
        
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(mainPath, 'wallet');
        const wallet = new S3Wallet(walletPath);

        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(user);
        if (userExists) {
            console.log(`An identity for the user ${user} already exists in the wallet`);
            return {
                success: false,
                result: FAIL_RESULT,
                msg: `An identity for the user ${user} already exists in the wallet`
            }
        }

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (!adminExists) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return {
                success: false,
                result: FAIL_RESULT,
                msg: 'An identity for the admin user "admin" does not exist in the wallet'
            }
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: false } });

        // Get the CA client object from the gateway for interacting with the CA.
        const adminIdentity = gateway.getCurrentIdentity();

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({ affiliation: 'org1.department1', maxEnrollments:0, enrollmentID: user, enrollmentSecret: pass, role: 'client', attrs: [{name:'medisot.username', value: user, ecert: true}] }, adminIdentity);
        const enrollment = await ca.enroll({ enrollmentID: user, enrollmentSecret: secret });
        const userIdentity = X509WalletMixin.createIdentity('medisotUsersMSP', enrollment.certificate, enrollment.key.toBytes());
        await wallet.import(user, userIdentity);

        console.log(`Successfully registered and enrolled user ${user} and imported it into the wallet`);
        return {
            success: true,
            result: SUCCESS_RESULT,
            msg: `Successfully registered and enrolled user ${user} and imported it into the wallet`
        }

    } catch (error) {
        console.error(`Failed to register user ${user}: ${error}`);
        return {
            success: false,
            result: FAIL_RESULT,
            msg: `Failed to register user ${user}: ${error}`
        }
    }
}

//registerUser('prasanths96', 'password');
module.exports = registerUser;
