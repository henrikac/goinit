package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/henrikac/goinit/cmd"
)

func TestNewCmd(t *testing.T) {
	root := cmd.NewRootCmd()
	newCmd := cmd.NewNewCmd()
	root.AddCommand(newCmd)

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INIT_PATH", pwd)
	if err != nil {
		t.Fatal(err)
	}

	buff := new(bytes.Buffer)

	root.SetOut(buff)
	root.SetArgs([]string{"new", "example"})
	err = root.Execute()
	if err != nil {
		t.Error("Unexpected error")
	}

	filenames := []string{"README.md", "LICENSE", ".gitignore", ".git", "go.mod"}

	dir := filepath.Join(pwd, "example")
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if len(files) != len(filenames) {
		t.Errorf("Expected a new project to contain %d files\nFound: %d\n", len(filenames), len(files))
	}

	for _, filename := range filenames {
		found := false
		for _, file := range files {
			if filename == file.Name() {
				found = true
			}
		}
		if !found {
			t.Errorf("Expected %s to be in a new project\n", filename)
		}
	}
}
