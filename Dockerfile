FROM fedora:29

MAINTAINER Jens Reimann <jreimann@redhat.com>
LABEL maintainer="Jens Reimann <jreimann@redhat.com>"

EXPOSE 8080

ENV \
    GOPATH=/go

RUN mkdir -p /go/src/github.com/ctron
ADD . /go/src/github.com/ctron/iot-simulator-console

RUN \
    dnf -y update && \
    dnf -y install nodejs golang && \
    go version && \
    cd /go/src/github.com/ctron/iot-simulator-console && \
    npm install && \
    npm run build && \
    cd cmd && \
    go build -o /iot-simulator-console . && \
    cd .. && \
    mv build / && \
    echo "Clean up" && \
    rm -Rf go && \
    dnf -y history undo last && dnf -y clean all && \
    true

ENTRYPOINT /iot-simulator-console
