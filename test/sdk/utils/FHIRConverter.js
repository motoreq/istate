
function makeItEasy(body) {
    var easyObject = {
        Patient: [],
        Practitioner: [],
        Organization: [],
        PractitionerRole: [],
        Appointment: [],
        Location: [],
        Observation: [],
    };


    if (body.resourceType === 'Bundle') {
        body.entry.forEach(function(i) {
            easyObject[i.resource.resourceType].push(i.resource);
            // console.log(easyObject[i.resource.resourceType])
            // console.log(i.resource.resourceType)
        });
    }
    else {
        easyObject[body.resourceType].push(body);
    }

    return easyObject;
}

const bookAppointment = function (body) {
    try {
        var args = {};
        var error = false;
        var errorMsg = '';

        var easyObject = makeItEasy(body);

        if(easyObject.Patient.length !== 1) {
            error = true;
            errorMsg = 'Expecting ONE patient resource';
        }
        console.log(easyObject.Patient[0].identifier[0].value)
        args.patientId = easyObject.Patient[0].identifier[0].value

        if(easyObject.Practitioner.length !== 1) {
            error = true;
            errorMsg = 'Expecting ONE Practitioner resource';
        }
        args.doctorId = easyObject.Practitioner[0].identifier[0].value;

        if(easyObject.Appointment.length !== 1) {
            error = true;
            errorMsg = 'Expecting ONE Appointment resource';
        }
        var epoch = new Date(easyObject.Appointment[0].start).getTime();
        args.startTimeSecs = epoch;

        var epoch = new Date(easyObject.Appointment[0].end).getTime();
        args.endTimeSecs = epoch;

        args.reason = easyObject.Appointment[0].reason[0].text;

        if(easyObject.Location.length !== 1) {
            error = true;
            errorMsg = 'Expecting ONE Location resource';
        }
        args.location = easyObject.Location[0].identifier[0].value;   


        if (error) {
            return {
                success: false,
                errorMsg: errorMsg,
            }
        }

        return {
            success: true,
            object: args,
        }
    } 
    catch (error) {
        console.log(error)
        return {
            success: false,
            errorMsg: error.toString(),
        }
    }
}

const addVitals = function (body) {
    try {
        var args = {}
        var error = false;
        var errorMsg = '';
        
        var easyObject = makeItEasy(body);

        if (easyObject.Observation.length != 1) {
            error = true;
            errorMsg = 'Expecting ONE Observation resource.';
        }
        
    
        easyObject.Observation[0].component.forEach(function(i){
            var code = i.valueQuantity.code.split('.')    
            if (code.length === 1) {
                args[code[0]] = i.valueQuantity.value;
            } 
            if (code.length === 2) {
                if (args[code[0]] === undefined) {
                    args[code[0]] = {};
                }
                args[code[0]][code[1]] = i.valueQuantity.value;
            } 
            
        });

        if (error) {
            return {
                success: false,
                errorMsg: errorMsg,
            }
        }

        return {
            success: true,
            object: args,
        }
    } catch (error) {
        console.log(error)
        return {
            success: false,
            errorMsg: error.toString(),
        }
    }
}


// Responses

const bookAppointmentResponse = function (req, res) {
    try {
        var response = {
            "resourceType": "Bundle",
            "id": '',
            "type": "transaction-response",
            "link": [
                {
                    "relation": "self",
                    "url": "http://hapi.fhir.org/baseDstu3"
                }
            ],
            "entry": [
                {
                    "response": {
                        "status": "201 Created",
                        "location": "Patient/1945825/_history/1",
                        "etag": "1",
                        "lastModified": "2019-05-26T18:03:56.567+00:00"
                    }
                },
                {
                    "response": {
                        "status": "201 Created",
                        "location": "Appointment/1945826/_history/1",
                        "etag": "1",
                        "lastModified": "2019-05-26T18:03:56.665+00:00"
                    }
                },
                {
                    "response": {
                        "status": "200 OK",
                        "location": "Location/1945824/_history/1",
                        "etag": "1"
                    }
                },
                {
                    "response": {
                        "status": "200 OK",
                        "location": "Practitioner/75844/_history/1",
                        "etag": "1"
                    }
                }
            ]
        }

        var parts = res.msg.split(' ');    
        var now = new Date(Date.now());
        response.id = parts[parts.length - 1];
        response.entry[0].response.location = `Patient/${req.patientId}/_history/1`;
        response.entry[0].response.lastModified = now;
        response.entry[1].response.location = `Appointment/${response.id}/_history/1`;
        response.entry[1].response.lastModified = now;
        response.entry[2].response.location = `Location/${req.location}/_history/1`;
        response.entry[3].response.location = `Practitioner/${req.doctorId}/_history/1`;

        console.log(response);
        return response;
    } 
    catch (error) {
        console.log(error);
        return createErrResponse(error.toString());
      
    }

}

