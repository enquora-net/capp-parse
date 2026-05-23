/*
 * internal/core/mode.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 */
package core

import "github.com/enquora-net/capp-parse/internal/types"

// accept reports whether path should be parsed under mode m.
// Wraps types.Mode.Accept for use within the core package.
func accept(m types.Mode, path string) bool {
	return m.Accept(path)
}
