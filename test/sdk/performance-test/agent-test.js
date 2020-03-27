var http = require("http");
var net = require("net");
var client = new net.Socket();
var agent = new http.Agent({
    keepAlive: true,
    maxSockets: 1000,
    keepAliveMsecs: 30000,
    maxFreeSockets:1000,
    timeout: 1000
});
// var host = 'localhost';
var host = '163.172.30.235';
var port = 3000;
var TxNum = 250;
var sendreg = function (i, host, port) {
    try{
    var startTime = Date.now();
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

    var reqBody = JSON.stringify(
        {"args": {"email":`performancetest${i}@gmail.com`}}
        );
    req.write(reqBody);
    // console.log('Sockets open: ',agent.sockets[`${host}:${port}:`].length)
    req.end();
    // agent.destroy(); 
    }
    catch(e){
        console.log(e);
    }
}

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


var reset = function (host, port) {
    try{
    var options = {
        // host: "13.58.177.170",
        // host: "localhost",
        host,
        port,
        path: "/api/reset",
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

async function main() {
    // for(var i=0; i<10; i++) {
    //     createConnection(host, port);
    // }
    for(var j=0; j<2; j++) {
        console.log("Round:",j+1);
        reset(host, port);
        await sleep(1000);
        for(var i=0; i<TxNum; i++){
            // sendreg(i, host, port);
            if(j!=0){
                sendreg(i, host, port);
            } else {
                connect(host, port);
            }
        }
        for(var i=0; i<5; i++){
            // console.log('========================================');
            await sleep(1000);
            console.log('Free Sockets open: ',agent)
        }
    }

}

main();


function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}
