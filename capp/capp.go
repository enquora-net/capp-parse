/*
 * capp/capp.go
 * cappuccino
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

/*
 * Package capp is the exported library interface for capp-parse.
 *
 * It provides both single-file and project-scale parsing of Objective-J
 * source trees using the tree-sitter grammar. The cmd layer and the
 * Cappuccino compiler are the two intended consumers.
 *
 * Tree lifetime: ParseProject returns a ProjectResult whose FileResult
 * entries carry live *sitter.Tree values. The caller must call
 * ProjectResult.Close() when the trees are no longer needed. Walk
 * manages tree lifetime internally — trees are closed after each
 * EmitFn call returns.
 */
package capp

import (
	"fmt"
	"time"
)

// BoundaryVersion identifies the interop contract version at the
// Lisette/Go boundary. Increment when the public API changes in a
// way that affects generated Go code.
const BoundaryVersion = 1

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

// ParseMode converts a flag string to a Mode.
func ParseMode(s string) (Mode, error) {
	return parseMode(s)
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
// The tree is not exposed here; for compiler use, call ParseFile or ParseProject.
type ParseResult struct {
	HasError   bool
	Errors     []ParseError
	NodeCount  int
	Sexp       string
	PrettySexp string
	Timing     *Timing // nil when Benchmark is false
}

// Parse parses a single file as described by cfg and returns a ParseResult
// suitable for CLI output. The tree is closed before returning.
func Parse(cfg ParseConfig) (ParseResult, error) {
	return runParse(cfg)
}

// ---------------------------------------------------------------------------
// FileResult — compiler-facing single-file result
// ---------------------------------------------------------------------------

// FileResult is the outcome of parsing a single file for compiler use.
// Tree and Source are live until ProjectResult.Close() is called.
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
// Call Close() when the compiler pass is complete.
type ProjectResult struct {
	Files      []*FileResult
	Duration   time.Duration
	FileCount  int
	ErrorCount int
	ByteCount  int64
}

// Close releases all live tree-sitter trees held by this result.
// Must be called exactly once when the compiler pass is complete.
func (pr *ProjectResult) Close() {
	closeProjectResult(pr)
}

// ParseProject walks the source tree described by cfg, parses every
// matching file, and returns a ProjectResult with all trees live.
// The caller must call ProjectResult.Close() when done.
func ParseProject(cfg ProjectConfig) (*ProjectResult, error) {
	return runParseProject(cfg)
}

// ParseFile parses a single file and returns a live FileResult.
// The caller must close the tree via ProjectResult.Close() or directly
// via the sitter.Tree embedded in FileResult.Tree.
func ParseFile(cfg ParseConfig) (*FileResult, error) {
	return runParseFile(cfg)
}

// ---------------------------------------------------------------------------
// Walk — streaming project walk for CLI use
// ---------------------------------------------------------------------------

// EmitFn is called once per file as results are produced.
// The tree is closed immediately after EmitFn returns.
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

// Walk executes the parallel walk described by cfg.
// Tree lifetime is managed internally; trees are closed after each EmitFn call.
func Walk(cfg WalkConfig) (WalkSummary, error) {
	return runWalk(cfg)
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

// Debug walks a file or directory, stopping at the first parse error.
func Debug(cfg DebugConfig) (DebugResult, error) {
	return runDebug(cfg)
}
