package cmd

import "github.com/spf13/cobra"

var rootCmd = NewRootCmd()

// NewRootCmd initializes a new root command.
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "goinit",
		Short: "A generator for new Go projects",
	}
}

// Execute the root cmd.
func Execute() error {
	return rootCmd.Execute()
}
