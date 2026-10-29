/*
 * host/host_darwin.go
 * cappuccino
 *
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package host

import "golang.org/x/sys/unix"

func detect(f *Facts) {
	f.OSName = "macOS"
	if v, err := unix.Sysctl("kern.osproductversion"); err == nil {
		f.OSVersion = v
		if build, err := unix.Sysctl("kern.osversion"); err == nil {
			f.OSVersion += " (" + build + ")"
		}
	}

	// hw.optional.arm64 is 1 on Apple silicon, including for a translated
	// process, and absent on Intel Macs.
	if v, err := unix.SysctlUint32("hw.optional.arm64"); err == nil && v == 1 {
		f.HostArch = "arm64"
	} else {
		f.HostArch = "amd64"
	}

	// sysctl.proc_translated is 1 when this process runs under Rosetta 2, and
	// absent on Intel Macs, which counts as native.
	f.TranslationKnown = true
	if v, err := unix.SysctlUint32("sysctl.proc_translated"); err == nil && v == 1 {
		f.Translated = true
		f.Translator = "Rosetta 2"
	}
}
