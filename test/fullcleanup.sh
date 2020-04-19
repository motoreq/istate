# Copyright 2020 Motoreq Infotech Pvt Ltd

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

