PROJECT_NAME := pr-reviewer-service
BINARY_NAME := pr-reviewer-service

.PHONY: run build test docker-build docker-run docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

docker-build:
	docker build -t $(PROJECT_NAME) .

docker-run: docker-build
	docker run --rm -p 8080:8080 $(PROJECT_NAME)

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
