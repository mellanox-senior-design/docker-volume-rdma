# Demo for accelerating containers over RDMA
FROM golang:1.7.5

WORKDIR /go/src/app
ENTRYPOINT ["go", "run", "main.go", "-logtostderr=true"]
CMD []

COPY . /go/src/app

RUN go get -v -t ./...
RUN go test ./... -cover
