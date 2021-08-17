FROM registry.trendyol.com/platform/base/image/golang:1.16.0-alpine3.13 as build

WORKDIR /kink-workspace

RUN apk add --no-cache git make

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN make cosign

FROM gcr.io/distroless/static:nonroot-amd64

WORKDIR /kink-workspace

COPY --from=build --chown=nonroot:nonroot /kink-workspace/kink /usr/local/bin/kink

USER nonroot
ENTRYPOINT ["kink"]


