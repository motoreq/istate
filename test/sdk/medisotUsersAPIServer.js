

/*
 * Project Name : Medisot
 * Created by Abhijit Ghosh on 01/08/18
 * EMAIL: xxx
 */


const express = require('express');
const path = require('path');
const bodyParser = require('body-parser');
const cors = require('cors');
const passport = require('passport');

const app = express();
const server = require('http').createServer(app);
server.keepAliveTimeout=120*1000;
server._events.request.setMaxListeners(1);
// console.log(server);

app.setMaxListeners(0);
// console.log("Prev Max listeners: ", app.getMaxListeners());
const api = require('./routes/api');

const port = 5000;

app.use(cors());


var config = require('./config/config');
var log4js = require('log4js');
var logger = log4js.getLogger('JWT-logs');
logger.level = 'debug';



//********************* END of GWT ***********/

//set static folder
app.use(express.static(path.join(__dirname, 'public')));

//app.use(bodyParser.json());

app.use(bodyParser.json({limit: '10mb', extended: true}))
app.use(bodyParser.urlencoded({limit: '10mb', extended: true}))

app.use('/api',api);

//body-parser Middleware
app.get('/', (req, res)=> {
   res.send("Invalid end point")
});

server.listen(port, ()=> {
    console.log("server started on port",port);
})
