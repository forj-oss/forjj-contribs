# This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

# So, update it for your need.
FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

COPY docker_files/*.sh /bin/

RUN apk update && \
    apk add --no-cache ca-certificates sudo && \
    update-ca-certificates --fresh && \
    rm -f /var/cache/apk/*tar.gz && \
    adduser devops devops -D && \
    chmod +xs /bin/update_user.sh && \
    chmod +x /bin/entrypoint.sh

# Required for DooD
RUN echo "devops ALL=(root:root) NOPASSWD:/bin/docker" >> /etc/sudoers.d/docker && \
    chmod 600 /etc/sudoers.d/docker

COPY templates/ /templates/

COPY jenkins /bin/

ENTRYPOINT ["/bin/entrypoint.sh"]

CMD ["--help"]
