package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = NewVersionCmd()

// NewVersionCmd initializes a new version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of goinit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.2.0")
		},
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
