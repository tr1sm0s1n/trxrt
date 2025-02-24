run:
	@go run .

build:
	@go build -o ./bin/wallet-x

fmt:
	@gofmt -s -w ./..
