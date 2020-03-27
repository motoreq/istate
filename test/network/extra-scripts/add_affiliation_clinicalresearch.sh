
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.clinicalresearch.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.clinicalresearch.com-cert.pem 
fabric-ca-client affiliation add clinicalresearch  -u https://admin:adminpw@ca.clinicalresearch.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.clinicalresearch.com-cert.pem 
