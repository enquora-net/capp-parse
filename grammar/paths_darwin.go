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

import "os"

// searchPaths returns the ordered list of directories to search for the
// grammar dynamic library on macOS.
//
// Search order:
//
//	~/Library/Application Support/dev.cappuccino/grammar/   (user managed)
//	/Library/Application Support/dev.cappuccino/grammar/    (system managed)
//	/usr/local/lib                                           (unmanaged)
func searchPaths() []string {
	var paths []string

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, home+"/Library/Application Support/dev.cappuccino/grammar")
	}

	paths = append(paths,
		"/Library/Application Support/dev.cappuccino/grammar",
		"/usr/local/lib",
	)

	return paths
}