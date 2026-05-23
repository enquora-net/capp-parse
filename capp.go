/*
 * capp.go
 * capp-parse
 *
 * Public facade for the capp-parse library.
 * Import path: github.com/enquora-net/capp-parse
 *
 * Type aliases re-export all public types from internal/types.
 * Function wrappers delegate to internal/core.
 * No CGo, no go-tree-sitter imports at this level.
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
	"github.com/enquora-net/capp-parse/internal/core"
	"github.com/enquora-net/capp-parse/internal/types"
)

// BoundaryVersion identifies the interop contract version at the
// Lisette/Go boundary. Increment when the public API changes in a
// way that affects generated Go code.
const BoundaryVersion = 1

// ---------------------------------------------------------------------------
// Type aliases — re-export all public types
// ---------------------------------------------------------------------------

type Mode = types.Mode
type Format = types.Format
type ParseError = types.ParseError
type Timing = types.Timing
type ParseConfig = types.ParseConfig
type ParseResult = types.ParseResult
type FileResult = types.FileResult
type ProjectConfig = types.ProjectConfig
type ProjectResult = types.ProjectResult
type EmitFn = types.EmitFn
type WalkConfig = types.WalkConfig
type WalkSummary = types.WalkSummary
type BenchReport = types.BenchReport
type BenchEvent = types.BenchEvent
type DebugConfig = types.DebugConfig
type DebugResult = types.DebugResult
type DebugProfile = types.DebugProfile
type DebugEvent = types.DebugEvent

// ---------------------------------------------------------------------------
// Mode constants
// ---------------------------------------------------------------------------

const (
	ModeAuto = types.ModeAuto
	ModeObjj = types.ModeObjj
	ModeJS   = types.ModeJS
	ModeBoth = types.ModeBoth
)

// ---------------------------------------------------------------------------
// Format constants
// ---------------------------------------------------------------------------

const (
	FormatDefault = types.FormatDefault
	FormatSexp    = types.FormatSexp
	FormatJSON    = types.FormatJSON
	FormatAST     = types.FormatAST
)

// ---------------------------------------------------------------------------
// Mode and format parsing
// ---------------------------------------------------------------------------

// ParseMode converts a flag string to a Mode.
func ParseMode(s string) (Mode, error) { return types.ParseMode(s) }

// ParseFormat converts a flag string to a Format.
func ParseFormat(s string) (Format, error) { return types.ParseFormat(s) }

// ---------------------------------------------------------------------------
// Single-file parse
// ---------------------------------------------------------------------------

// Parse parses a single file as described by cfg and returns a ParseResult
// suitable for CLI output. The tree is closed before returning.
func Parse(cfg ParseConfig) (ParseResult, error) { return core.RunParse(cfg) }

// ParseFile parses a single file and returns a live FileResult.
// The caller must call Close() on the owning ProjectResult when done.
func ParseFile(cfg ParseConfig) (*FileResult, error) { return core.RunParseFile(cfg) }

// ---------------------------------------------------------------------------
// Project-scale parse
// ---------------------------------------------------------------------------

// ParseProject walks the source tree described by cfg, parses every
// matching file, and returns a ProjectResult with all trees live.
// The caller must call ProjectResult.Close() when done.
func ParseProject(cfg ProjectConfig) (*ProjectResult, error) { return core.RunParseProject(cfg) }

// ---------------------------------------------------------------------------
// Walk — streaming project walk for CLI use
// ---------------------------------------------------------------------------

// Walk executes the parallel walk described by cfg.
func Walk(cfg WalkConfig) (WalkSummary, error) { return core.RunWalk(cfg) }

// ---------------------------------------------------------------------------
// Debug
// ---------------------------------------------------------------------------

// Debug walks a file or directory, stopping at the first parse error.
func Debug(cfg DebugConfig) (DebugResult, error) { return core.RunDebug(cfg) }
