FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/wind-watcher/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]