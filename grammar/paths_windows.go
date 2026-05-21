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

import "os"

// searchPaths returns the ordered list of directories to search for the
// grammar dynamic library on Windows.
//
// Search order:
//
//	%LOCALAPPDATA%\github.com/enquora-net/\grammar\   (per-user managed)
//	%PROGRAMDATA%\github.com/enquora-net/\grammar\    (system managed)
//
// There is no unmanaged fallback on Windows. Use --grammar to specify
// an explicit path when the library is not in a managed location.
func searchPaths() []string {
	var paths []string

	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		paths = append(paths, localAppData+`\github.com/enquora-net/\grammar`)
	}

	if programData := os.Getenv("PROGRAMDATA"); programData != "" {
		paths = append(paths, programData+`\github.com/enquora-net/\grammar`)
	}

	return paths
}
