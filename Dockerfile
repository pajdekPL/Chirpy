# Build stage
FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chirpy

# Final runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/chirpy .
COPY --from=builder /app/static ./static
COPY --from=builder /app/docs ./static


EXPOSE 8080

CMD ["./chirpy"]