const getAppointmentResponse = function (body) {
    try{
        let response;
        let now = new Date(Date.now());
        let obj = JSON.parse(body);

        response = {
        "resourceType": "Bundle",
        "id": "",
        "type": "searchset",
        "total": 1,
        "entry": []
        }

        let sampleentry = [
            {
                "fullUrl": "medisot.com/Appointment/1945826",
                "resource": {
                "resourceType": "Appointment",
                "id": "1945826",
                "meta": {
                    "versionId": "1",
                    "lastUpdated": "2019-05-26T18:46:11.358+00:00"
                },
                "identifier": [
                    {
                    "value": "AppointmentID_0192"
                    }
                ],
                "status": "proposed",
                "reason": [
                    {
                    "coding": [
                        {
                        "system": "http://snomed.info/sct",
                        "code": "109006",
                        "display": "Anxiety disorder of childhood OR adolescence (disorder)"
                        }
                    ],
                    "text": "X-RAY ANKLE 3+ VW"
                    }
                ],
                "description": "",
                "start": "2005-03-29T22:44:11.123+05:30",
                "end": "2005-03-29T22:44:11.123+05:30",
                "participant": [
                    {
                    "type": [
                        {
                        "coding": [
                            {
                            "display": "Patient"
                            }
                        ]
                        }
                    ],
                    "actor": {
                        "reference": "Patient/1945827",
                        "display": ""
                    },
                    "required": "required",
                    "status": "needs-action"
                    },
                    {
                    "type": [
                        {
                        "coding": [
                            {
                            "display": "Location of Appointment"
                            }
                        ]
                        }
                    ],
                    "actor": {
                        "reference": "Location/1945824",
                        "display": "facility"
                    },
                    "required": "required",
                    "status": "needs-action"
                    },
                    {
                    "type": [
                        {
                        "coding": [
                            {
                            "display": "Doctor"
                            }
                        ]
                        }
                    ],
                    "actor": {
                        "reference": "Practitioner/75844",
                        "display": ""
                    },
                    "required": "required",
                    "status": "needs-action"
                    }
                ]
                }
            },
            {
                "fullUrl": "medisot.com/Location/1945824",
                "resource": {
                "resourceType": "Location",
                "id": "1945824",
                "meta": {
                    "versionId": "1",
                    "lastUpdated": "2019-05-26T18:46:11.358+00:00"
                },
                "identifier": [
                    {
                    "type": {
                        "text": "Clinic"
                    },
                    "value": "Clinic_or_DoctorRoom_ID"
                    }
                ],
                "status": "active",
                "name": "Name of Clinic"
                },
                "request": {
                "method": "POST",
                "url": "Location",
                "ifNoneExist": "identifier=Clinic_or_DoctorRoom_ID"
                }
            },
            {
                "fullUrl": "medisot.com/Practitioner/75844",
                "resource": {
                "resourceType": "Practitioner",
                "id": "75844",
                "meta": {
                    "versionId": "1",
                    "lastUpdated": "2019-05-26T18:46:11.358+00:00"
                },
                "identifier": [
                    {
                    "type": {
                        "text": "Doctor ID"
                    },
                    "value": "1069"
                    }
                ],
                "name": [
                    {
                    "family": "",
                    "given": [
                        ""
                    ],
                    "prefix": [
                        ""
                    ]
                    }
                ],
                "gender": "unknown"
                },
                "request": {
                "method": "POST",
                "url": "Practitioner",
                "ifNoneExist": "identifier=1069"
                }
            },
            {
                "fullUrl": "medisot.com/PractitionerRole/unknown",
                "resource": {
                "resourceType": "PractitionerRole",
                "id": "unknown",
                "meta": {
                    "versionId": "1",
                    "lastUpdated": "2019-05-26T18:46:11.358+00:00"
                },
                "identifier": [
                    {
                    "type": {
                        "text": "speciality"
                    },
                    "value": ""
                    }
                ],
                "practitioner": {
                    "reference": "Practitioner/1069"
                },
                "location": [
                    {
                    "reference": "Location/Clinic_or_DoctorRoom_ID"
                    }
                ]
                },
                "request": {
                "method": "POST",
                "url": "PractitionerRole",
                "ifNoneExist": "identifier=unkown"
                }
            }];
        response.meta.lastUpdated = now;
        response.total = obj.length;
        
        for(var i = 0; i < obj.length; i++) {
            var appointment = JSON.parse(JSON.stringify(sampleentry[0]));
            
            appointment.fullUrl = `medisot.com/Appointment/${obj[i].appointmentId}`;
            appointment.resource.id = obj[i].appointmentId;
            appointment.resource.identifier[0].value = obj[i].appointmentId;
            
            appointment.resource.status = obj[i].status;
            appointment.resource.reason[0].text = obj[i].reason;
            
            appointment.resource.start = obj[i].startTimeSecs;
            appointment.resource.end = obj[i].endTimeSecs;
            appointment.resource.participant[0].actor.reference =`Patient/${obj[i].patientId}`
            appointment.resource.participant[1].actor.reference = `Location/${obj[i].location}`
            appointment.resource.participant[2].actor.reference = `Practitioner/${obj[i].doctorId}`


            var location = JSON.parse(JSON.stringify(sampleentry[1]));
            location.fullUrl = `medisot.com/Location/${obj[i].location}`
            location.resource.id = obj[i].location
            location.resource.identifier[0].value = obj[i].location
            location.request.ifNoneExist = `identifier=${obj[i].location}`

            var practitioner = JSON.parse(JSON.stringify(sampleentry[2]));
            practitioner.fullUrl = `medisot.com/Practitioner/${obj[i].doctorId}`
            practitioner.resource.id = obj[i].doctorId
            practitioner.resource.identifier[0].value = obj[i].doctorId
            practitioner.request.ifNoneExist = obj[i].doctorId
        
            var pracRole = JSON.parse(JSON.stringify(sampleentry[3]));
            pracRole.resource.practitioner.reference = `Practitioner/${obj[i].doctorId}`
            pracRole.resource.location[0].reference = obj[i].location

            
            response.entry.push(appointment);
            response.entry.push(location);
            response.entry.push(practitioner);
            response.entry.push(pracRole);
        }

        return response;
    }
    catch (error) {
        console.log(error.toString());
        return createErrResponse(error.toString());
    }
  
}
 

