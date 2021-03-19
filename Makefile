default: build

build: test
	mkdir -p bin
	go build ./cmd/gifwrap/ -o bin/gifwrap

build-travis: test
	mkdir -p bin/linux-amd64/gifwrap
	mkdir -p bin/darwin-amd64/gifwrap
	mkdir -p bin/darwin-arm64/gifwrap
	GOOS=linux  GOARCH=amd64 go build -o bin/linux-amd64/gifwrap  ./cmd/gifwrap
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/gifwrap ./cmd/gifwrap
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/gifwrap ./cmd/gifwrap

test:
	go vet ./...
	go test -v ./...

.PHONY: build test
