//go:build darwin || linux

/*
 * grammar/load_unix.go
 * cappuccino
 *
 * Created by David Richardson on Wednesday, May 20, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package grammar

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

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
