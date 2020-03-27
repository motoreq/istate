var body = "{\"bloodPressure\":\"[{\\\"systolic\\\":40,\\\"diastolic\\\":80,\\\"timeStamp\\\":1560357882}]\",\"ecg\":\"[{\\\"value\\\":50,\\\"timeStamp\\\":1560357882}]\",\"heartRate\":\"[{\\\"value\\\":80,\\\"timeStamp\\\":1560357882}]\",\"oxygenSaturation\":\"[{\\\"value\\\":90,\\\"timeStamp\\\":1560357882}]\",\"pulseRate\":\"[{\\\"value\\\":90,\\\"timeStamp\\\":1560357882}]\"}"

console.log(body)
res = {
    "result": "Transaction has been submitted.",
    "msg": "Successfully added appointment app_ceac75d6f01d6990ee8ef450850bb45e6c19d4454443ccbc11cf98e3d4521ff9"
}


  //console.log(body.entry[0])

  
//console.log(epoch);

var converter = require('./FHIRConverter')
const util = require('util')
//var req = converter.addVitals(body)
//var myObject = converter.getAppointmentResponse("[{\"docType\":\"appointment\",\"appointmentId\":\"app_8c3b8b72a1fa049f878c48269d42c5d14ef65ac4e13e271486b663da038c908d\",\"status\":\"BOOKED\",\"startTimeSecs\":1560238001,\"endTimeSecs\":1560239000,\"patientId\":\"prasanths96\",\"doctorId\":\"drrajeev\",\"location\":\"Bang2\",\"reason\":\"Cough\"},{\"docType\":\"appointment\",\"appointmentId\":\"app_9fcfd1db22939bf6a3908a24341c2894a25e51873924977b1e699a2c601da015\",\"status\":\"BOOKED\",\"startTimeSecs\":1,\"endTimeSecs\":5,\"patientId\":\"prasanths96\",\"doctorId\":\"drrajeev\",\"location\":\"Bang2\",\"reason\":\"Cough\"},{\"docType\":\"appointment\",\"appointmentId\":\"app_b7e8b2e4d433ccaf78e324d9f29fa01efd22c182633d80feb3be19bf64f5d4a4\",\"status\":\"BOOKED\",\"startTimeSecs\":1560237542,\"endTimeSecs\":1560238000,\"patientId\":\"prasanths96\",\"doctorId\":\"drrajeev\",\"location\":\"Bang2\",\"reason\":\"Cough\"}]");
//console.log(util.inspect(myObject.entry, false, null, true /* enable colors */))
// console.log(converter.addVitalsResponse(body));

console.log(converter.viewVitalsResponse('prasanths96', body).entry[0].resource.component);

