'use strict'
var firstStartTime
var startTime
var count = 0
const { FileSystemWallet, Gateway, DefaultEventHandlerStrategies } = require('fabric-network')
// const Channel = require('../node_modules/fabric-client/lib/Channel')
const Channel = require('./modified-sdk/Channel')
// const Transaction = require('../node_modules/fabric-network/lib/transaction')

const fs = require('fs')
const path = require('path')
const S3Wallet = require(path.resolve(__dirname, './modified-sdk/Wallet'))
const { makeAbsolutePath } = require('./makeAbsolutePath')
const correctBuffer = require('./jsonBufferCorrection')

const ccpPath = path.resolve(__dirname, 'connection.json')
const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
var ccp = JSON.parse(ccpJSON)

// Important Constants
const CHANNEL_NAME = 'medisotuserschannel'
const CHAINCODE_NAME = 'medisotUsers'
const SUCCESS_RESULT = 'Transaction has been submitted.'
const FAIL_RESULT = 'Transaction failed.'
const SUCCESS_RESULT_QUERY = 'Query Successful.'
const FAIL_RESULT_QUERY = 'Query failed.'
const SINGED_PROPOSAL = 'SIGNED_PROPOSAL'
const PROPOSAL_SUCCESS = 'Proposal generation successful.'
const PROPOSAL_FAIL = 'Proposal generation failed.'

const mainPath = path.resolve(__dirname)
// ccp = makeAbsolutePath(ccp)

var ccps = []
var totNumPeers = Object.keys(ccp.channels.medisotuserschannel.peers).length
for(var i = 0; i < totNumPeers; i ++) {
    var tempccp = JSON.parse(JSON.stringify(ccp,0 ,4))
    for(var j = 0; j < totNumPeers - 1; j ++) {
        var peerToDelete = (i + j) % totNumPeers
        delete tempccp.channels.medisotuserschannel.peers[`peer${peerToDelete}.users.medisot.com`]
        delete tempccp.channels.tokenchannel.peers[`peer${peerToDelete}.users.medisot.com`]
        delete tempccp.organizations.MedisotUsers.peers[`peer${peerToDelete}.users.medisot.com`]
        delete tempccp.peers[`peer${peerToDelete}.users.medisot.com`]        
    }    
    ccps.push(tempccp)
}

var currentCCP = 0

// Create a new file system based wallet for managing identities.
// var walletPath = path.join(mainPath, 'wallet')
// var wallet = new FileSystemWallet(walletPath)
var wallet = new S3Wallet()
// console.log(`Wallet path: ${walletPath}`)

var conObj = []

async function Preload() {
    try {
        var exists = await wallet.exists('admin')
        console.log(`Inside preload(): Exists: `, exists)
        if(!exists) {
            throw new Error("User: admin does not exist in wallet")
        }
        for(var i = 0; i < ccps.length; i ++) {        
            var gateway = new Gateway()
            console.log("Connection:", i)
            await gateway.connect( ccps[i], 
                                    { 
                                        wallet,  
                                        identity: 'admin', 
                                        discovery: { 
                                            enabled: false 
                                        }, 
                                        eventHandlerOptions: {
                                            commitTimeout: 100,
                                            strategy: DefaultEventHandlerStrategies.MSPID_SCOPE_ANYFORTX
                                        } 
                                    })

            // Get the network (channel) our contract is deployed to.
            //const network = await gateway.getNetwork('mychannel')
            var network = await gateway.getNetwork(CHANNEL_NAME)

            // Get the contract from the network.
            // const contract = network.getContract('mycc')
            var contract = network.getContract(CHAINCODE_NAME)
            //await gateway.disconnect()

            var client = gateway.getClient()
            conObj.push({
                network,
                client,
                contract
            })
        }
        return {
            success: true,
            result: `Gateway opening sucessful.`,
            msg: `Successfully opened gateway for Admin.`
        }
    }
    catch (e) {
        return {
            success: false,
            result: `Gateway opening failed.`,
            msg: e
        }
    }
}

// preload()

var userObjList = {}

const _setCurrentCCP = function() {
    currentCCP = (currentCCP + 1) % totNumPeers
}

const _getCurrentClient = function() {
    return conObj[currentCCP].client
}

