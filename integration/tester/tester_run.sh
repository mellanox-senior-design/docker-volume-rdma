#! /bin/bash

while ! nc -z plugin 8080; do
  echo "Waiting for the Plugin to launch (8080)."
  sleep 1
done
echo "The Plugin is up!"

sleep 5

echo "Getting depending packages"
cd $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma
go get -v github.com/docker/docker/pkg/namesgenerator
go get -t -v ./...
echo "Running integration tests"
go test -tags integration
