.PHONY: build clean test help default

BIN_NAME=berekbek

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

default: build

help:
	@echo 'Management commands for berekbek:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs dep ensure, mostly used for ci.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"

	go build -ldflags "-X github.com/gobliggg/berekbek/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/gobliggg/berekbek/version.BuildDate=${BUILD_DATE}" -o bin/${BIN_NAME}

get-deps:
	go get
clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...

