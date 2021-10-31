package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var license = "MIT"

const readme = `# %s

TODO: Write description here

## Usage

TODO: Write usage instructions here

## Contributing

1. Fork it (<https://github.com/your-github-user/%s/fork>)
2. Create your feature branch (` + "`git checkout -b my-new-feature`" + `)
3. Commit your changes (` + "`git commit -am 'Add some feature'`" + `)
4. Push to the branch (` + "`git push origin my-new-feature`" + `)
5. Create a new Pull Request

## Contributors

- [%s](https://github.com/your-github-user) - creator and maintainer
`

type userInfo struct {
	name string
}

var newCmd = NewNewCmd()

// NewNewCmd initializes a new new command.
func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project name]",
		Short: "Create a new Go project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			location, err := projectLocation()
			if err != nil {
				return err
			}
			projName := args[0]
			projFolder := filepath.Join(location, projName)
			if _, err := os.Stat(projFolder); !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("%s already exist", projName)
			}
			validLicense := isValidLicense()
			if !validLicense {
				return fmt.Errorf("unknown license: %s", strings.ToLower(license))
			}
			err = generateProject(projFolder)
			if err != nil {
				return fmt.Errorf("%s", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created %s\n", projName)
			return nil
		},
	}
	cmd.Flags().StringVarP(&license, "license", "l", "MIT", "Which license to add to project")
	return cmd
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func isValidLicense() bool {
	if license == "" {
		return true
	}
	filename := filepath.Join("license-templates", fmt.Sprintf("%s.txt", strings.ToLower(license)))
	_, err := licenses.ReadFile(filepath.ToSlash(filename))
	if err != nil {
		return false
	}
	return true
}

func generateProject(name string) error {
	err := os.Mkdir(name, 0777)
	if err != nil {
		return err
	}
	ui, err := getUserInfo()
	if err != nil {
		return err
	}
	err = generateReadme(name, ui)
	if err != nil {
		return err
	}
	if license != "" {
		err = generateLicense(name, ui)
		if err != nil {
			return err
		}
	}
	err = generateGitignore(name)
	if err != nil {
		return err
	}
	err = os.Chdir(name)
	if err != nil {
		return err
	}
	err = gitInit()
	if err != nil {
		return err
	}
	err = goModInit()
	if err != nil {
		return err
	}
	return nil
}

func goModInit() error {
	err := exec.Command("go", "mod", "init").Run()
	if err != nil {
		return err
	}
	fmt.Println("Initialized new Go module")
	return nil
}

func gitInit() error {
	err := exec.Command("git", "init").Run()
	if err != nil {
		return err
	}
	fmt.Println("Initialized empty Git repository")
	return nil
}

func generateReadme(path string, ui *userInfo) error {
	projName := filepath.Base(path)
	filename := filepath.Join(path, "README.md")
	err := os.WriteFile(filename, []byte(fmt.Sprintf(readme, projName, projName, ui.name)), 0666)
	if err != nil {
		return err
	}
	fmt.Printf("Created %s\n", filename)
	return nil
}

func generateLicense(path string, ui *userInfo) error {
	filename := filepath.Join(path, "LICENSE")
	licenseFilename := filepath.Join("license-templates", fmt.Sprintf("%s.txt", strings.ToLower(license)))
	file, err := licenses.ReadFile(filepath.ToSlash(licenseFilename))
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, []byte(fmt.Sprintf(string(file), time.Now().Year(), ui.name)), 0666)
	if err != nil {
		return err
	}
	fmt.Printf("Created %s\n", filename)
	return nil
}

func generateGitignore(path string) error {
	filename := filepath.Join(path, ".gitignore")
	err := os.WriteFile(filename, []byte{}, 0666)
	if err != nil {
		return err
	}
	fmt.Printf("Created %s\n", filename)
	return nil
}

func getUserInfo() (*userInfo, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	name, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		return nil, err
	}
	ui := &userInfo{
		name: strings.TrimSpace(string(name)),
	}
	return ui, nil
}

func projectLocation() (string, error) {
	if path, set := os.LookupEnv("GO_INIT_PATH"); set {
		if path != "" {
			return path, nil
		}
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		filename := filepath.Join(home, "go")
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(filename, 0777)
			if err != nil {
				return "", err
			}
		}
		gopath = filename
	}
	return gopath, nil
}
