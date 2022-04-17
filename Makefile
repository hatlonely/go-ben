NAME ?= ben
REGISTRY_ENDPOINT ?= docker.io
REGISTRY_NAMESPACE ?= hatlonely
IMAGE_TAG ?= $(shell git describe --tags | awk '{print(substr($$0,2,length($$0)))}')

export GOPROXY=https://goproxy.cn

define BUILD_VERSION
  version: $(shell git describe --tags)
gitremote: $(shell git remote -v | grep fetch | awk '{print $$2}')
   commit: $(shell git rev-parse HEAD)
 datetime: $(shell date '+%Y-%m-%d %H:%M:%S')
 hostname: $(shell hostname):$(shell pwd)
goversion: $(shell go version)
endef
export BUILD_VERSION

.PHONY: build
build: cmd/main.go $(wildcard internal/*/*.go) Makefile vendor
	mkdir -p build/bin
	go build -ldflags "-X 'main.Version=$$BUILD_VERSION'" -o build/bin/${NAME} cmd/main.go

.PHONY: clean
clean:
	rm -rf build

vendor: go.mod go.sum
	go mod tidy
	go mod vendor
