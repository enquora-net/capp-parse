//go:build windows

/*
 * grammar/load_windows.go
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

	"golang.org/x/sys/windows"
)

// openLibrary loads the library by its full path. The system DLL search order
// is never used, so a library planted elsewhere cannot be loaded instead.
func openLibrary(path string) (uintptr, error) {
	handle, err := windows.LoadLibraryEx(path, 0, windows.LOAD_WITH_ALTERED_SEARCH_PATH)
	if err != nil {
		return 0, fmt.Errorf("LoadLibrary %s: %w", path, err)
	}
	return uintptr(handle), nil
}

// lookup returns the address of an exported symbol, and whether it exists.
func lookup(handle uintptr, name string) (uintptr, bool) {
	sym, err := windows.GetProcAddress(windows.Handle(handle), name)
	if err != nil || sym == 0 {
		return 0, false
	}
	return sym, true
}
