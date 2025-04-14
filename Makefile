build:
	go build -o bin/ ./cmd/...

fmt:
	gofmt -w cmd internal
