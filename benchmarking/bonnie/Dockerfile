FROM python:2-onbuild

RUN mkdir /test

RUN apt-get update && \
    apt-get install -y bonnie++ && \
    apt-get autoremove && \
    apt-get clean

COPY run.sh run.sh

CMD ["./run.sh"]
