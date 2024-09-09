check: lint test

test:
	@go test ./...

lint:
	@golangci-lint run

build:
	@GOOS=darwin GOARCH=arm64 go build -o wordle
