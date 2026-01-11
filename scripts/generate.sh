#!/bin/bash

cd "$(dirname "$0")/.." || exit

mkdir -p ./internal/pb
mkdir -p ./docs/swagger

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

protoc -I ./api \
  -I ./api/google/api \
  --openapiv2_out=./docs/swagger \
  --openapiv2_opt logtostderr=true \
  ./api/user_api/user.proto ./api/task_api/task.proto ./api/activity_api/activity.proto

