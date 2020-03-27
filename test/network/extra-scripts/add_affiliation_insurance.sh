
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.insurance.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.insurance.com-cert.pem 
fabric-ca-client affiliation add insurance  -u https://admin:adminpw@ca.insurance.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.insurance.com-cert.pem 
