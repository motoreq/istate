const { InMemoryWallet, Wallet } = require('fabric-network')
const S3 = require('./lib/s3-client')

class S3Wallet {
    constructor(){
        this.inMemWallet = new InMemoryWallet()
    }
    async import(label, identity) {
        try {
            var result = await S3.Upload(label, JSON.stringify(identity))
        } catch (e) {
            throw new Error(e)
        }
    }

    async export(label) {
        console.log("Inside export....")
        if (!this.inMemWallet.exists(label)) {
            var identity = await S3.Get(label)
            await this.inMemWallet.import(label, JSON.parse(identity))
        }
        return await this.inMemWallet.export(label)
    }

    async exists(label) {
        try {
            if (!await this.inMemWallet.exists(label)) {
                var response = await S3.Get(label)
                var identity = response.Body.toString()
                await this.inMemWallet.import(label, JSON.parse(identity))
            }
            return true
        } catch (e) {
            if(e.message.includes('does not exist.')) {
                return false
            }
            else throw new Error(e)
        }
    }

    async setUserContext(client, label) {
        return await this.inMemWallet.setUserContext(client, label)
    }
}

module.exports = S3Wallet