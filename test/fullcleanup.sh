docker rm -f $(docker ps -aq --filter "name=dev.*.users")
docker rm -f $(docker ps -aq --filter "name=peer.*.users")
docker rm -f $(docker ps -aq --filter "name=orderer.users.*")
docker rm -f $(docker ps -aq --filter "name=ca.users.*")
docker rm -f $(docker ps -aq --filter "name=cliuser")
docker system prune -f
docker image prune -f
docker network prune -f
docker rmi -f $(docker images | grep dev.*.users | awk '{print $1}')
docker volume prune -f 
cd ./sdk/
npm stop
cd ../

