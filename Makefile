check: lint test

test:
	@go test ./...

lint:
	@golangci-lint run

run: mod
	@go run main.go

build: mod
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o ./bin/wordle
	@chmod +x ./bin/wordle

mod:
	@go mod download
