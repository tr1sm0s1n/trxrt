.DEFAULT_GOAL := help
GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
CURRENT_DIR = $(shell pwd)

.PHONY: run build fyne-cmd tidy fmt

run:
	@$(GO) run .

build:
	@$(GO) build -o ./bin/wallet-x

fyne-cmd:
	@GOBIN=$(CURRENT_DIR)/bin $(GO) install fyne.io/fyne/v2/cmd/fyne@latest

tidy:
	@$(GO) mod tidy

fmt:
	@$(GOFMT) -w $(GOFILES)
