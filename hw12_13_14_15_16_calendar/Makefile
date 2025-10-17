BIN := "./bin"
DOCKER_IMG="calendar:develop"
PROTO_DIR := "./api/proto"
OUT_DIR := api

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB="postgresql://postgres:123456@localhost:9999/test_db?sslmode=disable"

build_calendar:
	go build -v -o "$(BIN)/calendar" -ldflags "$(LDFLAGS)" ./cmd/calendar

build_scheduler:
	go build -v -o "$(BIN)/calendar_scheduler" -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

build_sender:
	go build -v -o "$(BIN)/calendar_sender" -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

build: build_calendar build_scheduler build_sender

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
	go test -race $(shell go list ./... | grep -v /integration_tests)

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

# разворачивание и запуск
up: 
	docker compose up --build --force-recreate

up-calendar:
	docker compose up calendar -d --build

up-sender:
	docker compose up sender -d --build

up-scheduler:
	docker compose up scheduler -d --build

down:
	docker compose down

integration-tests:
	docker compose -f ./docker-compose-test.yaml up integration-tests --build && docker compose -f ./docker-compose-test.yaml down --volumes
