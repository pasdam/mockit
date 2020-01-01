BUILD_DIR ?= .build
PROJECT_NAME ?= "app"

include scripts/makefiles/third_party/pasdam/makefiles/docker.mk
include scripts/makefiles/third_party/pasdam/makefiles/go.mk
include scripts/makefiles/third_party/pasdam/makefiles/go.mod.mk
include scripts/makefiles/third_party/pasdam/makefiles/help.mk

.DEFAULT_GOAL := help

## clean: Remove all artifacts
.PHONY: clean
clean: go-clean docker-clean
	@rm -rf .build generated

## build: Build the artifact
.PHONY: build
build: | proto-build go-build

## gitlab-ci-test: Run the stages locally to verify that they execute correctly
.PHONY: gitlab-ci-test
gitlab-ci-test:
	@gitlab-runner exec docker inspect
	@gitlab-runner exec docker build

## install-dep: Install the required tools (proto-gen-go)
.PHONY: install-dep
install-dep:
	@go get -u github.com/golang/protobuf/protoc-gen-go

## proto-build: Generate code from protobuf definitions
.PHONY: proto-build
proto-build: generated/pkg
	@protoc -I=proto --go_out=generated/pkg proto/forexpb/*.proto

generated/pkg:
	@mkdir -p generated/pkg
