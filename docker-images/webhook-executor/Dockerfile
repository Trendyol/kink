FROM golang:alpine3.10 AS build
WORKDIR /go/src/github.com/adnanh/webhook
ARG WEBHOOK_VERSION=2.6.10
RUN apk add --update -t build-deps curl libc-dev gcc libgcc
RUN curl -Lsf https://github.com/adnanh/webhook/archive/$WEBHOOK_VERSION.tar.gz | tar xzv --strip 1 -C . && \
    go get -d && \
    go build -o /usr/local/bin/webhook && \
    apk del --purge build-deps && \
    rm -rf /var/cache/apk/* && \
    rm -rf /go

FROM alpine:3.10
RUN apk add --update curl tini
COPY --from=build /usr/local/bin/webhook /usr/local/bin/webhook
RUN curl -Lso /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.15.5/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl
VOLUME ["/etc/webhook"]
EXPOSE 9000
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/webhook"]
