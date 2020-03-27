. setpeer.sh Patient peer0
export CHANNEL_NAME="medisotchannel"

peer channel fetch 0 medisotchannel.block -o orderer.medisot.net:7050 -c medisotchannel --tls --cafile $ORDERER_CA
peer channel join -b medisotchannel.block
peer channel fetch 0 tokenchannel.block -o orderer.medisot.net:7050 -c tokenchannel --tls --cafile $ORDERER_CA
peer channel join -b tokenchannel.block

peer lifecycle chaincode package medisotcc.tar.gz --path github.com/medisot --lang golang --label medisotcc_7 

. setpeer.sh Patient peer0
export CHANNEL_NAME="medisotchannel"
peer lifecycle chaincode install medisotcc.tar.gz
