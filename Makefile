# Go tools.
GO ?= go
GO_MD2MAN ?= go-md2man

# Paths.
PROJECT := github.com/openSUSE/helm-mirror
CMD := .
HELM_HOME_MIRROR := $(HELM_HOME)/plugins/helm-mirror
HELM_HOME_MIRROR_BIN := $(HELM_HOME_MIRROR)/bin

# We use Docker because Go is just horrific to deal with.
MIRROR_IMAGE := mirror_dev
DOCKER_RUN := docker run --rm -it --security-opt apparmor:unconfined --security-opt label:disable -v ${PWD}:/go/src/${PROJECT}

# Output directory.
BUILD_DIR ?= ./bin

# Release information.
GPG_KEYID ?=

# Version information.
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
COMMIT_NO := $(shell git rev-parse HEAD 2> /dev/null || true)
COMMIT := $(if $(shell git status --porcelain --untracked-files=no),"${COMMIT_NO}-dirty","${COMMIT_NO}")

# Get current Version changelog
CHANGE := $(shell sed -e '1,/v$(VERSION)/d;/v.*/Q' ./CHANGELOG.md)

BUILD_FLAGS ?=

BASE_FLAGS := ${BUILD_FLAGS} -tags "${BUILDTAGS}"

BASE_LDFLAGS := -X $(PROJECT)/cmd.version=$(VERSION)
BASE_LDFLAGS += -X $(PROJECT)/cmd.gitCommit=$(COMMIT)

DYN_BUILD_FLAGS := ${BASE_FLAGS} -buildmode=pie -ldflags "${BASE_LDFLAGS}"
TEST_BUILD_FLAGS := ${BASE_FLAGS} -buildmode=pie -ldflags "${BASE_LDFLAGS} -X ${PROJECT}/pkg/testutils.binaryType=test"
STATIC_BUILD_FLAGS := ${BASE_FLAGS} -ldflags "${BASE_LDFLAGS} -extldflags '-static'"

.DEFAULT: mirror

GO_SRC = $(shell find . -name \*.go)

# NOTE: If you change these make sure you also update local-validate-build.

mirror: $(GO_SRC)
	$(GO) build ${DYN_BUILD_FLAGS} -o $(BUILD_DIR)/helm-mirror ${CMD}

mirror.static: $(GO_SRC)
	env CGO_ENABLED=0 $(GO) build ${STATIC_BUILD_FLAGS} -o $(BUILD_DIR)/helm-mirror ${CMD}

install: $(GO_SRC)
	$(GO) install -v ${DYN_BUILD_FLAGS} ${CMD}

install.static: $(GO_SRC)
	$(GO) install -v ${STATIC_BUILD_FLAGS} ${CMD}

install.plugin: mirror
	mkdir -p $(HELM_HOME_MIRROR_BIN)
	cp bin/mirror $(HELM_HOME_MIRROR_BIN)
	cp plugin.yaml $(HELM_HOME_MIRROR)/

clean:
	rm -rf ./bin
	rm -f $(MANPAGES)

local-validate: local-validate-git local-validate-go local-validate-reproducible

EPOCH_COMMIT ?= 9ef2a655b2b071a3319892f9b249e2e8160eca10
local-validate-git:
	@type git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found." && false)
ifdef TRAVIS_COMMIT_RANGE
	git-validation -q -run DCO,short-subject
else
	git-validation -q -run DCO,short-subject -range $(EPOCH_COMMIT)..HEAD
endif

local-validate-go:
	@type gofmt    >/dev/null 2>/dev/null || (echo "ERROR: gofmt not found." && false)
	test -z "$$(gofmt -s -l . | grep -vE '^vendor/|^third_party/' | tee /dev/stderr)"
	@type golint   >/dev/null 2>/dev/null || (echo "ERROR: golint not found." && false)
	test -z "$$(golint $(PROJECT)/... | grep -vE '/vendor/|/third_party/' | tee /dev/stderr)"
	@go doc cmd/vet >/dev/null 2>/dev/null || (echo "ERROR: go vet not found." && false)
	test -z "$$($(GO) vet $$($(GO) list $(PROJECT)/... | grep -vE '/vendor/|/third_party/') 2>&1 | tee /dev/stderr)"

