/*
 * cmd/version.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 10, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the capp-parse version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "capp-parse %s\n", version)
		},
	}
}
