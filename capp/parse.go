/*
 * capp/parse.go
 * cappuccino
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package capp

import (
	"fmt"
	"time"

	"dev.cappuccino/capp-parse/grammar"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// runParseFile parses a single file and returns a live FileResult.
// The tree is owned by the caller; close via ProjectResult.Close()
// or directly on FileResult.Tree.(*sitter.Tree).
func runParseFile(cfg ParseConfig) (*FileResult, error) {
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
	fr := &FileResult{
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

// runParse parses a single file for CLI consumption.
// The tree is closed before returning; use runParseFile for compiler use.
func runParse(cfg ParseConfig) (ParseResult, error) {
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

	if cfg.Format == FormatSexp || cfg.Format == FormatDefault {
		result.Sexp = root.ToSexp()
		result.PrettySexp = FormatNode(root, cfg.Src, 0)
	}

	return result, nil
}
