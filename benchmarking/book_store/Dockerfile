FROM python:2-onbuild

COPY run.sh run.sh
COPY sql/db /usr/src/app/flyway-4.1.1/sql/db

RUN apt update  && apt install -y netcat && apt-get autoremove && apt-get clean

RUN curl https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/4.1.1/flyway-commandline-4.1.1-linux-x64.tar.gz >> fly.tar.gz && \
    tar -xvf fly.tar.gz && \
    rm fly.tar.gz


CMD ["./run.sh"]
