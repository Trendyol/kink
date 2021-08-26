# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Set version variables for LDFLAGS
GIT_TAG ?= dirty-tag
GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GIT_HASH ?= $(shell git rev-parse HEAD)
DATE_FMT = +'%Y-%m-%dT%H:%M:%SZ'
SOURCE_DATE_EPOCH ?= $(shell git log -1 --pretty=%ct)
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
GIT_TREESTATE = "clean"
DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(DIFF), 1)
    GIT_TREESTATE = "dirty"
endif

PKG=gitlab.trendyol.com/platform/base/poc/kink/cmd

LDFLAGS="-X $(PKG).GitVersion=$(GIT_VERSION) -X $(PKG).gitCommit=$(GIT_HASH) -X $(PKG).gitTreeState=$(GIT_TREESTATE) -X $(PKG).buildDate=$(BUILD_DATE)"

.PHONY: all kink release

all: kink

SRCS = $(shell find cmd -iname "*.go") $(shell find pkg -iname "*.go")

kink: $(SRCS)
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o kink

release:
	export GITLAB_TOKEN=$(GITLAB_TOKEN) && \
	export DOCKER_REGISTRY=$(DOCKER_REGISTRY) && \
	export DOCKER_USERNAME=$(DOCKER_USERNAME) && \
	export DOCKER_PASSWORD=$(DOCKER_PASSWORD) && \
	LDFLAGS=$(LDFLAGS) goreleaser release --rm-dist

release-local:
	docker container run --rm --privileged \
          -v $(shell pwd):/kink \
          -v /var/run/docker.sock:/var/run/docker.sock \
          -w /kink \
          --entrypoint='/bin/sh' \
          -e GITLAB_TOKEN=$(GITLAB_TOKEN) \
          -e DOCKER_REGISTRY=$(DOCKER_REGISTRY) \
          -e DOCKER_USERNAME=$(DOCKER_USERNAME) \
          -e DOCKER_PASSWORD=$(DOCKER_PASSWORD) \
          registry.trendyol.com/platform/base/image/gythialy/golang-cross:v1.17 -c "make release"
