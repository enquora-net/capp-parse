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
	sitter "github.com/tree-sitter/go-tree-sitter"
)

const (
	formatDefault = 0
	formatSexp    = 1
)

// RunParseFile parses a single file and returns a live FileResult.
func RunParseFile(cfg ParseConfig) (*FileResult, error) {
	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return nil, err
	}

	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return nil, fmt.Errorf("setting language: %w", err)
	}

	if cfg.Src == nil {
		return nil, fmt.Errorf("ParseFile: Src must not be nil")
	}

	t0 := time.Now()
	tree := parser.Parse(cfg.Src, nil)
	elapsed := time.Since(t0)

	root := tree.RootNode()
	fr := &FileResult{
		Path:      cfg.Path,
		Source:    cfg.Src,
		Tree:      tree,
		NodeCount: countNodes(root),
		Duration:  elapsed,
	}

	if root.HasError() {
		collectErrors(root, cfg.Src, &fr.Errors)
	}

	return fr, nil
}

// RunParse parses a single file for CLI consumption.
// The tree is closed before returning.
func RunParse(cfg ParseConfig) (ParseResult, error) {
	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return ParseResult{}, err
	}

	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return ParseResult{}, fmt.Errorf("setting language: %w", err)
	}

	t0 := time.Now()
	tree := parser.Parse(cfg.Src, nil)
	elapsed := time.Since(t0)
	defer tree.Close()

	root := tree.RootNode()
	result := ParseResult{
		HasError:  root.HasError(),
		NodeCount: countNodes(root),
	}

	if cfg.Benchmark {
		result.Timing = &Timing{Elapsed: elapsed, Bytes: len(cfg.Src)}
	}

	if result.HasError {
		collectErrors(root, cfg.Src, &result.Errors)
	}

	if cfg.Format == formatSexp || cfg.Format == formatDefault {
		result.Sexp = root.ToSexp()
		result.PrettySexp = FormatNode(root, cfg.Src, 0)
	}

	return result, nil
}
