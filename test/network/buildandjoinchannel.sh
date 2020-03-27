
#!/bin/bash -e
	
	echo "Building channel: medisotuserschannel" 		
        
            . setpeer.sh MedisotUsers peer0
            export CHANNEL_NAME="medisotuserschannel"
			peer channel create -o orderer.users.medisot.com:7050 -c $CHANNEL_NAME -f ./${CHANNEL_NAME}.tx --tls true --cafile $ORDERER_CA -t 1000s
			peer channel join -b $CHANNEL_NAME.block
			# export CHANNEL_NAME="tokenchannel"
			# peer channel create -o orderer.users.medisot.com:7050 -c $CHANNEL_NAME -f ./tokenchannel.tx --tls true --cafile $ORDERER_CA -t 1000s
			# peer channel join -b $CHANNEL_NAME.block