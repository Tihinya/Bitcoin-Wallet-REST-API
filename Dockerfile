FROM golang:1.21.1-alpine3.17 as builder

WORKDIR /app

COPY . .


RUN go mod download cloud.google.com/go/storage

RUN go build -o main .

FROM alpine:3.17

WORKDIR /app
COPY  --from=builder /app/.env .
COPY  --from=builder /app/dev_config.json .

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]



