/*
 * cmd/install/install.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 24, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * Package install implements the install subcommand.
 *
 * Extracts the embedded Objective-J grammar dynamic library to
 * /usr/local/lib/libtree-sitter-objj.dylib. Requires write access
 * to /usr/local/lib; use sudo if necessary.
 */
package install

import (
	"fmt"
	"os"
	"path/filepath"

	"dev.cappuccino/capp-parse/internal/grammar"

	"github.com/spf13/cobra"
)

// NewInstallCmd constructs the install command.
func NewInstallCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the Objective-J grammar library",
		Long: `Extract the embedded Objective-J grammar dynamic library to
/usr/local/lib/libtree-sitter-objj.dylib.

Requires write access to /usr/local/lib. Use sudo if necessary:

  sudo capp-parse install

Use --force to overwrite an existing installation.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dest := grammar.InstallPath()
			dir := filepath.Dir(dest)

			if _, err := os.Stat(dest); err == nil {
				if !force {
					fmt.Fprintf(cmd.OutOrStdout(),
						"already installed: %s\nuse --force to overwrite\n", dest)
					return nil
				}
			}

			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating %s: %w", dir, err)
			}

			if err := os.WriteFile(dest, grammar.EmbeddedLibrary, 0755); err != nil {
				return fmt.Errorf("writing %s: %w\ntry: sudo capp-parse install", dest, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "installed: %s\n", dest)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing installation")
	return cmd
}
