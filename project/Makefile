SOURCE_FILE := "cmd/main.go"
BUILD_OUTPUT := "./bin/project-app"
DOCKER_FILE := "build/Dockerfile"
DOCKER_IMG := "project:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# Build commands
build:
	go build -o $(BUILD_OUTPUT) $(SOURCE_FILE)

run: build
	$(BUILD_OUTPUT)

test:
	go test -race ./internal/...

# Deploy commands
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f $(DOCKER_FILE) .

run-img: build-img
	docker run $(DOCKER_IMG)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up-dev: build-img
	docker-compose --env-file=./deployments/.env -f ./deployments/docker-compose.dev.yaml up

.PHONY: build run build-img run-img version test lint up-dev
