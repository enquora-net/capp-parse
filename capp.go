/*
 * capp.go
 * capp-parse
 *
 * Public facade for the capp-parse library.
 * Import path: github.com/enquora-net/capp-parse
 *
 * All public types are defined in types.go at the root package level.
 * Function implementations delegate to internal/core via field conversion.
 * No internal package references appear in this file or types.go.
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
 */
package capp

import "github.com/enquora-net/capp-parse/internal/core"

// BoundaryVersion identifies the interop contract version at the
// Lisette/Go boundary. Increment when the public API changes.
const BoundaryVersion = 1

// ---------------------------------------------------------------------------
// Single-file parse
// ---------------------------------------------------------------------------

// Parse parses a single file as described by cfg.
// The tree is closed before returning.
func Parse(cfg ParseConfig) (ParseResult, error) {
	r, err := core.RunParse(core.ParseConfig{
		Path:        cfg.Path,
		Src:         cfg.Src,
		Mode:        int(cfg.Mode),
		Format:      int(cfg.Format),
		GrammarPath: cfg.GrammarPath,
		Benchmark:   cfg.Benchmark,
	})
	if err != nil {
		return ParseResult{}, err
	}
	return convertParseResult(r), nil
}

// ParseFile parses a single file and returns a live FileResult.
// The caller must call Close() on the owning ProjectResult when done.
func ParseFile(cfg ParseConfig) (*FileResult, error) {
	r, err := core.RunParseFile(core.ParseConfig{
		Path:        cfg.Path,
		Src:         cfg.Src,
		Mode:        int(cfg.Mode),
		Format:      int(cfg.Format),
		GrammarPath: cfg.GrammarPath,
		Benchmark:   cfg.Benchmark,
	})
	if err != nil {
		return nil, err
	}
	return convertFileResult(r), nil
}

// ---------------------------------------------------------------------------
// Project-scale parse
// ---------------------------------------------------------------------------

// ParseProject walks the source tree described by cfg and returns a
// ProjectResult with all trees live. Call Close() when done.
func ParseProject(cfg ProjectConfig) (*ProjectResult, error) {
	r, err := core.RunParseProject(core.ProjectConfig{
		Root:        cfg.Root,
		Paths:       cfg.Paths,
		Mode:        int(cfg.Mode),
		GrammarPath: cfg.GrammarPath,
		Workers:     cfg.Workers,
	})
	if err != nil {
		return nil, err
	}
	pr := &ProjectResult{
		Duration:   r.Duration,
		FileCount:  r.FileCount,
		ErrorCount: r.ErrorCount,
		ByteCount:  r.ByteCount,
	}
	for _, f := range r.Files {
		pr.Files = append(pr.Files, convertFileResult(f))
	}
	pr.SetCloseFn(func() { core.CloseProjectResult(r) })
	return pr, nil
}

// ---------------------------------------------------------------------------
// SourcePaths
// ---------------------------------------------------------------------------

// SourcePaths returns the paths of all parseable source files under root that
// match mode, skipping directories named in skip.  This is the entry point
// for build tools that need a filtered file list without parsing.
func SourcePaths(root string, mode Mode, skip []string) ([]string, error) {
	return core.SourcePaths(root, int(mode), skip)
}

// ---------------------------------------------------------------------------
// Walk
// ---------------------------------------------------------------------------

// Walk executes the parallel walk described by cfg.
func Walk(cfg WalkConfig) (WalkSummary, error) {
	var emitFn core.EmitFn
	if cfg.EmitFn != nil {
		fn := cfg.EmitFn
		emitFn = func(path string, r core.ParseResult) {
			fn(path, convertParseResult(r))
		}
	}
	r, err := core.RunWalk(core.WalkConfig{
		Root:        cfg.Root,
		Paths:       cfg.Paths,
		Mode:        int(cfg.Mode),
		GrammarPath: cfg.GrammarPath,
		Workers:     cfg.Workers,
		FailFast:    cfg.FailFast,
		Quiet:       cfg.Quiet,
		Benchmark:   cfg.Benchmark,
		Stdout:      cfg.Stdout,
		Stderr:      cfg.Stderr,
		EmitFn:      emitFn,
	})
	if err != nil {
		return WalkSummary{}, err
	}
	return convertWalkSummary(r), nil
}

// ---------------------------------------------------------------------------
// Debug
// ---------------------------------------------------------------------------

// Debug walks a file or directory, stopping at the first parse error.
func Debug(cfg DebugConfig) (DebugResult, error) {
	r, err := core.RunDebug(core.DebugConfig{
		Path:         cfg.Path,
		Mode:         int(cfg.Mode),
		GrammarPath:  cfg.GrammarPath,
		Profile:      cfg.Profile,
		ContextLines: cfg.ContextLines,
		NoXcode:      cfg.NoXcode,
	})
	if err != nil {
		return DebugResult{}, err
	}
	result := DebugResult{HasError: r.HasError}
	if r.Profile != nil {
		p := &DebugProfile{Total: r.Profile.Total}
		for _, e := range r.Profile.Events {
			p.Events = append(p.Events, DebugEvent{
				Path:    e.Path,
				Elapsed: e.Elapsed,
				OK:      e.OK,
			})
		}
		result.Profile = p
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

func convertParseResult(r core.ParseResult) ParseResult {
	result := ParseResult{
		HasError:   r.HasError,
		NodeCount:  r.NodeCount,
		Sexp:       r.Sexp,
		PrettySexp: r.PrettySexp,
	}
	if r.Timing != nil {
		result.Timing = &Timing{Elapsed: r.Timing.Elapsed, Bytes: r.Timing.Bytes}
	}
	for _, e := range r.Errors {
		result.Errors = append(result.Errors, ParseError{
			Row:     e.Row,
			Column:  e.Column,
			Message: e.Message,
		})
	}
	return result
}

func convertFileResult(r *core.FileResult) *FileResult {
	if r == nil {
		return nil
	}
	fr := &FileResult{
		Path:      r.Path,
		Source:    r.Source,
		Tree:      r.Tree,
		NodeCount: r.NodeCount,
		Duration:  r.Duration,
	}
	for _, e := range r.Errors {
		fr.Errors = append(fr.Errors, ParseError{
			Row:     e.Row,
			Column:  e.Column,
			Message: e.Message,
		})
	}
	return fr
}

func convertWalkSummary(r core.WalkSummary) WalkSummary {
	s := WalkSummary{
		TotalFiles: r.TotalFiles,
		OKFiles:    r.OKFiles,
		ErrorFiles: r.ErrorFiles,
		TotalBytes: r.TotalBytes,
		Elapsed:    r.Elapsed,
	}
	if r.Bench != nil {
		b := &BenchReport{
			WallTime:   r.Bench.WallTime,
			ParseTotal: r.Bench.ParseTotal,
			Files:      r.Bench.Files,
			Bytes:      r.Bench.Bytes,
		}
		for _, e := range r.Bench.Events {
			b.Events = append(b.Events, BenchEvent{
				Path:    e.Path,
				Elapsed: e.Elapsed,
				Bytes:   e.Bytes,
			})
		}
		s.Bench = b
	}
	return s
}
