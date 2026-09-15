// Package commands provides the Cobra root command and subcommand registration.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root command.
var RootCmd = &cobra.Command{
	Use:   "linmo",
	Short: "Linmo CLI",
	Long:  "Linmo - AI agent gateway and CLI.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "linmo: %v\n", err)
		os.Exit(1)
	}
}
