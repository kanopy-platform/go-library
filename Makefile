GO_MODULE := $(shell git config --get remote.origin.url | grep -o 'github\.com[:/][^.]*' | tr ':' '/')
REPO_SLUG := $(shell echo ${GO_MODULE} | cut -d/ -f2-)
CMD_NAME := $(shell basename ${GO_MODULE})
DEFAULT_APP_PORT ?= 8080
GIT_COMMIT := $(shell git rev-parse HEAD)

RUN ?= .*
PKG ?= ./...
.PHONY: test
test: ## Run tests in local environment
	golangci-lint run --timeout=5m $(PKG)
	go test -cover -run=$(RUN) $(PKG)

.PHONY: license-check
license-check:
	licensed cache
	licensed status


LDFLAGS = "-X 'github.com/${REPO_SLUG}/internal/version.version=${VERSION}' -X 'github.com/${REPO_SLUG}/internal/version.gitCommit=${GIT_COMMIT}'"
LINUX = $(CMD_NAME)-linux-$(VERSION)
MACOS_AMD = $(CMD_NAME)-macos-amd64-$(VERSION)
MACOS_ARM = $(CMD_NAME)-macos-arm64-$(VERSION)

.PHONY: dist
dist: dist-linux dist-darwin-amd64 dist-darwin-arm64  ## Cross compile binaries into ./dist/

dist-setup:
	mkdir -p ./bin ./dist

notarize-setup:
# Cannot be run locally, macnotary server is IP restricted
ifdef NOTARY_URI
	apt-get update && apt-get install unzip
	curl -LO $(NOTARY_BINARY_URL)
	unzip linux_$(BUILD_PIPELINE_ARCH).zip
	mv ./linux_$(BUILD_PIPELINE_ARCH)/macnotary /usr/local/bin/macnotary
else
	$(info Skipping notarize setup)
endif

notarize:
ifdef NOTARY_URI
	macnotary -f $(DIST_TAR) -m notarizeAndSign -u $(NOTARY_URI) -s $(NOTARY_SECRET) -k $(NOTARY_KEY_ID) -o $(DIST_TAR) -b mongodb.com
else
	$(info Skipping notarize)
endif

.PHONY: dist-linux
dist-linux: dist-setup
	GOOS=linux GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o ./bin/$(LINUX) .
	tar -zcvf dist/$(LINUX).tgz ./bin/$(LINUX) README.md

.PHONY: dist-darwin-amd64
dist-darwin-amd64: dist-setup notarize-setup
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o ./bin/$(MACOS_AMD) .
	tar -zcvf dist/$(MACOS_AMD).tgz ./bin/$(MACOS_AMD) README.md
	$(MAKE) DIST_TAR=dist/$(MACOS_AMD).tgz notarize

.PHONY: dist-darwin-arm64
dist-darwin-arm64: dist-setup notarize-setup
	GOOS=darwin GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o ./bin/$(MACOS_ARM) .
	tar -zcvf dist/$(MACOS_ARM).tgz ./bin/$(MACOS_ARM) README.md
	$(MAKE) DIST_TAR=dist/$(MACOS_ARM).tgz notarize


.PHONY: clean
clean: ## Clean up release artifacts
	rm -rf ./bin ./dist

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
