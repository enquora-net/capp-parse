/*
 * cmd/uninstall/uninstall.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 24, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package uninstall implements the uninstall subcommand.
 *
 * Removes the Objective-J grammar dynamic library from its installed
 * location. Only removes files placed by the install command.
 */
package uninstall

import (
	"fmt"
	"os"

	"capp-parse/internal/grammar"

	"github.com/spf13/cobra"
)

// NewUninstallCmd constructs the uninstall command.
func NewUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the installed Objective-J grammar library",
		Long: `Remove the Objective-J grammar dynamic library from
/usr/local/lib/libtree-sitter-objj.dylib.

Requires write access to /usr/local/lib. Use sudo if necessary:

  sudo capp-parse uninstall`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dest := grammar.InstallPath()

			if _, err := os.Stat(dest); os.IsNotExist(err) {
				fmt.Fprintf(cmd.OutOrStdout(), "not installed: %s\n", dest)
				return nil
			}

			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("removing %s: %w\ntry: sudo capp-parse uninstall", dest, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "removed: %s\n", dest)
			return nil
		},
	}
}
