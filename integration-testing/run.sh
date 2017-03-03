#!/bin/bash -e
# Exit immediately if a pipeline [...] returns a non-zero status.

# change to docker-volume-rdma directory
cd ..

# download packages named by import path along with dependencies
# then install named packages like 'go install'
# '...' matches all subdirectories
go get -t -d -v ./...

# build project
go build -v

# move executable binary to Dockerfile's context
mv docker-volume-rdma ./integration-testing/DockerVolumeRDMA

# change directory to integration-testing
cd integration-testing

# build Docker image named docker-volume-rdma-integration-test
docker build -t docker-volume-rdma-integration-test-image .

# start container with port 8080 exposed
docker run --rm -p 8080:8080 --name docker-volume-rdma-integration-test docker-volume-rdma-integration-test-image

# run tests that are tagged with integration
# go test -v -tags integration ./...
