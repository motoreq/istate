
#!/bin/bash
export ORDERER_CA=/opt/ws/crypto-config/ordererOrganizations/users.medisot.com/msp/tlscacerts/tlsca.users.medisot.com-cert.pem

if [ $# -lt 2 ];then
	echo "Usage : . setpeer.sh MedisotUsers|CommonUsers| <peerid>"
fi
export peerId=$2

if [[ $1 = "MedisotUsers" ]];then
	echo "Setting to organization MedisotUsers peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.users.medisot.com:7051
	export CORE_PEER_LOCALMSPID=medisotUsersMSP
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/users.medisot.com/peers/$peerId.users.medisot.com/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/users.medisot.com/peers/$peerId.users.medisot.com/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/users.medisot.com/peers/$peerId.users.medisot.com/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/users.medisot.com/users/Admin@users.medisot.com/msp
fi

if [[ $1 = "CommonUsers" ]];then
	echo "Setting to organization CommonUsers peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.users.common.com:7051
	export CORE_PEER_LOCALMSPID=commonUsersMSP
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/users.common.com/peers/$peerId.users.common.com/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/users.common.com/peers/$peerId.users.common.com/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/users.common.com/peers/$peerId.users.common.com/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/users.common.com/users/Admin@users.common.com/msp
fi
