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

	"github.com/ebitengine/purego"
)

func openLibrary(path string) (uintptr, error) {
	handle, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, fmt.Errorf("dlopen %s: %w", path, err)
	}
	return handle, nil
}

// lookup returns the address of an exported symbol, and whether it exists.
func lookup(handle uintptr, name string) (uintptr, bool) {
	sym, err := purego.Dlsym(handle, name)
	if err != nil || sym == 0 {
		return 0, false
	}
	return sym, true
}
