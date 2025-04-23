.DEFAULT_GOAL := help
GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
CURRENT_DIR = $(shell pwd)
BIN_DIR = $(CURRENT_DIR)/bin

.PHONY: run build fyne-cmd lint tidy fmt

run:
	@$(GO) run .

build:
	@$(GO) build -o ./bin/trxrt

fyne-cmd:
	@GOBIN=$(BIN_DIR) $(GO) install fyne.io/fyne/v2/cmd/fyne@latest

lint:
	@if [ ! -f "$(BIN_DIR)/golangci-lint" ]; then \
		echo "Installing golangci-lint..."; \
		GOBIN=$(BIN_DIR) $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
	fi
	@$(BIN_DIR)/golangci-lint run

tidy:
	@$(GO) mod tidy

fmt:
	@$(GOFMT) -w $(GOFILES)
