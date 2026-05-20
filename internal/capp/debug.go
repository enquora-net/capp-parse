/*
 * internal/capp/debug.go
 * capp-parse
 *
 * Created by David Richardson on Thursday, April 23, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package capp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"dev.cappuccino/capp-parse/internal/grammar"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

const (
	divider = "──────────────────────────────────────────────────────────────────────"
	header  = "══════════════════════════════════════════════════════════════════════"

	ansiGreen = "\033[32m"
	ansiRed   = "\033[31m"
	ansiReset = "\033[0m"

	iconOK  = ansiGreen + "✓" + ansiReset
	iconErr = ansiRed + "✗" + ansiReset
)

func runDebug(cfg DebugConfig) (DebugResult, error) {
	lang, err := grammar.Language(cfg.GrammarPath)
	if err != nil {
		return DebugResult{}, err
	}

	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return DebugResult{}, fmt.Errorf("setting language: %w", err)
	}

	contextLines := cfg.ContextLines
	if contextLines <= 0 {
		contextLines = 3
	}

	info, err := os.Stat(cfg.Path)
	if err != nil {
		return DebugResult{}, fmt.Errorf("accessing %s: %w", cfg.Path, err)
	}

	var events []DebugEvent
	var totalElapsed time.Duration
	hasError := false

	if info.IsDir() {
		paths, err := collectPaths(cfg.Path, cfg.Mode)
		if err != nil {
			return DebugResult{}, err
		}
		for _, path := range paths {
			ok, elapsed := debugFile(path, parser, cfg.Path, contextLines, cfg.NoXcode)
			totalElapsed += elapsed
			if cfg.Profile {
				events = append(events, DebugEvent{Path: path, Elapsed: elapsed, OK: ok})
			}
			if !ok {
				hasError = true
				break
			}
		}
	} else {
		ok, elapsed := debugFile(cfg.Path, parser, "", contextLines, cfg.NoXcode)
		totalElapsed += elapsed
		if cfg.Profile {
			events = append(events, DebugEvent{Path: cfg.Path, Elapsed: elapsed, OK: ok})
		}
		if !ok {
			hasError = true
		}
	}

	result := DebugResult{HasError: hasError}
	if cfg.Profile {
		result.Profile = &DebugProfile{
			Events: events,
			Total:  totalElapsed,
		}
	}
	return result, nil
}

func debugFile(path string, parser *sitter.Parser, root string, contextLines int, noXcode bool) (ok bool, elapsed time.Duration) {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\nEXCEPTION in %s\n%s\n%v\n", header, path, header, err)
		return false, 0
	}

	t0 := time.Now()
	tree := parser.Parse(src, nil)
	elapsed = time.Since(t0)
	defer tree.Close()

	root2 := tree.RootNode()
	lines := strings.Split(string(src), "\n")

	errorNode, chain := FindFirstError(root2, nil)
	if errorNode == nil {
		displayPath := path
		if root != "" {
			if rel, err := filepath.Rel(root, path); err == nil {
				displayPath = rel
			}
		}
		fmt.Printf("%s  %s\n", iconOK, displayPath)
		return true, elapsed
	}

	sp := errorNode.StartPosition()

	if !noXcode {
		openInXcode(path, sp.Row)
	}

	fmt.Printf("\n%s\nPARSE ERROR in %s\n%s\n", header, path, header)
	fmt.Printf("First error at line %d, column %d\n", sp.Row+1, sp.Column+1)
	fmt.Printf("Total lines: %d\n", len(lines))

	showSourceContext(errorNode, lines, contextLines)
	showParentChain(chain)

	fmt.Println("\n" + divider)
	fmt.Println("FULL CONCRETE SYNTAX TREE")
	fmt.Println(divider)
	fmt.Println(FormatNode(root2, src, 0))

	return false, elapsed
}

func showSourceContext(errorNode *sitter.Node, lines []string, contextLines int) {
	startLine := int(errorNode.StartPosition().Row)
	endLine := int(errorNode.EndPosition().Row)
	first := max(0, startLine-contextLines)
	last := min(len(lines)-1, endLine+contextLines)

	fmt.Println("\n" + divider)
	fmt.Println("SOURCE CONTEXT")
	fmt.Println(divider)

	for i := first; i <= last; i++ {
		lineText := strings.TrimRight(lines[i], "\r")
		lineNum := i + 1
		if startLine <= i && i <= endLine {
			fmt.Printf(">>> %4d | %s\n", lineNum, lineText)
			if startLine == endLine {
				col := int(errorNode.StartPosition().Column)
				width := max(1, int(errorNode.EndPosition().Column)-col)
				fmt.Printf("%s%s\n",
					strings.Repeat(" ", 3+1+4+3+col),
					strings.Repeat("^", width))
			}
		} else {
			fmt.Printf("    %4d | %s\n", lineNum, lineText)
		}
	}
}

func showParentChain(chain []*sitter.Node) {
	if len(chain) == 0 {
		return
	}
	fmt.Println("\n" + divider)
	fmt.Println("PARENT CONTEXT")
	fmt.Println(divider)
	for i, node := range chain {
		sp := node.StartPosition()
		fmt.Printf("%s↓ %s at %d:%d\n",
			strings.Repeat("  ", i), node.Kind(), sp.Row+1, sp.Column+1)
	}
}

func openInXcode(path string, row uint) {
	cmd := exec.Command("xed", "--line", fmt.Sprintf("%d", row+1), path)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "xed: %v\n", err)
	}
}
