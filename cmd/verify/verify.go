/*
 * cmd/verify/verify.go
 * cappuccino
 *
 * Created by David Richardson on Friday, April 24, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

/*
 * Package verify implements the verify subcommand.
 *
 * Checks whether capp-parse prerequisites are in place and reports their
 * status. When stdout is not a TTY, output is JSON for machine consumption.
 *
 * On failure, the full ordered list of searched paths is reported so that
 * the operator — or the orchestration application — can diagnose the missing
 * asset without consulting documentation.
 */
package verify

import (
	"encoding/json"
	"fmt"
	"os"

	"dev.cappuccino/capp-parse/grammar"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewVerifyCmd constructs the verify command.
func NewVerifyCmd() *cobra.Command {
	var grammarFlag string

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify capp-parse prerequisites",
		Long: `Check whether capp-parse prerequisites are installed and locatable.

Exits 0 if all prerequisites are satisfied, 1 otherwise.

When stdout is not a terminal, output is JSON:
  {"ok": true,  "grammar": {"found": true,  "path": "/usr/local/lib/libtree-sitter-objj.dylib"}}
  {"ok": false, "grammar": {"found": false, "searched": [...]}}`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			isTTY := term.IsTerminal(int(os.Stdout.Fd()))
			return runVerify(cmd, grammarFlag, isTTY)
		},
	}

	cmd.Flags().StringVarP(&grammarFlag, "grammar", "g", "", "explicit path to grammar dylib, overrides search")

	return cmd
}

func runVerify(cmd *cobra.Command, grammarFlag string, isTTY bool) error {
	searched := grammar.SearchPaths()
	found := grammar.FindLibrary()
	if grammarFlag != "" {
		if _, err := os.Stat(grammarFlag); err == nil {
			found = grammarFlag
		} else {
			found = ""
		}
		searched = []string{grammarFlag}
	}

	grammarOK := found != ""

	if isTTY {
		return emitTTY(cmd, grammarOK, found, searched)
	}
	return emitJSON(cmd, grammarOK, found, searched)
}

func emitTTY(cmd *cobra.Command, grammarOK bool, found string, searched []string) error {
	out := cmd.OutOrStdout()

	if grammarOK {
		fmt.Fprintf(out, "✓  grammar library: %s\n", found)
	} else {
		fmt.Fprintf(out, "✗  grammar library not found\n")
		fmt.Fprintf(out, "   searched:\n")
		for _, p := range searched {
			fmt.Fprintf(out, "     %s\n", p)
		}
		fmt.Fprintf(out, "   The grammar library is managed by the Cappuccino orchestration\n")
		fmt.Fprintf(out, "   application, or may be placed manually in any of the paths above.\n")
		fmt.Fprintf(out, "   Use --grammar to specify an explicit path.\n")
	}

	if !grammarOK {
		return fmt.Errorf("one or more prerequisites missing")
	}
	return nil
}

func emitJSON(cmd *cobra.Command, grammarOK bool, found string, searched []string) error {
	type grammarResult struct {
		Found    bool     `json:"found"`
		Path     string   `json:"path,omitempty"`
		Searched []string `json:"searched,omitempty"`
	}
	type result struct {
		OK      bool          `json:"ok"`
		Grammar grammarResult `json:"grammar"`
	}

	gr := grammarResult{Found: grammarOK}
	if grammarOK {
		gr.Path = found
	} else {
		gr.Searched = searched
	}

	line, _ := json.Marshal(result{OK: grammarOK, Grammar: gr})
	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", line)

	if !grammarOK {
		return fmt.Errorf("one or more prerequisites missing")
	}
	return nil
}
