/*
 * internal/grammar/grammar.go
 * capp-parse
 *
 * Created by David Richardson on Thursday, April 23, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package grammar loads a tree-sitter grammar dynamic library via purego.
 *
 * The grammar dynamic library is a first-class independently deployable
 * artifact available for use by third-party tooling in any language.
 *
 * Search order:
 *   /usr/local/lib                    (primary, idiomatic)
 *   /opt/local/lib                    (MacPorts)
 *   ~/Library/tree-sitter             (fallback, tree-sitter tooling)
 *   /usr/local/lib/tree-sitter        (fallback, tree-sitter tooling)
 *   /opt/local/lib/tree-sitter        (fallback, MacPorts tree-sitter tooling)
 *
 * Library name: libtree-sitter-<lang>.dylib / libtree-sitter-<lang>.so
 *
 * The canonical install location is /usr/local/lib.
 */
package grammar

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

const languageName = "objj"

var searchPaths = []string{
	"/usr/local/lib",
	"/opt/local/lib",
	filepath.Join(os.Getenv("HOME"), "Library/tree-sitter"),
	"/usr/local/lib/tree-sitter",
	"/opt/local/lib/tree-sitter",
}

var (
	langOnce sync.Once
	lang     *sitter.Language
	langErr  error
)

// Language returns the singleton tree-sitter Language for Objective-J.
// Safe for concurrent use; initialisation is performed at most once.
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
				"grammar library for %q not found\n"+
					"searched: %v\n"+
					"run: capp-parse install\n"+
					"or use --grammar to specify the path explicitly",
				languageName, searchPaths,
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
// string if not found. Exposed for use by the install and verify commands.
func FindLibrary() string {
	ext := libExt()
	candidates := []string{
		fmt.Sprintf("libtree-sitter-%s.%s", languageName, ext),
		fmt.Sprintf("tree-sitter-%s.%s", languageName, ext),
	}
	for _, base := range searchPaths {
		for _, name := range candidates {
			p := filepath.Join(base, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

// InstallPath returns the canonical install path for the grammar library.
func InstallPath() string {
	return filepath.Join("/usr/local/lib",
		fmt.Sprintf("libtree-sitter-%s.%s", languageName, libExt()))
}

func libExt() string {
	if runtime.GOOS == "darwin" {
		return "dylib"
	}
	return "so"
}

func load(path string) (*sitter.Language, error) {
	lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, fmt.Errorf("dlopen %s: %w", path, err)
	}
	var languageFunc func() uintptr
	purego.RegisterLibFunc(&languageFunc, lib, "tree_sitter_"+languageName)
	ptr := languageFunc()
	if ptr == 0 {
		return nil, fmt.Errorf("tree_sitter_%s returned nil", languageName)
	}
	return sitter.NewLanguage(unsafe.Pointer(ptr)), nil //nolint:unsafeptr
}
