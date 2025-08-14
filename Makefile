.DEFAULT_GOAL := help
GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
CURRENT_DIR = $(shell pwd)
BIN_DIR = $(CURRENT_DIR)/bin
RPC_URL = "http://127.0.0.1:8545"

.PHONY: install run build fyne-cmd lint tidy fmt help

#? install: Install prerequisites.
install:
	@sudo apt-get install gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev

#? run: Run TrXrT.
run:
	@env "RPC_URL=$(RPC_URL)" $(GO) run .

#? build: Build TrXrT.
build:
	@env "RPC_URL=$(RPC_URL)" $(GO) build -o ./bin/trxrt

#? fyne-cmd: Install fyne cmd.
fyne-cmd:
	@GOBIN=$(BIN_DIR) $(GO) install fyne.io/fyne/v2/cmd/fyne@latest

#? lint: Lint with golangci-lint.
lint:
	@if [ ! -f "$(BIN_DIR)/golangci-lint" ]; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s latest; \
	fi
	@$(BIN_DIR)/golangci-lint run

#? tidy: Tidy module.
tidy:
	@$(GO) mod tidy

#? fmt: Format files.
fmt:
	@$(GOFMT) -w $(GOFILES)

help: Makefile
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sort | sed -e 's/^/ /'
