var invoke = require('../routes/invoke');
loopCount = 1000;
async function main(func, args, user) {
    const result = await invoke(func, args, user);
    console.log(result.success);
    gatewayDisconnector(result.gateway);
}
startTime = 0;
function loop(num) {
    var func = 'RegisterEmail';
    var user = 'admin';
    startTime = Date.now();
    for(var i = 0; i < num; i++) {
        console.log(`performancetest${i}@gmail.com`);
        var args = JSON.stringify(
                {"args": {"email":`performancetest${i}@gmail.com`}}
                );
        main(func, args, user)
    }

    console.log("DONE....")
}



process.on('exit', calculateTime);

function calculateTime() {
    endTime = Date.now();
    console.log("TIME TAKEN: ", (endTime - startTime)/1000)
}

count = 0;
async function gatewayDisconnector(gateway) {
    count ++;
    if(count == loopCount) {
        await gateway.disconnect();
    }
}

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

async function run() {
    count = 0;
    await sleep(5000);
    loop(loopCount);
}

run();