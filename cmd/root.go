/*
 * cmd/root.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package cmd assembles the capp-parse command tree.
 * When this binary is absorbed into the toolchain, this file and version.go
 * are discarded; cmd/parse/ moves into the toolchain's own cmd/ tree.
 */
package cmd

import (
	"capp-parse/cmd/parse"

	"github.com/spf13/cobra"
)

// NewRootCmd constructs the root command.
func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:          "capp-parse",
		Short:        "Cappuccino Objective-J source parser",
		SilenceUsage: true,
	}

	root.AddCommand(parse.NewParseCmd())
	root.AddCommand(newVersionCmd(version))

	return root
}