# Make sure that our builds are reproducible even if you wait between them and
# the modified time of the files is different.
local-validate-reproducible:
	mkdir -p .tmp-validate
	make -B mirror && cp $(BUILD_DIR)/mirror .tmp-validate/mirror.a
	@echo sleep 10s
	@sleep 10s && touch $(GO_SRC)
	make -B mirror && cp $(BUILD_DIR)/mirror .tmp-validate/mirror.b
	diff -s .tmp-validate/mirror.{a,b}
	sha256sum .tmp-validate/mirror.{a,b}
	rm -r .tmp-validate/mirror.{a,b}

local-validate-build:
	$(GO) build ${DYN_BUILD_FLAGS} -o /dev/null ${CMD}
	env CGO_ENABLED=0 $(GO) build ${STATIC_BUILD_FLAGS} -o /dev/null ${CMD}
	$(GO) test -run nothing ${DYN_BUILD_FLAGS} $(PROJECT)/...

# Used for tests.
DOCKER_IMAGE :=opensuse/amd64:tumbleweed

mirrorimage:
	docker build -t $(MIRROR_IMAGE) .


test.unit: mirrorimage
	$(DOCKER_RUN) $(MIRROR_IMAGE) make test

test: local-validate-go
	rm -rf /tmp/mirror
	$(GO) test -v ./...

cover:
	bash <scripts/cover.sh

bootstrap:
	dep ensure

dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist:
	rm -rf build/mirror/* release/*
	mkdir -p build/mirror/bin release/
	cp README.md LICENSE plugin.yaml build/mirror
	GOOS=linux GOARCH=amd64 go build -o build/mirror/bin/helm-mirror -ldflags="$(BASE_LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-mirror-linux.tgz mirror/
	GOOS=darwin GOARCH=amd64 go build -o build/mirror/bin/helm-mirror -ldflags="$(BASE_LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-mirror-macos.tgz mirror/
	rm build/mirror/bin/helm-mirror
	GOOS=windows GOARCH=amd64 go build -o build/mirror/bin/helm-mirror.exe -ldflags="$(BASE_LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-mirror-windows.tgz mirror/

release: dist
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif
	github-release release -u openSUSE -r helm-mirror --tag v$(VERSION)  --name v$(VERSION) -s $(GITHUB_TOKEN) -d "$(CHANGE)"
	github-release upload -u openSUSE -r helm-mirror --tag v$(VERSION)  --name helm-mirror-linux.tgz --file release/helm-mirror-linux.tgz -s $(GITHUB_TOKEN)
	github-release upload -u openSUSE -r helm-mirror --tag v$(VERSION)  --name helm-mirror-macos.tgz --file release/helm-mirror-macos.tgz -s $(GITHUB_TOKEN)
	github-release upload -u openSUSE -r helm-mirror --tag v$(VERSION)  --name helm-mirror-windows.tgz --file release/helm-mirror-windows.tgz -s $(GITHUB_TOKEN)

MANPAGES_MD := $(wildcard doc/man/*.md)
MANPAGES    := $(MANPAGES_MD:%.md=%)

doc/man/%.1: doc/man/%.1.md
	$(GO_MD2MAN) -in $< -out $@

doc: $(MANPAGES)

.PHONY: mirror \
	mirror.static \
	install \
	install.static \
	install.plugin \
	clean \
	local-validate \
	local-validate-git \
	local-validate-go \
	local-validate-reproducible \
	local-validate-build \
	mirrorimage \
	test.unit
	test \
	cover \
	bootstrap \
	dist \
	release \
	doc