# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /out/api ./cmd/api

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates wget && update-ca-certificates \
	&& adduser -D -u 10001 appuser

COPY --from=builder /out/api /app/api
COPY --from=builder /src/docs /app/docs

ENV APP_PORT=8080
EXPOSE 8080

USER appuser
ENTRYPOINT ["/app/api"]

