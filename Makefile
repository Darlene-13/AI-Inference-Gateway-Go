# Short cut for common development commands.
.PHONY: help run test test-race cover lint fmt build tidy

help: ## list the commands
	@grep -E '^[a-z-]+:.*## ' Makefile | awk -F ':.*## ' '{printf "  %-10s %s\n", $$1, $$2}'

run: ## start the gateway with your local config
	go run ./cmd/gateway -config internal/config/local.yaml

test: ## run all unit tests
	go test ./...

test-race: ## run all tests with the race detector, like CI
	go test -race -count=1 ./...

cover: ## test coverage per function
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

lint: ## run the linters configured in .golangci.yml
	golangci-lint run

fmt: ## format the code
	golangci-lint fmt

build: ## build the binary into bin/
	go build -o bin/gateway ./cmd/gateway

tidy: ## add missing and remove unused dependencies
	go mod tidy
