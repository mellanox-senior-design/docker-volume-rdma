#! /bin/bash -xe
if [ -z "$GOPATH" ]; then
    echo 'GOPATH is not set' >&2
    exit 2
fi

cd "$GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma" || exit 1
go run -ldflags -s main.go -logtostderr=true "$@"
