#!/bin/bash

GIT_BRANCH=$(
	git rev-parse --abbrev-ref HEAD |
	tr '[:upper:]' '[:lower:]' |
	sed 's#[^a-z0-9._-]#-#'
)

docker build --tag docker-volume-rdma:"$GIT_BRANCH" .
docker rmi docker-volume-rdma:"$GIT_BRANCH"
./benchmarking/scenarios/test.sh

