targets="windows/amd64,windows/386,darwin/arm64,darwin/amd64,darwin/386,linux/amd64,linux/arm64,linux/386"
pkg="github.com/elwinar/rambler"
version=$(shell git describe --tags)
ldflags="-X main.VERSION=${version}"

default: build
all: build test

.PHONY: build
build: ## Build the binary for the local architecture
	go build --ldflags=${ldflags}

.PHONY: help
help: ## Get help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: release
release: ## Build the release files
	mkdir -p release/github.com/elwinar
	xgo --dest release --targets=$(targets) --ldflags=$(ldflags) .
	docker buildx build --output type=tar --build-arg VERSION=${version} --platform linux/amd64 --file build.Dockerfile . | tar -xO rambler >release/github.com/elwinar/rambler-alpine-amd64
	docker buildx build --output type=tar --build-arg VERSION=${version} --platform linux/arm64 --file build.Dockerfile . | tar -xO rambler >release/github.com/elwinar/rambler-alpine-arm64
	chmod +x release/github.com/elwinar/rambler-alpine-amd64 release/github.com/elwinar/rambler-alpine-arm64

.PHONY: test
test: ## Test the project
	go test ./...
