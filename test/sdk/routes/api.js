

/*
 * Project Name : Medisot
 * Created by Abhijit Ghosh on 01/08/18
 */

'use strict'
// This will automatically get loaded and enrolls admin when start
const enrollAdmin = require('./enrollAdmin'); 
const express = require('express');
const router = express.Router();
const app = express();
const bodyParser = require('body-parser')

const { InvokeHandler, Preload, Reset } = require('./invoke');

const enrollUser = require('./enrollUser');
const checkAuth = require('./checkAuth');
var registerUser = require('./registerUser');
// var query = require('./query');
var newEMR = true;
var fhirToMedisotSaveEMR = require('../utils/fhirToMedisotSaveEMR');
var convertFHIRToMedisot = require('../utils/convertFHIRToMedisot');
// Pras Converter
var fhirconverter = require('../utils/FHIRConverter');
// Request speed checker
var requestNum = 0;
var timeCounter = Date.now();
var currentRate = 0;

const path = require('path');
// const healthcoin = require('../fabtoken/healthcoin');
const http = require('./httpreq');

// Important Initializations
async function init() {
    // Enroll Admin
    var response = await enrollAdmin.Enroll();
    console.log("Init result:")
    console.log(response.result, ':', response.msg);
    if(!response.success) {
        process.exit(1);
    }

    response = await Preload();
    console.log(response.result, ':', response.msg);
    if(!response.success) {
        process.exit(1);
    }

}
init();


// Configuration variables
app.set('secret', 'secret_this_should_be_longer');
app.set('expirationTimeInSeconds', 3600);

// CORS Access
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));
// Prerequisite to avoid CORS errors. (Cross-Origin-Resource-Sharing)
app.use((req, res, next) => {
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader(
      "Access-Control-Allow-Headers",
      "Origin, X-Requested-With, Content-Type, Accept, Authorization", "Username"
    );
    res.setHeader(
      "Access-Control-Allow-Methods",
      "GET, POST, PATCH, DELETE, OPTIONS"
    );
    next();
});


router.use((req, res, next) => {
    requestNum ++;
    var difference = (Date.now() - timeCounter)/1000;
    if (difference > 1) {
        timeCounter = Date.now();
        currentRate = (currentRate + requestNum) / (2 * difference);
        requestNum = 0;
    }
    console.log(`=====================================================================`);
    console.log(`Average No. of request per sec received: ${currentRate}`);
    console.log(`=====================================================================`);
    next();
});

router.post("/connect", (req, res, next) => {
    res.status(200).json({
        msg: 'Connection success. Save this Socket for further use.'
    });
    res.end();
});

router.post("/registerEmail", (req, res, next) => {
    const func = 'RegisterEmail';
    const args = JSON.stringify(req.body.args);
    var user = 'admin';

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, args, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
        res.end();
    }
    main();
});


router.post("/confirmEmail", (req, res, next) => {

    const func = 'ConfirmEmail';
    const args = JSON.stringify(req.body.args);
    var user = 'admin';

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, args, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/registerUser", (req, res, next) => {

    const func = 'RegisterUser';
    const args = JSON.stringify(req.body.args);
    const username = req.body.args[0].username;
    const password = req.body.password;
    const user = 'admin';


    var pass = Buffer.from(`{"${username}":"${password}"}`).toString('base64');
    var transientMap = {
         PASSWORD: pass
    }


    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        var result = await InvokeHandler(func, args, user, transientMap);
        // result = await registerUser(username, password);
        if (result.success) {
            
            var result = await registerUser(username, password);
            if (!result.success) {                
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                });  

            }
            // await sleep(1000);
            // var issueArgs = {
            //     user: username,
            //     quantity: 50,
            // }
            // http.post('/api/issueToken', issueArgs, returnFun);
            // function returnFun(result) {
            //     // result = await healthcoin.issue(username, ["MED", "50"]);
            //     if (!result.success) {                
            //         res.status(404).json({
            //             result: result.result,
            //             msg: result.msg
            //         });  

            //     } else {
                    res.status(200).json({
                        result: result.result,
                        msg: result.msg
                    });      
            //     } 
            // }
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});


