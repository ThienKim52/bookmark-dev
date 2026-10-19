.PHONY: all build run test test-unit test-integration clean swagger fmt lint deps help docker-build docker-run docker-compose-up
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_TAG := latest
IMG_NAME=kimthanthien/bookmark
ifneq ($(GIT_TAG),)
   IMG_TAG := $(GIT_TAG)
endif
export IMG_TAG

# Variables
APP_NAME=health-check
CMD_DIR=./cmd/api
BUILD_DIR=./build
BINARY_NAME=$(BUILD_DIR)/$(APP_NAME)
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*' -not -path './build/*')

# Default target
all: swagger run

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Build complete: $(BINARY_NAME)"

## run: Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(CMD_DIR)/main.go

## test: Run all tests
test-all:
	go test ./... -coverprofile="coverage.out"
	go tool cover -html=coverage.out -o coverage.html

test:
	@echo "Running all tests..."
	@go test -v -coverprofile="coverage.out" ./...
	go tool cover -html=coverage.out -o coverage.html


## test-unit: Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v ./internal/... -run "^((?!integration).)*$$"

## test-integration: Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -v ./internal/integration_test/...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@go clean
	@echo "Clean complete"

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

## deps-verify: Verify dependencies
deps-verify:
	@echo "Verifying dependencies..."
	@go mod verify

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

COVERAGE_EXCLUDE=mocks|main.go|docs|test|config.go
COVERAGE_THRESHOLD = 50
test-cover:
	go test ./... -coverprofile="coverage.tmp" -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	  echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	  exit 1; \
	else \
	  echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi


## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(IMG_NAME):$(IMG_TAG) .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 bookmark-service:dev

COVERAGE_FOLDER=./test-output
docker-test:
	mkdir -p $(COVERAGE_FOLDER) 
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --build-arg COVERAGE_THRESHOLD="$(COVERAGE_THRESHOLD)" --progress=plain --target test -t test:test --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	  echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	  exit 1; \
	else \
	  echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

docker-compose-up:
	docker compose up --build

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker-push:
	docker push $(IMG_NAME):$(IMG_TAG)

generate-rsa-key:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private.pem -out public.pem
