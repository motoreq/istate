
#!/bin/bash -e
export PWD=`pwd`

export FABRIC_CFG_PATH=$PWD
export ARCH=$(uname -s)
export CRYPTOGEN=$PWD/bin/cryptogen
export CONFIGTXGEN=$PWD/bin/configtxgen
export CHANNEL_NAME="medisotuserschannel"

function generateArtifacts() {
	
	echo " *********** Generating artifacts ************ "
	echo " *********** Deleting old certificates ******* "
	
        rm -rf ./crypto-config
	
        echo " ************ Generating certificates ********* "
	
        $CRYPTOGEN generate --config=$FABRIC_CFG_PATH/crypto-config.yaml
        
        echo " ************ Generating tx files ************ "
	
		$CONFIGTXGEN -profile OrdererGenesis -outputBlock ./genesis.block
		
		$CONFIGTXGEN -profile $CHANNEL_NAME -outputCreateChannelTx ./$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
		
		echo "Generating anchor peers tx files for  Users.Medisot"
		$CONFIGTXGEN -profile $CHANNEL_NAME -outputAnchorPeersUpdate  ./${CHANNEL_NAME}medisotUsersMSPAnchor.tx -channelID $CHANNEL_NAME -asOrg medisotUsersMSP
		
		echo "Generating anchor peers tx files for  Patient"
		$CONFIGTXGEN -profile $CHANNEL_NAME -outputAnchorPeersUpdate  ./${CHANNEL_NAME}commonUsersMSPAnchor.tx -channelID $CHANNEL_NAME -asOrg commonUsersMSP
		

}
function generateDockerComposeFile(){
	OPTS="-i"
	if [ "$ARCH" = "Darwin" ]; then
		OPTS="-it"
	fi
	cp  docker-compose-template.yaml  docker-compose.yaml
	
	
	cd  crypto-config/peerOrganizations/users.medisot.com/ca
	PRIV_KEY=$(ls *_sk)
	cd ../../../../
	sed $OPTS "s/MEDISOTUSERS_PRIVATE_KEY/${PRIV_KEY}/g"  docker-compose.yaml
	
	
}
generateArtifacts 
cd $PWD
generateDockerComposeFile
cd $PWD


