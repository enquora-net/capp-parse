# capp-parse — Architecture

capp-parse is the Objective-J and JavaScript parser component of the
Cappuccino 2.0 toolchain. It replaces ad-hoc invocations of the Node-hosted
tree-sitter CLI with a self-contained binary and importable Go library,
providing a stable interface to the Objective-J grammar for the toolchain,
editor integrations, and third-party consumers.

---

## Governing principle

The library is the product. The binary is a convenience wrapper. Every
capability is expressed as an importable Go package before it is exposed as a
CLI command, and the command handlers carry no logic of their own.

This is not convention. It is the constraint that makes `capp-parse` a
composable component rather than a terminal artifact: the parent toolchain
absorbs the library packages without modification; the binary-specific
scaffolding (`cmd/root.go`, `cmd/version.go`, `cmd/capp-parse/main.go`) is
discarded at that point.

---

## Component boundaries

### Grammar loading (`grammar/`)

The tree-sitter Objective-J grammar is delivered as a platform-specific
dynamic library, distributed independently via
[tree-sitter-objj](https://github.com/enquora-net/tree-sitter-objj). It is
loaded at runtime via `purego`, crossing the C ABI without CGo.

This is the sole C boundary in the package. The dynamic library model makes
the grammar consumable by editor integrations and third-party tooling in any
language — not only Go — and decouples grammar versioning from toolchain
versioning.

`grammar/` owns path resolution and library loading. The public library
surface (`Parse`, `Walk`, `Debug`) accepts an optional explicit path;
`grammar.FindLibrary()` handles the ordered search over standard installation
locations when none is supplied.

### Core implementation (`internal/core/`)

All parse, walk, debug, and emit logic. Not directly importable by external
consumers; exposed through the public surface in the root package.

### Public library surface (`*.go` at root)

`Parse`, `Walk`, `Debug`, `EmitResult`, `EmitWalkSummary`, `EmitDebugProfile`,
`MakeEmitFn`, and associated config and result types. This is the stable
interface consumed by the parent toolchain and available to third-party
importers.

### CLI commands (`cmd/`)

`parse`, `debug`, and `verify` are implemented as fold-in seams. Their
`New*Cmd()` constructors accept no binary-specific dependencies and wire
directly into any Cobra root command. When absorbed into the parent toolchain:

- `cmd/parse/`, `cmd/debug/`, `cmd/verify/` move across unchanged.
- `cmd/root.go`, `cmd/version.go`, `cmd/capp-parse/main.go` are discarded.

---

## CGo constraint

CGo is not used. The go-tree-sitter binding itself has a CGo dependency, but
no additional C crossings are introduced by this package. The `grammar/`
package uses `purego` for dynamic loading; all other code is pure Go.

---

## Dependency direction

```
cmd/parse, cmd/debug, cmd/verify
        │
        ▼
  root package (Parse, Walk, Debug, ...)
        │
        ▼
  internal/core
        │
        ▼
  grammar/              go-tree-sitter / purego
```

Dependencies flow inward and downward. The grammar loader and core are
ignorant of CLI concerns. The public surface is ignorant of command structure.

---

## Versioning

Version, commit hash, and build date are injected at link time via ldflags.
`VERSION.txt` is the canonical version source for the release workflow;
`make release` reads it and tags the GitHub release accordingly. The running
binary can report all three values via `capp-parse version`.

---

## Relation to the parent toolchain

capp-parse is described from the parent toolchain's perspective in
[cappuccino/ARCHITECTURE.md](https://github.com/enquora-net/cappuccino/blob/main/ARCHITECTURE.md).
The grammar dynamic library boundary and the fold-in seam pattern are the
two architectural facts that matter at the integration point; both are stable
and will not change when the absorption occurs.
