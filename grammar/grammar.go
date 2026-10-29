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
 * Search locations are platform-specific and defined in paths_darwin.go,
 * paths_linux.go, and paths_windows.go; they are the canonical reference for
 * supported install locations. Library loading is platform-specific and
 * defined in load_unix.go and load_windows.go.
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
	"unsafe"

	"github.com/ebitengine/purego"
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

// Scope describes who a search location serves.
type Scope string

const (
	ScopeUser      Scope = "user"      // the current user only
	ScopeSystem    Scope = "system"    // all users
	ScopeUnmanaged Scope = "unmanaged" // all users, outside any managed location
)

// Location is one directory in the search order.
type Location struct {
	Dir   string
	Scope Scope
}

// Locations returns the ordered list of directories searched for the grammar
// library on the current platform: the canonical supported install locations.
func Locations() []Location {
	return locations()
}

// LibraryNames returns the file names accepted in each location, in order of
// preference.
func LibraryNames() []string {
	return libCandidates()
}

// FindLibrary returns the path to the grammar dynamic library, or empty
// string if not found: the first accepted file name in the first location
// that has one.
func FindLibrary() string {
	if found := FindLibraries(); len(found) > 0 {
		return found[0]
	}
	return ""
}

// FindLibraries returns every grammar library present in the search
// locations, in search order. The first, if any, is the one Language loads.
func FindLibraries() []string {
	var found []string
	for _, loc := range locations() {
		for _, name := range libCandidates() {
			p := filepath.Join(loc.Dir, name)
			if fi, err := os.Stat(p); err == nil && fi.Mode().IsRegular() {
				found = append(found, p)
			}
		}
	}
	return found
}

// SearchPaths returns the ordered list of directories consulted by
// FindLibrary on the current platform. Intended for use by verify and
// diagnostic tooling.
func SearchPaths() []string {
	locs := locations()
	dirs := make([]string, len(locs))
	for i, loc := range locs {
		dirs[i] = loc.Dir
	}
	return dirs
}

// Info describes a grammar library as reported by the library itself.
type Info struct {
	Path string `json:"path"`
	// ABIVersion is the tree-sitter language ABI the library was generated for.
	ABIVersion int `json:"abi_version"`
	// HasRelease reports whether the library carries tree-sitter language
	// metadata (ABI 15 or later); Major, Minor and Patch are its release
	// version, from the grammar's tree-sitter.json.
	HasRelease bool `json:"has_release"`
	Major      int  `json:"major"`
	Minor      int  `json:"minor"`
	Patch      int  `json:"patch"`
	// HasBuild reports whether the library exports tree_sitter_objj_version;
	// Build is the string it returns, e.g. 2.0.0-beta.2+gd434afd2ab1c.
	// Libraries released before 2.0.0-beta.2 do not export it.
	HasBuild bool   `json:"has_build"`
	Build    string `json:"build"`
}

// Inspect loads the library at path and reports its versions. Loading runs
// the library's initialisers in this process, and a process can hold only
// one copy of a given library reliably: inspect libraries other than the
// active one in a separate process.
func Inspect(path string) (Info, error) {
	info := Info{Path: path}
	handle, lang, err := open(path)
	if err != nil {
		return info, err
	}
	info.ABIVersion = int(lang.AbiVersion())
	if m := lang.Metadata(); m != nil {
		info.HasRelease = true
		info.Major = int(m.MajorVersion)
		info.Minor = int(m.MinorVersion)
		info.Patch = int(m.PatchVersion)
	}
	if sym, ok := lookup(handle, "tree_sitter_"+languageName+"_version"); ok {
		var version func() string
		purego.RegisterFunc(&version, sym)
		info.HasBuild = true
		info.Build = version()
	}
	return info, nil
}

func load(path string) (*sitter.Language, error) {
	_, lang, err := open(path)
	return lang, err
}

// open loads the library at path and returns its handle and language.
func open(path string) (uintptr, *sitter.Language, error) {
	handle, err := openLibrary(path)
	if err != nil {
		return 0, nil, err
	}
	sym, ok := lookup(handle, "tree_sitter_"+languageName)
	if !ok {
		return 0, nil, fmt.Errorf("%s: tree_sitter_%s not exported", path, languageName)
	}
	var languageFunc func() uintptr
	purego.RegisterFunc(&languageFunc, sym)
	ptr := languageFunc()
	if ptr == 0 {
		return 0, nil, fmt.Errorf("tree_sitter_%s returned nil", languageName)
	}
	return handle, sitter.NewLanguage(unsafe.Pointer(ptr)), nil //nolint:unsafeptr
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
	paths := SearchPaths()
	lines := make([]string, len(paths))
	for i, p := range paths {
		lines[i] = "  " + p
	}
	return strings.Join(lines, "\n")
}
