.PHONY: xctrainer

test:
	go test ./...

check:
	go mod tidy
	go fmt ./...

xctrainer: check test
	go build -o build/xctrainer  ./cmd/xctrainer/main.go
