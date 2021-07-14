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

const license = `The MIT License (MIT)

Copyright (c) %d %s <%s>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.`

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
	name  string
	email string
}

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new Go project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location, err := projectLocation()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		projName := args[0]
		projFolder := filepath.Join(location, projName)
		if _, err := os.Stat(projFolder); !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "%s already exist\n", projName)
			os.Exit(1)
		}
		err = generateProject(projFolder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
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
	err = generateLicense(name, ui)
	if err != nil {
		return err
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
	return nil
}

func gitInit() error {
	err := exec.Command("git", "init").Run()
	if err != nil {
		return err
	}
	return nil
}

func generateReadme(path string, ui *userInfo) error {
	projName := filepath.Base(path)
	filename := filepath.Join(path, "README.md")
	err := os.WriteFile(filename, []byte(fmt.Sprintf(readme, projName, projName, ui.name)), 0666)
	if err != nil {
		return err
	}
	return nil
}

func generateLicense(path string, ui *userInfo) error {
	filename := filepath.Join(path, "LICENSE")
	err := os.WriteFile(filename, []byte(fmt.Sprintf(license, time.Now().Year(), ui.name, ui.email)), 0666)
	if err != nil {
		return err
	}
	return nil
}

func generateGitignore(path string) error {
	filename := filepath.Join(path, ".gitignore")
	err := os.WriteFile(filename, []byte{}, 0666)
	if err != nil {
		return err
	}
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
	email, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		return nil, err
	}
	ui := &userInfo{
		name:  strings.TrimSpace(string(name)),
		email: strings.TrimSpace(string(email)),
	}
	return ui, nil
}

func projectLocation() (string, error) {
	if path, set := os.LookupEnv("GO_INIT_PATH"); set {
		return path, nil
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
