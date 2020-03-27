
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.doctor.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.doctor.com-cert.pem 
fabric-ca-client affiliation add doctor  -u https://admin:adminpw@ca.doctor.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.doctor.com-cert.pem 
