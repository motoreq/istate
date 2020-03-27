#!/usr/bin/env bash

docker-compose down --volumes 
# docker rm -f $(docker ps -a -q)
# docker rmi -f `docker images dev* -aq`
# docker volume rm $(docker volume ls -q)


