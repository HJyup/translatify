FROM golang:1.22.2 AS builder
WORKDIR /app

COPY . .
COPY ../common ../common

WORKDIR /app/chat
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o chat ./cmd

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/chat .
EXPOSE 8080
CMD ["./chat"]