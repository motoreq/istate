'use strict';

const fhirToMedisotDocReg = function(args, user){
    try{
        console.log("args=============",args)
        var userName = user;
        var argsObj = JSON.parse(args)
        var entryDetails = argsObj.entry;
        var doctorRegistrationID = "";
        var associatedHospital = [];
        var speciality = [];
        //var finalObj = {};
        for(var i=0; i< entryDetails.length; i++){
           var resourceType = entryDetails[i].resource.resourceType;
           console.log("resourceType----------",resourceType)
           if(resourceType == "Practitioner"){
            doctorRegistrationID = entryDetails[i].resource.identifier[0].value
            console.log("doctorRegistrationID----------",doctorRegistrationID)
           }else if(resourceType == "Organization"){
              for(var j=0; j< entryDetails[i].resource.identifier.length; j++){
                 var associatedSingleHospital = entryDetails[i].resource.identifier[j].value
                 associatedHospital.push(associatedSingleHospital);
              }
           }else if(resourceType == "PractitionerRole"){
            for(var j=0; j< entryDetails[i].resource.identifier.length; j++){
                var singleSpeciality = entryDetails[i].resource.identifier[j].value
                speciality.push(singleSpeciality);
             }
           }
        }
        var finalObj = {userName:userName, doctorRegistrationID:doctorRegistrationID, speciality:speciality, associatedHospital:associatedHospital};
        return finalObj
    }catch(error){
        console.log("error", error) 
    }
}

module.exports = fhirToMedisotDocReg;