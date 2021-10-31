package cmd_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/henrikac/goinit/cmd"
)

func ExampleNewLicensesListCmd() {
	root := cmd.NewRootCmd()
	licenses := cmd.NewLicensesCmd()
	list := cmd.NewLicensesListCmd()
	licenses.AddCommand(list)
	root.AddCommand(licenses)

	buff := new(bytes.Buffer)

	root.SetOut(buff)
	root.SetArgs([]string{"licenses", "list"})
	root.Execute()

	fmt.Println(buff.String())
	// Output:
	// Apache
	// BSD2
	// BSD3
	// GPL2
	// GPL3
	// MIT
}

func TestLicensesRead(t *testing.T) {
	root := cmd.NewRootCmd()
	licenses := cmd.NewLicensesCmd()
	read := cmd.NewLicensesReadCmd()
	licenses.AddCommand(read)
	root.AddCommand(licenses)

	tests := []struct {
		LicenseName   string
		ExpectedTitle string
	}{
		{LicenseName: "mit", ExpectedTitle: "The MIT License (MIT)"},
		{LicenseName: "gpl3", ExpectedTitle: "GNU GENERAL PUBLIC LICENSE"},
		{LicenseName: "bsd2", ExpectedTitle: "BSD 2-Clause License"},
		{LicenseName: "apache", ExpectedTitle: "Apache License"},
		{LicenseName: "gpl2", ExpectedTitle: "The GNU General Public License (GPL-2.0)"},
		{LicenseName: "bsd3", ExpectedTitle: "BSD 3-Clause License"},
	}

	for _, test := range tests {
		buff := new(bytes.Buffer)

		root.SetOut(buff)
		root.SetArgs([]string{"licenses", "read", test.LicenseName})
		root.Execute()

		// trimming spaces to the left because some of the licenses
		// has their title centered
		trimBuff := strings.TrimLeft(buff.String(), " ")

		if !strings.HasPrefix(trimBuff, test.ExpectedTitle) {
			t.Errorf("Expected: %s\nGot: %s\n", test.ExpectedTitle, trimBuff)
		}
	}
}

func TestLicensesReadUnknownLicense(t *testing.T) {
	root := cmd.NewRootCmd()
	licenses := cmd.NewLicensesCmd()
	read := cmd.NewLicensesReadCmd()
	licenses.AddCommand(read)
	root.AddCommand(licenses)

	tests := []string{"", "unknown-license", "bsd", "gpl"}

	for _, test := range tests {
		buff := new(bytes.Buffer)

		root.SetOut(buff)
		root.SetErr(buff)
		root.SetArgs([]string{"licenses", "read", test})

		err := root.Execute()
		if err == nil {
			t.Errorf("Expected \"%s\" to return an error\n", test)
		}

		expectedErr := fmt.Sprintf("Error: unknown license: \"%s\"\n\n%s\n",
			strings.ToLower(test), read.UsageString(),
		)

		if strings.Compare(buff.String(), expectedErr) != 0 {
			t.Errorf("Expected: %s\nGot: %s\n", expectedErr, buff.String())
		}
	}
}
