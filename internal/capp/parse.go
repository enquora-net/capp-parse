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
	gotreesitter "github.com/odvcencio/gotreesitter"
)

func runParse(cfg ParseConfig) (ParseResult, error) {
	lang, err := grammar.Language()
	if err != nil {
		return ParseResult{}, err
	}

	t0 := time.Now()
	parser := gotreesitter.NewParser(lang)
	tree, err := parser.Parse(cfg.Src)
	elapsed := time.Since(t0)

	if err != nil {
		return ParseResult{}, fmt.Errorf("parse %s: %w", cfg.Path, err)
	}

	root := tree.RootNode()
	result := ParseResult{
		HasError:  root.HasError(),
		NodeCount: countNodes(root),
	}

	if cfg.Benchmark {
		result.Timing = &Timing{Elapsed: elapsed, Bytes: len(cfg.Src)}
	}

	if result.HasError {
		collectErrors(root, cfg.Src, lang, &result.Errors)
	}

	if cfg.Sexp {
		result.Sexp = root.SExpr(lang)
	}

	return result, nil
}
