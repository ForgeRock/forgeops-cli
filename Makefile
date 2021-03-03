.PHONY: build clean test pr-build snapshot-build

BIN_NAME=forgeops

VERSION := $(shell grep "var Version " pkg/version/version.go | awk -F " = " -F '"' '{ print $$2 }')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "forgerock/forgeops-cli"

PR_VERSION_NAME = ${VERSION}-pr.${PR_NUMBER}

all: test build docs

pr-tag:
	git tag "$(PR_VERSION_NAME)"

snapshot:
	TAG_NAME=${PR_VERSION_NAME} BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_DIRTY=${GIT_DIRTY} goreleaser --snapshot --rm-dist

pr-build: pr-tag snapshot
	@echo "PR Build Completed: gs://engineering-devops_cloudbuild/forgeops-cli-artifacts/$(PR_NUMBER)"

release:
	TAG_NAME=${TAG_NAME} BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_DIRTY=${GIT_DIRTY} goreleaser
	@echo "Released: VERSION=${VERSION} BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_DIRTY=${GIT_DIRTY}"

install-tools:
	./hack/install-goreleaser.sh
	./hack/install-linter.sh

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	@test ! -e dist || rm -r dist


deps:
	go get ./...

gen-mocks: deps
	mockery --dir internal/k8s --output internal/mock --all

test tests:
	@go test ./...
	@echo "tests completed"

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build forgeops binary
build: fmt vet test clean
	go build -o bin/forgeops main.go

# generates forgeops-cli docs
docs doc: vet fmt
	rm -rf docs/
	go run main.go docs
