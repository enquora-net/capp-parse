/*
 * cmd/parse/parse.go
 * cappuccino
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
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

	"github.com/enquora-net/capp-parse/capp"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewParseCmd constructs the parse command.
func NewParseCmd() *cobra.Command {
	var (
		grammar string
		mode    string
		format  string
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

 --format applies to single-file and --src modes; ignored for directory targets.

 Output is human-readable when stdout is a terminal; JSON per line otherwise.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := capp.ParseMode(mode)
			if err != nil {
				return err
			}

			f, err := capp.ParseFormat(format)
			if err != nil {
				return err
			}

			isTTY := term.IsTerminal(int(os.Stdout.Fd()))

			if src {
				if len(args) > 0 {
					return fmt.Errorf("--src reads from stdin; no path argument expected")
				}
				return runStdinSource(grammar, m, f, isTTY)
			}

			if len(args) == 1 {
				return runPath(args[0], grammar, m, f, workers, isTTY)
			}

			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return runStdinPaths(grammar, m, workers, isTTY)
			}

			return cmd.Usage()
		},
	}

	cmd.Flags().StringVarP(&grammar, "grammar", "g", "", "explicit path to grammar dylib, overrides search")
	cmd.Flags().StringVar(&mode, "mode", "auto", "language filter: auto | objj | js | both")
	cmd.Flags().StringVarP(&format, "format", "f", "", "output format for single-file mode: sexp | json | ast (default: sexp)")
	cmd.Flags().IntVar(&workers, "workers", 0, "parallel workers for tree and path-stream modes (0 = GOMAXPROCS)")
	cmd.Flags().BoolVar(&src, "src", false, "read source bytes from stdin rather than file paths")

	return cmd
}

func runPath(path, grammar string, m capp.Mode, f capp.Format, workers int, isTTY bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("accessing %s: %w", path, err)
	}

	if info.IsDir() {
		return runTree(path, grammar, m, workers, isTTY)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	result, err := capp.Parse(capp.ParseConfig{
		Path:        path,
		Src:         src,
		Mode:        m,
		Format:      f,
		GrammarPath: grammar,
	})
	if err != nil {
		return err
	}

	capp.EmitResult(path, result, f, isTTY)
	if result.HasError {
		os.Exit(2)
	}
	return nil
}

func runTree(root, grammar string, m capp.Mode, workers int, isTTY bool) error {
	summary, err := capp.Walk(capp.WalkConfig{
		Root:        root,
		Mode:        m,
		GrammarPath: grammar,
		Workers:     workers,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Benchmark:   isTTY,
		EmitFn:      capp.MakeEmitFn(isTTY),
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

func runStdinPaths(grammar string, m capp.Mode, workers int, isTTY bool) error {
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
		Paths:       paths,
		Mode:        m,
		GrammarPath: grammar,
		Workers:     workers,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Benchmark:   isTTY,
		EmitFn:      capp.MakeEmitFn(isTTY),
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

func runStdinSource(grammar string, m capp.Mode, f capp.Format, isTTY bool) error {
	src, err := os.ReadFile("/dev/stdin")
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	result, err := capp.Parse(capp.ParseConfig{
		Path:        "<stdin>",
		Src:         src,
		Mode:        m,
		Format:      f,
		GrammarPath: grammar,
	})
	if err != nil {
		return err
	}

	capp.EmitResult("<stdin>", result, f, isTTY)
	if result.HasError {
		os.Exit(2)
	}
	return nil
}
