FROM fedora:29

MAINTAINER Jens Reimann <jreimann@redhat.com>
LABEL maintainer="Jens Reimann <jreimann@redhat.com>"

EXPOSE 8080

RUN dnf -y update
RUN dnf -y install nodejs golang
RUN go version

ENV \
    GOPATH=/go

RUN mkdir -p /go/src/github.com/ctron
ADD . /go/src/github.com/ctron/iot-simulator-console

RUN \
    cd /go/src/github.com/ctron/iot-simulator-console && \
    GO111MODULE=on go mod vendor && \
    npm install && \
    npm run build && \
    cd cmd && \
    go build -o /iot-simulator-console . && \
    mv build && \
    true

ENTRYPOINT /iot-simulator-console
