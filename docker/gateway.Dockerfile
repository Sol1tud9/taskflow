FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/gateway/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /api-gateway .
COPY configs/gateway.yaml ./configs/
COPY docs/swagger ./docs/swagger

EXPOSE 8080

CMD ["./api-gateway"]