const addVitalsResponse = function(body) {
    try {
        body.id = body.identifier[0].value;
        body.meta = {};
        body.meta.versionId = '1';
        body.meta.lastUpdated = new Date(Date.now());
        return body;
    } catch (error) {
        console.log(error.toString());
        return createErrResponse(error.toString());
    }
    
}

const viewVitalsResponse = function(req, res) {
    try {
        res = JSON.parse(res);
        Object.keys(res).forEach(function(key) {
        res[key] = JSON.parse(res[key]);
        });

        let response = {
            "resourceType": "Bundle",
            "id": "8cf9b154-2898-43d7-bc32-f3d9a2a2def6",
            "type": "searchset",
            "total": 2,
            "entry": [
                {
                    "fullUrl": "http://medisot.com/Observation/1958067",
                    "resource": {
                        "resourceType": "Observation",
                        "id": "patientid",
                        "meta": {
                            "versionId": "1",
                            "lastUpdated": "2019-06-11T16:57:18.995+00:00"
                        },
                        "identifier": [
                            {
                                "value": "patientid"
                            }
                        ],
                        "subject": {
                            "reference": "Patient/1945825"
                        },
                        "component": []
                    },
                    "search": {
                        "mode": "match"
                    }
                }
            ]
        }

        let samplecomponent = [
            {
                "valueQuantity": {
                    "value": 80,
                    "code": "heartRate"
                }
            },
            {
                "valueQuantity": {
                    "value": 90,
                    "code": "pulseRate"
                }
            },
            {
                "valueQuantity": {
                    "value": 90,
                    "code": "oxygenSaturation"
                }
            },
            {
                "valueQuantity": {
                    "value": 50,
                    "code": "ecg"
                }
            },
            {
                "valueQuantity": {
                    "value": 40,
                    "code": "bloodPressure.systolic"
                }
            },
            {
                "valueQuantity": {
                    "value": 80,
                    "code": "bloodPressure.diastolic"
                }
            }
        ]

        var total = 0;
        response.entry[0].resource.id = req;
        response.entry[0].resource.identifier[0].value = req;
        response.entry[0].resource.subject.reference = `Patient/${req}`;
        response.entry[0].fullUrl = `http://medisot.com/Observation/${req}`;
        if (res.bloodPressure != 'null') {
            total = total + res.bloodPressure.length;
            for(var i=0; i<res.bloodPressure.length; i++) {
                var bps = JSON.parse(JSON.stringify(samplecomponent[4]));
                var bpd = JSON.parse(JSON.stringify(samplecomponent[5]));
                bps.valueQuantity.value = res.bloodPressure[i].systolic;
                bpd.valueQuantity.value = res.bloodPressure[i].diastolic;

                response.entry[0].resource.component.push(bps);
                response.entry[0].resource.component.push(bpd);

            }
        }

        if (res.ecg != 'null') {
            total = total + res.ecg.length;
            for(var i=0; i<res.ecg.length; i++) {
                var ecg = JSON.parse(JSON.stringify(samplecomponent[3]));
                ecg.valueQuantity.value = res.ecg[i].value;

                response.entry[0].resource.component.push(ecg);            
            }
        }

        if (res.heartRate != 'null') {
            total = total + res.heartRate.length;
            for(var i=0; i<res.heartRate.length; i++) {
                var heartRate = JSON.parse(JSON.stringify(samplecomponent[0]));
                heartRate.valueQuantity.value = res.heartRate[i].value;

                response.entry[0].resource.component.push(heartRate);            
            }
        }

        if (res.oxygenSaturation != 'null') {
            total = total + res.oxygenSaturation.length;
            for(var i=0; i<res.oxygenSaturation.length; i++) {
                var oxygenSaturation = JSON.parse(JSON.stringify(samplecomponent[2]));
                oxygenSaturation.valueQuantity.value = res.oxygenSaturation[i].value;

                response.entry[0].resource.component.push(oxygenSaturation);            
            }
        }

        if (res.pulseRate != 'null') {
            total = total + res.pulseRate.length;
            for(var i=0; i<res.pulseRate.length; i++) {
                var pulseRate = JSON.parse(JSON.stringify(samplecomponent[1]));
                pulseRate.valueQuantity.value = res.pulseRate[i].value;

                response.entry[0].resource.component.push(pulseRate);            
            }
        }

        if(total == 0) {
            response = {
                "resourceType": "Bundle",
                "id": "",
                "type": "searchset",
                "total": 0
              }
              response.id = req;

        } else {
            response.entry[0].search.mode = 'match';
        }
        response.total = total;

        return response;
    } catch(error) {
        console.log(error.toString());
        return createErrResponse(error.toString());
    }
}


const createErrResponse = function(text) {
    var errRes = {
        "resourceType": "OperationOutcome",
        "id": "exception",
        "issue": [
          {
            "severity": "error",
            "code": "exception",
            "details": {
              "text": "Backend Communication Error"
            }
          }
        ]
    }

    errRes.issue[0].details.text = text;
    return errRes;
}


module.exports = {
    bookAppointment,
    addVitals,
    bookAppointmentResponse,    
    getAppointmentResponse,
    createErrResponse,
    addVitalsResponse,
    viewVitalsResponse,
}