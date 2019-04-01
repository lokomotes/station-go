FROM golang:1.11-alpine

COPY . /go/.

RUN /bin/sh /go/install.sh

ENTRYPOINT [ "main" ]