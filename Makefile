.PHONY: build

BUILD_DIR = ./build
BINARY_NAME = server
MAIN_PACKAGE = ./cmd/http
FULL_BINARY_NAME = $(BUILD_DIR)/$(BINARY_NAME)

all: test build

build:
	go build -o $(FULL_BINARY_NAME) $(MAIN_PACKAGE)

test:
	go test -v ./...

run: build
	$(FULL_BINARY_NAME)

clean:
	go clean
	rm -f $(FULL_BINARY_NAME)

tidy:
	go mod tidy
