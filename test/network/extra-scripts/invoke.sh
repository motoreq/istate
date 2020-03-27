#/bin/bash
. setpeer.sh Patient peer0

peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["RegisterEmail","[{\"email\":\"jeetabhi151@gmail.com\"},{\"email\":\"prashant@gmail.com\"}]"]}'

# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["ConfirmEmail","[{\"email\":\"prashant@gmail.com\",\"otp\":\"9p0D\"}]"]}'
# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["ConfirmEmail","[{\"email\":\"jeetabhi151@gmail.com\",\"otp\":\"9j0D\"}]"]}'
PASSWORD=`echo -n '{"prashant":"password","jeetabhi1":"password"}' | openssl base64`
# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["RegisterUser","[{\"username\":\"prashant\",\"name\":\"Pras\",\"age\":\"22\",\"sex\":\"M\",\"dob\":\"DOB\",\"phone\":\"8876567789\",\"email\":\"prashant@gmail.com\",\"uniqueGovtID\":\"HGTGJU231\"}]"]}' --transient "{\"PASSWORD\":\"$PASSWORD\"}"
# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["RegisterUser","[{\"username\":\"jeetabhi1\",\"name\":\"Pras\",\"age\":\"22\",\"sex\":\"M\",\"dob\":\"DOB\",\"phone\":\"8876567789\",\"email\":\"jeetabhi151@gmail.com\",\"uniqueGovtID\":\"HGTGJU231\"}]"]}' --transient "{\"PASSWORD\":\"$PASSWORD\"}"

PASSWORD1=`echo -n 'password' | openssl base64`
# peer chaincode query -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["CheckAuth","prashant"]}' --transient "{\"PASSWORD\":\"$PASSWORD1\"}"
# peer chaincode query -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["CheckAuth","jeetabhi1"]}' --transient "{\"PASSWORD\":\"$PASSWORD1\"}"


# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["LabRegistration","[{\"username\":\"prashant\",\"labName\":\"LABNAME\",\"labID\":\"12123\",\"testingOffered\":[\"test1\",\"test2\"],\"duration\":\"Duration\",\"status\":\"status\"}]"]}'

# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["DoctorRegistration","[{\"username\":\"prashant\",\"doctorRegistrationID\":\"DoctorRegID\",\"speciality\":[\"speciality1\",\"speciality2\"],\"associatedHospital\":[\"hospital1\",\"hospital2\"]}]"]}'

# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["PharmacyRegistration","[{\"username\":\"prashant\",\"pharmacyName\":\"PHARMACYNAME\",\"pharmacyRegistrationNo\":\"12123\",\"drugLicenceNo\":\"Drug_1234\"}]"]}'

# peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["InsuranceRegistration","[{\"username\":\"prashant\",\"insuranceID\":\"IID642\",\"insuranceName\":\"Iname\",\"category\":\"I_Category\"}]"]}'

peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c  '{"args":["RegisterEmail","[{\"email\":\"vas@gmail.com\"},{\"email\":\"pras@gmail.com\"},{\"email\":\"vasa@gmail.com\"}]"]}' 

###################################################### SAMPLES ############################################################

# ## Notice OTP from this command:
# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["RegisterEmail","[{\"email\":\"vas@gmail.com\"},{\"email\":\"pras@gmail.com\"},{\"email\":\"vasa@gmail.com\"}]"]}' 

# ## Update OTP in these:
# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["ConfirmEmail","[{\"email\":\"pras@gmail.com\",\"otp\":\"2p0D\"}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["ConfirmEmail","[{\"email\":\"vas@gmail.com\",\"otp\":\"2v0D\"}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["ConfirmEmail","[{\"email\":\"vasa@gmail.com\",\"otp\":\"2v0D\"}]"]}'

