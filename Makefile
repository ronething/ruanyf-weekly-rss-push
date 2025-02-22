# Binary name
BINARY_NAME=rss-push

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Build directory
BUILD_DIR=./bin

# Main package path
MAIN_PACKAGE=./cmd/rss-push

# Build Flags
LDFLAGS=-s -w

.PHONY: all build clean test run

all: clean build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Testing..."
	@go test -v ./...

run: build
	@echo "Running..."
	@./$(BUILD_DIR)/$(BINARY_NAME)
