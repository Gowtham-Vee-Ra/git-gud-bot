.PHONY: all build test clean run docker-build docker-run

# Build variables
BINARY_NAME=codereview
DOCKER_IMAGE=codereview-api

all: clean build

build:
	go build -o bin/$(BINARY_NAME) cmd/api/main.go

test:
	go test -v ./...

clean:
	go clean
	rm -rf bin/

run:
	go run cmd/api/main.go

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker-compose up

lint:
	golangci-lint run

generate:
	go generate ./...

deps:
	go mod tidy
	go mod verify

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down