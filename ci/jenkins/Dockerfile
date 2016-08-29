# This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

# So, update it for your need.
FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

COPY templates/ /templates/

RUN apk update &&     apk add --no-cache ca-certificates &&     update-ca-certificates --fresh &&     rm -f /var/cache/apk/*tar.gz &&     adduser devops devops -D

COPY jenkins /bin/jenkins

USER devops

ENTRYPOINT ["/bin/jenkins"]

CMD ["--help"]
