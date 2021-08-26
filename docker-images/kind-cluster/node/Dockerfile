ARG KUBERNETES_VERSION="1.21.2"
FROM kindest/node:v${KUBERNETES_VERSION}

RUN mv /usr/local/bin/entrypoint /usr/local/bin/entrypoint-original
COPY entrypoint-wrapper.sh /usr/local/bin/entrypoint