# ## Password must be passed in the format {"username":"pass"} and MUST be converted to base64 string before passing
# ## PASSWORD=`echo -n '{"pras":"password","vas":"password"}' | openssl base64`
# PASSWORD='eyJwcmFzIjoicGFzc3dvcmQiLCJ2YXMiOiJwYXNzd29yZCIsInZhc2EiOiJwYXNzd29yZCJ9' # 3 users - pras, vas, vasa
# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["RegisterUser","[{\"username\":\"pras\",\"name\":\"Pras\",\"age\":\"22\",\"sex\":\"M\",\"dob\":\"DOB\",\"phone\":\"8876567789\",\"email\":\"pras@gmail.com\",\"uniqueGovtID\":\"HGTGJU231\"}]"]}' --transient "{\"PASSWORD\":\"$PASSWORD\"}"

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["RegisterUser","[{\"username\":\"vas\",\"name\":\"Vas\",\"age\":\"27\",\"sex\":\"M\",\"dob\":\"DOB\",\"phone\":\"8876567789\",\"email\":\"vas@gmail.com\",\"uniqueGovtID\":\"HGTGJU231\"}]"]}' --transient "{\"PASSWORD\":\"$PASSWORD\"}"

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["RegisterUser","[{\"username\":\"vasa\",\"name\":\"Vasa\",\"age\":\"27\",\"sex\":\"M\",\"dob\":\"DOB\",\"phone\":\"8876567789\",\"email\":\"vasa@gmail.com\",\"uniqueGovtID\":\"HGTGJU231\"}]"]}' --transient "{\"PASSWORD\":\"$PASSWORD\"}"

# ## Again, password must be passed as base64 string
# PASSWORD1=`echo -n "password" | openssl base64`
# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["CheckAuth","pras"]}' --transient "{\"PASSWORD\":\"$PASSWORD1\"}"


# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["LabRegistration","[{\"username\":\"pras\",\"labName\":\"LABNAME\",\"labID\":\"12123\",\"testingOffered\":[\"test1\",\"test2\"],\"duration\":\"Duration\",\"status\":\"status\"}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["DoctorRegistration","[{\"username\":\"pras\",\"doctorRegistrationID\":\"DoctorRegID\",\"speciality\":[\"speciality1\",\"speciality2\"],\"associatedHospital\":[\"hospital1\",\"hospital2\"]}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["PharmacyRegistration","[{\"username\":\"pras\",\"pharmacyName\":\"PHARMACYNAME\",\"pharmacyRegistrationNo\":\"12123\",\"drugLicenceNo\":\"Drug_1234\"}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["InsuranceRegistration","[{\"username\":\"pras\",\"insuranceID\":\"IID642\",\"insuranceName\":\"Iname\",\"category\":\"I_Category\"}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["ClinicalResearchRegistration","[{\"username\":\"pras\",\"crID\":\"CRID642\",\"crName\":\"CRname\"}]"]}'

# ## Creates new EMR for that particular user, who is invoking
# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["createEMR"]}'