router.post("/registerUser_FHIR", (req, res, next) => {

    const func = 'RegisterUser';
    const args = JSON.stringify(req.body.args);
    // ======================ADDED FOR FHIR TO MEDISOT=============================//
    //const args = JSON.stringify(req.body.args);
    //const originalArgs = convertFHIRToMedisot(args)
    //console.log("originalArgs-->>>>>>>>>",originalArgs)

    //var originalArgsJSON = JSON.stringify(originalArgs)

    //=============================================================================//

    const username = req.body.args[0].username;
    const password = req.body.password;
    const user = 'admin';


    var pass = Buffer.from(`{"${username}":"${password}"}`).toString('base64');
    var transientMap = {
         PASSWORD: pass
    }


    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        var result = await InvokeHandler(func, args, user, transientMap);
        // result = await registerUser(username, password);
        if (result.success) {

            var finalObj = {
                "resourceType": "Person",
                "id": originalArgs[0].uniqueGovtID,
                "meta": {
                    "versionId": "1",
                    "lastUpdated": "2019-05-21T14:54:24.142+00:00"
                },
                "identifier": [
                    {
                        "value": originalArgs[0].uniqueGovtID,
                        "assigner": {
                            "display": "uniqueGovtID"
                        }
                    }
                ],
                "name": [
                    {
                        "given": [
                            originalArgs[0].username
                        ]
                    }
                ],
                "telecom": [
                    {
                        "system": "phone",
                        "value": originalArgs[0].phone
                    },
                    {
                        "system": "email",
                        "value": originalArgs[0].email
                    }
                ],
                "gender": originalArgs[0].sex,
                "birthDate": originalArgs[0].dob
            }
            
            const result = await registerUser(username, password);
            if (!result.success) {                
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                });                 
            }

            result = await healthcoin.issue(username, ["MED", "50"]);
            if (!result.success) {                
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                });  

            }
        
            res.status(200).json({
                result: result.result,
                msg: result.msg,
                finalObj : finalObj
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/checkAuth", (req, res, next) => {

    const func = 'CheckAuth';
    const user = req.body.username;
    // const password = req.body.password;

    // var pass = Buffer.from(`${password}`).toString('base64');
    // var transientMap = {
    //      PASSWORD: pass
    // }
    const args = JSON.stringify(req.body);

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        var result = await enrollUser(user, req.body.password);  
        setTimeout(function() { 
            timedOut(result);
        }, 100);
        async function timedOut(result) {
            if (result.success) {
                var result = await InvokeHandler(func, args, user, null, true); // transientMap = null, Query = true
                if (!result.success) {                
                    res.status(404).json({
                        result: result.result,
                        msg: result.msg
                    });                 
                }        
                res.status(200).json({
                    result: result.result,
                    msg: result.msg
                });       
            } else {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                }); 
            }
        }
    }
    main();
});


router.post("/invoke", (req, res, next) => {

    const func = req.body.func;
    var args = JSON.stringify(req.body.args);
    var user = req.body.user;
    console.log(user);
    // if (func === 'shareEMR') {
    //     console.log("inside shareEMR");
    //     var argObj = JSON.parse(args)[0];
    //     var jsonReadFields = Buffer.from(argObj.jsonReadFields).toString('base64');
    //     var jsonWriteFields = Buffer.from(argObj.jsonWriteFields).toString('base64');

    //     console.log(jsonReadFields);
    //     console.log(jsonWriteFields);
    //     argObj.jsonReadFields = jsonReadFields;
    //     argObj.jsonWriteFields = jsonWriteFields;
    //     args = JSON.stringify([argObj]);
    // } else if (func === 'writeEMR') {
    //     var argObj = JSON.parse(args)[0];
    //     var jsonData = Buffer.from(argObj.jsonData).toString('base64');

    //     console.log(jsonData);
    //     argObj.jsonData = jsonData;
    //     args = JSON.stringify([argObj]);
    // }

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, args, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});


// Application API
router.post("/query", (req, res, next) => {

    const func = req.body.func;
    var args = req.body.args;
    const user = req.body.user;
    if(typeof args !== "string") {
        args = JSON.stringify(args)
    }
    console.log(user);
    console.log(`args: ${args}, func: ${func}`);
    async function main() {
        const result = await InvokeHandler(func, args, user, null, true);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/getAccessibleUserInfo", (req, res, next) => {

    const func = "getAccessibleUserInfo";
    var args = "";
    const user = req.body.user;
    if(typeof args !== "string") {
        args = JSON.stringify(args)
    }
    console.log(user);
    console.log(`args: ${args}, func: ${func}`);
    async function main() {
        const result = await query(func, args, user);
        if (result.success) { 
            var output = JSON.parse(result.msg)
            Object.keys(output).forEach(function(key){
                if( output[key].Read != null) {
                    if(output[key].Read.length != 0) {
                        var part3 = output[key].Read[0].split("_")[2]
                        if(part3 == "x") {
                            output[key].Read.shift();
                            if(output[key].Read.length == 0) {
                                output[key].Read = null;
                            }
                        }
                    }
                }
            });
            result.msg = JSON.stringify(output);
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/getSignedProposal", (req, res, next) => {

    const func = req.body.func;
    var args = JSON.stringify(req.body.args);
    var user = req.body.user;
    var transientMap = req.body.transientMap;
    var chaincodeId = req.body.chaincodeId;
    var channelId = req.body.channelId;
    console.log(user);
    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await GetSignedProposal(func, args, user, transientMap, chaincodeId, channelId);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});


const SINGED_PROPOSAL = 'SIGNED_PROPOSAL'

router.post("/queryBySignedProposal", (req, res, next) => {

    var signedProposal = req.body.signedProposal;
    async function main() {
        const result = await InvokeHandler(SINGED_PROPOSAL, signedProposal);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});


// Token

router.get("/listToken", (req, res, next) => {
    const user = req.get('Username');

    console.log(user);
    async function main() {
        var headers = {
            Username: user
        }
        http.get('/api/listToken', headers, returnFun);
        function returnFun(result) { 
            if (result.success) { 
                res.status(200).json({
                    result: result.result,
                    msg: result.msg
                });       
            } else {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                }); 
            }
        }
    }
    main();
});

router.get("/viewTokenBalance", (req, res, next) => {
    const user = req.get('Username');

    console.log(user);
    async function main() {
        var headers = {
            Username: user
        }
        http.get('/api/listToken', headers, returnFun);
        function returnFun(result) { 
            if (result.success) { 
                var tokenList = JSON.parse(result.msg);
                var tokenBalance = 0;
                tokenList.forEach(function(token) {
                    if(token.type === "MED"){
                        tokenBalance +=  parseInt(token.quantity, 10);
                    }
                });
                console.log('Token Balance: ',tokenBalance)
                var returnObj = {};
                returnObj.tokenType = 'MED';
                returnObj.balance = tokenBalance;
                res.status(200).json({
                    result: result.result,
                    msg: JSON.stringify(returnObj)
                });       
            } else {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                }); 
            }
        }
    }
    main();
});

