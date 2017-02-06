#! /bin/bash -xe
if [ -z "$GOPATH" ]; then
    echo 'GOPATH is not set' >&2
    exit 2
fi

cd $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma
go run main.go -logtostderr=true $@
