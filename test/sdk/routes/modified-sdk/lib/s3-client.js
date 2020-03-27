const AWS = require('aws-sdk')

const BUCKET_NAME = 'users.medisot'
const IAM_USER_KEY = 'AKIARIQSQ4B6IOL55DJU'
const IAM_USER_SECRET = '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'


let s3 = new AWS.S3({
    accessKeyId: IAM_USER_KEY,
    secretAccessKey: IAM_USER_SECRET,
    Bucket: BUCKET_NAME
})


function Upload (key, data) {
    var promise = new Promise((resolve, reject) => {
        var params = {
            Bucket: BUCKET_NAME, 
            Key: key, Body: data
        }
        s3.upload(params, function(err, data) {
                if(err) {
                        reject(err)
                }
                if (!data) {
                    reject("No response from S3. Data is null.")
                }
                resolve(data)
        })
    })
    return promise
}

function Get (key) {
    var promise = new Promise((resolve, reject) => {
        var params = {
            Bucket: BUCKET_NAME, 
            Key: key
        }
        s3.getObject(params, function(err, data) {
            if (err) { 
                // console.log("S3 getObject Error:", err)
                reject(err) // an error occurred
            } 
            if (!data) {
                reject("No response from S3. Data is null.")
            }
            resolve(data)        // successful response  
            /*
            data = {
            AcceptRanges: "bytes", 
            ContentLength: 10, 
            ContentRange: "bytes 0-9/43", 
            ContentType: "text/plain", 
            ETag: "\"0d94420ffd0bc68cd3d152506b97a9cc\"", 
            LastModified: <Date Representation>, 
            Metadata: {
            }, 
            VersionId: "null"
            }
            */
        })
    })
    return promise
}    

function DeleteKeys (keys) {
    var promise = new Promise((resolve, reject) => {
        var params = {
            Bucket: BUCKET_NAME, 
            Delete: {
                Objects: keys, 
                Quiet: false
            }
        }
        s3.deleteObjects(params, function(err, data) {
            if (err) { 
                reject(err) // an error occurred
            } 
            if (!data) {
                reject("No response from S3. Data is null.")
            }  
            resolve(data)        // successful response  
            /*
            data = {
             Deleted: [
                {
               DeleteMarker: true, 
               DeleteMarkerVersionId: "A._w1z6EFiCF5uhtQMDal9JDkID9tQ7F", 
               Key: "objectkey1"
              }, 
                {
               DeleteMarker: true, 
               DeleteMarkerVersionId: "iOd_ORxhkKe_e8G8_oSGxt2PjsCZKlkt", 
               Key: "objectkey2"
              }
             ]
            }
            */
          });
    })
    return promise
}

async function DeleteAllKeys () {
    try {
        result = await ListAllKeys(null, null)
        keys = [];
        if (result.length == 0) {
            throw new Error("No keys available to delete.")
        }
        for(var i = 0; i < result.length; i ++) {
            keys.push({Key: result[i].Key})
        }
        console.log("Keys: ", keys)
        return await DeleteKeys(keys)
    } catch(e) {
        return e
    }
}


async function ListAllKeys (token, allKeys) {
    var allKeys = [];
    var params = {
            Bucket: BUCKET_NAME
        }
    var promise = new Promise((resolve, reject) => {
        try {
            if(token) {
                params.ContinuationToken = token;
            }

            s3.listObjectsV2(params, async function(err, data){
                if (err) { 
                    // console.log("S3 getObject Error:", err)
                    reject(err) // an error occurred
                } else if (!data) {
                    reject("No response from S3. Data is null.")
                } else {
                    allKeys = data.Contents
                    if(data.IsTruncated) {
                        var tempKeys = await listAllKeys(data.NextContinuationToken)// sample tempKeys
                                                                                    // [{   Key: 'admin',
                                                                                    //      LastModified: 2020-02-20T19:00:31.000Z,
                                                                                    //      ETag: '"9e9b2891c7707923e8e2fc79f3fa7b52"',
                                                                                    //      Size: 1104,
                                                                                    //      StorageClass: 'STANDARD' }]
                        allKeys.concat(tempKeys);
                    } else {
                        resolve(allKeys)
                    }    
                }        
            })
        } catch (e) {
            reject(e)
        }
    })
    return promise
}

module.exports = {
    Upload,
    Get,
    ListAllKeys,
    DeleteKeys,
    DeleteAllKeys
}