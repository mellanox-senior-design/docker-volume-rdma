FROM alpine

COPY tester_run.sh /
RUN apk update && apk add netcat-openbsd bash mysql-client

CMD ["./tester_run.sh"]
