/*
 * capp/tree.go
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
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

func countNodes(n *sitter.Node) int {
	count := 1
	for i := range n.ChildCount() {
		count += countNodes(n.Child(i))
	}
	return count
}

func collectErrors(n *sitter.Node, src []byte, errs *[]ParseError) {
	if n.IsError() || n.IsMissing() {
		start := n.StartPosition()
		*errs = append(*errs, ParseError{
			Row:     uint32(start.Row),
			Column:  uint32(start.Column),
			Message: errorContext(n, src),
		})
	}
	for i := range n.ChildCount() {
		collectErrors(n.Child(i), src, errs)
	}
}

func errorContext(n *sitter.Node, src []byte) string {
	start := n.StartByte()
	end := n.EndByte()
	srcLen := uint(len(src))
	if end > srcLen {
		end = srcLen
	}
	snippet := src[start:end]
	if len(snippet) > 40 {
		snippet = snippet[:40]
	}
	if n.IsMissing() {
		return fmt.Sprintf("missing %s", n.Kind())
	}
	return fmt.Sprintf("ERROR %q", snippet)
}

// FormatNode renders a node and its descendants as an indented tree.
// Leaf node text is shown truncated to maxTextLen runes.
const maxTextLen = 50

// FormatNode renders a parse tree node as an indented, human-readable tree.
func FormatNode(node *sitter.Node, source []byte, indent int) string {
	var sb strings.Builder
	prefix := strings.Repeat("  ", indent)
	kind := node.Kind()
	isError := kind == "ERROR"
	isMissing := node.IsMissing()

	label := kind
	if isError {
		sp := node.StartPosition()
		ep := node.EndPosition()
		label = fmt.Sprintf("❌ ERROR [%d:%d-%d:%d]",
			sp.Row+1, sp.Column+1, ep.Row+1, ep.Column+1)
	} else if isMissing {
		sp := node.StartPosition()
		label = fmt.Sprintf("⚠️ MISSING %s [%d:%d]", kind, sp.Row+1, sp.Column+1)
	}

	childCount := node.ChildCount()

	if childCount == 0 {
		if !isError && !isMissing {
			text := string(source[node.StartByte():node.EndByte()])
			text = strings.ReplaceAll(text, "\n", `\n`)
			text = strings.ReplaceAll(text, "\t", `\t`)
			if len(text) > maxTextLen {
				text = text[:maxTextLen-3] + "..."
			}
			fmt.Fprintf(&sb, "%s(%s %q)", prefix, label, text)
		} else {
			fmt.Fprintf(&sb, "%s(%s)", prefix, label)
		}
		return sb.String()
	}

	fmt.Fprintf(&sb, "%s(%s\n", prefix, label)
	for i := range childCount {
		child := node.Child(i)
		fmt.Fprintf(&sb, "%s\n", FormatNode(child, source, indent+1))
	}
	fmt.Fprintf(&sb, "%s)", prefix)
	return sb.String()
}

// FindFirstError performs a depth-first search for the first ERROR or MISSING node.
// Returns the node and the ancestor chain from root to its parent.
func FindFirstError(node *sitter.Node, ancestors []*sitter.Node) (*sitter.Node, []*sitter.Node) {
	if node.Kind() == "ERROR" || node.IsMissing() {
		return node, ancestors
	}
	for i := range node.ChildCount() {
		child := node.Child(i)
		if found, chain := FindFirstError(child, append(ancestors, node)); found != nil {
			return found, chain
		}
	}
	return nil, nil
}
