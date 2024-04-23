FROM golang:1.22.0 as build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build

FROM debian:12-slim
LABEL org.opencontainers.image.source https://github.com/sikalabs/github-apps-pull-secret-sync
COPY --from=build /build/github-apps-pull-secret-sync /usr/local/bin/github-apps-pull-secret-sync
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
