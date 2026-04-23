/*
 * main.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package main

import (
	"fmt"
	"os"

	"capp-parse/cmd"
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
