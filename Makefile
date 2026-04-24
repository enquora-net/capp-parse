# Makefile
# capp-parse

BINARY       := capp-parse
MODULE       := capp-parse
VERSION      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS      := -ldflags "-X main.version=$(VERSION)"

DARWIN_ARM   := $(BINARY)-darwin-arm64
DARWIN_AMD   := $(BINARY)-darwin-amd64
LINUX_ARM    := $(BINARY)-linux-arm64
LINUX_AMD    := $(BINARY)-linux-amd64
WINDOWS_ARM  := $(BINARY)-windows-arm64.exe
WINDOWS_AMD  := $(BINARY)-windows-amd64.exe

.PHONY: all clean darwin linux windows arm64 amd64

all: $(DARWIN_ARM) $(DARWIN_AMD) $(LINUX_ARM) $(LINUX_AMD) $(WINDOWS_ARM) $(WINDOWS_AMD)

darwin: $(DARWIN_ARM) $(DARWIN_AMD)

linux: $(LINUX_ARM) $(LINUX_AMD)

windows: $(WINDOWS_ARM) $(WINDOWS_AMD)

arm64: $(DARWIN_ARM) $(LINUX_ARM) $(WINDOWS_ARM)

amd64: $(DARWIN_AMD) $(LINUX_AMD) $(WINDOWS_AMD)

$(DARWIN_ARM):
    GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $@ .

$(DARWIN_AMD):
    GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $@ .

$(LINUX_ARM):
    GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o $@ .

$(LINUX_AMD):
    GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $@ .

$(WINDOWS_ARM):
    GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $@ .

$(WINDOWS_AMD):
    GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $@ .

clean:
    rm -f $(DARWIN_ARM) $(DARWIN_AMD) $(LINUX_ARM) $(LINUX_AMD) $(WINDOWS_ARM) $(WINDOWS_AMD)
