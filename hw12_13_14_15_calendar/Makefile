BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

### Docker ###
PROJECT_PATH := $(shell dirname $(shell pwd))
USER_ID := $(shell id -u)
GROUP_ID := $(shell id -g)
COMPOSE := docker-compose
DOCKER_COMPOSE = PROJECT_PATH=$(PROJECT_PATH) USER_ID=$(USER_ID) GROUP_ID=$(GROUP_ID) $(COMPOSE)

.PHONY: up
up: docker-build docker-up

.PHONY: down
down: docker-down

.PHONY: docker-build
docker-build:
	$(DOCKER_COMPOSE) build calendar

.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) up -d --remove-orphans

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) down --remove-orphans
### Docker ###

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
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
