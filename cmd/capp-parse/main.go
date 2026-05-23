/*
 * cmd/capp-parse/main.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package main

import (
	"fmt"
	"os"

	"github.com/enquora-net/capp-parse/cmd"
)

// version is set at link time:
//
//	go build -ldflags "-X main.version=v1.2.3"
var version = "dev"

func main() {
	if err := cmd.NewRootCmd(version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
