BIN_DIR := bin

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Install dependencies.
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go version
	go mod tidy
	go mod vendor

# Lint the code.
.PHONY: lint
lint:
	golangci-lint run

# Format the code.
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w .

# Build the application.
.PHONY: build
build: deps
	@echo "Building morpher-agent..."
	go build -o $(BIN_DIR)/morpher-agent main.go
	@echo "Done"

# Run all tests.
.PHONY: test
test:
	@echo "Running tests..."

# Clean up test artifacts.
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

.DEFAULT_GOAL := build
