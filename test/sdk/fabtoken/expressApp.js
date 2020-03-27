

/*
 * Project Name : Medisot
 * Created by Abhijit Ghosh on 01/08/18
 */

'use strict'

const express = require('express');
const router = express.Router();
const app = express();
const bodyParser = require('body-parser');

const healthcoin = require('./healthcoin');

// Configuration variables
app.set('secret', 'secret_this_should_be_longer');
app.set('expirationTimeInSeconds', 3600);

// CORS Access
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));
// Prerequisite to avoid CORS errors. (Cross-Origin-Resource-Sharing)
app.use((req, res, next) => {
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader(
      "Access-Control-Allow-Headers",
      "Origin, X-Requested-With, Content-Type, Accept, Authorization"
    );
    res.setHeader(
      "Access-Control-Allow-Methods",
      "GET, POST, PATCH, DELETE, OPTIONS"
    );
    next();
});

// Token

router.get("/listToken", (req, res, next) => {
    const user = req.get('Username');

    console.log(user);
    async function main() {
        const result = await healthcoin.list(user);
        if (result.success) { 
            res.status(200).json({
                success: true,
                result: result.result,
                msg: result.msg
            });       
        } else {
            res.status(404).json({
                success: false,
                result: result.result,
                msg: result.msg
            }); 
        }
    }
    main();
});

router.post("/issueToken", (req, res, next) => {

    const user = req.body.user;
    var quantity = req.body.quantity;

    if(typeof quantity === 'number') {
        quantity = quantity.toString();
    }
    async function main() {
        const result = await healthcoin.issue(user, ["MED", quantity]);
        if (result.success) { 
            res.status(200).json(result);       
        } else {
            res.status(404).json(result); 
        }
    }
    main();
});


// Util

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}



module.exports = router;
