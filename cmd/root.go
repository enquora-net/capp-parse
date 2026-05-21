/*
 * cmd/root.go
 * cappuccino
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

/*
 * Package cmd assembles the capp-parse command tree.
 * When this binary is absorbed into the toolchain, this file and version.go
 * are discarded; cmd/parse/, cmd/debug/, cmd/install/, cmd/verify/, and
 * cmd/uninstall/ move into the toolchain's own cmd/ tree.
 */
package cmd

import (
	"dev.cappuccino/capp-parse/cmd/debug"
	"dev.cappuccino/capp-parse/cmd/parse"
	"dev.cappuccino/capp-parse/cmd/verify"

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
	root.AddCommand(debug.NewDebugCmd())
	root.AddCommand(verify.NewVerifyCmd())
	root.AddCommand(newVersionCmd(version))

	return root
}