# ## jsonReadFields and jsonWriteFields has the json structure as String and this string must be converted to base64 before passing 
# ## Eg:Command to get base64 string in linux :  echo -n "{\"emrID\":\"\"}" | openssl base64
# ## Output : eyJlbXJJRCI6IiJ9  <--- This string must be sent as string in jsonReadFields/ jsonWriteFields arg.
# ## All EMR IDs are in format: emr_{username}_{#num} 
# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["shareEMR","[{\"emrID\":\"emr_pras_0\",\"userToShare\":\"vas\",\"jsonReadFields\":\"eyJlbXJJRCI6IiIsImRvY3RvclVzZXJOYW1lIjoiIiwidml0YWxzIjp7ImV4YW1pbmVyVXNlck5hbWUiOiIiLCJwdWxzZSI6IiIsImhlYXJ0UmF0ZSI6IiIsIkJQIjp7InN5c3RvbGljIjowLCJkaWFzdG9saWMiOjB9LCJveHlnZW5QZXJjZW50IjowLjAsIndlaWdodCI6MC4wLCJoZWlnaHQiOjAuMCwiZGF0ZSI6IiJ9LCJjaGllZkNvbXBsYWludCI6IiIsImNoaWVmQ29tcGxhaW50SGlzdG9yeSI6eyJmYW1pbHlIaXN0b3J5IjoiIiwicGVyc29uYWxIaXN0b3J5IjoiIiwicGFzdEhpc3RvcnkiOiIiLCJjdXJyZW50TWVkaWNhdGlvbiI6IiIsImFsbGVyZ2llcyI6IiJ9LCJwaHlzaWNhbEV4YW1pbmF0aW9uIjp7ImV4YW1pbmVyVXNlck5hbWUiOiIiLCJQIjoiIiwiSSI6IiIsIkMiOiIiLCJLIjoiIiwiTCI6IiIsIkUiOiIiLCJkYXRlIjoiIn0sInN5c3RlbWljRXhhbWluYXRpb24iOnsiZXhhbWluZXJVc2VyTmFtZSI6IiIsIlJTIjoiIiwiQ1MiOiIiLCJHSVQiOiIiLCJOUyI6IiIsImRhdGUiOiIifSwiZGlhZ25vc3RpY1Rlc3RzQWR2aXNlZCI6eyJhZHZpc29yVXNlck5hbWUiOiIiLCJwYXRob2xvZ3kiOiIiLCJyYWRpb2xvZ3kiOiIiLCJkYXRlIjoiIn0sImRpYWdub3N0aWNSZXBvcnRzIjp7InBhdGhvbG9neSI6eyJwYXRob2xvZ2lzdFVzZXJOYW1lIjoiIiwiUkJDIjoiIiwiV0JDIjoiIiwiSGIiOiIiLCJibG9vZFN1Z2FyIjoiIiwiY2hvbGVzdHJvbCI6IiIsInVyZWFDcmVhdGluaW5lIjoiIiwiRVNSIjoiIiwiZGF0ZSI6IiJ9LCJyYWRpb2xvZ3kiOnsicmFkaW9sb2dpc3RVc2VyTmFtZSI6IiIsInJhZGlvbG9neVJlcG9ydCI6IiIsInJhZGlvbG9neUltYWdlSGFzaCI6IiIsImRhdGUiOiIifX0sImRydWdzUHJlc2NyaWJlZCI6IiIsImRhdGUiOiIiLCJ3cml0ZXJVc2VyTmFtZSI6IiJ9\",\"jsonWriteFields\":\"eyJlbXJJRCI6IiJ9\",\"durationInSeconds\":100000}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["shareEMR","[{\"emrID\":\"emr_pras_0\",\"userToShare\":\"vas\",\"jsonReadFields\":\"eyJlbXJJRCI6IiJ9\",\"jsonWriteFields\":\"eyJlbXJJRCI6IiJ9\",\"durationInSeconds\":20}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["unshareEMR","[{\"emrID\":\"emr_pras_0\",\"userToUnshare\":\"vas\",\"readOrWrite\":\"read\"}]"]}'

