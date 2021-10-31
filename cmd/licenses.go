package cmd

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

//go:embed license-templates
var licenses embed.FS

var (
	licensesCmd     = NewLicensesCmd()
	licensesListCmd = NewLicensesListCmd()
	licensesReadCmd = NewLicensesReadCmd()
)

// NewLicensesCmd initializes a new licenses command.
func NewLicensesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "licenses [command]",
		Short: "Get information about available licenses",
		Args:  cobra.MaximumNArgs(1),
	}
}

// NewLicensesListCmd initializes a new license list command.
func NewLicensesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all available licenses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := getLicenseNames()
			if err != nil {
				return errors.New("unable to read licenses")
			}
			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}
}

// NewLicensesReadCmd initializes a new license read command.
func NewLicensesReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "read [license name]",
		Short: "Prints a license",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := filepath.Join("license-templates", fmt.Sprintf("%s.txt", strings.ToLower(args[0])))
			file, err := licenses.ReadFile(filepath.ToSlash(filename))
			if err != nil {
				errMsg := fmt.Sprintf("unknown license: \"%s\"\n", strings.ToLower(args[0]))
				return errors.New(errMsg)
			}
			if needsFormatting(file) {
				ui, err := getUserInfo()
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), string(file), time.Now().Year(), ui.name)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), string(file))
			}
			return nil
		},
	}
}

func init() {
	licensesCmd.AddCommand(licensesListCmd)
	licensesCmd.AddCommand(licensesReadCmd)
	rootCmd.AddCommand(licensesCmd)
}

func needsFormatting(b []byte) bool {
	return bytes.Contains(b, []byte("%s"))
}

func getLicenseNames() ([]string, error) {
	var licenseNames []string
	files, err := licenses.ReadDir("license-templates")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		split := strings.Split(file.Name(), ".")
		name := strings.ToUpper(split[0])
		if name == "APACHE" {
			name = strings.Title(strings.ToLower(name))
		}
		licenseNames = append(licenseNames, name)
	}
	return licenseNames, nil
}
