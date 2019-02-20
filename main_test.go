package vcsview

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

const (
	gitRepositoryPath = "testdata"+string(os.PathSeparator)+"git"
	hgRepositoryPath = "testdata"+string(os.PathSeparator)+"hg"
)

// Check git testing repository
// panic if repository not found
func checkGitTestingRepo() error {
	notFoundErr := func(path string) error {
		return fmt.Errorf("Git repository path not found: %s\nPlease, install repository using test_install.sh script.\n", path)
	}

	if stat, err := os.Stat(gitRepositoryPath); err != nil || !stat.IsDir() {
		return notFoundErr(gitRepositoryPath)
	}

	gitPath := gitRepositoryPath+string(os.PathSeparator)+".git"
	if stat, err := os.Stat(gitPath); err != nil || !stat.IsDir() {
		return notFoundErr(gitPath)
	}

	cmd := exec.Command("git", "status")
	cmd.Dir = gitRepositoryPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Git repository not found at %s: %v.\nPlease, re-install repository using test_install.sh script.\n", gitRepositoryPath, err)
	}

	return nil
}

func TestMain(m *testing.M) {
	if err := checkGitTestingRepo(); err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}