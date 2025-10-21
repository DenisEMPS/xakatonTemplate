# syntax=docker/dockerfile:1.7
FROM golang:1.25-alpine AS builder

ARG BUILD_PACKAGE=./cmd/app

ENV GOPRIVATE=git.a7ru.app \
    GONOSUMDB=git.a7ru.app \
    GIT_TERMINAL_PROMPT=0 \
    GOPROXY=direct \
    CGO_ENABLED=0 \
    GOOS=linux

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
COPY vendor/ ./vendor/

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/


RUN mkdir -p /build/app/config && \
    go build -mod=vendor -trimpath -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty)" \
    -o /build/app/app "$BUILD_PACKAGE"

FROM alpine:3.20

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates && \
    mkdir -p /app /migrations /app/data && \
    chown -R root:root /app /migrations && \
    chmod -R 555 /app /migrations && \
    chown -R appuser:appgroup /app/data && \
    chmod -R 770 /app/data

COPY --from=builder --chown=root:root --chmod=555 /build/app /app
COPY --chown=root:root --chmod=555 migrations /migrations

USER appuser

WORKDIR /app

EXPOSE 3001

CMD ["./app"]