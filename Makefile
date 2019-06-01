GO        ?= go
BINDIR    := $(CURDIR)/bin
LDFLAGS   := -w -s
TESTS     := ./...
TESTFLAGS :=

TARGETS   ?= darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS = find * -type d -exec

SHELL = /bin/bash

BUILD_PATH = github.com/belitre/kctx-manager/cmd/kctx-manager
GITHUB_REPO = https://github.com/belitre/kctx-manager.git
NAME = kctx-manager

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif

BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X $(BUILD_PATH)/version.Version=${BINARY_VERSION}
endif

LDFLAGS += -X $(BUILD_PATH)/version.GitCommit=${GIT_COMMIT}

info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"

.PHONY: build
build: tidy
	GOBIN=$(BINDIR) $(GO) install -ldflags '$(LDFLAGS)' $(BUILD_PATH)

# usage: make clean build-cross dist VERSION=v0.2-alpha
.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/{{.Dir}}" -osarch='$(TARGETS)' -ldflags '$(LDFLAGS)' $(BUILD_PATH)

.PHONY: dist
dist:
	( \
		cd _dist && \
		$(DIST_DIRS) tar -zcf $(NAME)-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r $(NAME)-${VERSION}-{}.zip {} \; \
	)

.PHONY: test
test: build
test: TESTFLAGS += -v
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	$(GO) test $(TESTS) $(GOFLAGS) $(TESTFLAGS)

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist

.PHONY: semantic-release
semantic-release:
	npm install
	npm audit fix
	npx semantic-release -r $(GITHUB_REPO)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif
ifndef HAS_GIT
	$(error You must install Git)
endif

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor: tidy
	go mod vendor