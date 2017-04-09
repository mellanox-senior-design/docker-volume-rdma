#! /bin/bash

echo "Building..."
(cd .. && rm -f docker-volume-rdma && GOOS=linux go build)

rm plugin/docker-volume-rdma
cp ../docker-volume-rdma plugin

rm tester/*.go
rm -r tester/db
rm -r tester/drivers
cp ../main.go tester
cp ../main_integration_test.go tester
cp -r ../db  tester
cp -r ../drivers tester

docker-compose up --abort-on-container-exit --build
docker-compose down --volumes
