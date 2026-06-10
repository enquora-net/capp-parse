/*
 * node.go
 * capp-parse
 *
 * CSTNode — public concrete syntax tree node for compiler and tooling consumers.
 *
 * All tree-sitter interactions are confined to internal/core.  This file
 * imports only internal/core, keeping the root package free of any sitter
 * dependency and preserving bindgen's ability to type-check the facade.
 *
 * Created by David Richardson on Tuesday, May 26, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package capp

import "github.com/enquora-net/capp-parse/internal/core"

// Grammar node-kind and field-name constants live in kinds.go.

// ---------------------------------------------------------------------------
// CSTNode
// ---------------------------------------------------------------------------

// CSTNode is a node in the concrete syntax tree produced by capp-parse.
// The inner field holds a *sitter.Node at runtime but is typed as interface{}
// so the root package never references sitter directly.  All sitter
// interactions are delegated to internal/core.
type CSTNode struct {
	inner  interface{}
	source []byte
}

func (n CSTNode) IsNull() bool { return n.inner == nil }

func (n CSTNode) Kind() string { return core.NodeKind(n.inner) }

func (n CSTNode) Text() string { return core.NodeText(n.inner, n.source) }

func (n CSTNode) ChildCount() int { return core.NodeChildCount(n.inner) }

func (n CSTNode) Child(i int) CSTNode {
	return CSTNode{inner: core.NodeChild(n.inner, i), source: n.source}
}

// NamedChildCount returns the number of named children of n, excluding
// anonymous tokens (punctuation, keywords).  Extras such as comments are
// named nodes and ARE counted; consumers must still skip them by kind.
func (n CSTNode) NamedChildCount() int { return core.NodeNamedChildCount(n.inner) }

// NamedChild returns the i-th named child of n.  Iterating named children
// skips punctuation and keyword tokens but not comments, which are named
// extras; consumers reconstructing structure skip those by kind.
func (n CSTNode) NamedChild(i int) CSTNode {
	return CSTNode{inner: core.NodeNamedChild(n.inner, i), source: n.source}
}

// IsNamed reports whether n is a named node (a grammar rule) rather than an
// anonymous token.
func (n CSTNode) IsNamed() bool { return core.NodeIsNamed(n.inner) }

func (n CSTNode) ChildByFieldName(name string) (CSTNode, bool) {
	child, ok := core.NodeChildByFieldName(n.inner, name)
	if !ok {
		return CSTNode{source: n.source}, false
	}
	return CSTNode{inner: child, source: n.source}, true
}

func (n CSTNode) HasError() bool    { return core.NodeHasError(n.inner) }
func (n CSTNode) IsError() bool     { return core.NodeIsError(n.inner) }
func (n CSTNode) IsMissing() bool   { return core.NodeIsMissing(n.inner) }
func (n CSTNode) StartRow() uint32  { return core.NodeStartRow(n.inner) }
func (n CSTNode) StartColumn() uint32 { return core.NodeStartColumn(n.inner) }

// ---------------------------------------------------------------------------
// FileResult.Root
// ---------------------------------------------------------------------------

// Root returns the root CSTNode of this file's syntax tree.
// The CSTNode is valid only while the owning FileResult or ProjectResult
// remains live and has not been closed.
func (f *FileResult) Root() CSTNode {
	if f == nil || f.Tree == nil {
		return CSTNode{}
	}
	return CSTNode{inner: core.FileRootNode(f.Tree), source: f.Source}
}
