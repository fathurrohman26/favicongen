# favicongen

.PHONY: all build test clean

all: build

build:
	@go run ./tools/build/main.go -o ./dist

build-all:
	@go run ./tools/build/main.go -o ./dist --all

test:
	@go test -v -race ./...

clean:
	@rm -rf ./dist