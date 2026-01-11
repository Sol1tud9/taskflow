FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /activity-service ./cmd/activity-service/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /activity-service .
COPY configs/activity.yaml ./configs/

EXPOSE 50051 8080

CMD ["./activity-service"]

