package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

//go:embed license-templates
var licenses embed.FS

var licensesCmd = &cobra.Command{
	Use:   "licenses [command]",
	Short: "Get information about available licenses",
	Args:  cobra.MaximumNArgs(1),
}

var licensesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available licenses",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		names, err := getLicenseNames()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, name := range names {
			fmt.Println(name)
		}
	},
}

var licensesReadCmd = &cobra.Command{
	Use:   "read [license name]",
	Short: "Prints a license",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := filepath.Join("license-templates", fmt.Sprintf("%s.txt", strings.ToLower(args[0])))
		file, err := licenses.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unknown license: %s\n", strings.ToLower(args[0]))
			fmt.Println("try: goinit licenses list")
			os.Exit(1)
		}
		if bytes.Contains(file, []byte("%s")) {
			ui, err := getUserInfo()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf(string(file), time.Now().Year(), ui.name)
		} else {
			fmt.Printf(string(file))
		}
	},
}

func init() {
	licensesCmd.AddCommand(licensesListCmd)
	licensesCmd.AddCommand(licensesReadCmd)
	rootCmd.AddCommand(licensesCmd)
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
