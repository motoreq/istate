docker rm -f peer0.patient.com
. setenv.sh
docker-compose up -d peer0.patient.com
docker exec -it cli bash -e ./joinExistingChannel.sh