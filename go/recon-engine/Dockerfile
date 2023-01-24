FROM golang:1.19.3-alpine AS builder

RUN apk update && \
    apk add build-base git && \
    mkdir -p /app

WORKDIR /app
COPY ./ ./

RUN go clean -modcache && \
    go mod tidy && \
    go build -o /app/recon /app/cmd/recon

FROM alpine:3.16.3

ARG UID=1000
ARG GID=1000

RUN apk update && \
    apk add shadow

RUN mkdir -p /app && \
    useradd -d /app -r -u ${UID} app && \
    chown app:${GID} /app

WORKDIR /app

COPY --from=builder /app/recon /usr/local/bin/recon

ENTRYPOINT ["recon"]
    
