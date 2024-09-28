FROM golang:1.23.1-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-service ./cmd/main.go

FROM alpine:latest

# RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/todo-service .

COPY config.toml .

CMD ["./todo-service"]

EXPOSE 8080
