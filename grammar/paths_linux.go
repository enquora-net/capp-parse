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
	"strings"
)

// searchPaths returns the ordered list of directories to search for the
// grammar dynamic library on Linux.
//
// Search order follows the XDG Base Directory Specification:
//
//	$XDG_DATA_HOME/dev.cappuccino/grammar/          (default: ~/.local/share/...)
//	$XDG_DATA_DIRS/dev.cappuccino/grammar/          (default: /usr/local/share/... then /usr/share/...)
//	/usr/local/lib                                   (unmanaged)
func searchPaths() []string {
	var paths []string

	// XDG_DATA_HOME: per-user managed location
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		if home, err := os.UserHomeDir(); err == nil {
			dataHome = home + "/.local/share"
		}
	}
	if dataHome != "" {
		paths = append(paths, dataHome+"/dev.cappuccino/grammar")
	}

	// XDG_DATA_DIRS: system managed locations
	dataDirs := os.Getenv("XDG_DATA_DIRS")
	if dataDirs == "" {
		dataDirs = "/usr/local/share:/usr/share"
	}
	for _, dir := range strings.Split(dataDirs, ":") {
		dir = strings.TrimRight(dir, "/")
		if dir != "" {
			paths = append(paths, dir+"/dev.cappuccino/grammar")
		}
	}

	// Unmanaged fallback
	paths = append(paths, "/usr/local/lib")

	return paths
}