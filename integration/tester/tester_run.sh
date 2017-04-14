#! /bin/bash

while ! nc -z plugin 8080; do
  echo "Waiting for the Plugin to launch (8080)."
  sleep 1
done
echo "The Plugin is up!"

sleep 5

cd $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma
echo "Running integration tests"
go get -v github.com/docker/docker/pkg/namesgenerator
go get -t -v ./...
go test -tags integration
