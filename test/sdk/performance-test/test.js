//var host = "localhost";
var host = "163.172.30.235";
//var host = "52.15.197.224";
var port = 3000;
var TPS = 10000000;
var performanceTest = require('./performanceTest');
async function check(){
    performanceTest.beginTest(0, TPS, TPS, callback, host, port); // function(inputTimeLimit, inputStartTPS, maxTPS, callback, host, port)
    // performanceTest.customLoadgen(50, 10000, host, port) // (rate, numReq, host, port)
}

// async function main(){
// for(var i=0; i<1; i++){
//         check();
// //      await sleep(160000);
// }
// }

// main();
check();
function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
}
function callback(res){
    console.log(res);
}

