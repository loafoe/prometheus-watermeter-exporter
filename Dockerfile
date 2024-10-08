FROM golang:1.23.2 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN git rev-parse --short HEAD
RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 go build -o prometheus-watermeter-exporter -ldflags "-X main.GitCommit=${GIT_COMMIT}"

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /
COPY --from=builder /build/prometheus-watermeter-exporter /usr/bin
EXPOSE 8880