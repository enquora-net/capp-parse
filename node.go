/*
 * node.go
 * capp-parse
 *
 * CSTNode — public concrete syntax tree node for compiler and tooling consumers.
 *
 * Wraps *sitter.Node and retains the source bytes so Text() requires no
 * caller-supplied buffer; the CST is equivalent to the source text and carries
 * it implicitly.  All method signatures return representable Go types so that
 * Lisette's bindgen can consume this API without encountering unrepresentable
 * internal types.
 *
 * FileResult.Root() is defined here because it produces a CSTNode; the
 * method is on a type declared in types.go, which is valid within the same
 * package.
 *
 * Grammar node-kind and field-name constants are co-located here so that
 * consumers pattern-match against named identifiers rather than raw string
 * literals, remaining resilient to grammar evolution.
 *
 * Created by David Richardson on Tuesday, May 26, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package capp

import sitter "github.com/tree-sitter/go-tree-sitter"

// ---------------------------------------------------------------------------
// Grammar node-kind constants
// ---------------------------------------------------------------------------

// NodeKind constants are the values returned by CSTNode.Kind() for
// nodes relevant to import resolution.  Consumers match against these
// rather than raw string literals.
const (
	// NodeKindObjjImport is the node produced by an @import directive.
	// Its path field is either NodeKindSystemLibString or NodeKindString.
	NodeKindObjjImport = "objj_import"

	// NodeKindSystemLibString is the angle-bracket form: <Foundation/CPObject.j>
	// The text includes the surrounding < and > delimiters.
	NodeKindSystemLibString = "system_lib_string"

	// NodeKindString is the quoted form: "MyClass.j"
	// The text includes the surrounding quotation marks.
	NodeKindString = "string"
)

// ---------------------------------------------------------------------------
// Grammar field-name constants
// ---------------------------------------------------------------------------

// FieldImportPath is the named field on NodeKindObjjImport that holds the
// import path child (either system_lib_string or string).
const FieldImportPath = "path"

// ---------------------------------------------------------------------------
// CSTNode
// ---------------------------------------------------------------------------

// CSTNode is a node in the concrete syntax tree produced by capp-parse.
//
// The node and source fields are unexported; bindgen sees only the method
// signatures, all of which return representable Go types.  The source bytes
// are retained alongside the node so Text() can be called without a
// caller-supplied buffer.
//
// A zero-value CSTNode is null.  Call IsNull before traversing children or
// reading text on a node returned by Child or ChildByFieldName, as either
// may return a null node when the requested index or field does not exist.
type CSTNode struct {
	node   *sitter.Node
	source []byte
}

// IsNull reports whether this node has no underlying tree-sitter node.
// Child and ChildByFieldName return a null CSTNode when the requested
// index is out of range or the named field is absent.
func (n CSTNode) IsNull() bool { return n.node == nil }

// Kind returns the grammar rule name for this node, e.g. "objj_import",
// "system_lib_string", "string", "identifier".  Returns an empty string
// for a null node.
func (n CSTNode) Kind() string {
	if n.node == nil {
		return ""
	}
	return n.node.Kind()
}

// Text returns the source text spanned by this node, derived directly from
// the retained source bytes.  For a null node Text returns an empty string.
func (n CSTNode) Text() string {
	if n.node == nil {
		return ""
	}
	start := n.node.StartByte()
	end := n.node.EndByte()
	if end > uint(len(n.source)) {
		end = uint(len(n.source))
	}
	return string(n.source[start:end])
}

// ChildCount returns the total number of children, named and anonymous.
// Returns 0 for a null node.
func (n CSTNode) ChildCount() int {
	if n.node == nil {
		return 0
	}
	return int(n.node.ChildCount())
}

// Child returns the i-th child (named or anonymous, zero-based).
// Returns a null CSTNode when i is out of range or the node is null.
func (n CSTNode) Child(i int) CSTNode {
	if n.node == nil || i < 0 || uint(i) >= n.node.ChildCount() {
		return CSTNode{source: n.source}
	}
	return CSTNode{node: n.node.Child(uint(i)), source: n.source}
}

// ChildByFieldName returns the child associated with the named grammar field
// and true, or a null CSTNode and false when no such field exists on this
// node.  Use the FieldImportPath and similar constants rather than raw
// string literals.
func (n CSTNode) ChildByFieldName(name string) (CSTNode, bool) {
	if n.node == nil {
		return CSTNode{source: n.source}, false
	}
	child := n.node.ChildByFieldName(name)
	if child == nil {
		return CSTNode{source: n.source}, false
	}
	return CSTNode{node: child, source: n.source}, true
}

// HasError reports whether this node or any descendant contains a parse error.
func (n CSTNode) HasError() bool {
	if n.node == nil {
		return false
	}
	return n.node.HasError()
}

// IsError reports whether this node itself is an ERROR node inserted by
// tree-sitter's error recovery.
func (n CSTNode) IsError() bool {
	if n.node == nil {
		return false
	}
	return n.node.IsError()
}

// IsMissing reports whether this node is a MISSING node — a zero-width node
// synthesised by error recovery to satisfy a grammar rule.
func (n CSTNode) IsMissing() bool {
	if n.node == nil {
		return false
	}
	return n.node.IsMissing()
}

// StartRow returns the zero-based source row of the first byte of this node.
func (n CSTNode) StartRow() uint32 {
	if n.node == nil {
		return 0
	}
	return uint32(n.node.StartPosition().Row)
}

// StartColumn returns the zero-based source column of the first byte of
// this node.
func (n CSTNode) StartColumn() uint32 {
	if n.node == nil {
		return 0
	}
	return uint32(n.node.StartPosition().Column)
}

// ---------------------------------------------------------------------------
// FileResult.Root
// ---------------------------------------------------------------------------

// Root returns the root CSTNode of this file's syntax tree.
//
// The returned CSTNode is only valid while the owning FileResult or
// ProjectResult remains live and has not been closed.  Callers must not
// retain CSTNode values beyond the Close() call of the owning result.
func (f *FileResult) Root() CSTNode {
	if f == nil || f.Tree == nil {
		return CSTNode{}
	}
	tree := f.Tree.(*sitter.Tree)
	return CSTNode{node: tree.RootNode(), source: f.Source}
}
