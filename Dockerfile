# Demo for accelerating containers over RDMA
FROM golang:1.7.5

WORKDIR /go/src/github.com/mellanox-senior-design/docker-volume-rdma
ENTRYPOINT ["go", "run", "main.go", "-logtostderr=true"]
CMD []

COPY . /go/src/github.com/mellanox-senior-design/docker-volume-rdma

RUN go get -v -t ./...
RUN go test ./... -cover
