var http = require("http");
var agent = new http.Agent({
    keepAlive: true,
    maxSockets: 1000,
    keepAliveMsecs: 30000,
    maxFreeSockets: 1000,
    timeout: 10000
});
var activateAgent = false;

var connect = function (host, port) {
    try{
    var options = {
        // host: "13.58.177.170",
        // host: "localhost",
        host,
        port,
        path: "/api/connect",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer token",
            'connection': 'keep-alive',
        },
        agent: agent
    };
    var req = http.request(options, function (res) {   
        try{
        var responseString = "";
        res.on("data", function (data) {
            responseString += data;
            // save all the data from response
            //console.log(responseString + data)
            //responses.push(data);
        });
        res.on("end", function () {
            // print to console when response ends
            console.log(responseString);
        });
        }
        catch(e){
            console.log(e);
        }
    });
    req.end();
    // agent.destroy(); 
    }
    catch(e){
        console.log(e);
    }
}

var sendreg = function (i, totalTPS, cluster, host, port) {
    var startTime = Date.now();
    // var startTime;
    var options = {
        // host: "13.58.177.170",
        // host: "localhost",
        host,
        port,
        path: "/api/registerEmail",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer token",
            'connection': 'keep-alive',
        },
        agent: agent
    };
    if(!activateAgent) {
        options.agent = false;
    }   
    var req = http.request(options, function (res) {   
        var responseString = "";
        res.on("data", function (data) {
            responseString += data;
            // save all the data from response
            //console.log(responseString + data)
            //responses.push(data);
        });
        res.on("end", function () {
            // console.log(responseString);
            console.log("Latency: ", (Date.now() - startTime) / 1000);
            callback(responseString, i, res.statusCode, totalTPS, cluster, startTime, Date.now(), null);
            // print to console when response ends
        });
    });

    var reqBody = JSON.stringify(
        {"args": [{"email":`performancetest${i}@gmail.com`}]}
        );
    // startTime = Date.now();
    req.write(reqBody);
    req.end();
}

var deletereg = function (i, totalTPS, cluster, host, port) {
    var startTime = Date.now();
    // var startTime;
    var options = {
        // host: "13.58.177.170",
        // host: "localhost",
        host,
        port,
        path: "/api/invoke",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer token",
            'connection': 'keep-alive',
        },
        agent: agent
    };
    if(!activateAgent) {
        options.agent = false;
    }   
    var req = http.request(options, function (res) {   
        var responseString = "";
        res.on("data", function (data) {
            responseString += data;
            // save all the data from response
            //console.log(responseString + data)
            //responses.push(data);
        });
        res.on("end", function () {
            // console.log(responseString);
            console.log("Latency: ", (Date.now() - startTime) / 1000);
            callback(responseString, i, res.statusCode, totalTPS, cluster, startTime, Date.now(), null);
            // print to console when response ends
        });
    });

    var reqBody = JSON.stringify(
        {
            "func": "deleteEmail",
            "args": [{"email":`performancetest${i}@gmail.com`}],
            "user": "admin"
        }
        );
    // startTime = Date.now();
    req.write(reqBody);
    req.end();
}

var sendSimpleStore = function (i, totalTPS, cluster, host, port) {
    var startTime = Date.now();
    var options = {
        // host: "52.14.31.90",
        // host: "localhost",
        host,
        port,
        path: "/api/testPerformance",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer token",
        },
        agent: agent,
        // timeout: 600000 // This is connect time out, not read timeout
    };
    if(!activateAgent) {
        options.agent = false;
    }   
    var req = http.request(options, function (res) {   
        var responseString = "";
        res.on("data", function (data) {
            responseString += data;
            // save all the data from response
            //console.log(responseString + data)
            //responses.push(data);
        });
        res.on("end", function () {
            // console.log(responseString);
            console.log("Latency: ", (Date.now() - startTime) / 1000);
            callback(responseString, i, res.statusCode, totalTPS, cluster, startTime, Date.now(), null);
            // print to console when response ends
        });
    });

    var reqBody = JSON.stringify(
        {
            "func": "simpleStore",
            "args":{"test":`performancetest-${i}`},
            "user": "admin"
        }
        );
    // req.on('socket', (s) => { s.setTimeout(300000, () => { s.destroy(); })});
    // req.setTimeout(0, function () {
    //     req.abort();
    //     console.log("timeout");
    // });
    req.on("error", function(e){
        // console.log(`Caught error: ${e}`);
        callback(null, i, null, totalTPS, cluster,null, null, e);
    });
    req.write(reqBody);

    req.end(); 
}

