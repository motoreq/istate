'use strict';

var methods = {}
methods.fhirToMedisotSaveEMR = function(args, callback){
    try{
       // console.log("args=============",args);
       // var userName = user;
        var argsObj = JSON.parse(args)
        var entryDetails = argsObj.entry;
        var entryDetails1 = argsObj.entry;
        var doctorRegistrationID = "";
        var associatedHospital = [];
        var speciality = [];
        var chiefComplaint = [];
        var chiefComplaintHistory = {};
        var diagnosticTestsAdvised = [];
        var report = [];
        var allergies = [];
        var currentMedication = [];
        var familyHistory = [];
        var pastHistory = [];    
        var personalHistory = []; // same above
       // var diagnosticTestsAdvised = [];
        var test = [];
        var drugsPrescribed = [];
        var physicalExamination = {};
        var vitals = {};
        var systemicExamination = {};
        var date = "";
        var writerUserName = "";
        var empID = args.empID
        var doctorUserName = ""
        var finalVitals = []
        var medicationRequest = [];
        var fullURL = ""
        var fullURLDR = ""
        var referenceURL = []
        var displayValTest = []
        var refUrlVal = []
        var test1 = ""
        var test2 = ""
        var referenceNumber = ""

        var testingDone = ""
        var referenceNumberURL = []
        var diagnosis = ""
        var location = ""
       
        //var finalObj = {};
        for(var i=0; i< entryDetails.length; i++){
           var resourceType = entryDetails[i].resource.resourceType;
           console.log("resourceType----------",resourceType)
           if(resourceType == "Condition"){
             var clinicalStatus = entryDetails[i].resource.clinicalStatus
             if(clinicalStatus == "active"){
                if(entryDetails[i].resource.code.hasOwnProperty('coding')){
                    for(var j=0; j< entryDetails[i].resource.code.coding.length; j++){
                        var chiefComplaintCondition = entryDetails[i].resource.code.coding[j].display
                        chiefComplaint.push(chiefComplaintCondition);
                    }
               }
             }else if(clinicalStatus == "resolved" || clinicalStatus == "Inactive" || clinicalStatus == "remission" || clinicalStatus == "recurrence"){
                if(entryDetails[i].resource.code.hasOwnProperty('coding')){
                    for(var j=0; j< entryDetails[i].resource.code.coding.length; j++){
                        var pastHistoryCondition = entryDetails[i].resource.code.coding[j].display
                        pastHistory.push(pastHistoryCondition);
                    }
               }
            }else{
                var chiefComplaintCondition = entryDetails[i].resource.code.text
                chiefComplaint.push(chiefComplaintCondition);
             }
           }else if(resourceType == "AllergyIntolerance"){
            if(entryDetails[i].resource.code.hasOwnProperty('coding')){
                for(var j=0; j< entryDetails[i].resource.code.coding.length; j++){
                    var chiefComplaintAllergyIntolerance = entryDetails[i].resource.code.coding[j].display
                    allergies.push(chiefComplaintAllergyIntolerance);
                }
               
            }else{
                var chiefComplaintAllergyIntolerance = entryDetails[i].resource.code.text
                allergies.push(chiefComplaintAllergyIntolerance);
            }
            
           } else if(resourceType == "FamilyMemberHistory"){
            
                // if(entryDetails[i].resource.code.hasOwnProperty('coding')){
                //     for(var j=0; j< entryDetails[i].resource.code.coding.length; j++){
                //         var chiefComplaintFamilyMemberHistory = entryDetails[i].resource.code.coding[j].display
                //         familyHistory.push(chiefComplaintFamilyMemberHistory);
                //     }
                // }
                if(entryDetails[i].resource.hasOwnProperty('condition')){
                    for(var j=0; j< entryDetails[i].resource.condition.length; j++){
                        var chiefComplaintFamilyMemberHistory = entryDetails[i].resource.condition[j].code.text
                        familyHistory.push(chiefComplaintFamilyMemberHistory);
                    }
                }
           }else if(resourceType == "ProcedureRequest"){ 
               fullURL = entryDetails[i].fullUrl

               console.log("fullURL------",fullURL)
            
            if(entryDetails[i].resource.code.hasOwnProperty('coding')){
            for(var j=0; j< entryDetails[i].resource.code.coding.length; j++){
                testingDone = entryDetails[i].resource.code.coding[j].display
             }
            }else{
                testingDone = entryDetails[i].resource.code.text
            }
           }else if(resourceType == "DiagnosticReport"){ 
               if(entryDetails[i].resource.hasOwnProperty('result')){
                for(var j=0; j< entryDetails[i].resource.result.length; j++){
                    var referenceURL1 = entryDetails[i].resource.result[j].reference
                    referenceURL.push(referenceURL1)
                    console.log("referenceURL----",referenceURL)
                    displayValTest.push(entryDetails[i].resource.result[j].display)
                    refUrlVal.push(entryDetails[i].resource.result[j])
                }
               }
               if(entryDetails[i].resource.hasOwnProperty('basedOn')){
                    if (entryDetails[i].resource.basedOn.reference == fullURL) {
                      //  var reportDiagonostics = DiagnosticReport.result.display  // Need to discuss with chidamber

                       if (entryDetails[i].resource.hasOwnProperty('result')) {
                          var resultLength = entryDetails[i].resource.length
                          for(var k=0; k<resultLength; k++){
                             // referenceNumberURL = entryDetails[i].resource.result[resultLength].reference
                             referenceNumberURL.push(entryDetails[i].resource.result[resultLength].reference)
                             console.log("referenceNumberURL---",referenceNumberURL)

                          }
                       }
                    }
                }
            
           }else if(resourceType == "Observation"){ 
            if(entryDetails[i].resource.hasOwnProperty('status')){
                 var status = entryDetails[i].resource.status;
                 if(status == "preliminary"){
                    diagnosis = entryDetails[i].resource.code.coding[0].display;
                }
            }
             for(var j=0; j<referenceNumberURL.length; j++){
                if (entryDetails[i].fullUrl === referenceNumberURL[j]){
                    if(entryDetails[i].resource.valueQuantity.hasOwnProperty('value')){
                    var observationValue = entryDetails[i].resource.valueQuantity.value
                    var observationValueString = observationValue.toString()
                    var testReport1 = {"test":testingDone, "report":observationValueString}
                    diagnosticTestsAdvised.push(testReport1)
                    console.log("diagnosticTestsAdvised-xxxxxx-----",diagnosticTestsAdvised)
                    }
                }
              }

               console.log("referenceURL.length----",refUrlVal.length)
               for(var m=0; m<refUrlVal.length; m++){
                   
                   for(var k=0; k< entryDetails1.length; k++){
                    if (entryDetails1[k].fullUrl === refUrlVal[m].reference){
                            var singleDisplay = refUrlVal[m].display
                            if(entryDetails1[k].resource.hasOwnProperty('valueQuantity')){
                                var observationValue = entryDetails1[k].resource.valueQuantity.value
                                var observationValueString = observationValue.toString()
                                var testReport2 = {"test":singleDisplay, "report":observationValueString}
                                diagnosticTestsAdvised.push(testReport2)
                                console.log("diagnosticTestsAdvised>>>>>>>",diagnosticTestsAdvised)
                           }
                        }
                    }
                }
            if(entryDetails[i].resource.hasOwnProperty('performer')){
               doctorUserName = entryDetails[i].resource.performer.reference
            }
            
            if(entryDetails[i].resource.code.hasOwnProperty('coding')){
                for(var j=0; j< entryDetails[i].resource.code.coding.length; j++){
                    var display = entryDetails[i].resource.code.coding[j].display
                    if (display == "Body Height"){
                        var value = entryDetails[i].resource.valueQuantity.value
                        var unit = entryDetails[i].resource.valueQuantity.unit
                        vitals = {"bodyHeight":value+unit}
                    }else if(display == "Body Weight"){
                      //  var value = entryDetails[i].resource.valueQuantity.value  //Need to discuss chidamer
                       // var unit = entryDetails[i].resource.valueQuantity.unit
                      //  vitals = {"bodyWeight":value+unit}
                    }else if(display == "Body Mass Index"){
                        var value = entryDetails[i].resource.valueQuantity.value
                        var unit = entryDetails[i].resource.valueQuantity.unit
                        vitals = {"bmi":value+unit}
                    }else if(display == "Diastolic Blood Pressure"){
                        var value = entryDetails[i].resource.valueQuantity.value
                        var unit = entryDetails[i].resource.valueQuantity.unit
                        vitals = {"dbp":value+unit}
                    }else if(display == "Systolic Blood Pressure"){
                        var value = entryDetails[i].resource.valueQuantity.value
                        var unit = entryDetails[i].resource.valueQuantity.unit
                        vitals = {"sbp":value+unit}
                    }else if(display == "Oral temperature"){
                        var value = entryDetails[i].resource.valueQuantity.value
                        var unit = entryDetails[i].resource.valueQuantity.unit
                        vitals = {"oralTemp":value+unit}
                    }else if(display == "Quality adjusted life years"){
                        var value = entryDetails[i].resource.valueQuantity.value
                        var unit = entryDetails[i].resource.valueQuantity.unit
                        vitals = {"qaly":value+unit}
                    }
                    finalVitals.push(vitals)

                }
            }
           }
           else if(resourceType == "Procedure"){
            var previousSurgeries = [];
            if(entryDetails[i].resource.hasOwnProperty('status')){
                if(entryDetails[i].resource.status == "completed"){
                   var previousSurgery = entryDetails[i].resource.code.text;
                   previousSurgeries.push(previousSurgery);

                }
            }
            if(entryDetails[i].resource.code.hasOwnProperty('text')){
              var systemicExaminationVal = entryDetails[i].resource.code.text
              systemicExamination = {"CE":systemicExaminationVal}
            }
           }else if(resourceType == "MedicationRequest"){
               if(entryDetails[i].resource.status == "active"){
                if(entryDetails[i].resource.medicationCodeableConcept.hasOwnProperty('coding')){
                    for(var j=0; j< entryDetails[i].resource.medicationCodeableConcept.coding.length; j++){
                    var displayM = entryDetails[i].resource.medicationCodeableConcept.coding[j].display
                    medicationRequest.push(displayM)
               }
               if(entryDetails[i].resource.hasOwnProperty('dosageInstruction')){
                for(var j=0; j< entryDetails[i].resource.dosageInstruction.length; j++){
                var doseQuantity = entryDetails[i].resource.dosageInstruction[j].doseQuantity.value
                var doseQuantityString = doseQuantity.toString()
                medicationRequest.push(doseQuantityString)
               }
              }
             }
           }
          }else if(resourceType == "Location"){
              // location = entryDetails[i].resource.location
              location = entryDetails[i].resource.address.city
          }
           chiefComplaintHistory ={allergies,currentMedication,familyHistory,pastHistory,personalHistory}
          // diagnosticTestsAdvised = {test, report}
           date = "" // hardcoded

        }
        // diagnosticTestsAdvised = diagnosticTestsAdvised.filter( function( item, index, inputArray ) {
        //     return inputArray.indexOf(item) == index;
        // });

        diagnosticTestsAdvised = diagnosticTestsAdvised.filter((thing,index) => {
            return index === diagnosticTestsAdvised.findIndex(obj => {
              return JSON.stringify(obj) === JSON.stringify(thing);
            });
          });

        console.log("diagnosticTestsAdvised---------",diagnosticTestsAdvised);
 
        var finalObj = {chiefComplaint:chiefComplaint, chiefComplaintHistory:chiefComplaintHistory, date:date, diagnosticTestsAdvised:diagnosticTestsAdvised, doctorUserName:doctorUserName, drugsPrescribed:medicationRequest,empID:empID, physicalExamination:physicalExamination, systemicExamination:systemicExamination, vitals:finalVitals, writerUserName:writerUserName,previousSurgeries:previousSurgeries,diagnosis:diagnosis,location:location};
        console.log("finalObj---------",finalObj);
        callback(finalObj)
    }catch(error){
        console.log("error", error) 
    }
}

module.exports = methods;

