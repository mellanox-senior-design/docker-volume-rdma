#! /bin/bash

while ! nc -z plugin 8080; do
  echo "Waiting for the Plugin to launch (8080)."
  sleep 1
done
echo "The Plugin is up!"

sleep 5

echo "running integration tests"
cd $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma
go get -t -v ./...
go test -tags integration
