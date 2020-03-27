'use strict'

const express = require('express');
const router = express.Router();
const app = express();
const bodyParser = require('body-parser');
var performanceTester = require('../performance-test/performanceTest');
var host = "localhost";
var port = 3000;
var maxTPS = 1000;
var testTimeout = 120000;
var requestsToSend;
var cluster = require('cluster');
if (cluster.isMaster){
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
        "Origin, X-Requested-With, Content-Type, Accept, Authorization"
        );
        res.setHeader(
        "Access-Control-Allow-Methods",
        "GET, POST, PATCH, DELETE, OPTIONS"
        );
        next();
    });

    // Used by frontend
    router.post("/checkPerformance", (req, res, next) => {
        req.setTimeout(0);
        requestsToSend = req.body.requestsToSend;
        console.log(`Testing performace... No. Of Requests: ${requestsToSend}`);
        
        performanceTester.beginTest(testTimeout, requestsToSend, requestsToSend, callback, host, port); // function(inputTimeLimit, inputStartTPS, maxTPS, callback, host, port)    
        
        function callback(result){
            if (result.success) { 
                res.status(200).json(result);       
            } else {
                res.status(404).json(result); 
            }
        }

    });




    // Util

    function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

} else {
    // console.log(`express worker`)
    performanceTester.beginTest(testTimeout, requestsToSend, requestsToSend, callback, host, port); // function(inputTimeLimit, inputStartTPS, maxTPS, callback, host, port)    
    function callback(result){
        if (result.success) { 
            res.status(200).json(result);       
        } else {
            res.status(404).json(result); 
        }
    }
}


module.exports = router;
