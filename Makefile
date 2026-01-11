.PHONY: build test lint proto docker-up docker-down clean

build:
	go build -o bin/user-service ./cmd/user-service/main.go
	go build -o bin/task-service ./cmd/task-service/main.go
	go build -o bin/activity-service ./cmd/activity-service/main.go
	go build -o bin/api-gateway ./cmd/gateway/main.go

test:
	go test ./... -v

test-cover:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

proto:
	@if not exist internal\pb mkdir internal\pb
	protoc -I ./api \
		-I ./api/google/api \
		--go_out=./internal/pb --go_opt=paths=source_relative \
		--go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
		./api/user_api/user.proto ./api/task_api/task.proto ./api/activity_api/activity.proto ./api/models/*.proto
	protoc -I ./api \
		-I ./api/google/api \
		--grpc-gateway_out=./internal/pb \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt logtostderr=true \
		./api/user_api/user.proto ./api/task_api/task.proto ./api/activity_api/activity.proto

swagger:
	@if not exist docs\swagger mkdir docs\swagger
	protoc -I ./api \
		-I ./api/google/api \
		--openapiv2_out=./docs/swagger \
		--openapiv2_opt logtostderr=true \
		./api/user_api/user.proto ./api/task_api/task.proto ./api/activity_api/activity.proto

proto-all: proto swagger

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

deps:
	go mod download
	go mod tidy

run-user:
	CONFIG_PATH=./configs/user.yaml go run ./cmd/user-service/main.go

run-task:
	CONFIG_PATH=./configs/task.yaml go run ./cmd/task-service/main.go

run-activity:
	CONFIG_PATH=./configs/activity.yaml go run ./cmd/activity-service/main.go

run-gateway:
	CONFIG_PATH=./configs/gateway.yaml go run ./cmd/gateway/main.go

migrate-up:
	migrate -path ./migrations/user -database "postgres://admin:admin@localhost:5432/user_db?sslmode=disable" up
	migrate -path ./migrations/task -database "postgres://admin:admin@localhost:5433/task_db?sslmode=disable" up
	migrate -path ./migrations/activity -database "postgres://admin:admin@localhost:5434/activity_db_shard_0?sslmode=disable" up
	migrate -path ./migrations/activity -database "postgres://admin:admin@localhost:5435/activity_db_shard_1?sslmode=disable" up

migrate-down:
	migrate -path ./migrations/user -database "postgres://admin:admin@localhost:5432/user_db?sslmode=disable" down -all
	migrate -path ./migrations/task -database "postgres://admin:admin@localhost:5433/task_db?sslmode=disable" down -all
	migrate -path ./migrations/activity -database "postgres://admin:admin@localhost:5434/activity_db_shard_0?sslmode=disable" down -all
	migrate -path ./migrations/activity -database "postgres://admin:admin@localhost:5435/activity_db_shard_1?sslmode=disable" down -all

infra-up:
	docker compose up redis postgres-user postgres-task postgres-activity-shard-0 postgres-activity-shard-1 kafka -d

infra-down:
	docker compose down redis postgres-user postgres-task postgres-activity-shard-0 postgres-activity-shard-1 kafka

help:
	@echo "Available targets:"
	@echo "  build         - Build all services"
	@echo "  test          - Run all tests"
	@echo "  test-cover    - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  proto         - Generate proto files"
	@echo "  swagger       - Generate Swagger/OpenAPI documentation"
	@echo "  proto-all     - Generate proto files and Swagger docs"
	@echo "  docker-up     - Start all services in Docker"
	@echo "  docker-down   - Stop all Docker services"
	@echo "  docker-logs   - View Docker logs"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  run-user      - Run user service locally"
	@echo "  run-task      - Run task service locally"
	@echo "  run-activity  - Run activity service locally"
	@echo "  run-gateway   - Run API gateway locally"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  infra-up      - Start infrastructure only"
	@echo "  infra-down    - Stop infrastructure"

