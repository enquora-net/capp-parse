/*
 * internal/core/node.go
 * capp-parse
 *
 * Node accessor functions accepting interface{} so the root package never
 * imports go-tree-sitter directly.  Bindgen processes the root package;
 * keeping sitter confined to internal/core preserves the CGo-free facade.
 *
 * Created by David Richardson on Tuesday, May 26, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package core

import sitter "github.com/tree-sitter/go-tree-sitter"

// NodeKind returns the grammar rule name of n, or "" if n is nil.
func NodeKind(n interface{}) string {
	if n == nil {
		return ""
	}
	return n.(*sitter.Node).Kind()
}

// NodeText returns the source text spanned by n.
func NodeText(n interface{}, source []byte) string {
	if n == nil {
		return ""
	}
	node := n.(*sitter.Node)
	start := node.StartByte()
	end := node.EndByte()
	if end > uint(len(source)) {
		end = uint(len(source))
	}
	return string(source[start:end])
}

// NodeChildCount returns the total child count of n, or 0 if n is nil.
func NodeChildCount(n interface{}) int {
	if n == nil {
		return 0
	}
	return int(n.(*sitter.Node).ChildCount())
}

// NodeChild returns the i-th child of n as interface{}, or nil if out of range.
func NodeChild(n interface{}, i int) interface{} {
	if n == nil {
		return nil
	}
	node := n.(*sitter.Node)
	if i < 0 || uint(i) >= node.ChildCount() {
		return nil
	}
	return node.Child(uint(i))
}

// NodeChildByFieldName returns the child with the given field name, or nil.
func NodeChildByFieldName(n interface{}, name string) (interface{}, bool) {
	if n == nil {
		return nil, false
	}
	child := n.(*sitter.Node).ChildByFieldName(name)
	if child == nil {
		return nil, false
	}
	return child, true
}

// NodeHasError reports whether n or any descendant contains a parse error.
func NodeHasError(n interface{}) bool {
	if n == nil {
		return false
	}
	return n.(*sitter.Node).HasError()
}

// NodeIsError reports whether n is itself an ERROR node.
func NodeIsError(n interface{}) bool {
	if n == nil {
		return false
	}
	return n.(*sitter.Node).IsError()
}

// NodeIsMissing reports whether n is a MISSING node.
func NodeIsMissing(n interface{}) bool {
	if n == nil {
		return false
	}
	return n.(*sitter.Node).IsMissing()
}

// NodeStartRow returns the zero-based source row of n.
func NodeStartRow(n interface{}) uint32 {
	if n == nil {
		return 0
	}
	return uint32(n.(*sitter.Node).StartPosition().Row)
}

// NodeStartColumn returns the zero-based source column of n.
func NodeStartColumn(n interface{}) uint32 {
	if n == nil {
		return 0
	}
	return uint32(n.(*sitter.Node).StartPosition().Column)
}

// FileRootNode returns the root *sitter.Node of a *sitter.Tree held as
// interface{} in FileResult.Tree, alongside the source bytes.
func FileRootNode(tree interface{}) interface{} {
	if tree == nil {
		return nil
	}
	return tree.(*sitter.Tree).RootNode()
}
