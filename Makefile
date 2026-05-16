APP_NAME = plane
GO = go
LDFLAGS = -ldflags="-s -w"

.PHONY: build build-all clean test

build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/plane

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 ./cmd/plane

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-arm64 ./cmd/plane

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 ./cmd/plane

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 ./cmd/plane

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

clean:
	rm -rf bin/

test:
	$(GO) test ./...
