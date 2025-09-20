.PHONY: all build test clean run migrate

# Variables
BINARY_NAME=opus
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_CLEAN=$(GO_CMD) clean
GO_MOD=$(GO_CMD) mod
GO_RUN=$(GO_CMD) run

# Targets
all: build

build:
	$(GO_BUILD) -o bin/$(BINARY_NAME) ./cmd/server

build-cli:
	$(GO_BUILD) -o bin/$(BINARY_NAME)-cli ./cmd/cli

test:
	$(GO_TEST) -v ./modules/...

test-cover:
	$(GO_TEST) -coverprofile=coverage.out ./...
	$(GO_CMD) tool cover -html=coverage.out

clean:
	$(GO_CLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_NAME)-cli

run:
	$(GO_RUN) cmd/server/main.go

tidy:
	$(GO_MOD) tidy

lint:
	~/go/bin/golangci-lint run

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Integration tests
test-integration:
	$(GO_TEST) -tags=integration -v ./modules/...
	