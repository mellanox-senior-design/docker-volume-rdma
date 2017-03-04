#!/bin/bash -e
# -e : Exit immediately if a pipeline [...] returns a non-zero status.

if [ -z "$GOPATH" ]; then
    echo 'GOPATH is not set' >&2
    exit 2
fi

# change to docker-volume-rdma directory
cd ..

# download packages named by import path along with dependencies
# then install named packages like 'go install'
# '...' matches all subdirectories
go get -t -v ./...

GOOS=linux go build -v

# move executable binary to Dockerfile's context
mv docker-volume-rdma ./integration-testing

# change directory to integration-testing
cd integration-testing

# build Docker image named docker-volume-rdma-integration-test
docker build -t docker-volume-rdma-integration-test-image .

# remove executable binary
rm docker-volume-rdma

# start container with port 8080 exposed
docker run --rm -p 8080:8080 --name docker-volume-rdma-integration-test docker-volume-rdma-integration-test-image

# run tests that are tagged with integration
# go test -v -tags integration ../...
