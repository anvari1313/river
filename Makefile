BASH_PATH:=$(shell which bash)
SHELL=$(BASH_PATH)
ROOT := $(shell realpath $(dir $(lastword $(MAKEFILE_LIST))))
APP := river
APP_IMPORT_PATH := github.com/anvari1313/river
BUILD_PATH ?= ".build/app"
BUILD_DATE ?= $(shell date +'%Y-%m-%dT%H:%M:%S%z')
GIT_HEAD_REF := $(shell cat .git/HEAD | cut -d' ' -f2)
GIT_BRANCH ?= $(shell echo $(GIT_HEAD_REF) | cut -d'/' -f3)
GIT_SHA ?= $(shell cat .git/$(GIT_HEAD_REF) | head -c 8)
GIT_TAG ?= $(shell git describe --exact-match $GIT_SHA 2&>/dev/null || echo "v0.0.0")
LDFLAGS := "-w -s \
	-X $(APP_IMPORT_PATH)/cmd.BuildDate=$(BUILD_DATE)\
	-X $(APP_IMPORT_PATH)/cmd.GitCommit=$(GIT_SHA) \
	-X $(APP_IMPORT_PATH)/cmd.GitRef=$(GIT_BRANCH) \
	-X $(APP_IMPORT_PATH)/cmd.GitTag=$(GIT_TAG)"

all: format lint-ci build-static

############################################################
# Build & Run
############################################################
build:
	go build -v -race .

build-static:
	CGO_ENABLED=0 go build -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

build-static-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -v -o $(APP) -installsuffix cgo -ldflags $(LDFLAGS) .

docker:
	docker build \
	  --build-arg=BUILD_PATH=$(BUILD_PATH) \
	  --build-arg=BUILD_DATE=$(BUILD_DATE) \
	  --build-arg=GIT_BRANCH=$(GIT_BRANCH) \
	  --build-arg=GIT_SHA=$(GIT_SHA) \
	  --build-arg=GIT_TAG=$(GIT_TAG) \
	  -t $(APP):$(GIT_BRANCH) .

install:
	cp $(APP) $(GOPATH)/bin

run:
	go run -race .

############################################################
# Test & Coverage
############################################################
test: check-gotestsum
	gotestsum --junitfile-testcase-classname short --junitfile .report.xml -- -gcflags 'all=-N -l' -mod vendor ./...

coverage:
	gotestsum -- -gcflags 'all=-N -l' -mod vendor -v -coverprofile=.testCoverage.txt ./...
	GOFLAGS=-mod=vendor go tool cover -func=.testCoverage.txt

coverage-report: coverage
	GOFLAGS=-mod=vendor go tool cover -html=.testCoverage.txt -o testCoverageReport.html
	gocover-cobertura < .testCoverage.txt > .cobertura.xml

############################################################
# Format & Lint
############################################################
check-goimport:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-goimport
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-golint:
	which golint || (go get -u golang.org/x/lint/golint)

lint: check-golint
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R golint -set_exit_status R

check-golangci-lint:
	which golangci-lint || (go get -u github.com/golangci/golangci-lint/cmd/golangci-lint)

lint-ci: check-golangci-lint
	golangci-lint run -c .golangci.yml ./...

.PHONY: build install run ci-test
