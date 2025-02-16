PROJECT_NAME := "timble"
REST_DIR := cmd/rest
MAIN_REST := "$(CURDIR)/$(REST_DIR)"
PKG := "$(PROJECT_NAME)"
OUTPUT_DIR := "deploy/_output"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /mocks/)
SDIRS = $(shell ls scripts)
export VAR_SCRIPTS  ?= $(SDIRS:/=)

.PHONY: prepare fetch vet unit-test race-test msan-test coverage coverhtml vendor compile

prepare: fetch vendor

fetch:
	go mod tidy
	go mod download
	go mod verify

vet: ## Run analyze
	go vet ${PKG_LIST}

unit-test: ## Run unit tests
	go test -short ${PKG_LIST}

race-test: ## Run data race detector
	go test -race -short ${PKG_LIST}

vendor: ## Run download dependencies
	go mod vendor

coverage: ## Generate global code coverage report
	./tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./tools/coverage.sh html;

compile: ## Run build go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -p 1 -o $(OUTPUT_DIR)/rest/timble $(MAIN_REST)
	$(foreach script, $(VAR_SCRIPTS), \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(OUTPUT_DIR)/scripts/$(script) scripts/$(script)/main.go;)

compile_osx: ## Run build go for mac
	go build -ldflags=-s -o $(OUTPUT_DIR)/rest/timble $(REST_DIR)/main.go
	$(foreach script, $(VAR_SCRIPTS), \
		go build -ldflags=-s -o $(OUTPUT_DIR)/scripts/$(script) scripts/$(script)/main.go;)

pretty:
	gofmt -w .
	goimports -w .

run-rest:
	go run cmd/rest/main.go

