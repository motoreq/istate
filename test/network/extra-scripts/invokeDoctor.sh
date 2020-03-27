#/bin/bash
. setpeer.sh Doctor peer0
peer chaincode invoke -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["DoctorRegistration","[{\"doctorRegistrationID\":\"Doc001\",\"doctorName\":\"Rajeev\",\"speciality\":[\"ENT\",\"NEURO\"],\"associatedHospital\":[\"DESUN\",\"TATA MEDICAL\",\"COLUMBIA ASIA\"],\"sex\":\"M\"},{\"doctorRegistrationID\":\"Doc002\",\"doctorName\":\"LALIMA\",\"speciality\":[\"SURGEON\",\"PHYCHO\"],\"associatedHospital\":[\"RUBY\",\"COLUMBIA ASIA\",\"MEDICA\"],\"sex\":\"F\"}]"]}'
