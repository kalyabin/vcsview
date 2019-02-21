package vcsview

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

const (
	noRepositoryPath = "testdata"
	gitRepositoryPath = "testdata"+string(os.PathSeparator)+"git"
	hgRepositoryPath = "testdata"+string(os.PathSeparator)+"hg"
)

// Check testing repository
// repoType - string identifier of repository type (hg, git)
// repoPath - path of repository (not a .git or .hg path)
// returns error if directory or repository not found or not initialized
func checkRepo(repoType string, repoPath string) error {
	repoTypeName := ""
	vcsPath := repoPath
	cmdName := ""

	switch repoType {
	case "hg":
		repoTypeName = "Mercurial"
		vcsPath += string(os.PathSeparator) + ".hg"
		cmdName = "hg"
	case "git":
		repoTypeName = "GIT"
		vcsPath += string(os.PathSeparator) + ".git"
		cmdName = "git"
	}

	notFoundErr := func(path string) error {
		return fmt.Errorf("%s repository path not found: %s\nPlease, install repository using test_install.sh script.\n", repoTypeName, path)
	}

	if stat, err := os.Stat(repoPath); err != nil || !stat.IsDir() {
		return notFoundErr(repoPath)
	}

	if stat, err := os.Stat(vcsPath); err != nil || !stat.IsDir() {
		return notFoundErr(vcsPath)
	}

	cmd := exec.Command(cmdName, "status")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s repository not found at %s: %v.\nPlease, re-install repository using test_install.sh script.\n", repoTypeName, repoPath, err)
	}

	return nil
}

func TestMain(m *testing.M) {
	if err := checkRepo("git", gitRepositoryPath); err != nil {
		panic(err)
	}
	if err := checkRepo("hg", hgRepositoryPath); err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}