/*
 * grammar/grammar.go
 * cappuccino
 *
 * Created by David Richardson on Thursday, April 23, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

/*
 * Package grammar loads the Objective-J tree-sitter grammar dynamic library.
 *
 * The grammar library is a versioned, independently deployable artifact
 * managed by the Cappuccino orchestration application. This package provides
 * runtime location and loading; it does not embed or install the library.
 *
 * Search paths are platform-specific and defined in paths_darwin.go,
 * paths_linux.go, and paths_windows.go. Library loading is platform-specific
 * and defined in load_unix.go and load_windows.go.
 *
 * Use --grammar to override the search with an explicit path.
 *
 * Library name: libtree-sitter-objj.dylib (macOS), libtree-sitter-objj.so (Linux),
 *               tree-sitter-objj.dll (Windows)
 */
package grammar

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

const languageName = "objj"

var (
	langOnce sync.Once
	lang     *sitter.Language
	langErr  error
)

// Language returns the singleton tree-sitter Language for Objective-J.
// Safe for concurrent use; initialisation is performed at most once.
//
// When grammarPath is non-empty it is used directly, bypassing the search.
// Subsequent calls ignore grammarPath — the language is initialised once.
// The --grammar flag must therefore be resolved before any concurrent parse
// begins; cobra command execution guarantees this in normal usage.
func Language(grammarPath string) (*sitter.Language, error) {
	langOnce.Do(func() {
		path := grammarPath
		if path == "" {
			path = FindLibrary()
		}
		if path == "" {
			langErr = fmt.Errorf(
				"grammar library for %q not found\nsearched:\n%s\nuse --grammar to specify the path explicitly",
				languageName,
				formatSearchPaths(),
			)
			return
		}
		lang, langErr = load(path)
	})
	if langErr != nil {
		return nil, langErr
	}
	return lang, nil
}

// FindLibrary returns the path to the grammar dynamic library, or empty
// string if not found. Searches the platform-specific managed and unmanaged
// locations in order.
func FindLibrary() string {
	candidates := libCandidates()
	for _, base := range searchPaths() {
		for _, name := range candidates {
			p := filepath.Join(base, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

// SearchPaths returns the ordered list of directories consulted by
// FindLibrary on the current platform. Intended for use by verify and
// diagnostic tooling.
func SearchPaths() []string {
	return searchPaths()
}

func libCandidates() []string {
	ext := libExt()
	return []string{
		fmt.Sprintf("libtree-sitter-%s.%s", languageName, ext),
		fmt.Sprintf("tree-sitter-%s.%s", languageName, ext),
	}
}

func libExt() string {
	switch runtime.GOOS {
	case "darwin":
		return "dylib"
	case "windows":
		return "dll"
	default:
		return "so"
	}
}

func formatSearchPaths() string {
	paths := searchPaths()
	lines := make([]string, len(paths))
	for i, p := range paths {
		lines[i] = "  " + p
	}
	return strings.Join(lines, "\n")
}
