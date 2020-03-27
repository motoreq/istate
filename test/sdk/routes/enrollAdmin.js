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

// const mainPath = path.resolve(__dirname);
// ccp = makeAbsolutePath(ccp);

const SUCCESS_RESULT = 'Enrollment Successful.'
const FAIL_RESULT = 'Enrollment failed.'

const Enroll = async function() {
    try {
        console.log("wihtin in enrolladmin")
        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.users.medisot.com'].url;
        const ca = new FabricCAServices(caURL);

        // Create a new file system based wallet for managing identities.
        // const walletPath = path.join(mainPath, 'wallet');
        // const wallet = new FileSystemWallet(walletPath);
        const wallet = new S3Wallet();
        // await wallet.import("teest", "daataa")

        var msg;
        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (adminExists) {
            // console.log('An identity for the admin user "admin" already exists in the wallet');
            msg = 'An identity for the admin user "admin" already exists in the wallet';
            return {
                success: true,
                result: SUCCESS_RESULT,
                msg
            }
        }
        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        const identity = X509WalletMixin.createIdentity('medisotUsersMSP', enrollment.certificate, enrollment.key.toBytes());
        await wallet.import('admin', identity);
        // console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
        msg = 'Successfully enrolled admin user "admin" and imported it into the wallet'
        return {
            success: true,
            result: SUCCESS_RESULT,
            msg
        }

    } catch (error) {
        // console.error(`Failed to enroll admin user "admin": ${error}`);
        // process.exit(1);
        var msg = `Failed to enroll admin user "admin": ${error}`
        return {
            success: false,
            result: FAIL_RESULT,
            msg
        }
    }
}

// enrollAdmin();
module.exports = { 
    Enroll
};
