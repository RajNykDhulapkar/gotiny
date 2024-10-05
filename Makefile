BINARY_NAME=gotiny
PACKAGE=github.com/RajNykDhulapkar/gotiny
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=1

# Default target
all: run

# Run the application
run:
	@echo "Starting the GoTiny application..."
	@go run main.go

# Build the application
build:
	@echo "Building the GoTiny application..."
	@go build -o ./bin/$(BINARY_NAME) main.go
	@echo "Build complete! Binary generated: ./bin/$(BINARY_NAME)"

# Test the application
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean up generated files and binaries
clean:
	@echo "Cleaning up binaries and temporary files..."
	@go clean
	@rm -f $(BINARY_NAME)

# Docker compose to start Redis (Assumes you have a docker-compose.yml)
start-redis:
	@echo "Starting Redis container with Docker..."
	@docker run --name gotiny-redis -p $(REDIS_PORT):$(REDIS_PORT) -d redis

stop-redis:
	@echo "Stopping Redis container..."
	@docker stop gotiny-redis
	@docker rm gotiny-redis

# Lint the codebase
lint:
	@echo "Linting the codebase..."
	@golangci-lint run ./...

# Format the Go code
format:
	@echo "Formatting Go code..."
	@go fmt ./...

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod tidy

# Rebuild everything
rebuild: clean build

# Help
help:
	@echo "Makefile commands:"
	@echo "  run          - Run the application"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean up binaries and temporary files"
	@echo "  start-redis  - Start a Redis instance using Docker"
	@echo "  stop-redis   - Stop the Redis Docker container"
	@echo "  lint         - Lint the codebase"
	@echo "  format       - Format the Go code"
	@echo "  deps         - Install dependencies"
	@echo "  rebuild      - Clean and rebuild the project"
	@echo "  help         - Show this help message"
