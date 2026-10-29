/*
 * host/host_windows.go
 * cappuccino
 *
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package host

import (
	"fmt"

	"golang.org/x/sys/windows"
)

const (
	imageFileMachineI386  = 0x014c
	imageFileMachineAMD64 = 0x8664
	imageFileMachineARM64 = 0xaa64
)

func detect(f *Facts) {
	v := windows.RtlGetVersion()
	f.OSName = "Windows"
	if v.MajorVersion == 10 {
		f.OSName = "Windows 10"
		if v.BuildNumber >= 22000 {
			f.OSName = "Windows 11"
		}
	}
	f.OSVersion = fmt.Sprintf("%d.%d.%d", v.MajorVersion, v.MinorVersion, v.BuildNumber)

	// IsWow64Process2 reports the machine's native architecture. An x64
	// process emulated on ARM64 is not a WOW64 process, so its own process
	// machine reads as unknown; comparing the binary's architecture with the
	// native one detects it.
	var process, native uint16
	if err := windows.IsWow64Process2(windows.CurrentProcess(), &process, &native); err != nil {
		return
	}
	switch native {
	case imageFileMachineAMD64:
		f.HostArch = "amd64"
	case imageFileMachineARM64:
		f.HostArch = "arm64"
	case imageFileMachineI386:
		f.HostArch = "386"
	default:
		f.HostArch = fmt.Sprintf("machine 0x%04x", native)
	}
	f.TranslationKnown = true
	if f.HostArch != f.BinaryArch {
		f.Translated = true
		f.Translator = "Windows " + f.HostArch + " emulation of " + f.BinaryArch
	}
}
