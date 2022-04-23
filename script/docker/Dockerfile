FROM golang:1.16-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

EXPOSE 8080

WORKDIR /app/
COPY --from=builder /app/bin/bot /app/bot
COPY --from=builder /app/config /app/config

ENTRYPOINT ["./bot"]
