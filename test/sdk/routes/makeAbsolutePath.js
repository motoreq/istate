const path = require('path');

const makeAbsolutePath = function(ccp) {
    // Change relative paths in connection to absolute paths
    var orderers = Object.keys(ccp.orderers);
    var peers = Object.keys(ccp.peers);
    // Changing orderer relative paths
    for (var i = 0; i < orderers.length; i++) {
        var relativePath = ccp.orderers[orderers[i]].tlsCACerts.path;
        ccp.orderers[orderers[i]].tlsCACerts.path = path.resolve(__dirname, relativePath);
    }
    // Changing orderer relative paths
    for (var i = 0; i < peers.length; i++) {
        var relativePath = ccp.peers[peers[i]].tlsCACerts.path;
        ccp.peers[peers[i]].tlsCACerts.path = path.resolve(__dirname, relativePath);
    }

    return ccp;
}

module.exports = {
    makeAbsolutePath
}