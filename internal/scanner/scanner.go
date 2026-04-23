/*
 * internal/scanner/scanner.go
 * capp-parse
 *
 * Created by David Richardson on Sunday, April 12, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package scanner is a pure-Go port of the Cappuccino Objective-J tree-sitter
 * external scanner (scanner.c).
 *
 * The scanner is stateless: Create returns nil, Serialize writes 0 bytes,
 * Deserialize and Destroy are no-ops. All logic lives in Scan.
 *
 * Token enum order MUST match the externals array in grammar.js:
 *
 *	0  _automatic_semicolon
 *	1  _template_chars
 *	2  _ternary_qmark
 *	3  html_comment
 *	4  ||  (LOGICAL_OR)
 *	5  escape_sequence
 *	6  regex_pattern
 *	7  jsx_text
 */
package scanner

import (
	"unicode"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// tokenType mirrors the C enum; ordinal == index in validSymbols.
type tokenType int

const (
	automaticSemicolon tokenType = iota
	templateChars
	ternaryQmark
	htmlComment
	logicalOr
	escapeSequence
	regexPattern
	jsxText
)

// sym converts a tokenType to the gotreesitter.Symbol type.
func sym(t tokenType) gotreesitter.Symbol {
	return gotreesitter.Symbol(t)
}

// Scanner implements gotreesitter.ExternalScanner.
// The payload threading through all methods is nil — the scanner is stateless.
type Scanner struct{}

// New returns a Scanner satisfying gotreesitter.ExternalScanner.
func New() gotreesitter.ExternalScanner {
	return &Scanner{}
}

func (s *Scanner) Create() any                   { return nil }
func (s *Scanner) Destroy(_ any)                 {}
func (s *Scanner) Serialize(_ any, _ []byte) int { return 0 }
func (s *Scanner) Deserialize(_ any, _ []byte)   {}

// Scan is the main dispatch, mirroring tree_sitter_objj_external_scanner_scan.
func (s *Scanner) Scan(_ any, lexer *gotreesitter.ExternalLexer, valid []bool) bool {
	// fmt.Printf("Scan called: valid=%v lookahead=%q\n", valid, lexer.Lookahead())
	if valid[templateChars] {
		if valid[automaticSemicolon] {
			return false
		}
		return scanTemplateChars(lexer)
	}

	if valid[jsxText] && scanJSXText(lexer) {
		return true
	}

	if valid[automaticSemicolon] {
		scannedComment := false
		ok := scanAutomaticSemicolon(lexer, !valid[logicalOr], &scannedComment)
		if !ok && !scannedComment && valid[ternaryQmark] && lexer.Lookahead() == '?' {
			return scanTernaryQmark(lexer)
		}
		return ok
	}

	if valid[ternaryQmark] {
		return scanTernaryQmark(lexer)
	}

	if valid[htmlComment] && !valid[logicalOr] && !valid[escapeSequence] && !valid[regexPattern] {
		return scanHTMLComment(lexer)
	}

	return false
}

// ---------------------------------------------------------------------------
// Template characters
// ---------------------------------------------------------------------------

func scanTemplateChars(lexer *gotreesitter.ExternalLexer) bool {
	lexer.SetResultSymbol(sym(templateChars))
	hasContent := false
	for {
		lexer.MarkEnd()
		switch lexer.Lookahead() {
		case '`':
			return hasContent
		case 0:
			return false
		case '$':
			lexer.Advance(false)
			if lexer.Lookahead() == '{' {
				return hasContent
			}
		case '\\':
			return hasContent
		default:
			lexer.Advance(false)
		}
		hasContent = true
	}
}

// ---------------------------------------------------------------------------
// Whitespace / comment scanner
// ---------------------------------------------------------------------------

type whitespaceResult int

const (
	wsReject    whitespaceResult = iota
	wsNoNewline                  // ambiguous, keep scanning
	wsAccept
)

func scanWhitespaceAndComments(lexer *gotreesitter.ExternalLexer, scannedComment *bool, consume bool) whitespaceResult {
	sawBlockNewline := false

	for {
		for unicode.IsSpace(lexer.Lookahead()) {
			lexer.Advance(true)
		}

		if lexer.Lookahead() != '/' {
			return wsAccept
		}
		lexer.Advance(true)

		switch lexer.Lookahead() {
		case '/':
			lexer.Advance(true)
			for {
				ch := lexer.Lookahead()
				if ch == 0 || ch == '\n' || ch == 0x2028 || ch == 0x2029 {
					break
				}
				lexer.Advance(true)
			}
			*scannedComment = true

		case '*':
			lexer.Advance(true)
			for lexer.Lookahead() != 0 {
				if lexer.Lookahead() == '*' {
					lexer.Advance(true)
					if lexer.Lookahead() == '/' {
						lexer.Advance(true)
						*scannedComment = true
						if lexer.Lookahead() != '/' && !consume {
							if sawBlockNewline {
								return wsAccept
							}
							return wsNoNewline
						}
						break
					}
				} else if ch := lexer.Lookahead(); ch == '\n' || ch == 0x2028 || ch == 0x2029 {
					sawBlockNewline = true
					lexer.Advance(true)
				} else {
					lexer.Advance(true)
				}
			}

		default:
			return wsReject
		}
	}
}

// ---------------------------------------------------------------------------
// Automatic semicolon insertion
// ---------------------------------------------------------------------------

func scanAutomaticSemicolon(lexer *gotreesitter.ExternalLexer, commentCondition bool, scannedComment *bool) bool {
	lexer.SetResultSymbol(sym(automaticSemicolon))
	lexer.MarkEnd()

	for {
		if lexer.Lookahead() == 0 {
			return true
		}

		if lexer.Lookahead() == '/' {
			result := scanWhitespaceAndComments(lexer, scannedComment, false)
			if result == wsReject {
				return false
			}
			if result == wsAccept && commentCondition &&
				lexer.Lookahead() != ',' && lexer.Lookahead() != '=' {
				return true
			}
		}

		if lexer.Lookahead() == '}' {
			return true
		}

		ch := lexer.Lookahead()
		if ch == '\n' || ch == 0x2028 || ch == 0x2029 {
			break
		}

		if !unicode.IsSpace(ch) {
			return false
		}

		lexer.Advance(true)
	}

	lexer.Advance(true)

	if scanWhitespaceAndComments(lexer, scannedComment, true) == wsReject {
		return false
	}

	switch lexer.Lookahead() {
	case '`', ',', ':', ';', '*', '%', '>', '<', '=', '(', '?', '^', '|', '&', '/':
		return false

	case '.':
		lexer.Advance(true)
		return unicode.IsDigit(lexer.Lookahead())

	case '+':
		lexer.Advance(true)
		return lexer.Lookahead() == '+'

	case '-':
		lexer.Advance(true)
		return lexer.Lookahead() == '-'

	case '!':
		lexer.Advance(true)
		return lexer.Lookahead() != '='

	case 'i':
		lexer.Advance(true)
		if lexer.Lookahead() != 'n' {
			return true
		}
		lexer.Advance(true)
		if !unicode.IsLetter(lexer.Lookahead()) {
			return false
		}
		for _, expected := range "stanceof" {
			if lexer.Lookahead() != expected {
				return true
			}
			lexer.Advance(true)
		}
		if !unicode.IsLetter(lexer.Lookahead()) {
			return false
		}
	}

	return true
}

// ---------------------------------------------------------------------------
// Ternary ?
// ---------------------------------------------------------------------------

func scanTernaryQmark(lexer *gotreesitter.ExternalLexer) bool {
	for unicode.IsSpace(lexer.Lookahead()) {
		lexer.Advance(true)
	}

	if lexer.Lookahead() != '?' {
		return false
	}
	lexer.Advance(false)

	if lexer.Lookahead() == '?' {
		return false
	}

	lexer.MarkEnd()
	lexer.SetResultSymbol(sym(ternaryQmark))

	if lexer.Lookahead() == '.' {
		lexer.Advance(false)
		return unicode.IsDigit(lexer.Lookahead())
	}
	return true
}

// ---------------------------------------------------------------------------
// HTML comment  <!-- ... -->
// ---------------------------------------------------------------------------

func scanHTMLComment(lexer *gotreesitter.ExternalLexer) bool {
	for ch := lexer.Lookahead(); unicode.IsSpace(ch) || ch == 0x2028 || ch == 0x2029; ch = lexer.Lookahead() {
		lexer.Advance(true)
	}

	const commentStart = "<!--"
	const commentEnd = "-->"

	switch lexer.Lookahead() {
	case '<':
		for _, expected := range commentStart {
			if lexer.Lookahead() != expected {
				return false
			}
			lexer.Advance(false)
		}
	case '-':
		for _, expected := range commentEnd {
			if lexer.Lookahead() != expected {
				return false
			}
			lexer.Advance(false)
		}
	default:
		return false
	}

	for {
		ch := lexer.Lookahead()
		if ch == 0 || ch == '\n' || ch == 0x2028 || ch == 0x2029 {
			break
		}
		lexer.Advance(false)
	}

	lexer.SetResultSymbol(sym(htmlComment))
	lexer.MarkEnd()
	return true
}

// ---------------------------------------------------------------------------
// JSX text
// ---------------------------------------------------------------------------

func scanJSXText(lexer *gotreesitter.ExternalLexer) bool {
	sawText := false
	atNewline := false

	for {
		ch := lexer.Lookahead()
		if ch == 0 || ch == '<' || ch == '>' || ch == '{' || ch == '}' || ch == '&' {
			break
		}
		isWspace := unicode.IsSpace(ch)
		if ch == '\n' {
			atNewline = true
		} else {
			atNewline = atNewline && isWspace
			if !atNewline {
				sawText = true
			}
		}
		lexer.Advance(false)
	}

	lexer.SetResultSymbol(sym(jsxText))
	return sawText
}
