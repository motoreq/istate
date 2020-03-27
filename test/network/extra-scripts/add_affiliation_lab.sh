
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.lab.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.lab.com-cert.pem 
fabric-ca-client affiliation add lab  -u https://admin:adminpw@ca.lab.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.lab.com-cert.pem 
