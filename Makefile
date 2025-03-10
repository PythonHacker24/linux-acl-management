# Makefile for building Go backend server

# Variables
APP_NAME = backend-server
BUILD_DIR = build
SRC_DIR = backend-server/main.go

# Go commands
GO = go
GO_BUILD = $(GO) build

# Paths
BINARY = $(BUILD_DIR)/$(APP_NAME)

.PHONY: all clean build run

all: build

# Build the Go application
build:
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) -o $(BINARY) $(SRC_DIR)
	@echo "Build completed: $(BINARY)"

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	@echo "Cleaned build artifacts."

# Run the compiled binary
run: build
	$(BINARY)
