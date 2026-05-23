/*
 * cmd/debug/debug.go
 * cappuccino
 *
 * Created by David Richardson on Thursday, April 23, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

/*
 * Package debug implements the debug subcommand.
 *
 * The debug command walks a file or directory tree, stops at the first parse
 * error, displays the full concrete syntax tree and source context, and opens
 * the offending file in Xcode at the error line. This workflow dramatically
 * accelerates grammar and compiler development.
 *
 * Like cmd/parse, this package is designed as a fold-in seam: NewDebugCmd()
 * wires directly into the parent toolchain's root command unchanged when
 * capp-parse is absorbed into cappuccino.
 */
package debug

import (
	"os"

	capp "github.com/enquora-net/capp-parse"

	"github.com/spf13/cobra"
)

// NewDebugCmd constructs the debug command.
func NewDebugCmd() *cobra.Command {
	var (
		grammar string
		mode    string
		profile bool
		context int
		noXcode bool
	)

	cmd := &cobra.Command{
		Use:   "debug <path>",
		Short: "Walk source, stop at first error, show context and open in Xcode",
		Long: `Walk a file or directory tree and stop at the first parse error.

On error, the following are displayed:
  - Source context around the error (configurable line count)
  - Parent node chain from root to the error
  - Full concrete syntax tree

The offending file is opened in Xcode at the error line via xed.

This workflow is the recommended approach during grammar and compiler
development. A single command replaces the cycle of manual jumps between
terminal, source editor, and grammar editor.

Examples:
  capp-parse debug Frameworks/AppKit
  capp-parse debug --profile --no-xcode Frameworks/
  capp-parse debug --context 5 src/AppController.j`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := capp.ParseMode(mode)
			if err != nil {
				return err
			}

			result, err := capp.Debug(capp.DebugConfig{
				Path:         args[0],
				Mode:         m,
				GrammarPath:  grammar,
				Profile:      profile,
				ContextLines: context,
				NoXcode:      noXcode,
			})
			if err != nil {
				return err
			}

			if result.Profile != nil {
				capp.EmitDebugProfile(result.Profile)
			}

			if result.HasError {
				os.Exit(2)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&grammar, "grammar", "g", "", "explicit path to grammar dylib, overrides search")
	cmd.Flags().StringVarP(&mode, "mode", "m", "auto", "language filter: auto | objj | js | both")
	cmd.Flags().BoolVarP(&profile, "profile", "p", false, "print per-file timing and aggregate summary")
	cmd.Flags().IntVarP(&context, "context", "c", 3, "source lines of context around error")
	cmd.Flags().BoolVar(&noXcode, "no-xcode", false, "suppress xed invocation")

	return cmd
}
