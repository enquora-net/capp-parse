/*
 * cmd/root.go
 * capp-parse
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
 * are discarded; cmd/parse/, cmd/debug/, cmd/verify/ move into the
 * toolchain's own cmd/ tree unchanged.
 */
package cmd

import (
	"github.com/enquora-net/capp-parse/cmd/debug"
	"github.com/enquora-net/capp-parse/cmd/parse"
	"github.com/enquora-net/capp-parse/cmd/verify"

	"github.com/spf13/cobra"
)

// NewRootCmd constructs the root command.
func NewRootCmd(version, commit, date string) *cobra.Command {
	root := &cobra.Command{
		Use:          "capp-parse",
		Short:        "Cappuccino Objective-J source parser",
		SilenceUsage: true,
	}

	root.AddCommand(parse.NewParseCmd())
	root.AddCommand(debug.NewDebugCmd())
	root.AddCommand(verify.NewVerifyCmd())
	root.AddCommand(newVersionCmd(version, commit, date))

	return root
}
