/*
 * host/host.go
 * cappuccino
 *
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

/*
 * Package host reports facts about the machine and the running process, for
 * diagnostics and support: the operating system and its version, the
 * architecture the binary was built for, the machine's own architecture, and
 * whether the process runs under architecture translation, such as an amd64
 * binary under Rosetta 2 on Apple silicon or under x64 emulation on Windows
 * on ARM.
 *
 * Detection is platform-specific and defined in host_darwin.go,
 * host_linux.go, and host_windows.go.
 */
package host

import "runtime"

// Facts describes the host and the running process.
type Facts struct {
	OS        string `json:"os"`         // runtime.GOOS
	OSName    string `json:"os_name"`    // e.g. "macOS", "Ubuntu 24.04.1 LTS", "Windows 11"
	OSVersion string `json:"os_version"` // e.g. "26.1 (25B78)", "6.8.0-45-generic", "10.0.26100"
	// BinaryArch is the architecture the binary was built for (runtime.GOARCH).
	BinaryArch string `json:"binary_arch"`
	// HostArch is the machine's own architecture, in GOARCH terms. On Linux it
	// is what the kernel reports, which an emulator can disguise.
	HostArch string `json:"host_arch"`
	// TranslationKnown reports whether translation can be detected on this
	// platform. It cannot on Linux: qemu-user reports the emulated
	// architecture, and nothing reliably reveals it to the process.
	TranslationKnown bool `json:"translation_known"`
	// Translated reports whether the process runs under architecture
	// translation; meaningful only when TranslationKnown.
	Translated bool `json:"translated"`
	// Translator names the translation layer when Translated, e.g. "Rosetta 2".
	Translator string `json:"translator"`
}

// Detect returns the facts for this host and process.
func Detect() Facts {
	f := Facts{OS: runtime.GOOS, BinaryArch: runtime.GOARCH, HostArch: runtime.GOARCH}
	detect(&f)
	return f
}
