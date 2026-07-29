.PHONY: build test lint run docker-build tidy

BINARY_NAME=user-service
DOCKER_IMAGE=ghcr.io/release-orchestrator/user-service
VERSION?=dev

build:
	go build -o bin/$(BINARY_NAME) ./cmd

test:
	go test -v ./...

lint:
	golangci-lint run ./...

run: build
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/user_db?sslmode=disable" ./bin/$(BINARY_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

tidy:
	go mod tidy