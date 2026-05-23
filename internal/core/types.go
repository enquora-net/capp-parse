/*
 * internal/core/types.go
 * capp-parse
 *
 * Internal type definitions for the core package.
 * Parallel to the root public types; conversion is performed at the facade boundary.
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 */
package core

import "time"

// ParseConfig is the internal specification for a single-file parse.
// Mode and Format are carried as int to avoid dependency on root package types.
type ParseConfig struct {
	Path        string
	Src         []byte
	Mode        int
	Format      int
	GrammarPath string
	Benchmark   bool
}

// ParseError is a structured parse error with source position.
type ParseError struct {
	Row     uint32
	Column  uint32
	Message string
}

// Timing carries diagnostic timing for one parse operation.
type Timing struct {
	Elapsed time.Duration
	Bytes   int
}

// ParseResult is the outcome of a single-file parse.
type ParseResult struct {
	HasError   bool
	Errors     []ParseError
	NodeCount  int
	Sexp       string
	PrettySexp string
	Timing     *Timing
}

// FileResult is the outcome of parsing a single file for compiler use.
type FileResult struct {
	Path      string
	Source    []byte
	Tree      interface{}
	Errors    []ParseError
	NodeCount int
	Duration  time.Duration
}

// ProjectConfig is the internal specification for a project-scale parse.
type ProjectConfig struct {
	Root        string
	Paths       []string
	Mode        int
	GrammarPath string
	Workers     int
}

// ProjectResult is the aggregate outcome of a project-scale parse.
type ProjectResult struct {
	Files      []*FileResult
	Duration   time.Duration
	FileCount  int
	ErrorCount int
	ByteCount  int64
}

// EmitFn is called once per file during a walk.
type EmitFn func(path string, result ParseResult)

// WalkConfig is the internal specification for a source-tree walk.
type WalkConfig struct {
	Root        string
	Paths       []string
	Mode        int
	GrammarPath string
	Workers     int
	FailFast    bool
	Quiet       bool
	Benchmark   bool
	Stdout      interface{ Write([]byte) (int, error) }
	Stderr      interface{ Write([]byte) (int, error) }
	EmitFn      EmitFn
}

// WalkSummary is the aggregate outcome of a source-tree walk.
type WalkSummary struct {
	TotalFiles int
	OKFiles    int
	ErrorFiles int
	TotalBytes int
	Elapsed    time.Duration
	Bench      *BenchReport
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

// DebugConfig is the internal specification for a debug walk.
type DebugConfig struct {
	Path         string
	Mode         int
	GrammarPath  string
	Profile      bool
	ContextLines int
	NoXcode      bool
}

// DebugResult is the outcome of a debug walk.
type DebugResult struct {
	HasError bool
	Profile  *DebugProfile
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
