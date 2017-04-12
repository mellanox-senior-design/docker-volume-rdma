# Demo for accelerating containers over RDMA
FROM golang:1.7.5

WORKDIR /go/src/github.com/mellanox-senior-design/docker-volume-rdma
ENTRYPOINT ["go", "run", "main.go", "-logtostderr=true"]
CMD []

RUN go get -v \
    github.com/docker/docker/pkg/namesgenerator \
    github.com/docker/go-plugins-helpers/volume \
    github.com/mattn/go-sqlite3

COPY . /go/src/github.com/mellanox-senior-design/docker-volume-rdma

RUN go get -v -t ./...
RUN go test ./... -cover
