const Channel = require('../../node_modules/fabric-client/lib/Channel');
const client_utils = require('../../node_modules/fabric-client/lib/client-utils');
// const Channel = require('./fabric-client/lib/Channel');
// const client_utils = require('./fabric-client/lib/client-utils');

class ModChannel extends Channel {
    constructor(name, clientContext) {
      super(name, clientContext);
    }
    /*
	 * Internal static method to allow transaction proposals to be called without
	 * creating a new channel
	 */
	static async sendTransactionProposalBySignedProposal(targets, proposal, timeout) {
		const method = 'sendTransactionProposal(static)';
		// logger.debug('%s - start', method);

        let errorMsg;
        if (!targets || targets.length < 1) {
			errorMsg = 'Missing peer objects in Transaction proposal';
		}

		if (errorMsg) {
			// logger.error('%s error %s', method, errorMsg);
			throw new Error(errorMsg);
		}

		const responses = await client_utils.sendPeersProposal(targets, proposal.signed, timeout);
		return [responses, proposal.source];
    }
    
    /**
	 * Sends a proposal to one or more endorsing peers that will be handled by the chaincode.
	 * There is no difference in how the endorsing peers process a request
	 * to invoke a chaincode for transaction vs. to invoke a chaincode for query. All requests
	 * will be presented to the target chaincode's 'Invoke' method which must be implemented to
	 * understand from the arguments that this is a query request. The chaincode must also return
	 * results in the byte array format and the caller will have to be able to decode
	 * these results.
	 *
	 * If the request contains a <code>txId</code> property, that transaction ID will be used, and its administrative
	 * privileges will apply. In this case the <code>useAdmin</code> parameter to this function will be ignored.
	 *
	 * @param {ChaincodeQueryRequest} request
	 * @param {boolean} useAdmin - Optional. Indicates that the admin credentials should be used in making
	 *                  this call. Ignored if the <code>request</code> contains a <code>txId</code> property.
	 * @returns {Promise} A Promise for an array of byte array results returned from the chaincode
	 *                    on all Endorsing Peers
	 * @example
	 * <caption>Get the list of query results returned by the chaincode</caption>
	 * const responsePayloads = await channel.queryByChaincode(request);
	 * for (let i = 0; i < responsePayloads.length; i++) {
	 *     console.log(util.format('Query result from peer [%s]: %s', i, responsePayloads[i].toString('utf8')));
	 * }
	 */
	async queryByChaincodeBySignedProposal(targets, proposal, timeout) {
		// logger.debug('queryByChaincodeBySignedProposal - start');

		const proposalResults = await ModChannel.sendTransactionProposalBySignedProposal(targets, proposal, timeout);
		const responses = proposalResults[0];
		// logger.debug('queryByChaincode - results received');

		if (!responses || !Array.isArray(responses)) {
			throw new Error('Payload results are missing from the chaincode query');
		}

		const results = [];
		responses.forEach((response) => {
			if (response instanceof Error) {
				results.push(response);
			} else if (response.response && response.response.payload) {
				if (response.response.status === 200) {
					results.push(response.response.payload);
				} else {
					if (response.response.message) {
						results.push(new Error(response.response.message));
					} else {
						results.push(new Error(response));
					}
				}
			} else {
				// logger.error('queryByChaincode - unknown or missing results in query ::' + results);
				results.push(new Error(response));
			}
		});

		return results;
	}
}


module.exports = ModChannel;