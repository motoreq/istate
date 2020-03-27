'use strict';

const convertFHIRToMedisot = function(args) {
    try {
   console.log("within convertHFIRToMedisot---", args);
   var argsObj = JSON.parse(args)
   console.log("----------------",argsObj)
   var uniqueGovtID = argsObj.identifier[0].value
   console.log("uniqueGovtID---", uniqueGovtID);
   var name = argsObj.name[0].given[0]
   console.log("name---", name);
   var telecom = argsObj.telecom.length
   var sex = argsObj.gender;
   var dob = argsObj.birthDate;
   var var1 = "";
   var var2 = "";
   var phone = "";
   var email = "";
   for(var i=0; i<telecom; i++){
       console.log(telecom);
      var1 = argsObj.telecom[i].system;
      if(var1 == "phone"){
         phone = argsObj.telecom[i].value
      }else{
        email = argsObj.telecom[i].value
      }
   }
   console.log("phone",phone)
   console.log("email",email)
   console.log("sex",sex)
   console.log("dob",dob)
   var finalObj = {uniqueGovtID:uniqueGovtID, username:name, sex:sex, dob:dob, phone:phone, email:email}
   var finalRegisterUser = [];
   finalRegisterUser.push(finalObj);
   return finalRegisterUser
    }catch (error){
        console.log("error", error)
    }
}


module.exports = convertFHIRToMedisot;