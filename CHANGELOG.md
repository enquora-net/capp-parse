# Changelog

All notable changes to capp-parse are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## 2.0.0-beta2 — 2026-10-28

### Added

- `grammar.Locations()` — the supported grammar library install locations for
  the current platform, in search order, each marked as serving the current
  user, all users, or all users outside any managed location. These are the
  canonical list of supported locations.
- `grammar.FindLibraries()` — every grammar library present in those
  locations, in search order; the first is the one `grammar.Language()` loads.
- `grammar.Inspect()` — loads one grammar library and reports its tree-sitter
  ABI, its release version from tree-sitter's language metadata, and its build
  version from `tree_sitter_objj_version()`, when the library provides them.
- `host` package — the operating system and its version, the architecture the
  binary was built for, the machine's own architecture, and whether the
  process runs under architecture translation (Rosetta 2 on macOS, x64
  emulation on Windows on ARM; not detectable on Linux).

### Changed

- Grammar library install locations, searched in this order:
  - macOS: `~/Library/Application Support/dev.cappuccino/grammar`, then
    `/Library/Application Support/dev.cappuccino/grammar`, then
    `/usr/local/lib`
  - Linux: `$XDG_DATA_HOME/cappuccino/grammar` (default
    `~/.local/share/cappuccino/grammar`), then `/usr/local/lib`
  - Windows: `%LOCALAPPDATA%\Cappuccino\grammar`, then
    `%ProgramFiles%\Cappuccino\grammar`

- Locations no longer searched:
  - macOS: `~/Library/Application Support/github.com/enquora-net/grammar`
    and `/Library/Application Support/github.com/enquora-net/grammar`
  - Linux: `$XDG_DATA_HOME/github.com/enquora-net/grammar` and
    `$XDG_DATA_DIRS/github.com/enquora-net/grammar` (default
    `/usr/local/share/...`, then `/usr/share/...`); compiled libraries do not
    belong under `share`
  - Windows: `%LOCALAPPDATA%\github.com\enquora-net\grammar` and
    `%PROGRAMDATA%\github.com\enquora-net\grammar`; any user can create
    folders under `%PROGRAMDATA%`
- On Windows, the grammar library is loaded by its full path only, never
  through the system DLL search order.

---

## 2.0.0-beta1 — 2026-06-17

First public release of the Go-native Cappuccino Objective-J parser.

capp-parse replaces ad-hoc invocations of the Node-hosted tree-sitter CLI
with a self-contained binary and importable Go library. It is distributed as
a component of the Cappuccino 2.0 toolchain and is independently usable by
editor integrations, linters, and third-party pipelines in any language.

### Added

- Single binary distribution for macOS, Linux, and Windows on ARM64 and AMD64,
  cross-compiled via Zig from a macOS host
- `parse` — parse a single file, a directory tree, or a stream of file paths
  from stdin; human-readable output when stdout is a terminal, JSON per line
  otherwise; parallel workers configurable via `--workers`; exit code 2 on
  syntax errors
- `debug` — walk a source tree, stop at the first parse error, display source
  context, parent node chain, and full concrete syntax tree, and open the
  failing file in Xcode at the error line via `xed`; `--profile` for per-file
  timing; `--no-xcode` to suppress `xed`
- `verify` — check that the grammar dynamic library is installed and locatable;
  reports the full search path on failure; JSON output when stdout is not a
  terminal
- `version` — report version, commit hash, and build date
- Grammar dynamic library search across standard installation paths with
  explicit override via `--grammar`
- Importable Go library (`github.com/enquora-net/capp-parse`) providing
  `Parse`, `Walk`, `Debug`, and emit functions for use by the parent toolchain
  and third-party consumers
- Deterministic version, commit, and build date embedding via ldflags
- SHA-256 checksum files for all release artifacts
