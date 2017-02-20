FROM ubuntu

RUN apt update && apt install bonnie++ && mkdir -p /tmp

CMD bonnie++ -d /tmp -r 2048 -u root
