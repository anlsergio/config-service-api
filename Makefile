.PHONY: build

# Where the compiled binary will be created at.
BUILD_DIR = ./build
# The name of the compiled binary.
BINARY_NAME = server
# The package where the binary will bootstrap the application.
MAIN_PACKAGE = ./cmd/http
# The relative path to the compiled binary
FULL_BINARY_NAME = $(BUILD_DIR)/$(BINARY_NAME)

all: test build

##############
# Run/Compile
##############
build:
	go build -o $(FULL_BINARY_NAME) $(MAIN_PACKAGE)

run: build
	$(FULL_BINARY_NAME)

##############
# Housekeeping
##############
clean:
	go clean
	rm -f $(FULL_BINARY_NAME)

tidy:
	go mod tidy

##############
# Test
##############
test:
	go test -v ./...

test-cov:
	./scripts/generate-test-cov.sh

show-cov:
	go tool cover -html=/tmp/test.cov
