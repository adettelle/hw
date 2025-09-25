BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
PROTO_DIR := "./api/proto"
OUT_DIR := api

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB="postgresql://postgres:123456@localhost:9999/test_db?sslmode=disable"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

migrate-create:
	./migrate create -ext sql -dir internal/migrator/migration/ -seq $(name)

migrate-up:
	./migrate -path internal/migrator/migration/ -database $(DB) -verbose up

migrate-down:
	./migrate -path internal/migrator/migration/ -database $(DB) -verbose down

generate:
	protoc -I$(PROTO_DIR) \
		--go_out=paths=source_relative:$(OUT_DIR) \
  		--go-grpc_out=paths=source_relative:$(OUT_DIR) \
  		./api/proto/event.proto ./api/proto/event_service.proto  