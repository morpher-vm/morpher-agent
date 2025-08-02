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
