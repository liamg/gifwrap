default: build

build: test
	mkdir -p bin
	go build ./cmd/gifwrap/ -o bin/gifwrap

build-travis: test
	mkdir -p bin
	GOOS=linux  GOARCH=amd64 go build -o bin/gifwrap-linux-amd64  ./cmd/gifwrap
	GOOS=darwin GOARCH=amd64 go build -o bin/gifwrap-darwin-amd64 ./cmd/gifwrap
	GOOS=darwin GOARCH=arm64 go build -o bin/gifwrap-darwin-arm64 ./cmd/gifwrap

test:
	go vet ./...
	go test -v ./...

.PHONY: build test
