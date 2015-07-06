FROM golang:1.4.2
MAINTAINER YI-HUNG JEN <yihungjen@macrodatalab.com>

RUN go get gopkg.in/redis.v3
RUN go get github.com/streadway/amqp

COPY bigobject-registry/worker /go/src/github.com/yihungjen/bigobject-registry/worker
COPY registry-driver /go/src/github.com/yihungjen/registry-driver

RUN go install github.com/yihungjen/bigobject-registry/worker

ENV REDIS_NETLOC localhost:6379
ENV ACCEPT_COMMAND set,del,expired

ENV AMQP_URI amqp://localhost:5672

CMD ["worker"]
