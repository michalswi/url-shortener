GOLANG_VERSION := 1.13.4
ALPINE_VERSION := 3.10

GIT_REPO := github.com/michalswi/url-shortener
DOCKER_REPO := michalsw
APPNAME := url-shortener

VERSION ?= $(shell git describe --always)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
LAST_COMMIT_USER ?= $(shell git log -1 --format='%cn <%ce>')
LAST_COMMIT_HASH ?= $(shell git log -1 --format=%H)
LAST_COMMIT_TIME ?= $(shell git log -1 --format=%cd --date=format:'%Y-%m-%d %H:%M:%S')

SERVICE_ADDR := 8080
PPROF_ADDR := 5050
STORE_ADDR := 6379
DNS_NAME := localhost

.DEFAULT_GOAL := all
.PHONY: all test go-run go-build docker-build docker-run docker-stop docker-push release dockertest

all: test go-build

test:
	go test -v ./...

go-run:
	SERVICE_ADDR=$(SERVICE_ADDR) \
	PPROF_ADDR=$(PPROF_ADDR) \
	STORE_ADDR=$(STORE_ADDR) \
	DNS_NAME=$(DNS_NAME) \
	go run .	

go-build:
	CGO_ENABLED=0 \
	go build \
	-v \
	-ldflags "-s -w -X '$(GIT_REPO)/version.AppVersion=$(VERSION)' \
	-X '$(GIT_REPO)/version.BuildTime=$(BUILD_TIME)' \
	-X '$(GIT_REPO)/version.LastCommitUser=$(LAST_COMMIT_USER)' \
	-X '$(GIT_REPO)/version.LastCommitHash=$(LAST_COMMIT_HASH)' \
	-X '$(GIT_REPO)/version.LastCommitTime=$(LAST_COMMIT_TIME)'" \
	-o $(APPNAME)-$(VERSION) .

docker-build:
	docker build \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg APPNAME="$(APPNAME)" \
	--build-arg VERSION="$(VERSION)" \
	--build-arg BUILD_TIME="$(BUILD_TIME)" \
	--build-arg LAST_COMMIT_USER="$(LAST_COMMIT_USER)" \
	--build-arg LAST_COMMIT_HASH="$(LAST_COMMIT_HASH)" \
	--build-arg LAST_COMMIT_TIME="$(LAST_COMMIT_TIME)" \
	--label="build.version=$(VERSION)" \
	--tag="$(DOCKER_REPO)/$(APPNAME):latest" \
	--tag="$(DOCKER_REPO)/$(APPNAME):$(VERSION)" \
	.

docker-run:
	docker run -d --rm \
	--name $(APPNAME) \
	-p $(SERVICE_ADDR):$(SERVICE_ADDR) \
	-p $(PPROF_ADDR):$(PPROF_ADDR) \
	$(DOCKER_REPO)/$(APPNAME):latest

docker-stop:
	docker stop $(APPNAME)	

docker-push:
	docker push $(DOCKER_REPO)/$(NAME):latest
	docker push $(DOCKER_REPO)/$(NAME):$(VERSION)

dockertest: docker-build docker-run

release: docker-build docker-push