# peer chaincode invoke -C $CHANNEL_NAME -n mycc $ORDERER_CONN_ARGS -c '{"args":["writeEMR","[{\"emrID\":\"emr_pras_0\",\"jsonData\":\"eyJlbXJJRCI6IkNoYW5nZWQgRU1SIElEIiwiZG9jdG9yVXNlck5hbWUiOiJwcmFzX2RvY3RvciIsInZpdGFscyI6eyJleGFtaW5lclVzZXJOYW1lIjoiY3NkIiwicHVsc2UiOiJmZHNmIiwiaGVhcnRSYXRlIjoiIiwiQlAiOnsic3lzdG9saWMiOjAsImRpYXN0b2xpYyI6MH0sIm94eWdlblBlcmNlbnQiOjAuMCwid2VpZ2h0IjowLjAsImhlaWdodCI6MC4wLCJkYXRlIjoiIn0sImNoaWVmQ29tcGxhaW50IjoiIiwiY2hpZWZDb21wbGFpbnRIaXN0b3J5Ijp7ImZhbWlseUhpc3RvcnkiOiIiLCJwZXJzb25hbEhpc3RvcnkiOiIiLCJwYXN0SGlzdG9yeSI6IiIsImN1cnJlbnRNZWRpY2F0aW9uIjoiIiwiYWxsZXJnaWVzIjoiIn0sInBoeXNpY2FsRXhhbWluYXRpb24iOnsiZXhhbWluZXJVc2VyTmFtZSI6IiIsIlAiOiIiLCJJIjoiIiwiQyI6IiIsIksiOiIiLCJMIjoiIiwiRSI6IiIsImRhdGUiOiIifSwic3lzdGVtaWNFeGFtaW5hdGlvbiI6eyJleGFtaW5lclVzZXJOYW1lIjoiIiwiUlMiOiIiLCJDUyI6IiIsIkdJVCI6IiIsIk5TIjoiIiwiZGF0ZSI6IiJ9LCJkaWFnbm9zdGljVGVzdHNBZHZpc2VkIjp7ImFkdmlzb3JVc2VyTmFtZSI6IiIsInBhdGhvbG9neSI6IiIsInJhZGlvbG9neSI6IiIsImRhdGUiOiIifSwiZGlhZ25vc3RpY1JlcG9ydHMiOnsicGF0aG9sb2d5Ijp7InBhdGhvbG9naXN0VXNlck5hbWUiOiIiLCJSQkMiOiIiLCJXQkMiOiIiLCJIYiI6IiIsImJsb29kU3VnYXIiOiIiLCJjaG9sZXN0cm9sIjoiIiwidXJlYUNyZWF0aW5pbmUiOiIiLCJFU1IiOiIiLCJkYXRlIjoiIn0sInJhZGlvbG9neSI6eyJyYWRpb2xvZ2lzdFVzZXJOYW1lIjoiIiwicmFkaW9sb2d5UmVwb3J0IjoiIiwicmFkaW9sb2d5SW1hZ2VIYXNoIjoiIiwiZGF0ZSI6IiJ9fSwiZHJ1Z3NQcmVzY3JpYmVkIjoiIiwiZGF0ZSI6IiIsIndyaXRlclVzZXJOYW1lIjoiIn0=\"}]"]}'


# ## ACL is not implemented in this, need to convert this appropriate to specific function and add acl later. - currently helpful for testing
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"args":["GeneralQuery","pras@gmail.com"]}'

# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"args":["GeneralQuery","pras"]}'

# ## Displays only the accessible fields.
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"args":["viewEMR","emr_pras_0"]}'

# ## Displays all the users registered, groups can be passed in groups field.
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"args":["queryAllUserProfile","{\"groups\":[\"doctor\",\"insurance\"]}"]}'

