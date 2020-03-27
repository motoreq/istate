const S3 = require('./routes/modified-sdk/lib/s3-client')

async function deleteKeys() {
    try {
        result = await S3.DeleteAllKeys()
        console.log(result)
    } catch (e) {
        console.log("Error caught:")
        console.log(e)
    }
}

deleteKeys();