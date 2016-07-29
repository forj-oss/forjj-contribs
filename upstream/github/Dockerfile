FROM golang:1.6.2-alpine

RUN apk update && apk add --no-cache git && mkdir /go/src/github

ENV TERM vt100

COPY * /go/src/github/
COPY entrypoint.sh /tmp

WORKDIR /go/src/github

RUN go get -insecure && \
    cd $GOPATH/src/github.hpe.com/christophe-larsonneur/go-forjj/cmd/genflags/ && \
    go get && \
    cd - && \
    go generate && \
    go install

RUN adduser devops devops -D

USER devops

ENTRYPOINT [ "/tmp/entrypoint.sh" ]
