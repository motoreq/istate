var cluster = require('cluster')
async function main(){
    if (cluster.isMaster){
        for(var i=0; i<2; i++){
            var worker = cluster.fork();
            worker.on('message', function(msg) {
                console.log(msg);
            });
        }
        cluster.on('exit', function(worker){
            console.log(`${worker.id} - exited.`)
        })
    } else {
        await callback(cluster);
        cluster.worker.disconnect();
    }
}
main();

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

async function callback(cluster1) {
    await sleep(2000 * cluster1.worker.id);
    cluster1.worker.send(`Worker: ${cluster1.worker.id}`);
}