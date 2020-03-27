/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
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

const enrollUser = async function main(user, pass) {
    try {
        console.log("wihtin in enrolluser")
        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.users.medisot.com'].url;

        const ca = new FabricCAServices(caURL);
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(mainPath, 'wallet');
        const wallet = new S3Wallet();
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the  user.
        const userExists = await wallet.exists(user);
        var alreadyInWallet = false;
        if (userExists) {
            console.log(`An identity for the  user ${user} already exists in the wallet`);
            alreadyInWallet = true;
        }
        // Enroll the user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: user, enrollmentSecret: pass });
        
        if (!alreadyInWallet) {        
            // Enroll the user, and import the new identity into the wallet.
            const identity = X509WalletMixin.createIdentity('medisotUsersMSP', enrollment.certificate, enrollment.key.toBytes());
            await wallet.import(user, identity);
            console.log(`Successfully enrolled user ${user} and imported it into the wallet`);
        }
        return {
            success: true,
            enrollment: enrollment,
            result: SUCCESS_RESULT,
            msg: `Successfully enrolled user ${user} and imported it into the wallet`
        };

    } catch (error) {
        console.error(`Failed to enroll user ${user}: ${error}`);
        return {
            success: false,
            result: FAIL_RESULT,
            msg: `Failed to enroll user ${user}: ${error}`
        };
    }
}

//enrollUser('user3', 'userpw');
module.exports = enrollUser;