router.post("/issueToken", (req, res, next) => {

    const user = req.body.user;
    var quantity = req.body.quantity;

    if(typeof quantity === 'number') {
        quantity = quantity.toString();
    }
    async function main() {
        // const result = await healthcoin.issue(user, ["MED", quantity]);
        const args = req.body
        // let result;
        http.post('/api/issueToken', args, returnFun);
        function returnFun(result) { 
            if (result.success) { 
                res.status(200).json({
                    result: result.result,
                    msg: result.msg
                });       
            } else {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                }); 
            }
        }
    }
    main();
});

router.post("/claimReward", (req, res, next) => {

    const func = "evaluateReward";
    var args = JSON.stringify(req.body.args);
    var user = req.body.user;
    console.log(user);

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, args, user);
        if (result.success) { 
            let evaluationResult
            var flag = false;
            let respMsg;
            try {
                evaluationResult = JSON.parse(result.msg);
                console.log("EVALRESULT",evaluationResult)
                respMsg = result.msg
            }
            catch (e) {
                flag = true
            }
            if(!evaluationResult.eligible || flag) {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                });
            } else {
                var tokenArgs = {
                    user: user,
                    quantity: evaluationResult.tokens,
                }
                // console.log("TOKENARG", tokenArgs)
                http.post('/api/issueToken', tokenArgs, returnFun);
                async function returnFun(result) {
                    if(result.msg.status != 'SUCCESS' || result.msg.status == undefined) {
                        res.status(404).json({
                            result: result.result,
                            msg: respMsg
                        }); 
                    }
                    var txId = result.msg.txId;        
                    var args = {};
                    args.username = user;
                    args.tokenId = txId;
                    args = JSON.stringify(args);
                    const updateResult = await InvokeHandler("updateRewardTracker", args, user);
                    var statusCode = 404;
                    if(updateResult.success) {
                        statusCode = 200;
                    }
                    res.status(statusCode).json({
                        result: result.result,
                        msg: respMsg
                    });                         
                    
                } 
            }   
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/reclaimReward", (req, res, next) => {

    const func = "reclaimUnissuedReward";
    var args = req.body.args;
    var user = req.body.user;
    console.log(user);

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, args, user);
        if (result.success) { 
            let evaluationResult
            var flag = false;
            let respMsg;
            try {
                evaluationResult = JSON.parse(result.msg);
                console.log("EVALRESULT",evaluationResult)
                respMsg = result.msg
            }
            catch (e) {
                flag = true
            }
            if(!evaluationResult.eligible || flag) {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                });
            } else {
                var tokenArgs = {
                    user: user,
                    quantity: evaluationResult.tokens,
                }
                // console.log("TOKENARG", tokenArgs)
                http.post('/api/issueToken', tokenArgs, returnFun);
                async function returnFun(result) {
                    if(result.msg.status != 'SUCCESS' || result.msg.status == undefined) {
                        res.status(404).json({
                            result: result.result,
                            msg: respMsg
                        }); 
                    }
                    var txId = result.msg.txId;        
                    var args = {};
                    args.username = user;
                    args.tokenId = txId;
                    args = JSON.stringify(args);
                    const updateResult = await InvokeHandler("updateRewardTracker", args, user);
                    var statusCode = 404;
                    if(updateResult.success) {
                        statusCode = 200;
                    }
                    res.status(statusCode).json({
                        result: result.result,
                        msg: respMsg
                    });                         
                    
                } 
            }   
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

// Util

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}















