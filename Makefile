.PHONY: build clean test pr-build snapshot-build

BIN_NAME=forgeops

VERSION := $(shell grep "const Version " pkg/version/version.go | awk -F " = " '{ print $$2 }')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "forgerock/forgeops-cli"

PR_VERSION_NAME = "$(VERSION)-pr.$(PR_NUMBER)"

pr-tag:
	git tag "$(PR_VERSION_NAME)"

snapshot:
	BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_DIRTY=${GIT_DIRTY} goreleaser --snapshot

pr-build: pr-tag snapshot
	@echo "PR Build Completed: gs://engineering-devops_cloudbuild/forgeops-cli-artifacts/$(PR_NUMBER)"

release:
	BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_DIRTY=${GIT_DIRTY} goreleaser
	@echo "Released: VERSION=${VERSION} BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_DIRTY=${GIT_DIRTY}"

install-tools:
	./hack/install-goreleaser.sh

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	@test ! -e dist || rm -r dist

test:
	@go test ./...
	@echo "tests completed"
