/*
 * internal/core/mode.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 */
package core

import (
	"path/filepath"
	"strings"
)

// Mode int values mirror the root package constants (ModeAuto=0, ModeObjj=1, ModeJS=2, ModeBoth=3).
// accept reports whether path should be parsed under the given mode.
func accept(mode int, path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch mode {
	case 1: // ModeObjj
		return ext == ".j" || ext == ".sj"
	case 2: // ModeJS
		return ext == ".js"
	case 3: // ModeBoth
		return ext == ".j" || ext == ".sj" || ext == ".js"
	default: // ModeAuto = 0
		return ext == ".j" || ext == ".sj" || ext == ".js"
	}
}
