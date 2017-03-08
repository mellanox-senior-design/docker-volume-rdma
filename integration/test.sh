#! /bin/bash

echo "Building..."
(cd .. && rm -f docker-volume-rdma && GOOS=linux go build)

rm plugin/docker-volume-rdma
cp ../docker-volume-rdma plugin

docker-compose up --abort-on-container-exit --build
docker-compose down --volumes
