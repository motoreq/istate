#/bin/bash
. setpeer.sh Patient peer0
# peer chaincode query -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["GeneralQuery","prashant@gmail.com"]}'

# peer chaincode query -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["GeneralQuery","prashant"]}'

peer chaincode query -o orderer.medisot.net:7050  --tls --cafile $ORDERER_CA -C medisotchannel -n medisot -c '{"args":["GeneralQuery","EMR_prashant"]}'