// router.post('/registerUser', (req, res) => {
//     var postBody = req.body
//     registerUser.data.regisUserDetails(postBody, function(message){
//         console.log("message>>>>>>>>>>>>>>>",message)
//         res.json(message)
//     })
// });





// ------------------------ Decode : JWT ---------------------------- //
//var jwt = require('jsonwebtoken');
var expressJWT = require('express-jwt');

// router.use(expressJWT({
//     secret: 'secretKey'
// }));

// function decodeJWT(headers){
//     var retCreds = {
//         username: ""
//        // orgName: ""
//     };
//    // console.debug("headers================",headers)
//     jwt.verify(headers, 'secretKey', function(err, decoded){
//       //  console.debug("decoded------------",decoded);


//         retCreds.username = decoded.user_id;
//       //  retCreds.orgName = decoded.orgName;

//       //  console.debug("=======================", decoded.username, " ================== ", decoded.orgName);
//     });

//     return retCreds;
// }

// ------------------------------ Decode : JWT Ends -------------------// 

// router.get('/register',(req, res, next) =>{
//     res.send('REGISTER');
// });

// router.post('/InvokeHandler', (req, res) => {
//    var postBody = req.body
//    console.log("postBody>>>>>>>>>>>",postBody)
//     //let rootDoc = JSON.parse(JSON.stringify(postBody))
//    console.log("within InvokeHandler js");
//    InvokeHandler.data.submitIntoLedger(postBody, function(message){
//        console.log("message>>>>>>>>>",message);
//        res.json(message)
//    })     

// });


 





// router.post('/authenticate', (req, res) => {
//     var postBody = req.body
//     loginService.data.loginDetails(postBody, function(message){
//        res.json(message);
//        console.log('===================================>>>>> ', message)
//     })
// });


//================ A L L  F H I R  A D D E D  M E T H O D ======================//

router.post("/saveEMR_FHIR", (req, res, next) => {

    const func = "writeEMRAndCreateEncounter";
    //var args = req.body.args;
    var user = req.body.user;
    var patientIdOrEMRId = req.body.patientIdOrEMRId;
    var isDoctor = req.body.isDoctor;
    console.log(user);

    //var args = JSON.stringify(req.body.args)

   // console.log(`args: ${args}, func: ${func}`);

    const args1 = JSON.stringify(req.body.args);
    return new Promise((resolve, reject) => {
        fhirToMedisotSaveEMR.fhirToMedisotSaveEMR(args1, function(message){
            console.log("message>>>>>>>>>",message);
            //res.json(message)
            resolve(message)
            var queryMessage = message
            return queryMessage
        });
    }).then((medisotResult =>{
        console.log("result>>>>>======",medisotResult)
        var medisotJSON = JSON.stringify(medisotResult)
        var args = {patientIdOrEMRId : patientIdOrEMRId, isDoctor:isDoctor, jsonData: medisotJSON}
        console.log("args---->>>",args)
        //var emrID = args[0].emrID
       // console.log("args[0].jsonData---->>>",args[0].jsonData)
       // var jsonData = Buffer.from(args[0].jsonData).toString('base64');

       // args[0].jsonData = jsonData;
        var args1 = JSON.stringify(args)

        console.log(`args: ${args1}, func: ${func}`);

        async function main() {
            let result = await InvokeHandler(func, args1, user);
            if (result.success) {
                let parts = result.msg.split(' ')
                res.status(200).json({
                    result: result.result,
                    msg: result.msg,
                    emrId: parts[parts.length - 1]
                });       
            } else {
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                }); 
            }
        }
        main();

    }));

    
});

