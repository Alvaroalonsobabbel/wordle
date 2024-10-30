check: lint test

test:
	@go test ./...

lint:
	@golangci-lint run

build: mod
	@go build -o wordle
	chmod +x wordle

mod:
	@go mod download
