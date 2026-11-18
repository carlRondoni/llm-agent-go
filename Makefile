.PHONY: help local-init local-docker-start local-docker-build-and-start local-docker-post-start \
        local-complete-first-start local-docker-stop local-docker-restart local-docker-restart-and-build \
        docker-start docker-start-build docker-stop docker-restart docker-pull-model \
        test-unit test-integration test test-coverage

# same precedence as docker compose: environment, then .env, then llama3
OLLAMA_MODEL ?= $(shell sed -n 's/^OLLAMA_MODEL=//p' .env 2>/dev/null | tr -d '\r"')
ifeq ($(strip $(OLLAMA_MODEL)),)
OLLAMA_MODEL := llama3
endif
# the integration tests read it too
export OLLAMA_MODEL

help:
	@echo "Available commands:"
	@echo "  make local-complete-first-start  - Create .env if missing, start all Docker services and pull the model"
	@echo "  make docker-start                - Start Docker compose services in background"
	@echo "  make docker-start-build          - Rebuild the API image and start Docker compose services"
	@echo "  make docker-stop                 - Stop Docker compose services"
	@echo "  make docker-pull-model           - Pull OLLAMA_MODEL (environment, then .env, then llama3) into the Ollama container"
	@echo "  make test-unit                   - Run unit tests (//go:build unit)"
	@echo "  make test-integration            - Rebuild and start the stack, then run integration tests (//go:build integration)"
	@echo "  make test                        - Run all tests (unit + integration)"
	@echo "  make test-coverage               - Run unit tests with coverage report"

local-init:
	cp -n .env.example .env 2>/dev/null || true

local-docker-start: docker-start
local-docker-build-and-start: docker-start-build
local-docker-post-start: docker-pull-model
local-complete-first-start: local-init docker-start docker-pull-model
local-docker-stop: docker-stop
local-docker-restart: docker-restart
local-docker-restart-and-build: docker-stop docker-start-build

docker-start: local-init
	docker compose up -d

docker-start-build: local-init
	docker compose up --build -d

docker-pull-model:
	docker exec ollama ollama pull $(OLLAMA_MODEL)

docker-stop:
	docker compose down

docker-restart: docker-stop docker-start

test-unit:
	go test -v -tags=unit ./...

test-integration: docker-start-build docker-pull-model
	go test -v -p 1 -tags=integration ./...

test: docker-start-build docker-pull-model
	go test -v -p 1 -tags="unit,integration" ./...

test-coverage:
	go test -tags=unit -covermode=count -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
