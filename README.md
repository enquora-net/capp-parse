# capp-parse

Cappuccino 2.0 – Objective-J and JavaScript source parser for the
[Cappuccino Project's](https://cappuccino.dev) port of AppKit and Foundation
to the web.

capp-parse is a component of the [Cappuccino toolchain](https://github.com/enquora-net/cappuccino).
It wraps the [tree-sitter-objj](https://github.com/enquora-net/tree-sitter-objj) grammar
as a self-contained binary and importable Go library, providing single-file and
recursive tree-walking parse modes along with a debug output mode designed to
accelerate grammar and compiler development workflows.

The binary is consumed internally by the toolchain and is available independently
for use in pipelines, editor integrations, and third-party tooling in any language.

---

## Installation

Download the binary for your platform from the
[releases page](https://github.com/enquora-net/capp-parse/releases), make it
executable, and place it in your PATH.

The Objective-J Tree-sitter dynamic library must be installed separately.
Download the appropriate binary for your platform from the
[tree-sitter-objj releases page](https://github.com/enquora-net/tree-sitter-objj/releases)
and place it in `/usr/local/lib`.

On macOS, Gatekeeper quarantines downloaded files. Strip the attribute from
both before first use:

```sh
xattr -d com.apple.quarantine capp-parse
xattr -d com.apple.quarantine /usr/local/lib/libtree-sitter-objj.dylib
```

Use `capp-parse verify` to confirm the grammar library is locatable.

> **Pre-release.** This software is in developer preview. It should only be
> used in environments where rollback or recovery is in place.

---

## Prerequisites

The compiled grammar dynamic library must be installed and locatable. The
following paths are searched in order on macOS:

```
~/Library/Application Support/github.com/enquora-net/grammar/
/Library/Application Support/github.com/enquora-net/grammar/
/usr/local/lib
```

For this release, `/usr/local/lib` is the supported manual install location.
An explicit path can be supplied to any command via `--grammar`.

---

## Commands

### parse

Parse a file, directory tree, or stream of file paths from stdin.

```
capp-parse parse [flags] [path]

Flags:
  -g, --grammar  path              explicit grammar dylib path, overrides search
      --mode     auto|objj|js|both file type filter (default: auto)
      --workers  N                 parallel workers (default: GOMAXPROCS)
  -f, --format   sexp|json|ast     output format for single-file mode (default: sexp)
      --src                        read source bytes from stdin
```

```sh
capp-parse parse src/AppController.j
capp-parse parse Frameworks/
find . -name '*.j' | capp-parse parse
cat src/AppController.j | capp-parse parse --src --mode objj
```

Output is human-readable when stdout is a terminal; JSON per line otherwise.

Exit codes: `0` = OK, `2` = syntax errors found.

---

### debug

Walk a file or directory tree, stop at the first parse error, and display full
diagnostic output: source context, parent node chain, and full concrete syntax
tree. The failing file is opened in Xcode at the error line via `xed`.

```
capp-parse debug [flags] <path>

Flags:
  -g, --grammar   path              explicit grammar dylib path, overrides search
  -m, --mode      auto|objj|js|both file type filter (default: auto)
  -p, --profile                     print per-file timing and aggregate summary
  -c, --context   N                 source lines of context around error (default: 3)
      --no-xcode                    suppress xed invocation
```

```sh
capp-parse debug Frameworks/AppKit
capp-parse debug --profile --no-xcode Frameworks/
capp-parse debug --context 5 src/AppController.j
```

Exit codes: `0` = all files OK, `2` = parse error found.

---

### verify

Check that capp-parse prerequisites are installed and locatable. Reports the
full ordered search path on failure. Exits 0 if all prerequisites are
satisfied, 1 otherwise.

```
capp-parse verify [flags]

Flags:
  -g, --grammar  path   explicit grammar dylib path, overrides search
```

```sh
capp-parse verify
capp-parse verify --grammar /opt/local/lib/tree-sitter/libtree-sitter-objj.dylib
```

When stdout is not a terminal, output is JSON:

```json
{"ok": true,  "grammar": {"found": true,  "path": "/usr/local/lib/libtree-sitter-objj.dylib"}}
{"ok": false, "grammar": {"found": false, "searched": [...]}}
```

---

### version

```sh
capp-parse version
```

---

## Performance

A full cold walk of the Cappuccino AppKit corpus (225 files, 5.5 MB), Apple
Silicon M4, macOS 26:

```
files: 225   total: <300ms   mean: 1.2ms   max: 35ms
```

---

## Build

```sh
# Current platform
make build

# All platforms (requires Zig)
make build-all

# Darwin only
make darwin
```

Version is read from `VERSION.txt`. Cross-compilation uses Zig as the C
compiler for all targets; Darwin targets require a macOS host and Xcode
command-line tools.

---

## As a Library

```go
import capp "github.com/enquora-net/capp-parse"

result, err := capp.Parse(capp.ParseConfig{
    Path:        "src/AppController.j",
    Src:         src,
    Mode:        capp.ModeObjJ,
    GrammarPath: "",   // use default search
})
```

`cmd/parse/`, `cmd/debug/`, and `cmd/verify/` are fold-in seams: their
`New*Cmd()` constructors wire directly into the parent toolchain's root
command without modification when this package is absorbed into the omnibus
binary.

---

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for design rationale and component
boundaries.

---

## License

Copyright David Richardson. The binaries distributed here may be freely run
for any purpose. All other rights are reserved pending transfer to the
Cappuccino Project, at which point this software will be released under
AGPL-3.0.
