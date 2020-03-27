
const express = require('express');
const path = require('path');
const bodyParser = require('body-parser');
const cors = require('cors');
const passport = require('passport');

const app = express();

const api = require('./expressApp');

const port = 3001;

app.use(cors());


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

app.listen(port, ()=> {
    console.log("server started on port",port);
})