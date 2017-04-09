FROM golang:alpine
RUN apk update && apk add netcat-openbsd bash mysql-client git gcc g++

COPY tester_run.sh /
COPY main.go $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma/
COPY main_integration_test.go $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma/
COPY db $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma/db
COPY drivers $GOPATH/src/github.com/mellanox-senior-design/docker-volume-rdma/drivers

CMD ["/tester_run.sh"]
