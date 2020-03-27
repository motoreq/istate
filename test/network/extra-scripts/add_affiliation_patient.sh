
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.patient.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.patient.com-cert.pem 
fabric-ca-client affiliation add patient  -u https://admin:adminpw@ca.patient.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.patient.com-cert.pem 