router.post("/registerUser_FHIR", (req, res, next) => {

    const func = 'RegisterUser';
    //const args = JSON.stringify(req.body.args);
    // ======================ADDED FOR FHIR TO MEDISOT=============================//
    const args = JSON.stringify(req.body.args);
    const originalArgs = convertFHIRToMedisot(args)
    console.log("originalArgs-->>>>>>>>>",originalArgs)

    var originalArgsJSON = JSON.stringify(originalArgs)
   // console.log("originalArgsJSON-->>>>>>>>>",originalArgsJSON)

    //=============================================================================//

    //const username = req.body.args[0].username;
    const username = originalArgs[0].username;
    const password = req.body.password;
    const user = 'admin';

    console.log("username------",username);
    var pass = Buffer.from(`{"${username}":"${password}"}`).toString('base64');
    var transientMap = {
         PASSWORD: pass
    }


    console.log(`originalArgsJSON::::::: ${originalArgsJSON}, func: ${func}`);

    async function main() {
        var result = await InvokeHandler(func, originalArgsJSON, user, transientMap);
        // result = await registerUser(username, password);
        if (result.success) {
            console.log("with in success")

            var finalObj = {
                "resourceType": "Person",
                "id": originalArgs[0].uniqueGovtID,
                "meta": {
                    "versionId": "1",
                    "lastUpdated": "2019-05-21T14:54:24.142+00:00"
                },
                "identifier": [
                    {
                        "value": originalArgs[0].uniqueGovtID,
                        "assigner": {
                            "display": "uniqueGovtID"
                        }
                    }
                ],
                "name": [
                    {
                        "given": [
                            originalArgs[0].username
                        ]
                    }
                ],
                "telecom": [
                    {
                        "system": "phone",
                        "value": originalArgs[0].phone
                    },
                    {
                        "system": "email",
                        "value": originalArgs[0].email
                    }
                ],
                "gender": originalArgs[0].sex,
                "birthDate": originalArgs[0].dob
            }
            
            const result = await registerUser(username, password);
            if (!result.success) {                
                res.status(404).json({
                    result: result.result,
                    msg: result.msg
                });                 
            }
        
            res.status(200).json({
                result: result.result,
                msg: result.msg,
                finalObj : finalObj
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/shareDoctor_FHIR", (req, res, next) => {

    // {
    //     "args": {"userToShare": "drrajeev", "durationInSeconds": 36000},
    //     "newEMR": true,
    //     "user": "mraditya"
    // }

    const func = "shareEMR";
    var args = req.body.args;
    //var user = req.body.user;
    var userString = req.body.args.criteria;

    //var the_string = "sometext-20202";
    var parts = userString.split('?');
    var userType = parts[0]
    var userID = parts[1]

    var userIdentifer = ""

    if(userType == "Patient"){
        userIdentifer = userID.split('=');
    }

    var user = userIdentifer[1];
    console.log("user------",user)

    var userToShare = req.body.userToShare
    var durationInSeconds = req.body.durationInSeconds

    console.log("userToShare------",userToShare)
    console.log("durationInSeconds------",durationInSeconds)

    var dummyArgs = JSON.stringify([]);
    // Create new EMR
    async function main1() {
        const result = await InvokeHandler("createEMR", dummyArgs, user);
        if (!result.success)       
        {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    var duration = 0;
    if(req.body.allowCreation) {
        main1();
        //newEMR = false;
        duration = 3000;
    }

    // Retrieve user to get EMRKeys
    var EMRKeys = [];
    const main2 = async() => {
        return await query("getUserProfile", undefined, user);        
    }
    var str = "strggrgtgtgtgtgt"
    function timedOut() {
        main2().then(temp => result(temp)).catch(e => console.log(e));
    }
    setTimeout(timedOut, duration);

    function result(result) {
        if (!result.success)       
        {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
            var temp = JSON.parse(result.msg);
            EMRKeys = temp.EMRKeys;
            step2();
    }


    function step2() {
        console.log(EMRKeys);
        // var jsonReadFields = '{"emrID":"","doctorUserName":"","vitals":{"examinerUserName":"","pulse":"","heartRate":"","BP":{"systolic":0,"diastolic":0},"oxygenPercent":0.0,"weight":0.0,"height":0.0,"date":""},"chiefComplaint":"","chiefComplaintHistory":{"familyHistory":"","personalHistory":"","pastHistory":"","currentMedication":"","allergies":""},"physicalExamination":{"examinerUserName":"","P":"","I":"","C":"","K":"","L":"","E":"","date":""},"systemicExamination":{"examinerUserName":"","RS":"","CS":"","GIT":"","NS":"","date":""},"diagnosticTestsAdvised":{"advisorUserName":"","pathology":"","radiology":"","date":""},"diagnosticReports":{"pathology":{"pathologistUserName":"","RBC":"","WBC":"","Hb":"","bloodSugar":"","cholestrol":"","ureaCreatinine":"","ESR":"","date":""},"radiology":{"radiologistUserName":"","radiologyReport":"","radiologyImageHash":"","date":""}},"drugsPrescribed":"","date":"","writerUserName":""}';
        // var jsonWriteFields = '{"emrID":"","doctorUserName":"","vitals":{"examinerUserName":"","pulse":"","heartRate":"","BP":{"systolic":0,"diastolic":0},"oxygenPercent":0.0,"weight":0.0,"height":0.0,"date":""},"chiefComplaint":"","chiefComplaintHistory":{"familyHistory":"","personalHistory":"","pastHistory":"","currentMedication":"","allergies":""},"physicalExamination":{"examinerUserName":"","P":"","I":"","C":"","K":"","L":"","E":"","date":""},"systemicExamination":{"examinerUserName":"","RS":"","CS":"","GIT":"","NS":"","date":""},"diagnosticTestsAdvised":{"advisorUserName":"","pathology":"","radiology":"","date":""},"diagnosticReports":{"pathology":{"pathologistUserName":"","RBC":"","WBC":"","Hb":"","bloodSugar":"","cholestrol":"","ureaCreatinine":"","ESR":"","date":""},"radiology":{"radiologistUserName":"","radiologyReport":"","radiologyImageHash":"","date":""}},"drugsPrescribed":"","date":"","writerUserName":""}';
        
        var jsonReadFields = '{"emrID":"","doctorUserName":"","vitals":{"examinerUserName":"","pulse":"","heartRate":"","BP":{"systolic":0,"diastolic":0},"oxygenPercent":0.0,"weight":0.0,"height":0.0,"date":""},"chiefComplaint":"","chiefComplaintHistory":{"familyHistory":"","personalHistory":"","pastHistory":"","currentMedication":"","allergies":""},"physicalExamination":{"examinerUserName":"","P":"","I":"","C":"","K":"","L":"","E":"","date":""},"systemicExamination":{"examinerUserName":"","RS":"","CS":"","GIT":"","NS":"","date":""},"diagnosticTestsAdvised":"","drugsPrescribed":"","writerUserName":"","timeStamp":""}';
        var jsonWriteFields = '{"emrID":"","doctorUserName":"","vitals":{"examinerUserName":"","pulse":"","heartRate":"","BP":{"systolic":0,"diastolic":0},"oxygenPercent":0.0,"weight":0.0,"height":0.0,"date":""},"chiefComplaint":"","chiefComplaintHistory":{"familyHistory":"","personalHistory":"","pastHistory":"","currentMedication":"","allergies":""},"physicalExamination":{"examinerUserName":"","P":"","I":"","C":"","K":"","L":"","E":"","date":""},"systemicExamination":{"examinerUserName":"","RS":"","CS":"","GIT":"","NS":"","date":""},"diagnosticTestsAdvised":"","drugsPrescribed":"","writerUserName":""}';
        

        jsonReadFields = Buffer.from(jsonReadFields).toString('base64');
        jsonWriteFields = Buffer.from(jsonWriteFields).toString('base64');

        var jsonTempWriteFields = "";
        var newargs = [];

        EMRKeys.forEach(function(emrID, index) {
            if(index === EMRKeys.length - 1) {
                jsonTempWriteFields = jsonWriteFields
                console.log("CHANGED WRITE FIELDS", EMRKeys[index]);
            }
            console.log(jsonTempWriteFields);
            var temp = {
                emrID,
                userToShare: userToShare,
                jsonReadFields,
                jsonWriteFields: jsonTempWriteFields,
                durationInSeconds: durationInSeconds
            }
            newargs.push(temp);
        });

        newargs = JSON.stringify(newargs);

        console.log(newargs);
        console.log(`args: ${args}, func: ${func}`);

        
        main(newargs);
    }

    async function main(newargs) {
        const result = await InvokeHandler(func, newargs, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
});

router.post("/shareLab_FHIR", (req, res, next) => {

    // {
    //     "args": {"userToShare": "vas", "emrID":"emr_prasanths96_0", "durationInSeconds": 36000},
    //     "user": "prasanths96"
    // }

    const func = "shareEMR";
    var args = req.body.args;
    var user = req.body.user;

    var userString = req.body.args.criteria;

    //var the_string = "sometext-20202";
    var parts = userString.split('?');
    var userType = parts[0]
    var userID = parts[1]

    var userIdentifer = ""

    if(userType == "DigonosticReport"){
        userIdentifer = userID.split('=');
    }

    var emrID = userIdentifer[1];
    console.log("emrID------",emrID)

    var userToShare = req.body.userToShare
    var durationInSeconds = req.body.durationInSeconds

    console.log("userToShare------",userToShare)
    console.log("durationInSeconds------",durationInSeconds)
    console.log("user------",user)
    
    // var jsonReadFields = '{"emrID":"","doctorUserName":"","diagnosticTestsAdvised":{"advisorUserName":"","pathology":"","radiology":"","date":""}}';
    var jsonReadFields = '{"emrID":"","doctorUserName":"","diagnosticTestsAdvised":"","timeStamp":""}';
    var jsonWriteFields = '{"diagnosticTestsAdvised":""}';
     
    jsonReadFields = Buffer.from(jsonReadFields).toString('base64');
    jsonWriteFields = Buffer.from(jsonWriteFields).toString('base64');
    var newargs = [];

    var temp = {
        emrID: emrID,
        userToShare: userToShare,
        jsonReadFields,
        jsonWriteFields,
        durationInSeconds: durationInSeconds
    }
    newargs.push(temp);

    newargs = JSON.stringify(newargs);

    console.log(newargs);
    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, newargs, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/sharePharmacy_FHIR", (req, res, next) => {

    // {
    //     "args": {"userToShare": "vas", "emrID": "emr_prasanths96_0", "durationInSeconds": 36000},
    //     "user": "prasanths96"
    // }

    const func = "shareEMR";
    var args = req.body.args;
    var user = req.body.user;

    var userString = req.body.args.criteria;

    //var the_string = "sometext-20202";
    var parts = userString.split('?');
    var userType = parts[0]
    var userID = parts[1]

    var userIdentifer = ""

    if(userType == "Pharmacy"){
        userIdentifer = userID.split('=');
    }

    var emrID = userIdentifer[1];
    console.log("emrID------",emrID)

    var userToShare = req.body.userToShare
    var durationInSeconds = req.body.durationInSeconds

    console.log("userToShare------",userToShare)
    console.log("durationInSeconds------",durationInSeconds)
    console.log("user------",user)

    var jsonReadFields = '{"emrID":"","doctorUserName":"","drugsPrescribed":"","timeStamp":""}';
    var jsonWriteFields = ""    
    
    jsonReadFields = Buffer.from(jsonReadFields).toString('base64');
    jsonWriteFields = Buffer.from(jsonWriteFields).toString('base64');

    var newargs = [];

    var temp = {
        emrID: emrID,
        userToShare: userToShare,
        jsonReadFields,
        jsonWriteFields,
        durationInSeconds: durationInSeconds
    }
    newargs.push(temp);


    newargs = JSON.stringify(newargs);

    console.log(newargs);
    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, newargs, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/getAppointmentByID_FHIR", (req, res, next) => {

    const func = req.body.func;
    var args = req.body.args;
    const user = req.body.user;
    if(typeof args !== "string") {
        args = JSON.stringify(args)
    }
    console.log(user);
    console.log(`args: ${args}, func: ${func}`);
    async function main() {
        const result = await query(func, args, user);
        if (result.success) {
            console.log("result--->>>",result.msg);
            var resultParse = JSON.parse(result.msg);
            console.log("result.msg[0].appointmentId--->>>",resultParse[0].appointmentId);
            var minsDuration = resultParse[0].endTimeSecs - resultParse[0].startTimeSecs
            var finalOBJ = {
                "resourceType": "Bundle",
                "id": resultParse[0].appointmentId,
                "type": "searchset",
                "total": 1,
                "entry": [
                  {
                    "fullUrl": "{{host}}/Appointment/"+resultParse[0].appointmentId,
                    "resource": {
                      "resourceType": "Appointment",
                      "id": resultParse[0].appointmentId,
                      "identifier": [
                        {
                          "value": "123_mediSOT"
                        }
                      ],
                      "status": resultParse[0].status,
                      "reason": [
                        {
                        //   "coding": [
                        //     {
                        //       "system": "http://snomed.info/sct",   //Chidamber will send later
                        //       "code": "413095006"
                        //     }
                        //   ],
                          "text": resultParse[0].reason
                        }
                      ],
                      "description": resultParse[0].reason,
                      "start": resultParse[0].startTimeSecs,
                      "end": resultParse[0].endTimeSecs,
                      "minutesDuration": minsDuration,
                      "created": "2015-12-02",
                      "participant": [
                        {
                          "actor": {
                            "reference": "Patient/"+resultParse[0].patientId,
                            "display": resultParse[0].patientId
                          },
                          "status": resultParse[0].status
                        },
                        {
                          "actor": {
                            "reference": "Location/"+resultParse[0].location,
                            "display": resultParse[0].location
                          },
                          "status": resultParse[0].status
                        },
                        {
                          "actor": {
                            "reference": "Practitioner/"+resultParse[0].doctorId,
                            "display": resultParse[0].doctorId
                          },
                          "status": resultParse[0].status
                        }
                      ]
                    },
                    "search": {
                      "mode": "match"
                    }
                  }
                ]
              } 
            res.status(200).json({
                result: result.result,
                msg: finalOBJ
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/getAllAppointments_FHIR", (req, res, next) => {

    const func = req.body.func;
    var args = req.body.args;
    const user = req.body.user;
    if(typeof args !== "string") {
        args = JSON.stringify(args)
    }
    console.log(user);
    console.log(`args: ${args}, func: ${func}`);
    async function main() {
        const result = await query(func, args, user);
        if (result.success) {
            console.log("result--->>>",result.msg);
            var resultParse = JSON.parse(result.msg);
            var finalOBJArray = [];
            for(var i=0; i<resultParse.length; i++){
                var minsDuration = resultParse[i].endTimeSecs - resultParse[i].startTimeSecs
                var finalOBJ = {
                    "resourceType": "Bundle",
                    "id": resultParse[i].appointmentId,
                    "type": "searchset",
                    "total": 1,
                    "entry": [
                      {
                        "fullUrl": "{{host}}/Appointment/"+resultParse[i].appointmentId,
                        "resource": {
                          "resourceType": "Appointment",
                          "id": resultParse[i].appointmentId,
                          "identifier": [
                            {
                              "value": "123_mediSOT"
                            }
                          ],
                          "status": resultParse[i].status,
                          "reason": [
                            {
                            //   "coding": [
                            //     {
                            //       "system": "http://snomed.info/sct",  //Chidamber will share it later
                            //       "code": "413095006"
                            //     }
                            //   ],
                              "text": resultParse[i].reason
                            }
                          ],
                          "description": resultParse[i].reason,
                          "start": resultParse[i].startTimeSecs,
                          "end": resultParse[i].endTimeSecs,
                          "minutesDuration": minsDuration,
                          "created": "2015-12-02",
                          "participant": [
                            {
                              "actor": {
                                "reference": "Patient/"+resultParse[i].patientId,
                                "display": resultParse[i].patientId
                              },
                              "status": resultParse[i].status
                            },
                            {
                              "actor": {
                                "reference": "Location/"+resultParse[i].location,
                                "display": resultParse[i].location
                              },
                              "status": resultParse[i].status
                            },
                            {
                              "actor": {
                                "reference": "Practitioner/"+resultParse[i].doctorId,
                                "display": resultParse[i].doctorId
                              },
                              "status": resultParse[i].status
                            }
                          ]
                        },
                        "search": {
                          "mode": "match"
                        }
                      }
                    ]
                  } 
                  finalOBJArray.push(finalOBJ);
            }
            res.status(200).json({
                result: result.result,
                msg: finalOBJArray
            });       
        } else {q
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/changeAppointmentStatus_fhir", (req, res, next) => {
    const func = "changeAppointmentStatus";
    var args = req.body;
    var appointmentID = req.body.id;
    console.log("appointmentID----",appointmentID);

    var identiFier = req.body.identifier;
    var participant = req.body.participant;

    var emrID = "";
    var userID = "";
    var _userID = "";
    var status = "";
    var finalUserID = "";

    var finalArg = {};

    for(var i=0; i<identiFier.length; i++){
        emrID = identiFier[i].value
        console.log("emrID---",emrID);
    }
    for(var i=0; i<participant.length; i++){
        userID = participant[i].actor.reference
        console.log("userID---",userID)
        _userID = userID.split("/");
        console.log(" _userID[1]-----", _userID[1])
        finalUserID = _userID[1]
        status = participant[i].status
    }
    finalArg = {"appointmentId":appointmentID,"status":status,"emrId":emrID}
    var finalArgsString = JSON.stringify(finalArg);
    console.log("finalArgsString---",finalArgsString);
    console.log("func---",func);
    console.log("finalUserID---",finalUserID);
    async function main() {
        const result = await InvokeHandler(func, finalArgsString, finalUserID);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/getTransactionLog", (req, res, next) => {

    const func = "getHistoryByKey";
    var args = req.body.args;
    const user = req.body.user;
    if(typeof args !== "string") {
        args = JSON.stringify(args)
    }
    console.log(user);
    console.log(`args: ${args}, func: ${func}`);
    async function main() {
        const result = await query(func, args, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});


// Performance checker

router.post("/testPerformance-old", (req, res, next) => {
    req.setTimeout(0);
    const func = 'RegisterEmail';
    var user = 'admin';

    // console.log(`args: ${args}, func: ${func}`);

    async function main(temp) {
        const args = `[{"email":"performancetest${temp}@gmail.com"}]`;
        const result = await InvokeHandler(func, args, user);
        console.log(`result-${temp} came...`)
        callbackFun(result, temp);
    }
    let start;
    for(var i=0; i<100; i++) {
        start = Date.now();
        main(i)
    }
    //main();

    var outputList = [];
    function callbackFun(result, i){  
        // var millis = Date.now() - start;
        outputList.push(result);
        if(outputList.length == 2) {
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });  
        }
    }
});

router.post("/reset", (req, res, next) => {
    console.log("Inside reset");
    async function main() {
        const result = await Reset();
            res.status(200).json(result);       
    }
    main();
});

// Used by test script
router.post("/testPerformance", (req, res, next) => {
    req.setTimeout(0);
    const func = req.body.func;
    var args = JSON.stringify(req.body.args);
    var user = req.body.user;
    console.log(user);
    if (func === 'shareEMR') {
        console.log("inside shareEMR");
        var argObj = JSON.parse(args)[0];
        var jsonReadFields = Buffer.from(argObj.jsonReadFields).toString('base64');
        var jsonWriteFields = Buffer.from(argObj.jsonWriteFields).toString('base64');

        console.log(jsonReadFields);
        console.log(jsonWriteFields);
        argObj.jsonReadFields = jsonReadFields;
        argObj.jsonWriteFields = jsonWriteFields;
        args = JSON.stringify([argObj]);
    } else if (func === 'writeEMR') {
        var argObj = JSON.parse(args)[0];
        var jsonData = Buffer.from(argObj.jsonData).toString('base64');

        console.log(jsonData);
        argObj.jsonData = jsonData;
        args = JSON.stringify([argObj]);
    }

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await InvokeHandler(func, args, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/invoke-test", (req, res, next) => {

    const func = req.body.func;
    var args = JSON.stringify(req.body.args);
    var user = req.body.user;
    console.log(user);

    console.log(`args: ${args}, func: ${func}`);

    async function main() {
        const result = await invoketest(func, args, user);
        if (result.success) { 
            res.status(200).json({
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

module.exports = router;
