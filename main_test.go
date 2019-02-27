package vcsview

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	noRepositoryPath = "testdata"
	gitRepositoryPath = "testdata"+pathSeparator+"git"
	hgRepositoryPath = "testdata"+pathSeparator+"hg"
)

var (
	expectedGitBranches = []string{"master", "remotes/origin/branch1", "remotes/origin/branch2", "remotes/origin/master"}

	gitReadCommitTestCase = struct{
		repoPath string
		commitId string
		commit Commit
	}{
		repoPath: gitRepositoryPath,
		commitId: "747ad5712f0ddbb482ebb6e07eb779e70b94687f",
		commit: Commit{
			id: "747ad5712f0ddbb482ebb6e07eb779e70b94687f",
			date: time.Date(2016, time.Month(10), 4, 22, 7, 27, 0, time.FixedZone("UTC+3", 3*60*60)),
			author: Contributor{
				name: "Max Kalyabin",
				email: "maksim@kalyabin.ru",
			},
			parents: []string{"81cb0276737ca3345faaaec5a5df2a3e1ff5d775"},
			message: "random commit for random file",
		},
	}
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
		vcsPath += pathSeparator + ".hg"
		cmdName = "hg"
	case "git":
		repoTypeName = "GIT"
		vcsPath += pathSeparator + ".git"
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