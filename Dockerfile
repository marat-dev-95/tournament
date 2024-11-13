
FROM golang:1.23 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y gcc

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
#GOOS=linux GOARCH=amd64
RUN  go build -o /app/main /app/cmd/app

FROM debian:bookworm-slim

WORKDIR /build

COPY --from=builder /app/main ./

COPY --from=builder /app/db ./db

RUN chmod +x /build/main

CMD ["/build/main"]

