FROM alpine

COPY plugin_run.sh /
COPY docker-volume-rdma /
RUN apk update && apk add netcat-openbsd bash mysql-client

ENTRYPOINT ["./plugin_run.sh"]
