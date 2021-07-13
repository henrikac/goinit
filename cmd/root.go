package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "goinit",
	Short: "A generator for new Go projects",
}

func Execute() error {
	return rootCmd.Execute()
}
