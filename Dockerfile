FROM golang:1.16-alpine3.14 as builder

WORKDIR /app

COPY . .


RUN go mod download cloud.google.com/go/storage

RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY  --from=builder /app/db/migrations .
COPY  --from=builder /app/.env .
COPY  --from=builder /app/dev_config.json .


COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]