const _getCurrentContract = function() {
    return conObj[currentCCP].contract
}

const _getCurrentNetwork = function() {
    return conObj[currentCCP].network
}

const _setUserContextToClient = async function(user) {
    try {
        if(!user) {
            return {
                success: true
            }
        }
        if(userObjList[user] !== undefined) {
            _getCurrentClient().setUserContext(userObjList[user], true)
        } else {
            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(user)
            if (!userExists) {
                console.log(`An identity for the user ${user} does not exist in the wallet`)
                throw new Error('User does not exist in wallet')
            }
            var identity = await wallet.export(user)
            // console.log("IDENTITY")
            // console.log(identity)
            // Create a new gateway for connecting to our peer node.
            //const gateway = new Gateway()
            //await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: false } })
            // console.log("CLIENT:")
            // console.log(gateway.getClient())
            
            // console.log(createUserOpts)
            var createUserOpts = {
                username: user,
                mspid: identity.mspId,
                cryptoContent: {
                    privateKeyPEM: identity.privateKey,
                    signedCertPEM: identity.certificate
                },
                skipPersistence: true
            }
            var userObj = await _getCurrentClient().createUser(createUserOpts)
            userObjList[user] = userObj
        }
        var msg = `UserContext set successful.`
        return {
            success: true,
            result: SUCCESS_RESULT,
            msg
        }  
    } catch (e) {
        console.log("SetUserContext Error: ", e)
        return {
            success: false,
            result: FAIL_RESULT,
            msg: e.toString()
        }
    }
}

const _getTransactionObject = function(func, transientMap) {
    var transactionObj = _getCurrentContract().createTransaction(func)
    // Transient
    if (transientMap !== undefined) {
        transactionObj.setTransient(transientMap)
    }
    return transactionObj
}

// const _getTransactionObject = function(func, transientMap) {
//     var transactionObj = new Transaction(_getCurrentContract(), func)
//     console.log(`========================== NEW TRANSACTION ===========================`)
//     console.log(transactionObj)
//     console.log(`========================== FROM CONTRACT =============================`)
//     console.log(_getCurrentContract().createTransaction(func))
//     // Transient
//     if (transientMap !== undefined) {
//         transactionObj.setTransient(transientMap)
//     }
//     return transactionObj
// }

