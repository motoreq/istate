#/bin/bash
export CHANNEL_NAME="medisotuserschannel"
./bin/configtxgen -profile OrdererGenesis -channelID sys-channel -outputBlock genesis.block
./bin/configtxgen -profile OrdererGenesis -channelID sys-channel -outputBlock $CHANNEL_NAME.block
./bin/configtxgen -profile $CHANNEL_NAME -outputCreateChannelTx ./$CHANNEL_NAME.tx -channelID $CHANNEL_NAME

./bin/configtxgen -profile OrdererGenesis -channelID sys-channel -outputBlock tokenchannel.block
./bin/configtxgen -profile tokenChannel -outputCreateChannelTx ./tokenchannel.tx -channelID tokenchannel
