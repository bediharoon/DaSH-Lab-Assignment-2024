#!/bin/bash

set -e

mkdir -p ./server/bin

cp ./server/servDockerfile ./server/bin/Dockerfile
cp ./server/session.conf ./server/bin
cp ./server/start.sh ./server/bin

go build -o ./server/bin/main ./server/main.go ./server/server.go


mkdir -p ./client/bin

cp ./client/clientDockerfile ./client/bin/Dockerfile
cp ./client/session.conf ./client/bin
cp ./client/start.sh ./client/bin
cp ./client/input.txt ./client/bin

go build -o ./client/bin/main ./client/main.go


docker-compose up -d
