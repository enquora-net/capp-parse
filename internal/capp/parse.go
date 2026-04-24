/*
 * internal/capp/parse.go
 * capp-parse
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package capp

import (
	"fmt"
	"time"

	"capp-parse/internal/grammar"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

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
