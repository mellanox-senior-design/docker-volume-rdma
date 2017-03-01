FROM python:2-onbuild

COPY run.sh run.sh
COPY sql sql

ADD https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/4.1.1/flyway-commandline-4.1.1-linux-x64.tar.gz fly.tar.gz
RUN tar -xvf fly.tar.gz


CMD ["./run.sh"]
