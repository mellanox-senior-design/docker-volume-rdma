FROM golang:1.7.5

COPY run.sh run.sh
COPY bench.py bench.py
COPY sql/db /go/flyway-4.1.1/sql/db

RUN apt update && \
    apt install -y netcat && \
    apt-get autoremove && \
    apt-get clean && \
    go get -u github.com/rakyll/hey

RUN curl https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/4.1.1/flyway-commandline-4.1.1-linux-x64.tar.gz >> fly.tar.gz && \
    tar -xvf fly.tar.gz && \
    rm fly.tar.gz

CMD ["./run.sh"]
