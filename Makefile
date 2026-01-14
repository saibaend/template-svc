DEFAULT_GOAL=run-local

swag:
	swag init -g ./cmd/main.go --parseDependency

run-local:
	go run ./cmd/main.go

run-tests:
	gotestsum --format testname

lint:
	golangci-lint run ./...

gen:
	go generate ./...

test:
	go test -race -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out