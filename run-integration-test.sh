#!/bin/bash -xe
# -e : Exit immediately if a pipeline [...] returns a non-zero status.

if [ -z "$GOPATH" ]; then
    echo 'GOPATH is not set' >&2
    exit 2
fi

# download packages named by import path along with dependencies
# then install named packages like 'go install'
# '...' matches all subdirectories
go get -t -v ./...

GOOS=linux go build

# build Docker image named in-memory-test-image
docker build --file integration-in-memory-test.docker --tag in-memory-test-image .

# start container with port 8080 exposed
docker run --rm  -d -p 8080:8080 --name in-memory-test in-memory-test-image

#echo "TODO: add in-memory-test.go"
# run tests that are tagged with integration
# go test -v -tags in-memory-test ./...

# stop in-memory-test container
docker stop in-memory-test

docker run --rm -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=foo mysql

#  Wait for MySQL
#  http://stackoverflow.com/a/27601038/3259030
while ! nc -vz mysql 3306; do
    echo "Waiting for MySQL to launch (3306)."
    sleep 1
done
echo "MySQL is up!"

sleep 5

mysql -u root -p foo << EOF
  CREATE SCHEMA rdma-volumes;
EOF

# build Docker image named in-memory-test-image
docker build --file integration-mysql-test.docker --tag mysql-test-image .

# start container with port 8080 exposed // TODO: add -d later
docker run --rm -p 8080:8080 --link mysql --name mysql-test mysql-test-image

# echo "TODO: add mysql-test.go"
# run tests that are tagged with integration
# go test -v -tags mysql-test ./...

# stop mysql-test container
# docker stop mysql-test

# stop mysql container
# docker stop mysql

# remove executable binary
rm docker-volume-rdma