# echo -n "{\"emrID\":\"Changed EMR ID\",\"doctorUserName\":\"pras_doctor\",\"vitals\":{\"examinerUserName\":\"csd\",\"pulse\":\"fdsf\",\"heartRate\":\"\",\"BP\":{\"systolic\":0,\"diastolic\":0},\"oxygenPercent\":0.0,\"weight\":0.0,\"height\":0.0,\"date\":\"\"},\"chiefComplaint\":\"\",\"chiefComplaintHistory\":{\"familyHistory\":\"\",\"personalHistory\":\"\",\"pastHistory\":\"\",\"currentMedication\":\"\",\"allergies\":\"\"},\"physicalExamination\":{\"examinerUserName\":\"\",\"P\":\"\",\"I\":\"\",\"C\":\"\",\"K\":\"\",\"L\":\"\",\"E\":\"\",\"date\":\"\"},\"systemicExamination\":{\"examinerUserName\":\"\",\"RS\":\"\",\"CS\":\"\",\"GIT\":\"\",\"NS\":\"\",\"date\":\"\"},\"diagnosticTestsAdvised\":{\"advisorUserName\":\"\",\"pathology\":\"\",\"radiology\":\"\",\"date\":\"\"},\"diagnosticReports\":{\"pathology\":{\"pathologistUserName\":\"\",\"RBC\":\"\",\"WBC\":\"\",\"Hb\":\"\",\"bloodSugar\":\"\",\"cholestrol\":\"\",\"ureaCreatinine\":\"\",\"ESR\":\"\",\"date\":\"\"},\"radiology\":{\"radiologistUserName\":\"\",\"radiologyReport\":\"\",\"radiologyImageHash\":\"\",\"date\":\"\"}},\"drugsPrescribed\":\"\",\"date\":\"\",\"writerUserName\":\"\"}" | openssl base64
# eyJlbXJJRCI6IkNoYW5nZWQgRU1SIElEIiwiZG9jdG9yVXNlck5hbWUiOiJwcmFzX2RvY3RvciIsInZpdGFscyI6eyJleGFtaW5lclVzZXJOYW1lIjoiY3NkIiwicHVsc2UiOiJmZHNmIiwiaGVhcnRSYXRlIjoiIiwiQlAiOnsic3lzdG9saWMiOjAsImRpYXN0b2xpYyI6MH0sIm94eWdlblBlcmNlbnQiOjAuMCwid2VpZ2h0IjowLjAsImhlaWdodCI6MC4wLCJkYXRlIjoiIn0sImNoaWVmQ29tcGxhaW50IjoiIiwiY2hpZWZDb21wbGFpbnRIaXN0b3J5Ijp7ImZhbWlseUhpc3RvcnkiOiIiLCJwZXJzb25hbEhpc3RvcnkiOiIiLCJwYXN0SGlzdG9yeSI6IiIsImN1cnJlbnRNZWRpY2F0aW9uIjoiIiwiYWxsZXJnaWVzIjoiIn0sInBoeXNpY2FsRXhhbWluYXRpb24iOnsiZXhhbWluZXJVc2VyTmFtZSI6IiIsIlAiOiIiLCJJIjoiIiwiQyI6IiIsIksiOiIiLCJMIjoiIiwiRSI6IiIsImRhdGUiOiIifSwic3lzdGVtaWNFeGFtaW5hdGlvbiI6eyJleGFtaW5lclVzZXJOYW1lIjoiIiwiUlMiOiIiLCJDUyI6IiIsIkdJVCI6IiIsIk5TIjoiIiwiZGF0ZSI6IiJ9LCJkaWFnbm9zdGljVGVzdHNBZHZpc2VkIjp7ImFkdmlzb3JVc2VyTmFtZSI6IiIsInBhdGhvbG9neSI6IiIsInJhZGlvbG9neSI6IiIsImRhdGUiOiIifSwiZGlhZ25vc3RpY1JlcG9ydHMiOnsicGF0aG9sb2d5Ijp7InBhdGhvbG9naXN0VXNlck5hbWUiOiIiLCJSQkMiOiIiLCJXQkMiOiIiLCJIYiI6IiIsImJsb29kU3VnYXIiOiIiLCJjaG9sZXN0cm9sIjoiIiwidXJlYUNyZWF0aW5pbmUiOiIiLCJFU1IiOiIiLCJkYXRlIjoiIn0sInJhZGlvbG9neSI6eyJyYWRpb2xvZ2lzdFVzZXJOYW1lIjoiIiwicmFkaW9sb2d5UmVwb3J0IjoiIiwicmFkaW9sb2d5SW1hZ2VIYXNoIjoiIiwiZGF0ZSI6IiJ9fSwiZHJ1Z3NQcmVzY3JpYmVkIjoiIiwiZGF0ZSI6IiIsIndyaXRlclVzZXJOYW1lIjoiIn0=

