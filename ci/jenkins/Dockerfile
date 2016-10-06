# This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

# So, update it for your need.
FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

RUN apk update &&     apk add --no-cache ca-certificates sudo &&     update-ca-certificates --fresh &&     rm -f /var/cache/apk/*tar.gz &&     adduser devops devops -D

# Required for DooD
RUN echo "devops ALL=(root:root) NOPASSWD:/bin/docker,NOPASSWD:/bin/docker-config-update.sh" >> /etc/sudoers.d/docker && chmod 600 /etc/sudoers.d/docker

COPY templates/ /templates/

COPY jenkins docker-config-update.sh /bin/

USER devops

ENTRYPOINT ["/bin/jenkins"]

CMD ["--help"]
