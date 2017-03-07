FROM golang:1.7.5

COPY run.sh run.sh

RUN apt update && \
    apt install -y netcat && \
    apt-get autoremove && \
    apt-get clean && \
    go get -u github.com/rakyll/hey

CMD ["./run.sh"]
