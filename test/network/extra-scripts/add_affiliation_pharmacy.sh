
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.pharmacy.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.pharmacy.com-cert.pem 
fabric-ca-client affiliation add pharmacy  -u https://admin:adminpw@ca.pharmacy.com:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.pharmacy.com-cert.pem 
