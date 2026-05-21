/*
 * capp/mode.go
 * cappuccino
 *
 * Created by David Richardson on Saturday, April 11, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package capp

import (
	"fmt"
	"path/filepath"
	"strings"
)

func parseMode(s string) (Mode, error) {
	switch strings.ToLower(s) {
	case "auto", "":
		return ModeAuto, nil
	case "objj":
		return ModeObjj, nil
	case "js":
		return ModeJS, nil
	case "both":
		return ModeBoth, nil
	default:
		return 0, fmt.Errorf("unknown mode %q: use auto | objj | js | both", s)
	}
}

func (m Mode) accept(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch m {
	case ModeObjj:
		return ext == ".j" || ext == ".sj"
	case ModeJS:
		return ext == ".js"
	case ModeBoth:
		return ext == ".j" || ext == ".sj" || ext == ".js"
	default: // ModeAuto
		return ext == ".j" || ext == ".sj" || ext == ".js"
	}
}
