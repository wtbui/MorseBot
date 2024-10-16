# Define variables
PROJECT_DIR := $(CURDIR)
BIN_DIR := $(PROJECT_DIR)/bin
BINARY_NAME := Morsebot

# Default target
all: build

# Build target
build:
	@echo "Building the Go application..."
	@mkdir -p $(BIN_DIR)    # Ensure bin directory exists
	go build -o $(BIN_DIR)/$(BINARY_NAME)

# Clean target
clean:
	@echo "Cleaning up the binary..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)

# Run target (optional)
run: build
	@echo "Running the application..."
	$(BIN_DIR)/$(BINARY_NAME)

# Phony targets
.PHONY: all build clean run
