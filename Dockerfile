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

WORKDIR /go/src/github.com/ctron/iot-simulator-console

RUN GO111MODULE=on go mod vendor

RUN npm install
RUN npm run build
RUN cd cmd && go build -o /iot-simulator-console .
RUN mv build /

WORKDIR /

ENTRYPOINT /iot-simulator-console
