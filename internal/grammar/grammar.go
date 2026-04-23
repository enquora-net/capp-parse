/*
 * internal/grammar/grammar.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package grammar loads the Cappuccino Objective-J tree-sitter grammar
 * from an embedded binary blob produced by ts2go and wires in the
 * pure-Go external scanner.
 *
 * # Obtaining the blob
 *
 * Run ts2go against your parser.c once:
 *
 *	go run github.com/odvcencio/gotreesitter/cmd/ts2go \
 *	    -input  path/to/tree-sitter-objj/src/parser.c \
 *	    -out    internal/grammar/objj.bin \
 *	    -compact
 *
 * The resulting file is committed to the repository and embedded here
 * via go:embed.  Re-run ts2go whenever grammar.js / parser.c changes.
 *
 * # External scanner
 *
 * The grammar requires an external scanner (scanner.c in the grammar
 * source).  That logic is ported to Go in internal/scanner and is
 * registered with the language via gotreesitter.NewLanguageWithScanner.
 */

/*
* Package grammar registers the Cappuccino Objective-J tree-sitter language
* with the gotreesitter grammar registry.
*
* Registration happens in init() using the same mechanism as the built-in
* gotreesitter grammars: Register binds the LangEntry, RegisterExternalScanner
* binds the scanner by name. Language() retrieves the registered entry.
*
* # Obtaining the blob
*
* Run ts2go against src/parser.c (copied from tree-sitter-objj):
*
*	ts2go -input src/parser.c -out internal/grammar/objj.bin -compact
*
* Commit the resulting file. Re-run whenever grammar.js / parser.c changes.
package grammar
*/
package grammar

import (
	_ "embed"
	"fmt"
	"sync"

	"capp-parse/internal/scanner"
	gotreesitter "github.com/odvcencio/gotreesitter"
)

//go:embed objj.bin
var objjBlob []byte

var (
	langOnce sync.Once
	lang     *gotreesitter.Language
	langErr  error
)

// Language returns the singleton gotreesitter Language for Objective-J.
// Safe for concurrent use; initialisation is performed at most once.
func Language() (*gotreesitter.Language, error) {
	langOnce.Do(func() {
		l, err := gotreesitter.LoadLanguage(objjBlob)
		if err != nil {
			langErr = fmt.Errorf("loading objj grammar: %w", err)
			return
		}
		l.ExternalScanner = scanner.New()
        // fmt.Printf("ExternalTokenCount: %d\n", l.ExternalTokenCount)
        // fmt.Printf("ExternalScanner set: %v\n", l.ExternalScanner != nil)
        // fmt.Printf("ExternalSymbols: %v\n", l.ExternalSymbols)
        // fmt.Printf("ExternalLexStates: %d rows\n", len(l.ExternalLexStates))
		lang = l
	})
	if langErr != nil {
		return nil, langErr
	}
	return lang, nil
}
