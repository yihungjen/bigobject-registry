FROM golang:1.4.2
MAINTAINER YI-HUNG JEN <yihungjen@macrodatalab.com>

RUN go get gopkg.in/redis.v3

COPY bigobject-registry/registry /go/src/github.com/yihungjen/bigobject-registry/registry
COPY registry-driver /go/src/github.com/yihungjen/registry-driver

RUN go install github.com/yihungjen/bigobject-registry/registry

ENV REDIS_NETLOC localhost:6379
ENV ACCEPT_COMMAND append,expired

EXPOSE 80
CMD ["registry"]
