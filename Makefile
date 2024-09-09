check: lint test

test:
	@go test ./...

lint:
	@golangci-lint run

build:
	@go build -o wordle