const Invoke = async function(func, args, transientMap) {
    try {
        console.log('Inside Invoke')

        //client.setUserContext(userObj, true)
        // console.log(gateway.getClient())
    
        var transactionObj = _getTransactionObject(func, transientMap)
        var msg = await transactionObj.submit(args)
        msg = msg.toString()
        
        // Submit the specified transaction.
        // await contract.submitTransaction( func, args)   
        
        console.log('Transaction has been submitted')
        console.log(msg.toString())
        // var result = 'Transaction has been submitted'
        // Disconnect from the gateway.
        //await gateway.disconnect()
        
        return {
            success: true,
            result: SUCCESS_RESULT,
            msg
        }               

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`)
        return {
            success: false,
            result: FAIL_RESULT,
            msg: error.toString()
        }
    }
}

const Query = async function(func, args, transientMap) {
    try {
        console.log('Inside Query')

        //client.setUserContext(userObj, true)
        // console.log(gateway.getClient())
    
        var transactionObj = _getTransactionObject(func, transientMap)
        var msg;
        if (args === undefined) {
            console.log("ON TRACK OK!")
            msg = await transactionObj.evaluate('')
        } else {
            msg = await transactionObj.evaluate(args)
        }
        msg = msg.toString()
        
        // Submit the specified transaction.
        // await contract.submitTransaction( func, args)   
        
        console.log(SUCCESS_RESULT_QUERY)
        console.log(msg.toString())
        // var result = 'Transaction has been submitted'
        // Disconnect from the gateway.
        //await gateway.disconnect()
        
        return {
            success: true,
            result: SUCCESS_RESULT_QUERY,
            msg
        }               

    } catch (error) {
        console.error(`${FAIL_RESULT_QUERY}: ${error}`)
        return {
            success: false,
            result: FAIL_RESULT_QUERY,
            msg: error.toString()
        }
    }
}

const Reset = function() {
    firstStartTime = undefined
    count = 0
    return {
        msg: "DONE"
    }
}
// var func = "RegisterEmail"
// var args = "[{\"email\":\"prasa@gmail.com\"}]" // 9p0D
// var func = "ConfirmEmail"  
// var args = "[{\"email\":\"prasa@gmail.com\",\"otp\":\"9p0D\"}]"
// var func = "RegisterUser"
// var args = "[{\"username\":\"prasa\",\"name\":\"Pras\",\"age\":\"22\",\"sex\":\"M\",\"dob\":\"DOB\",\"phone\":\"8876567789\",\"email\":\"prasa@gmail.com\",\"uniqueGovtID\":\"HGTGJU231\"}]"
// var pass = Buffer.from('{"prasa":"password"}').toString('base64')
// var transientMap = {
//     PASSWORD: pass
// }

// invoke(func, args, 'user1')
// invoke(func, args, 'user1', transientMap)

const GetSignedProposal = async function(func, args, user, transientMap, chaincodeId, channelName) {
    try {
        var response = await _setUserContextToClient(user) 
        if(!response.success) {
            throw new Error(response.msg)
        }
        var transactionObj = _getTransactionObject(func, transientMap)
        var newargs = [args]
        var request = transactionObj._buildRequest(newargs)
        request.chaincodeId = chaincodeId
        var channel = new Channel(CHANNEL_NAME, _getCurrentClient())
        var proposal = Channel._buildSignedProposal(request, channelName, channel._clientContext)
        return {
            success: true,
            result: PROPOSAL_SUCCESS,
            msg: JSON.stringify(proposal),
        }
    } catch (e) {
        console.log(`GetSignedProposal Error:`, e)
        return {
            success: false,
            result: PROPOSAL_FAIL,
            msg: e.toString(),
        }
    }
}

const QueryBySignedProposal = async function(signedProposal) {
    try {
        console.log("Inside QueryBySignedProposal")
        var timeout = 5000
        var channel = new Channel(CHANNEL_NAME, _getCurrentClient())
        var result = await channel.queryByChaincodeBySignedProposal(_getCurrentNetwork().getChannel().getPeers(), signedProposal, timeout)
        console.log("Result:", result.toString())
        console.log("Status:", result[0].status)
        if(result[0].status && result[0].status != 200) {
            throw new Error(result[0].toString())
        }
        var msg = result[0].toString()
        return {
            success: true,
            result: SUCCESS_RESULT,
            msg: msg,
        }
    } catch (e) {
        console.log("QueryBySignedProposal Error: ", e)
        return {
            success: false,
            result: FAIL_RESULT,
            msg: e.toString(),
        }
    }
}

const InvokeHandler = async function(func, args, user, transientMap, query) {
    try {
        count++
        var now = Date.now()
        startTime = now
        if (firstStartTime == undefined || now - firstStartTime > 20000) {
            firstStartTime = now
        }

        _setCurrentCCP()
        var response = await _setUserContextToClient(user) 
        if(!response.success) {
            throw new Error(response.msg)
        }
        // Todo Handle
        if(query) {
            return await Query(func, args, transientMap)
        } else {
            switch(func) {
                case SINGED_PROPOSAL:
                    console.log("Inside Switch Signed Proposal")
                    console.log(args)
                    return await QueryBySignedProposal(correctBuffer(JSON.parse(args)))
                default:
                    return await Invoke(func, args, transientMap)     
                    // var signedProposal = await GetSignedProposal(func, args, user, transientMap)
                    // return await QueryBySignedProposal(signedProposal.msg)  
            }     
        }   
        // Handle ends

    } catch(e) {
        return {
            success: false,
            result: FAIL_RESULT,
            msg: e.toString()
        }
    } finally {
        var endTime = Date.now()
        console.log("===============================================")
        console.log("FirstStartTime: ", firstStartTime)
        console.log("StartTime: ", startTime)
        console.log("EndTime: ", endTime)
        console.log("Elapsed: ", (endTime - startTime)/1000)
        console.log("Elapsed to send requests alone: ", (startTime - firstStartTime)/1000)
        console.log("Num requests: ", count)
        console.log("Gateway cache size: ", conObj.length)
        console.log("===============================================")
    }
}

module.exports = { 
    Invoke,
    Query,
    Preload,
    Reset,
    GetSignedProposal,
    InvokeHandler,
}
