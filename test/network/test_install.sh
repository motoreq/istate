#!/bin/bash
. setpeer.sh MedisotUsers peer0 
export CHANNEL_NAME="medisotuserschannel"
peer chaincode install -n medisotUsers -v 1.0 -l golang -p  github.com/test/
peer chaincode instantiate -o orderer.users.medisot.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n medisotUsers -v 1.0 -c '{"Args":["init",""]}' -P " OR( 'medisotUsersMSP.member','commonUsersMSP.member') " 

