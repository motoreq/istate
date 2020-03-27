
#!/bin/bash
echo "Stopping and resetting the network"
. stop.sh
. setenv.sh
docker-compose up -d
sleep 15
docker exec -it cliuser bash -e ./buildandjoinchannel.sh
sleep 5
docker exec -it cliuser bash -e ./test_install.sh
# cd ../medisot_sdk
# rm wallet/* -r
# cp routes/connection.json .

# function installNodeModules() {
# 	echo
# 	if [ -d node_modules ]; then
# 		echo "============== node modules installed already ============="
# 	else
# 		echo "============== Installing node modules ============="
# 		npm install
# 	fi
# 	echo
# }

# installNodeModules

# node enrollAdmin.js
# cd ../medisotNetwork 