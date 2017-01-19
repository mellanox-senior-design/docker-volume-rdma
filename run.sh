#! /bin/bash -xe
cd $GOPATH/src/github.com/Jacobingalls/docker-volume-rdma
go run main.go -logtostderr=true $@
