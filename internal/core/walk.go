/*
 * internal/core/walk.go
 * capp-parse
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 */
package core

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/enquora-net/capp-parse/grammar"
	"github.com/enquora-net/capp-parse/internal/types"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ---------------------------------------------------------------------------
// Project-scale parse — compiler path
// ---------------------------------------------------------------------------

// RunParseProject walks the source tree, parses every matching file, and
// returns a ProjectResult with all trees live. Call Close() when done.
func RunParseProject(cfg types.ProjectConfig) (*types.ProjectResult, error) {
	workers := cfg.Workers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}

	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return nil, err
	}

	paths, err := resolvePaths(cfg.Root, cfg.Paths, cfg.Mode)
	if err != nil {
		return nil, err
	}

	type rawResult struct {
		path    string
		src     []byte
		tree    *sitter.Tree
		errors  []types.ParseError
		nodes   int
		elapsed time.Duration
		readErr error
	}

	work := make(chan string, workers*4)
	results := make(chan rawResult, workers*4)

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parser := sitter.NewParser()
			defer parser.Close()
			if err := parser.SetLanguage(lang); err != nil {
				return
			}
			for path := range work {
				src, readErr := os.ReadFile(path)
				if readErr != nil {
					results <- rawResult{path: path, readErr: readErr}
					continue
				}
				t0 := time.Now()
				tree := parser.Parse(src, nil)
				elapsed := time.Since(t0)
				root := tree.RootNode()
				r := rawResult{
					path:    path,
					src:     src,
					tree:    tree,
					nodes:   countNodes(root),
					elapsed: elapsed,
				}
				if root.HasError() {
					collectErrors(root, src, &r.errors)
				}
				results <- r
			}
		}()
	}

	go func() {
		for _, p := range paths {
			work <- p
		}
		close(work)
		wg.Wait()
		close(results)
	}()

	pr := &types.ProjectResult{}
	wallStart := time.Now()

	for r := range results {
		pr.FileCount++
		if r.readErr != nil {
			pr.ErrorCount++
			pr.Files = append(pr.Files, &types.FileResult{
				Path: r.path,
				Errors: []types.ParseError{{
					Message: fmt.Sprintf("read error: %v", r.readErr),
				}},
			})
			continue
		}
		pr.ByteCount += int64(len(r.src))
		if len(r.errors) > 0 {
			pr.ErrorCount++
		}
		pr.Files = append(pr.Files, &types.FileResult{
			Path:      r.path,
			Source:    r.src,
			Tree:      r.tree,
			Errors:    r.errors,
			NodeCount: r.nodes,
			Duration:  r.elapsed,
		})
	}

	pr.Duration = time.Since(wallStart)
	pr.SetCloseFn(func() { closeProjectResult(pr) })
	return pr, nil
}

// closeProjectResult releases all live trees in a ProjectResult.
func closeProjectResult(pr *types.ProjectResult) {
	for i := range pr.Files {
		if pr.Files[i].Tree != nil {
			if tree, ok := pr.Files[i].Tree.(*sitter.Tree); ok {
				tree.Close()
			}
			pr.Files[i].Tree = nil
		}
	}
}

// ---------------------------------------------------------------------------
// Streaming walk — CLI path
// ---------------------------------------------------------------------------

// RunWalk executes the parallel walk described by cfg.
func RunWalk(cfg types.WalkConfig) (types.WalkSummary, error) {
	stdout := writerOrStdout(cfg.Stdout)
	stderr := writerOrStderr(cfg.Stderr)

	workers := cfg.Workers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}

	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return types.WalkSummary{}, err
	}

	paths, err := resolvePaths(cfg.Root, cfg.Paths, cfg.Mode)
	if err != nil {
		return types.WalkSummary{}, err
	}

	type result struct {
		path       string
		elapsed    time.Duration
		bytes      int
		hasErr     bool
		errors     []types.ParseError
		nodeCount  int
		sexp       string
		prettySexp string
		readErr    error
	}

	work := make(chan string, workers*4)
	results := make(chan result, workers*4)

	var stopped atomic.Bool
	wallStart := time.Now()

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parser := sitter.NewParser()
			defer parser.Close()
			if err := parser.SetLanguage(lang); err != nil {
				return
			}
			for path := range work {
				if stopped.Load() {
					return
				}
				src, readErr := os.ReadFile(path)
				if readErr != nil {
					results <- result{path: path, readErr: readErr}
					if cfg.FailFast {
						stopped.Store(true)
					}
					continue
				}
				t0 := time.Now()
				tree := parser.Parse(src, nil)
				elapsed := time.Since(t0)
				root := tree.RootNode()
				r := result{
					path:       path,
					elapsed:    elapsed,
					bytes:      len(src),
					hasErr:     root.HasError(),
					nodeCount:  countNodes(root),
					sexp:       root.ToSexp(),
					prettySexp: FormatNode(root, src, 0),
				}
				if r.hasErr {
					collectErrors(root, src, &r.errors)
					if cfg.FailFast {
						stopped.Store(true)
					}
				}
				tree.Close()
				results <- r
			}
		}()
	}

	go func() {
		for _, p := range paths {
			if stopped.Load() {
				break
			}
			work <- p
		}
		close(work)
		wg.Wait()
		close(results)
	}()

	summary := types.WalkSummary{}
	bench := newBenchAccumulator(cfg.Benchmark)

	for r := range results {
		summary.TotalFiles++
		if r.readErr != nil {
			summary.ErrorFiles++
			fmt.Fprintf(stderr, "read error %s: %v\n", r.path, r.readErr)
			continue
		}
		summary.TotalBytes += r.bytes
		bench.record(r.path, r.elapsed, r.bytes)

		pr := types.ParseResult{
			HasError:   r.hasErr,
			Errors:     r.errors,
			NodeCount:  r.nodeCount,
			Sexp:       r.sexp,
			PrettySexp: r.prettySexp,
		}

		if cfg.EmitFn != nil {
			cfg.EmitFn(r.path, pr)
		} else if !cfg.Quiet {
			if r.hasErr {
				fmt.Fprintf(stderr, "%s: syntax errors\n", r.path)
			} else {
				fmt.Fprintf(stdout, "%s: OK\n", r.path)
			}
		}

		if r.hasErr {
			summary.ErrorFiles++
		} else {
			summary.OKFiles++
		}
	}

	summary.Elapsed = time.Since(wallStart)
	summary.Bench = bench.report(summary.Elapsed)
	return summary, nil
}

// ---------------------------------------------------------------------------
// Shared utilities
// ---------------------------------------------------------------------------

func resolvePaths(root string, explicit []string, mode types.Mode) ([]string, error) {
	if len(explicit) > 0 {
		return explicit, nil
	}
	return collectPaths(root, mode)
}

func collectPaths(root string, mode types.Mode) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && accept(mode, path) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", root, err)
	}
	return paths, nil
}

func writerOrStdout(w interface{ Write([]byte) (int, error) }) io.Writer {
	if w != nil {
		return w.(io.Writer)
	}
	return os.Stdout
}

func writerOrStderr(w interface{ Write([]byte) (int, error) }) io.Writer {
	if w != nil {
		return w.(io.Writer)
	}
	return os.Stderr
}
