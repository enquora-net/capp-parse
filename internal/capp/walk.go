/*
 * internal/capp/walk.go
 * capp-parse
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package capp

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

	"dev.cappuccino/capp-parse/internal/grammar"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

func runWalk(cfg WalkConfig) (WalkSummary, error) {
	stdout := writerOrStdout(cfg.Stdout)
	stderr := writerOrStderr(cfg.Stderr)

	workers := cfg.Workers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}

	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return WalkSummary{}, err
	}

	var paths []string
	if len(cfg.Paths) > 0 {
		paths = cfg.Paths
	} else {
		paths, err = collectPaths(cfg.Root, cfg.Mode)
		if err != nil {
			return WalkSummary{}, err
		}
	}

	type result struct {
		path      string
		elapsed   time.Duration
		bytes     int
		hasErr    bool
		errors    []string
		nodeCount int
		readErr   error
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
					path:      path,
					elapsed:   elapsed,
					bytes:     len(src),
					hasErr:    root.HasError(),
					nodeCount: countNodes(root),
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

	summary := WalkSummary{}
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

		pr := ParseResult{
			HasError:  r.hasErr,
			Errors:    r.errors,
			NodeCount: r.nodeCount,
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

func collectPaths(root string, mode Mode) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && mode.accept(path) {
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