var createEMR = function (i, totalTPS, cluster, host, port) {
    
    var options = {
        // host: "52.14.31.90",
        // host: "localhost",
        host,
        port,
        path: "/api/invoke",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer token",
        },
        agent: agent,
        // timeout: 600000 // This is connect time out, not read timeout
    };
    if(!activateAgent) {
        options.agent = false;
    }   
    var req = http.request(options, function (res) {   
        var responseString = "";
        res.on("data", function (data) {
            responseString += data;
            // save all the data from response
            //console.log(responseString + data)
            //responses.push(data);
        });
        res.on("end", function () {
            // console.log(responseString);
            callback(responseString, i, res.statusCode, totalTPS, cluster);
            // print to console when response ends
        });
    });

    var reqBody = JSON.stringify(
        {
            "func": "createEMR",
            "args": [],
            "user": "prasanths96"
        }
        );
    // req.on('socket', (s) => { s.setTimeout(300000, () => { s.destroy(); })});
    // req.setTimeout(0, function () {
    //     req.abort();
    //     console.log("timeout");
    // });
    req.on("error", function(e){
        // console.log(`Caught error: ${e}`);
        callback(null, i, null, totalTPS, cluster, e);
    });
    req.write(reqBody);

    req.end(); 
}


var query = function (i, totalTPS, cluster, host, port) {
    // {
    //     "func": "getAllTLogs",
    //     "user": "admin"
    // }
}

// Vars
var numSuccess = 0;
var numFail = 0;
var latencySum = 0;

// Callback
function callback(responseString, i, status, totalTPS, cluster, startTime, endTime, err) {
    // console.log(i + ": " + status + ": " + responseString);
    if(status == 200) {
        numSuccess++;
        latencySum += (endTime - startTime) / 1000;
    } else {
        numFail++;
    }
    // console.log(`Success:Fail = ${numSuccess} : ${numFail}`);

    if((numSuccess + numFail) === totalTPS) {
        var averageLatency = latencySum / numSuccess;
        cluster.worker.send({
            responseStat: true,
            numSuccess,
            numFail,
            latency: averageLatency
        });
        cluster.worker.disconnect();
    }    
    if(err) {
        cluster.worker.send({
            responseStat: true,
            numSuccess,
            numFail: totalTPS - numSuccess,
        });
        // process.exit();
    } 
    /*
    if((numSuccess + numFail) === currentTPS) {
        console.log("All response received... Recalculating TPS..");
        if(numFail > 0){
            if (currentTPS < lowestErrorTPS) {
                lowestErrorTPS = currentTPS;
            }
            currentTPS = currentTPS / 2;
            main();
        } else {
            if (currentTPS > highestSuccessTPS) {
                highestSuccessTPS = currentTPS;
            }
            main();
        }
        numSuccess = 0;
        numFail = 0;
    }
    */
}

// HTTP request sender
var sendHTTPRequest = function (options, reqBody, i, totalTPS, cluster) {
    var req = http.request(options, function (res) {   
        var responseString = "";
        res.on("data", function (data) {
            responseString += data;
            // save all the data from response
            //console.log(responseString + data)
            //responses.push(data);
        });
        res.on("end", function () {
            // console.log(responseString);
            callback(responseString, i, res.statusCode, totalTPS, cluster);
            // print to console when response ends
        });
    });
    // req.on('socket', (s) => { s.setTimeout(300000, () => { s.destroy(); })});
    // req.setTimeout(0, function () {
    //     req.abort();
    //     console.log("timeout");
    // });
    req.on("error", function(e){
        // console.log(`Caught error: ${e}`);
        callback(null, i, null, totalTPS, cluster, e);
    });
    req.write(reqBody);

    req.end(); 
}



module.exports = {
    sendreg,
    sendSimpleStore,
    createEMR,
    connect,
    deletereg,
}
