.EXPORT_ALL_VARIABLES:

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://calendar:calendar@localhost:15432/calendar
GOOSE_MIGRATION_DIR=./migrations

BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_STORER := "./bin/calendar_storer"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/config.release="develop" -X github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/config.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/config.gitHash=$(GIT_HASH)

migrate-prepare:
	go install github.com/pressly/goose/v3/cmd/goose@latest
migrate-status:
	goose status
migrate-up:
	goose up
migrate-down:
	goose down
migrate: migrate-prepare migrate-up

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_STORER) -ldflags "$(LDFLAGS)" ./cmd/calendar_storer

run: build
	$(BIN) -config ./configs/calendar_config.yml

run_scheduler: build
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yml

run_storer: build
	$(BIN_STORER) -config ./configs/storer_config.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

generate:
	go generate ./...
	
test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run --timeout 3m0s ./...

.PHONY: build run build-img run-img version test lint
