/*
 * internal/capp/tree.go
 * capp-parse
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package capp

import (
	"fmt"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

func countNodes(n *gotreesitter.Node) int {
	count := 1
	for i := 0; i < n.ChildCount(); i++ {
		count += countNodes(n.Child(i))
	}
	return count
}

func collectErrors(n *gotreesitter.Node, src []byte, lang *gotreesitter.Language, errs *[]string) {
	if n.IsError() || n.IsMissing() {
		start := n.StartPoint()
		*errs = append(*errs, fmt.Sprintf("line %d col %d: %s",
			start.Row+1, start.Column+1, errorContext(n, src, lang)))
	}
	for i := 0; i < n.ChildCount(); i++ {
		collectErrors(n.Child(i), src, lang, errs)
	}
}

func errorContext(n *gotreesitter.Node, src []byte, lang *gotreesitter.Language) string {
	start := n.StartByte()
	end := n.EndByte()
	srcLen := uint32(len(src))
	if end > srcLen {
		end = srcLen
	}
	snippet := src[start:end]
	if len(snippet) > 40 {
		snippet = snippet[:40]
	}
	if n.IsMissing() {
		return fmt.Sprintf("missing %s", n.Type(lang))
	}
	return fmt.Sprintf("ERROR %q", snippet)
}
