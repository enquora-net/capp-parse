/*
 * internal/capp/capp.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package capp is the sole internal API consumed by the cmd layer.
 * cmd imports nothing else from internal/.
 * Package capp is the sole internal API consumed by the cmd layer.
 * cmd imports nothing else from internal/.
 */
package capp

import "time"

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

// ---------------------------------------------------------------------------
// Parse
// ---------------------------------------------------------------------------

// ParseConfig is the complete specification for a single-file parse.
type ParseConfig struct {
	Path      string
	Src       []byte
	Mode      Mode
	Benchmark bool
	Sexp      bool
}

// ParseResult is the complete outcome of a single-file parse.
type ParseResult struct {
	HasError  bool
	Errors    []string
	NodeCount int
	Sexp      string
	Timing    *Timing // nil when Benchmark is false
}

// Timing carries diagnostic timing for one parse operation.
type Timing struct {
	Elapsed time.Duration
	Bytes   int
}

// Parse parses a single file as described by cfg.
func Parse(cfg ParseConfig) (ParseResult, error) {
	return runParse(cfg)
}

// ---------------------------------------------------------------------------
// Walk
// ---------------------------------------------------------------------------

// EmitFn is called once per file as results are produced.
// It is the caller's responsibility to format and write the output.
// When nil, Walk emits nothing per-file.
type EmitFn func(path string, result ParseResult)

// WalkConfig is the complete specification for a source-tree walk.
type WalkConfig struct {
	Root      string   // directory to walk; empty when Paths is set
	Paths     []string // explicit path list (stdin stream mode); overrides Root
	Mode      Mode
	Workers   int
	FailFast  bool
	Quiet     bool
	Benchmark bool
	Stdout    interface{ Write([]byte) (int, error) }
	Stderr    interface{ Write([]byte) (int, error) }
	EmitFn    EmitFn // called per file; nil = silent
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
func Walk(cfg WalkConfig) (WalkSummary, error) {
	return runWalk(cfg)
}

// ParseMode converts a flag string to a Mode.
func ParseMode(s string) (Mode, error) {
	return parseMode(s)
}
