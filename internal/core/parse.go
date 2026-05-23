/*
 * internal/core/parse.go
 * capp-parse
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 */
package core

import (
	"fmt"
	"time"

	"github.com/enquora-net/capp-parse/grammar"
	"github.com/enquora-net/capp-parse/internal/types"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// RunParseFile parses a single file and returns a live FileResult.
// The tree is owned by the caller; release via the owning ProjectResult.Close().
func RunParseFile(cfg types.ParseConfig) (*types.FileResult, error) {
	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return nil, err
	}

	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return nil, fmt.Errorf("setting language: %w", err)
	}

	src := cfg.Src
	if src == nil {
		return nil, fmt.Errorf("ParseFile: Src must not be nil; read the file before calling")
	}

	t0 := time.Now()
	tree := parser.Parse(src, nil)
	elapsed := time.Since(t0)

	root := tree.RootNode()
	fr := &types.FileResult{
		Path:      cfg.Path,
		Source:    src,
		Tree:      tree,
		NodeCount: countNodes(root),
		Duration:  elapsed,
	}

	if root.HasError() {
		collectErrors(root, src, &fr.Errors)
	}

	return fr, nil
}

// RunParse parses a single file for CLI consumption.
// The tree is closed before returning; use RunParseFile for compiler use.
func RunParse(cfg types.ParseConfig) (types.ParseResult, error) {
	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return types.ParseResult{}, err
	}

	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return types.ParseResult{}, fmt.Errorf("setting language: %w", err)
	}

	t0 := time.Now()
	tree := parser.Parse(cfg.Src, nil)
	elapsed := time.Since(t0)
	defer tree.Close()

	root := tree.RootNode()
	result := types.ParseResult{
		HasError:  root.HasError(),
		NodeCount: countNodes(root),
	}

	if cfg.Benchmark {
		result.Timing = &types.Timing{Elapsed: elapsed, Bytes: len(cfg.Src)}
	}

	if result.HasError {
		collectErrors(root, cfg.Src, &result.Errors)
	}

	if cfg.Format == types.FormatSexp || cfg.Format == types.FormatDefault {
		result.Sexp = root.ToSexp()
		result.PrettySexp = FormatNode(root, cfg.Src, 0)
	}

	return result, nil
}