# echo -n "{\"emrID\":\"\",\"doctorUserName\":\"\",\"vitals\":{\"examinerUserName\":\"\",\"pulse\":\"\",\"heartRate\":\"\",\"BP\":{\"systolic\":0,\"diastolic\":0},\"oxygenPercent\":0.0,\"weight\":0.0,\"height\":0.0,\"date\":\"\"},\"chiefComplaint\":\"\",\"chiefComplaintHistory\":{\"familyHistory\":\"\",\"personalHistory\":\"\",\"pastHistory\":\"\",\"currentMedication\":\"\",\"allergies\":\"\"},\"physicalExamination\":{\"examinerUserName\":\"\",\"P\":\"\",\"I\":\"\",\"C\":\"\",\"K\":\"\",\"L\":\"\",\"E\":\"\",\"date\":\"\"},\"systemicExamination\":{\"examinerUserName\":\"\",\"RS\":\"\",\"CS\":\"\",\"GIT\":\"\",\"NS\":\"\",\"date\":\"\"},\"diagnosticTestsAdvised\":{\"advisorUserName\":\"\",\"pathology\":\"\",\"radiology\":\"\",\"date\":\"\"},\"diagnosticReports\":{\"pathology\":{\"pathologistUserName\":\"\",\"RBC\":\"\",\"WBC\":\"\",\"Hb\":\"\",\"bloodSugar\":\"\",\"cholestrol\":\"\",\"ureaCreatinine\":\"\",\"ESR\":\"\",\"date\":\"\"},\"radiology\":{\"radiologistUserName\":\"\",\"radiologyReport\":\"\",\"radiologyImageHash\":\"\",\"date\":\"\"}},\"drugsPrescribed\":\"\",\"date\":\"\",\"writerUserName\":\"\"}" | openssl base64
# eyJlbXJJRCI6IiIsImRvY3RvclVzZXJOYW1lIjoiIiwidml0YWxzIjp7ImV4YW1pbmVyVXNlck5hbWUiOiIiLCJwdWxzZSI6IiIsImhlYXJ0UmF0ZSI6IiIsIkJQIjp7InN5c3RvbGljIjowLCJkaWFzdG9saWMiOjB9LCJveHlnZW5QZXJjZW50IjowLjAsIndlaWdodCI6MC4wLCJoZWlnaHQiOjAuMCwiZGF0ZSI6IiJ9LCJjaGllZkNvbXBsYWludCI6IiIsImNoaWVmQ29tcGxhaW50SGlzdG9yeSI6eyJmYW1pbHlIaXN0b3J5IjoiIiwicGVyc29uYWxIaXN0b3J5IjoiIiwicGFzdEhpc3RvcnkiOiIiLCJjdXJyZW50TWVkaWNhdGlvbiI6IiIsImFsbGVyZ2llcyI6IiJ9LCJwaHlzaWNhbEV4YW1pbmF0aW9uIjp7ImV4YW1pbmVyVXNlck5hbWUiOiIiLCJQIjoiIiwiSSI6IiIsIkMiOiIiLCJLIjoiIiwiTCI6IiIsIkUiOiIiLCJkYXRlIjoiIn0sInN5c3RlbWljRXhhbWluYXRpb24iOnsiZXhhbWluZXJVc2VyTmFtZSI6IiIsIlJTIjoiIiwiQ1MiOiIiLCJHSVQiOiIiLCJOUyI6IiIsImRhdGUiOiIifSwiZGlhZ25vc3RpY1Rlc3RzQWR2aXNlZCI6eyJhZHZpc29yVXNlck5hbWUiOiIiLCJwYXRob2xvZ3kiOiIiLCJyYWRpb2xvZ3kiOiIiLCJkYXRlIjoiIn0sImRpYWdub3N0aWNSZXBvcnRzIjp7InBhdGhvbG9neSI6eyJwYXRob2xvZ2lzdFVzZXJOYW1lIjoiIiwiUkJDIjoiIiwiV0JDIjoiIiwiSGIiOiIiLCJibG9vZFN1Z2FyIjoiIiwiY2hvbGVzdHJvbCI6IiIsInVyZWFDcmVhdGluaW5lIjoiIiwiRVNSIjoiIiwiZGF0ZSI6IiJ9LCJyYWRpb2xvZ3kiOnsicmFkaW9sb2dpc3RVc2VyTmFtZSI6IiIsInJhZGlvbG9neVJlcG9ydCI6IiIsInJhZGlvbG9neUltYWdlSGFzaCI6IiIsImRhdGUiOiIifX0sImRydWdzUHJlc2NyaWJlZCI6IiIsImRhdGUiOiIiLCJ3cml0ZXJVc2VyTmFtZSI6IiJ9

# echo -n "{\"emrID\":\"\"}" | openssl base64
# eyJlbXJJRCI6IiJ9
