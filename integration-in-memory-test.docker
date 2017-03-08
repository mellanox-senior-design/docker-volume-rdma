FROM scratch

COPY docker-volume-rdma /

ENTRYPOINT ["/docker-volume-rdma", "-logtostderr=true", "-db=in-memory"]
