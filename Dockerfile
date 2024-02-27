FROM golang:1.16-alpine3.14 as builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:lates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
