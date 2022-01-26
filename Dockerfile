# Copy from https://t.me/golangID/138064
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /main

FROM alpine:latest

WORKDIR /app

COPY assets ./assets
COPY data ./data
COPY --from=builder /main /main

EXPOSE 8080

ENTRYPOINT ["/main"]

