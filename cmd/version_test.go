package cmd_test

import (
	"bytes"

	"github.com/henrikac/goinit/cmd"
)

func ExampleVersion() {
	root := cmd.NewRootCmd()
	version := cmd.NewVersionCmd()
	root.AddCommand(version)

	buff := bytes.NewBufferString("")

	root.SetOut(buff)
	root.SetArgs([]string{"version"})
	root.Execute()

	// Output:
	// v0.2.0
}
