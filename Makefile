# Makefile
# capp-parse
#
# Cross-compilation via Zig. All targets buildable from any host.
# macOS SDK path resolved via xcrun for cross-architecture Darwin targets.
#
# Each build artifact is accompanied by a individual .sha256 checksum file.

BINARY       := capp-parse
MODULE       := github.com/enquora-net/capp-parse
VERSION      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS      := -ldflags "-X main.version=$(VERSION)"
MACOS_SDK    := $(shell xcrun --sdk macosx --show-sdk-path)
BUILD_DIR    := build

DARWIN_ARM   := $(BUILD_DIR)/$(BINARY)-darwin-arm64
DARWIN_AMD   := $(BUILD_DIR)/$(BINARY)-darwin-amd64
LINUX_ARM    := $(BUILD_DIR)/$(BINARY)-linux-arm64
LINUX_AMD    := $(BUILD_DIR)/$(BINARY)-linux-amd64
WINDOWS_ARM  := $(BUILD_DIR)/$(BINARY)-windows-arm64.exe
WINDOWS_AMD  := $(BUILD_DIR)/$(BINARY)-windows-amd64.exe

.PHONY: all clean darwin linux windows arm64 amd64 native ci-linux ci-darwin

all: $(DARWIN_ARM) $(DARWIN_AMD) $(LINUX_ARM) $(LINUX_AMD) $(WINDOWS_ARM) $(WINDOWS_AMD)

native:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/capp-parse

ci-linux: $(LINUX_ARM) $(LINUX_AMD) $(WINDOWS_ARM) $(WINDOWS_AMD)

ci-darwin: $(DARWIN_ARM) $(DARWIN_AMD)

darwin: $(DARWIN_ARM) $(DARWIN_AMD)

linux: $(LINUX_ARM) $(LINUX_AMD)

windows: $(WINDOWS_ARM) $(WINDOWS_AMD)

arm64: $(DARWIN_ARM) $(LINUX_ARM) $(WINDOWS_ARM)

amd64: $(DARWIN_AMD) $(LINUX_AMD) $(WINDOWS_AMD)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(DARWIN_ARM): $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin  GOARCH=arm64 \
	CC="zig cc -target aarch64-macos -isysroot $(MACOS_SDK) -L$(MACOS_SDK)/usr/lib -F$(MACOS_SDK)/System/Library/Frameworks" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse
	cd $(BUILD_DIR) && shasum -a 256 $(notdir $@) > $(notdir $@).sha256

$(DARWIN_AMD): $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin  GOARCH=amd64 \
	CC="zig cc -target x86_64-macos -isysroot $(MACOS_SDK) -L$(MACOS_SDK)/usr/lib -F$(MACOS_SDK)/System/Library/Frameworks" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse
	cd $(BUILD_DIR) && shasum -a 256 $(notdir $@) > $(notdir $@).sha256

$(LINUX_ARM): $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux   GOARCH=arm64 \
	CC="zig cc -target aarch64-linux-gnu.2.28" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse
	cd $(BUILD_DIR) && shasum -a 256 $(notdir $@) > $(notdir $@).sha256

$(LINUX_AMD): $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux   GOARCH=amd64 \
	CC="zig cc -target x86_64-linux-gnu.2.28" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse
	cd $(BUILD_DIR) && shasum -a 256 $(notdir $@) > $(notdir $@).sha256

$(WINDOWS_ARM): $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 \
	CC="zig cc -target aarch64-windows" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse
	cd $(BUILD_DIR) && shasum -a 256 $(notdir $@) > $(notdir $@).sha256

$(WINDOWS_AMD): $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
	CC="zig cc -target x86_64-windows" \
	go build $(LDFLAGS) -o $@ ./cmd/capp-parse
	cd $(BUILD_DIR) && shasum -a 256 $(notdir $@) > $(notdir $@).sha256

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY)
