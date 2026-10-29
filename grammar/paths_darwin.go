/*
 * grammar/paths_darwin.go
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
// grammar dynamic library on macOS. This list is the canonical reference
// for supported install locations.
//
// Search order:
//
//	~/Library/Application Support/dev.cappuccino/grammar   (current user)
//	/Library/Application Support/dev.cappuccino/grammar    (all users)
//	/usr/local/lib                                         (all users, unmanaged)
func locations() []Location {
	var locs []Location
	if home, err := os.UserHomeDir(); err == nil {
		locs = append(locs, Location{
			Dir:   filepath.Join(home, "Library", "Application Support", "dev.cappuccino", "grammar"),
			Scope: ScopeUser,
		})
	}
	locs = append(locs,
		Location{Dir: "/Library/Application Support/dev.cappuccino/grammar", Scope: ScopeSystem},
		Location{Dir: "/usr/local/lib", Scope: ScopeUnmanaged},
	)
	return locs
}
