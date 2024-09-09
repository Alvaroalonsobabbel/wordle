GOOS=darwin
GOARCH=arm64

check: lint test

test:
	@go test ./...

lint:
	@golangci-lint run

build:
	@go build -o wordle
