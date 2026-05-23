/*
 * internal/types/types.go
 * capp-parse
 *
 * All shared type definitions for the capp-parse public API.
 * Pure Go — no CGo, no go-tree-sitter imports.
 * Consumed by internal/core (implementation) and the root package (facade).
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package types

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Mode
// ---------------------------------------------------------------------------

// Mode controls which file extensions are accepted for parsing.
type Mode int

const (
	ModeAuto Mode = iota // infer from extension: .j .sj → objj, .js → js
	ModeObjj             // .j and .sj only
	ModeJS               // .js only
	ModeBoth             // .j .sj .js
)

// Accept reports whether path should be parsed under mode m.
func (m Mode) Accept(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch m {
	case ModeObjj:
		return ext == ".j" || ext == ".sj"
	case ModeJS:
		return ext == ".js"
	case ModeBoth:
		return ext == ".j" || ext == ".sj" || ext == ".js"
	default: // ModeAuto
		return ext == ".j" || ext == ".sj" || ext == ".js"
	}
}

// ParseMode converts a flag string to a Mode.
func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(s) {
	case "auto", "":
		return ModeAuto, nil
	case "objj":
		return ModeObjj, nil
	case "js":
		return ModeJS, nil
	case "both":
		return ModeBoth, nil
	default:
		return 0, fmt.Errorf("unknown mode %q: use auto | objj | js | both", s)
	}
}

// ---------------------------------------------------------------------------
// Format
// ---------------------------------------------------------------------------

// Format controls the output representation for single-file parse results.
type Format int

const (
	FormatDefault Format = iota // same as FormatSexp
	FormatSexp                  // tree-sitter s-expression
	FormatJSON                  // JSON tree (not yet implemented)
	FormatAST                   // typed IR (not yet implemented)
)

// ParseFormat converts a flag string to a Format.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "sexp", "":
		return FormatSexp, nil
	case "json":
		return FormatJSON, nil
	case "ast":
		return FormatAST, nil
	default:
		return 0, fmt.Errorf("unknown format %q: use sexp | json | ast", s)
	}
}

// ---------------------------------------------------------------------------
// ParseError
// ---------------------------------------------------------------------------

// ParseError is a structured parse error with source position.
type ParseError struct {
	Row     uint32
	Column  uint32
	Message string
}

// String returns a human-readable representation of the error.
func (e ParseError) String() string {
	return fmt.Sprintf("line %d col %d: %s", e.Row+1, e.Column+1, e.Message)
}

// ---------------------------------------------------------------------------
// Timing
// ---------------------------------------------------------------------------

// Timing carries diagnostic timing for one parse operation.
type Timing struct {
	Elapsed time.Duration
	Bytes   int
}

// ---------------------------------------------------------------------------
// Single-file parse
// ---------------------------------------------------------------------------

// ParseConfig is the complete specification for a single-file parse.
type ParseConfig struct {
	Path        string
	Src         []byte
	Mode        Mode
	Format      Format
	GrammarPath string // explicit dylib path; empty = search
	Benchmark   bool
}

// ParseResult is the outcome of a single-file parse for CLI consumption.
type ParseResult struct {
	HasError   bool
	Errors     []ParseError
	NodeCount  int
	Sexp       string
	PrettySexp string
	Timing     *Timing // nil when Benchmark is false
}

// ---------------------------------------------------------------------------
// FileResult — compiler-facing single-file result
// ---------------------------------------------------------------------------

// FileResult is the outcome of parsing a single file for compiler use.
// Tree and Source are live until Close() is called on the owning ProjectResult.
type FileResult struct {
	Path      string
	Source    []byte
	Tree      interface{} // *sitter.Tree — typed as interface{} at the boundary
	Errors    []ParseError
	NodeCount int
	Duration  time.Duration
}

// HasError reports whether the file had any parse errors.
func (f *FileResult) HasError() bool {
	return len(f.Errors) > 0
}

// ---------------------------------------------------------------------------
// Project-scale parse
// ---------------------------------------------------------------------------

// ProjectConfig is the complete specification for a project-scale parse.
type ProjectConfig struct {
	Root        string   // directory to walk; empty when Paths is set
	Paths       []string // explicit path list; overrides Root
	Mode        Mode
	GrammarPath string // explicit dylib path; empty = search
	Workers     int
}

// ProjectResult is the aggregate outcome of a project-scale parse.
// Call Close() when the compiler pass is complete to release live trees.
type ProjectResult struct {
	Files      []*FileResult
	Duration   time.Duration
	FileCount  int
	ErrorCount int
	ByteCount  int64
	closeFn    func() // set by internal/core
}

// SetCloseFn is called by internal/core to register the tree-release function.
func (pr *ProjectResult) SetCloseFn(fn func()) { pr.closeFn = fn }

// Close releases all live tree-sitter trees held by this result.
// Must be called exactly once when the compiler pass is complete.
func (pr *ProjectResult) Close() {
	if pr.closeFn != nil {
		pr.closeFn()
	}
}

// ---------------------------------------------------------------------------
// Walk — streaming project walk for CLI use
// ---------------------------------------------------------------------------

// EmitFn is called once per file as results are produced.
type EmitFn func(path string, result ParseResult)

// WalkConfig is the complete specification for a source-tree walk.
type WalkConfig struct {
	Root        string   // directory to walk; empty when Paths is set
	Paths       []string // explicit path list (stdin stream mode); overrides Root
	Mode        Mode
	GrammarPath string // explicit dylib path; empty = search
	Workers     int
	FailFast    bool
	Quiet       bool
	Benchmark   bool
	Stdout      interface{ Write([]byte) (int, error) }
	Stderr      interface{ Write([]byte) (int, error) }
	EmitFn      EmitFn // called per file; nil = silent
}

// WalkSummary is the aggregate outcome of a source-tree walk.
type WalkSummary struct {
	TotalFiles int
	OKFiles    int
	ErrorFiles int
	TotalBytes int
	Elapsed    time.Duration
	Bench      *BenchReport // nil when Benchmark is false
}

// BenchReport is the aggregate timing data collected during a walk.
type BenchReport struct {
	WallTime   time.Duration
	ParseTotal time.Duration
	Files      int
	Bytes      int
	Events     []BenchEvent
}

// BenchEvent is a single timed parse observation.
type BenchEvent struct {
	Path    string
	Elapsed time.Duration
	Bytes   int
}

// ---------------------------------------------------------------------------
// Debug
// ---------------------------------------------------------------------------

// DebugConfig is the complete specification for a debug walk.
type DebugConfig struct {
	Path         string
	Mode         Mode
	GrammarPath  string // explicit dylib path; empty = search
	Profile      bool
	ContextLines int
	NoXcode      bool
}

// DebugResult is the outcome of a debug walk.
type DebugResult struct {
	HasError bool
	Profile  *DebugProfile // nil when Profile is false
}

// DebugProfile is the timing data collected during a debug walk.
type DebugProfile struct {
	Events []DebugEvent
	Total  time.Duration
}

// DebugEvent is a single timed parse observation from a debug walk.
type DebugEvent struct {
	Path    string
	Elapsed time.Duration
	OK      bool
}
