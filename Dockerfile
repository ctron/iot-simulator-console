FROM fedora:29

MAINTAINER Jens Reimann <jreimann@redhat.com>
LABEL maintainer="Jens Reimann <jreimann@redhat.com>"

EXPOSE 8080

RUN dnf -y update
RUN dnf -y install nodejs golang

RUN mkdir -p /src
ADD . /

RUN go mod download

RUN \
    npm install && \
    npm run build && \
    go build -o /iot-simulator-console .

ENTRYPOINT /iot-simulator-console
