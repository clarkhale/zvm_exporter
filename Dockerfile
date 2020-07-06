ARG ARCH="s390x"
ARG OS="linux"
FROM        quay.io/prometheus/busybox:latest
LABEL maintainer="Clark Hale <chale@redhat.com>"

ARG ARCH="s390x"
ARG OS="linux"
COPY ./zvm_exporter /bin/zvm_exporter

USER nobody
ENTRYPOINT ["/bin/zvm_exporter"]
EXPOSE     9100