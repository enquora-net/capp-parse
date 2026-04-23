/*
 * cmd/parse/parse.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/* Package parse implements the parse subcommand.
 *
 * This package is the fold-in seam: when capp-parse is absorbed into the
 * parent toolchain binary, NewParseCmd() is wired into the toolchain's own
 * root command without modification. cmd/root.go and cmd/version.go are
 * discarded at that point; this package and internal/ move across unchanged.
 */

package parse

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"capp-parse/internal/capp"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewParseCmd constructs the parse command.
func NewParseCmd() *cobra.Command {
	var (
		mode    string
		workers int
		src     bool
	)

	cmd := &cobra.Command{
		Use:   "parse [path]",
		Short: "Parse Objective-J and JavaScript source",
		Long: `Parse a file, a directory tree, or a stream of file paths from stdin.

 Scope is determined from the argument:
   capp-parse parse src/AppController.j      single file
   capp-parse parse Frameworks/              directory tree
   find . -name '*.j' | capp-parse parse     path stream from stdin

 To parse source bytes directly from stdin:
   cat src/AppController.j | capp-parse parse --src --mode objj

 Output is human-readable with performance data when stdout is a terminal.
 Output is JSON per line when stdout is not a terminal.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := capp.ParseMode(mode)
			if err != nil {
				return err
			}

			isTTY := term.IsTerminal(int(os.Stdout.Fd()))

			// stdin source mode
			if src {
				if len(args) > 0 {
					return fmt.Errorf("--src reads from stdin; no path argument expected")
				}
				return runStdinSource(m, isTTY)
			}

			// path argument mode
			if len(args) == 1 {
				return runPath(args[0], m, workers, isTTY)
			}

			// stdin path-stream mode
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return runStdinPaths(m, workers, isTTY)
			}

			return cmd.Usage()
		},
	}

	cmd.Flags().StringVar(&mode, "mode", "auto", "language filter: auto | objj | js | both")
	cmd.Flags().IntVar(&workers, "workers", 0, "parallel workers for tree and path-stream modes (0 = GOMAXPROCS)")
	cmd.Flags().BoolVar(&src, "src", false, "read source bytes from stdin rather than file paths")

	return cmd
}

func runPath(path string, m capp.Mode, workers int, isTTY bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("accessing %s: %w", path, err)
	}

	if info.IsDir() {
		return runTree(path, m, workers, isTTY)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	result, err := capp.Parse(capp.ParseConfig{
		Path: path,
		Src:  src,
		Mode: m,
	})
	if err != nil {
		return err
	}

	capp.EmitResult(path, result, isTTY)
	if result.HasError {
		os.Exit(2)
	}
	return nil
}

func runTree(root string, m capp.Mode, workers int, isTTY bool) error {
	summary, err := capp.Walk(capp.WalkConfig{
		Root:      root,
		Mode:      m,
		Workers:   workers,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Benchmark: isTTY,
		EmitFn:    capp.MakeEmitFn(isTTY),
	})
	if err != nil {
		return err
	}

	if isTTY {
		capp.EmitWalkSummary(summary)
	}

	if summary.ErrorFiles > 0 {
		os.Exit(2)
	}
	return nil
}

func runStdinPaths(m capp.Mode, workers int, isTTY bool) error {
	var paths []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			paths = append(paths, line)
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	summary, err := capp.Walk(capp.WalkConfig{
		Paths:     paths,
		Mode:      m,
		Workers:   workers,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Benchmark: isTTY,
		EmitFn:    capp.MakeEmitFn(isTTY),
	})
	if err != nil {
		return err
	}

	if isTTY {
		capp.EmitWalkSummary(summary)
	}

	if summary.ErrorFiles > 0 {
		os.Exit(2)
	}
	return nil
}

func runStdinSource(m capp.Mode, isTTY bool) error {
	src, err := os.ReadFile("/dev/stdin")
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	result, err := capp.Parse(capp.ParseConfig{
		Path: "<stdin>",
		Src:  src,
		Mode: m,
	})
	if err != nil {
		return err
	}

	capp.EmitResult("<stdin>", result, isTTY)
	if result.HasError {
		os.Exit(2)
	}
	return nil
}
