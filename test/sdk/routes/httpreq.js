const http = require('http')
const HOST = 'localhost'
const PORT = 3001

const FAIL_RESULT = 'Error when sending request.'

const post = async function(path, data, callback) {
    var options = {
      hostname: HOST,
      port: PORT,
      path: path,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Content-Length': 0
      }
    }
    var dataStr = JSON.stringify(data)
    options.headers['Content-Length'] = dataStr.length
  
    const req = http.request(options, (res) => {
        console.log(`statusCode: ${res.statusCode}`)
      
        res.on('data', (d) => {
          var dStr = d.toString()
          console.log(dStr)
          callback(JSON.parse(dStr))
        })
      })
      
      req.on('error', (error) => {
        console.error(error)
        return {
          success: false,
          result: FAIL_RESULT,
          msg: error.toString()
        }
      })
      
      req.write(dataStr)
      req.end()
}

const get = async function(path, headers, callback) {
  var options = {
    hostname: HOST,
    port: PORT,
    path: path,
    method: 'GET',
    headers:{}
  }
  options.headers = headers;

  const req = http.request(options, (res) => {
    console.log(`statusCode: ${res.statusCode}`)
    res.on('data', (d) => {
      var dStr = d.toString()
      console.log(dStr)
      callback(JSON.parse(dStr))
    })
  })

  req.on('error', (error) => {
    console.error(error)
    return {
      success: false,
      result: FAIL_RESULT,
      msg: error.toString()
    }
  })

  req.end()
}

// post('/api/issueToken', {
// 	"user": "admin",
// 	"quantity": "100"
// })

module.exports = {
    post,
    get,
}