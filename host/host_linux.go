/*
 * host/host_linux.go
 * cappuccino
 *
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package host

import (
	"bufio"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func detect(f *Facts) {
	f.OSName = osReleaseName()
	var uts unix.Utsname
	if err := unix.Uname(&uts); err == nil {
		f.OSVersion = unix.ByteSliceToString(uts.Release[:])
		f.HostArch = goarch(unix.ByteSliceToString(uts.Machine[:]))
	}
	// Translation cannot be detected: qemu-user reports the emulated
	// architecture to uname, and binfmt_misc entries show only that emulation
	// is installed, not that this process uses it. WSL is not translation.
	f.TranslationKnown = false
}

// osReleaseName returns PRETTY_NAME from /etc/os-release, e.g.
// "Ubuntu 24.04.1 LTS", or "Linux" when it is unavailable.
func osReleaseName() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Linux"
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if value, ok := strings.CutPrefix(scanner.Text(), "PRETTY_NAME="); ok {
			return strings.Trim(value, `"'`)
		}
	}
	return "Linux"
}

// goarch maps a uname machine name to its GOARCH name.
func goarch(machine string) string {
	switch machine {
	case "x86_64":
		return "amd64"
	case "aarch64", "arm64":
		return "arm64"
	default:
		return machine
	}
}
