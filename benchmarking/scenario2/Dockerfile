# FROM ubuntu
FROM golang:1.7.5

# RUN apt-get -y update && apt-get install -y golang && go get -u github.com/rakyll/hey
RUN go get -u github.com/rakyll/hey

# CMD ["hey" "-m=GET" "-disable-compression" "http://localhost:4200/?p=20"]
CMD ["hey", "-m=GET", "-disable-compression", "http://google.com"]
