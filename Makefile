GO        ?= go
BINDIR    := $(CURDIR)/bin
LDFLAGS   := -w -s
TESTS     := ./...
TESTFLAGS :=

TARGETS   ?= darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS = find * -type d -exec

SHELL = /bin/bash

BUILD_PATH = github.com/belitre/kctx-manager/cmd/kctx-manager
NAME = kctx-manager

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)

HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)
GOLANGCI_LINT_VERSION := v1.23.6

GOLANGCI_VERSION_CHECK := $(shell golangci-lint --version | grep -oh $(GOLANGCI_LINT_VERSION);)

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif

BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X $(BUILD_PATH)/version.Version=${BINARY_VERSION}
endif

LDFLAGS += -X $(BUILD_PATH)/version.GitCommit=${GIT_COMMIT}

.PHONY: info
info:
	@echo "Git status:        $(GIT_DIRTY)"
	@echo "Version:           ${VERSION}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"

.PHONY: build
build: bootstrap info tidy fmt 
	GOBIN=$(BINDIR) $(GO) install -ldflags '$(LDFLAGS)' $(BUILD_PATH)

# usage: make clean build-cross dist VERSION=v0.2-alpha
.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross: bootstrap info tidy fmt
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
	npm ci
	npx semantic-release 

.PHONY: semantic-release-dry-run
semantic-release-dry-run:
	npm install
	npx semantic-release -d

.PHONY: lint
lint: bootstrap-lint 
	@echo "lint target..."
	@golangci-lint run --enable-all --disable lll,nakedret,funlen,gochecknoglobals ./...

.PHONY: bootstrap-lint
bootstrap-lint:
	@echo "bootstrap lint..."
ifndef HAS_GOLANGCI_LINT
	@echo "golangci-lint $(GOLANGCI_LINT_VERSION) not found..."
	@GOPROXY=direct GOSUMDB=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
else
	@echo "golangci-lint found, checking version..."
ifeq ($(GOLANGCI_VERSION_CHECK), )
	@echo "found different version, installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@GOPROXY=direct GOSUMDB=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
else
	@echo "golangci-lint version $(GOLANGCI_VERSION_CHECK) found!"
endif
endif

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

.PHONY: fmt
fmt:
	@echo "fmt target..."
	@gofmt -l -w -s $(SRC)

.PHONY: install-npm-check-updates
install-npm-check-updates:
	npm install npm-check-updates

.PHONY: update-dependencies
update-dependencies: install-npm-check-updates
	ncu -u
	npm install