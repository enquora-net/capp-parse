# capp-parse

Pure-Go Cappuccino Objective-J source parser built on
[gotreesitter](https://github.com/odvcencio/gotreesitter).  No CGo.  No C
toolchain.  Cross-compiles to any `GOOS/GOARCH` target Go supports, including
`wasip1`.

---

## Prerequisites

```
go get github.com/odvcencio/gotreesitter@latest
go get github.com/spf13/cobra@latest
```

---

## One-time grammar blob generation

gotreesitter reads the same parse-table format the C runtime uses, extracted
from `parser.c` by the `ts2go` tool and stored as a compressed binary blob.
This step is a build-time prerequisite; the blob is committed to the repo so
downstream consumers need only `go build`.

```sh
go run github.com/odvcencio/gotreesitter/cmd/ts2go \
    -input  path/to/tree-sitter-objj/src/parser.c \
    -out    internal/grammar/objj.bin \
    -compact
```

Re-run whenever `grammar.js` / `parser.c` changes.

---

## Build

```sh
go build -o capp-parse .

# With version stamped in binary:
go build -ldflags "-X github.com/cappuccino/capp-parse/cmd.Version=$(git describe --tags)" \
    -o capp-parse .
```

---

## Commands

### parse — single file

```
capp-parse parse [flags] <file>

Flags:
  --mode   auto | objj | js        file type filter (default: auto)
  --bench                          print diagnostic timing
  --sexp                           print the full S-expression tree
```

```sh
capp-parse parse src/AppController.j
capp-parse parse --sexp --bench Foundation/CPObject.j
```

Exit codes: 0 = OK, 2 = syntax errors found.

### walk — source tree

```
capp-parse walk [flags] <directory>

Flags:
  --mode     both | objj | js      file type filter (default: both)
  --bench                          print aggregate timing summary
  --workers  N                     parallel workers (default: GOMAXPROCS)
  --fail-fast                      stop at first error
  --quiet                          suppress per-file OK lines
```

```sh
capp-parse walk --bench Frameworks/
capp-parse walk --mode objj --quiet --bench .
```

Exit codes: 0 = all files OK, 2 = one or more files with syntax errors.

### version

```
capp-parse version
```

---

## External scanner

The grammar's `scanner.c` is ported verbatim to Go in
`internal/scanner/scanner.go`.  The scanner is stateless (create returns nil,
serialize returns 0 bytes) so the port is a pure translation of the C logic
with no state-management complexity.

Token enum order matches `externals` in `grammar.js` exactly:

| Index | Token                |
|-------|----------------------|
| 0     | `_automatic_semicolon` |
| 1     | `_template_chars`    |
| 2     | `_ternary_qmark`     |
| 3     | `html_comment`       |
| 4     | `\|\|` (LOGICAL_OR)  |
| 5     | `escape_sequence`    |
| 6     | `regex_pattern`      |
| 7     | `jsx_text`           |

---

## Performance notes

gotreesitter full-parse throughput is roughly 10–11× slower than the C runtime
on the benchmark in the upstream README (19 ms vs 1.75 ms on a 19 KB Go file).
For the Cappuccino source tree the expectation is:

- **C runtime sub-second cold walk** → **gotreesitter ~10–15 s cold walk**
- Incremental reparse (keystroke-at-a-time) is ~104× *faster* than the C
  runtime due to Go-managed tree memory and zero-copy subtree reuse.

The `--bench` flag reports wall time, total parse-CPU time, per-file breakdown,
and throughput in MB/s, giving a reproducible baseline for the evaluation.

---

## Evaluation criteria

This project exists to answer one question: is gotreesitter's cold-walk
throughput acceptable for `capp-parse` as a standalone CLI tool given the
advantages of a pure-Go, CGo-free, WASM-compilable binary?

If the answer is no, the fallback is the native tree-sitter C runtime accessed
via CGo or a dynamic-library subprocess — the same path the existing parser
uses.  The `--bench` output from a full Cappuccino source walk is the deciding
datum.
