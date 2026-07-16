# Прокси для сборки. По умолчанию пустые (прямое соединение).
# На сервере с Squid задайте через docker compose build --build-arg или в .env.
ARG HTTP_PROXY=
ARG HTTPS_PROXY=
ARG http_proxy=
ARG https_proxy=

FROM golang:1.24-alpine AS builder

# Прокси доступны только на этапе сборки.
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG http_proxy
ARG https_proxy
ENV HTTP_PROXY=$HTTP_PROXY \
    HTTPS_PROXY=$HTTPS_PROXY \
    http_proxy=$http_proxy \
    https_proxy=$https_proxy

WORKDIR /src

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /out/api ./cmd

FROM alpine:3.20

# Прокси нужен только для apk install в финальном образе. После — не сохраняется.
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG http_proxy
ARG https_proxy
ENV HTTP_PROXY=$HTTP_PROXY \
    HTTPS_PROXY=$HTTPS_PROXY \
    http_proxy=$http_proxy \
    https_proxy=$https_proxy

WORKDIR /app

RUN apk add --no-cache ca-certificates wget && update-ca-certificates \
	&& adduser -D -u 10001 appuser

# Сбрасываем прокси в финальном образе, чтобы рантайм ходил в сеть напрямую.
ENV HTTP_PROXY= \
    HTTPS_PROXY= \
    http_proxy= \
    https_proxy=

COPY --from=builder /out/api /app/api
COPY --from=builder /src/docs /app/docs

ENV APP_PORT=8080
EXPOSE 8080

USER appuser
ENTRYPOINT ["/app/api"]
