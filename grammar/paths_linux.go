/*
 * grammar/paths_linux.go
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
// grammar dynamic library on Linux. This list is the canonical reference
// for supported install locations.
//
// Search order:
//
//	$XDG_DATA_HOME/cappuccino/grammar   (current user; default ~/.local/share/...)
//	/usr/local/lib                      (all users)
//
// XDG_DATA_DIRS is deliberately not searched: its defaults (/usr/local/share,
// /usr/share) hold architecture-independent files, not compiled libraries.
func locations() []Location {
	var locs []Location
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		if home, err := os.UserHomeDir(); err == nil {
			dataHome = filepath.Join(home, ".local", "share")
		}
	}
	if dataHome != "" {
		locs = append(locs, Location{
			Dir:   filepath.Join(dataHome, "cappuccino", "grammar"),
			Scope: ScopeUser,
		})
	}
	locs = append(locs, Location{Dir: "/usr/local/lib", Scope: ScopeSystem})
	return locs
}
