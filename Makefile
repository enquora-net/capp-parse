# Makefile
# capp-parse
#
# Created by David Richardson on 2026-04-10.
#
# Development workflow:
#   make build        Build for the current host
#   make run          Build and run (no arguments)
#   make test         Run the full test suite
#   make install      Install to $GOPATH/bin
#
# Release workflow:
#   1. Update VERSION.txt and CHANGELOG.md.
#   2. Commit and push.
#   3. make release
#   Requires: gh (GitHub CLI), authenticated with repo write access.
#
# Cross-compilation:
#   capp-parse has a CGo dependency (tree-sitter via go-tree-sitter).
#   All six platform targets are produced using Zig as the C cross-compiler.
#   The macOS SDK path is resolved via xcrun; Darwin targets must be built
#   from a macOS host.
#
#   Linux targets pin glibc 2.28 for broad distribution compatibility.
#   Windows targets produce .exe binaries via the Zig mingw-w64 toolchain.

BINARY     := capp-parse
BUILD_DIR  := build
MODULE     := github.com/enquora-net/capp-parse
VERSION    := $(shell head -1 VERSION.txt)
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
MACOS_SDK  := $(shell xcrun --sdk macosx --show-sdk-path 2>/dev/null)

LDFLAGS := -ldflags "\
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(BUILD_DATE)"

DARWIN_ARM := $(BUILD_DIR)/$(BINARY)_darwin_arm64_$(VERSION)
DARWIN_AMD := $(BUILD_DIR)/$(BINARY)_darwin_amd64_$(VERSION)
LINUX_ARM  := $(BUILD_DIR)/$(BINARY)_linux_arm64_$(VERSION)
LINUX_AMD  := $(BUILD_DIR)/$(BINARY)_linux_amd64_$(VERSION)
WIN_ARM    := $(BUILD_DIR)/$(BINARY)_windows_arm64_$(VERSION).exe
WIN_AMD    := $(BUILD_DIR)/$(BINARY)_windows_amd64_$(VERSION).exe

.PHONY: all build run test clean install deps build-all checksums release \
        darwin linux windows arm64 amd64 ci-linux ci-darwin

all: build

# ---------------------------------------------------------------------------
# Development
# ---------------------------------------------------------------------------

build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/capp-parse

run: build
	$(BUILD_DIR)/$(BINARY)

test:
	go test -race -v ./...

install:
	CGO_ENABLED=1 go install $(LDFLAGS) ./cmd/capp-parse

deps:
	go mod download
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)

# ---------------------------------------------------------------------------
# Release builds — all six platform targets via Zig CC
# ---------------------------------------------------------------------------

build-all: $(DARWIN_ARM) $(DARWIN_AMD) $(LINUX_ARM) $(LINUX_AMD) $(WIN_ARM) $(WIN_AMD)

ci-darwin: $(DARWIN_ARM) $(DARWIN_AMD)

ci-linux: $(LINUX_ARM) $(LINUX_AMD) $(WIN_ARM) $(WIN_AMD)

darwin: $(DARWIN_ARM) $(DARWIN_AMD)

linux: $(LINUX_ARM) $(LINUX_AMD)

windows: $(WIN_ARM) $(WIN_AMD)

arm64: $(DARWIN_ARM) $(LINUX_ARM) $(WIN_ARM)

amd64: $(DARWIN_AMD) $(LINUX_AMD) $(WIN_AMD)

$(DARWIN_ARM):
	@mkdir -p $(BUILD_DIR)
	@echo "building darwin/arm64 ..."
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
	CC="zig cc -target aarch64-macos -isysroot $(MACOS_SDK) -L$(MACOS_SDK)/usr/lib -F$(MACOS_SDK)/System/Library/Frameworks" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse

$(DARWIN_AMD):
	@mkdir -p $(BUILD_DIR)
	@echo "building darwin/amd64 ..."
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
	CC="zig cc -target x86_64-macos -isysroot $(MACOS_SDK) -L$(MACOS_SDK)/usr/lib -F$(MACOS_SDK)/System/Library/Frameworks" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse

$(LINUX_ARM):
	@mkdir -p $(BUILD_DIR)
	@echo "building linux/arm64 ..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
	CC="zig cc -target aarch64-linux-gnu.2.28" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse

$(LINUX_AMD):
	@mkdir -p $(BUILD_DIR)
	@echo "building linux/amd64 ..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
	CC="zig cc -target x86_64-linux-gnu.2.28" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse

$(WIN_ARM):
	@mkdir -p $(BUILD_DIR)
	@echo "building windows/arm64 ..."
	@CGO_ENABLED=1 GOOS=windows GOARCH=arm64 \
	CC="zig cc -target aarch64-windows" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse

$(WIN_AMD):
	@mkdir -p $(BUILD_DIR)
	@echo "building windows/amd64 ..."
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
	CC="zig cc -target x86_64-windows" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse

checksums:
	@echo "computing checksums ..."
	@cd $(BUILD_DIR) && for f in $(BINARY)_*; do \
		[ -f "$$f" ] && shasum -a 256 "$$f" | awk '{print $$1}' > "$$f.sha256" \
		&& echo "  $$f.sha256"; \
	done || true
	@echo "done."

# ---------------------------------------------------------------------------
# Release — requires gh authenticated with repo write access
# ---------------------------------------------------------------------------

release: build-all checksums
	@if ! gh auth status > /dev/null 2>&1; then \
		echo "error: gh is not authenticated — run 'gh auth login' first" >&2; exit 1; \
	fi
	@tag="v$(VERSION)"; \
	notes=$$(awk "/^## $(VERSION)/{found=1; next} found && /^## /{exit} found{print}" CHANGELOG.md); \
	[ -z "$$notes" ] && notes="Release $(VERSION)"; \
	prerelease=""; \
	case "$(VERSION)" in *-*) prerelease="--prerelease" ;; esac; \
	echo "creating release $$tag ..."; \
	gh release create "$$tag" \
		--title "$$tag" \
		--notes "$$notes" \
		$$prerelease; \
	echo "uploading artifacts ..."; \
	for f in $(BUILD_DIR)/$(BINARY)_* $(BUILD_DIR)/*.sha256; do \
		[ -f "$$f" ] || continue; \
		gh release upload "$$tag" "$$f"; \
		echo "  uploaded: $$(basename $$f)"; \
	done; \
	echo ""; \
	echo "release $$tag published."; \
	gh release view "$$tag" --web

.DEFAULT_GOAL := build
