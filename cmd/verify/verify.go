/*
 * cmd/verify/verify.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 24, 2026.
  * Copyright (c) 2026 David Richardson. All rights reserved.
 *
*/

/*
 * Package verify implements the verify subcommand.
 *
 * Checks whether capp-parse prerequisites are in place and reports
 * their status. Currently verifies the Objective-J grammar library.
 */
package verify

import (
	"fmt"

	"capp-parse/internal/grammar"

	"github.com/spf13/cobra"
)

// NewVerifyCmd constructs the verify command.
func NewVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify capp-parse prerequisites",
		Long:  `Check whether capp-parse prerequisites are installed and locatable.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			allOK := true

			// Grammar library
			path := grammar.FindLibrary()
			if path == "" {
				fmt.Fprintf(out, "✗  grammar library not found\n")
				fmt.Fprintf(out, "   run: capp-parse install\n")
				allOK = false
			} else {
				fmt.Fprintf(out, "✓  grammar library: %s\n", path)
			}

			if !allOK {
				return fmt.Errorf("one or more prerequisites missing")
			}
			return nil
		},
	}
}
