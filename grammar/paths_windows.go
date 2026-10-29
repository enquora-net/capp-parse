/*
 * grammar/paths_windows.go
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
	"os"
	"path/filepath"
)

// locations returns the ordered list of directories to search for the
// grammar dynamic library on Windows. This list is the canonical reference
// for supported install locations.
//
// Search order:
//
//	%LOCALAPPDATA%\Cappuccino\grammar   (current user)
//	%ProgramFiles%\Cappuccino\grammar   (all users; writable by administrators only)
//
// %ProgramData% is deliberately not used: ordinary users can create folders
// there, so any user could plant a library that would load for every user.
// There is no unmanaged fallback on Windows. Use --grammar to specify an
// explicit path when the library is not in a supported location.
func locations() []Location {
	var locs []Location
	if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
		locs = append(locs, Location{
			Dir:   filepath.Join(dir, "Cappuccino", "grammar"),
			Scope: ScopeUser,
		})
	}
	if dir := os.Getenv("ProgramFiles"); dir != "" {
		locs = append(locs, Location{
			Dir:   filepath.Join(dir, "Cappuccino", "grammar"),
			Scope: ScopeSystem,
		})
	}
	return locs
}
