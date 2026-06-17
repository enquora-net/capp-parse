# Developer Guide — capp-parse

## Prerequisites

**Go 1.26 or later** — required. See [go.dev/dl](https://go.dev/dl) or install
via MacPorts (`sudo port install go`) or Homebrew (`brew install go`).

**Zig** — required for cross-platform release builds. Install via MacPorts
(`sudo port install zig`) or Homebrew (`brew install zig`). Not required for
native development builds.

**tree-sitter-objj dynamic library** — required at runtime. Install from the
[tree-sitter-objj releases page](https://github.com/enquora-net/tree-sitter-objj/releases)
to `/usr/local/lib`. Confirm with `capp-parse verify`.

**gh (GitHub CLI)** — required for `make release` only.

---

## Build

```sh
# Native development build
make build

# Run tests
make test

# Install to $GOPATH/bin
make install

# All six platform targets (requires Zig, macOS host for Darwin targets)
make build-all

# Checksums only
make checksums
```

---

## Release workflow

1. Update `VERSION.txt` with the new version string.
2. Add a matching `## <version> — <date>` section to `CHANGELOG.md`.
3. Commit and push.
4. `make release`

The `release` target builds all platforms, computes checksums, creates a
GitHub release tagged `v<VERSION>`, extracts release notes from `CHANGELOG.md`,
and uploads all artifacts. Pre-release versions (any version containing a
hyphen) are marked as pre-releases on GitHub automatically.

---

## Code standards

- `gofmt -s -w .` before every commit.
- No CGo beyond the go-tree-sitter dependency. The grammar is loaded via
  purego; no additional C ABI crossings are permitted.
- No capability in command handlers. All logic lives in the library packages;
  commands are consumers only.
- `make test` must pass clean under `-race` before any push.

---

## Repository layout

```
capp-parse/
├── cmd/
│   ├── capp-parse/     entry point (main.go)
│   ├── debug/          debug subcommand
│   ├── parse/          parse subcommand
│   ├── verify/         verify subcommand
│   ├── root.go         root command (discarded on toolchain fold-in)
│   └── version.go      version command (discarded on toolchain fold-in)
├── grammar/            grammar dynamic library loading and path resolution
├── internal/
│   └── core/           parse, walk, debug, and emit implementation
├── *.go                public library surface (Parse, Walk, Debug, emit fns)
├── ARCHITECTURE.md
├── CHANGELOG.md
├── DEVELOPERS.md
├── LICENSE
├── Makefile
├── README.md
└── VERSION.txt
```
