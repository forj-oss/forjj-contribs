FROM alpine:latest

COPY github /bin/github

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates --fresh && \
    rm -f /var/cache/apk/*tar.gz && \
    adduser devops devops -D

USER devops

ENTRYPOINT ["/bin/jenkins"]

CMD ["--help"]
