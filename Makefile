default: build

build: test
	mkdir -p bin
	go build -o bin/gifwrap ./cmd/gifwrap/

demo: build
	mkdir -p bin
	./bin/gifwrap -f https://media.giphy.com/media/QMHoU66sBXqqLqYvGO/giphy.gif

build-travis: test
	mkdir -p bin
	GOOS=linux  GOARCH=amd64 go build -o bin/gifwrap-linux-amd64  ./cmd/gifwrap
	GOOS=darwin GOARCH=amd64 go build -o bin/gifwrap-darwin-amd64 ./cmd/gifwrap
	GOOS=darwin GOARCH=arm64 go build -o bin/gifwrap-darwin-arm64 ./cmd/gifwrap

test:
	go vet ./...
	go test -v ./...

.PHONY: build test
