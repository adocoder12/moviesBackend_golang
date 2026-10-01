APP_NAME := api
BIN_DIR  := bin


#--- commands ---
.PHONY: run dev build test vet fmt lint tidy check clean

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build         Build the binary"
	@echo "  dev           Run the app with live reload"
	@echo "  run           Run the app directly"
	@echo "  lint          Run golangci-lint"
	@echo "  migrate-up    Run all up migrations"
	@echo "  migrate-down  Rollback the last migration"
	@echo "  docker-up     Run up docker             "
	@echo "  docker-down     Run down docker         "

## run: start the API
run:
	go run ./cmd/api/main.go

## dev: start with live reload (uses .air.toml)
dev:
	air
## build: compile the binary to bin/api
build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

## test: run all tests with the race detector
test:
	go test -race ./...
## fmt: format all code
fmt:
	gofmt -w .
## vet: catch common mistakes
vet:
	go vet ./...
## lint: run golangci-lint (must be installed)
lint:
	golangci-lint run

## tidy: clean up go.mod and go.sum
tidy:
	go mod tidy

## check: everything CI will run
check: vet test build

## clean: remove build output
clean:
	rm -rf $(BIN_DIR)
