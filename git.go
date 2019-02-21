package vcsview

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

var (
	versionPattern = regexp.MustCompile(`([\d]+\.?([\d]+)?\.([\d]+)?)`)
)

type Git struct {
	Cli
}

// add specific params to execution
func (g *Git) Execute(dir string, out io.Writer, params ...string) error {
	params = append([]string{"--no-pager"}, params...)
	return g.Cli.Execute(dir, out, params...)
}

// Returns repository settings pathname
// like .git, .hg, etc.
func (g Git) RepositoryPathname() string {
	return ".git"
}

// Check Git version
// returns error if git command not found, or it hasn't version arguments
func (g Git) Version() (string, error) {
	buf := new(bytes.Buffer)

	if err := g.Execute(".", buf, "--version"); err != nil {
		return "", err
	}

	return versionPattern.FindString(buf.String()), nil
}

// Check project repository
// ProjectPath is absolute path to project path
// Returns error if repository not found at provided projectPath
// Returns nil if repository found
func (g Git) CheckRepository(projectPath string) error {
	repoPath := projectPath+string(os.PathSeparator)+g.RepositoryPathname()

	stats, err := os.Stat(repoPath)

	if err != nil {
		return err
	}

	if !stats.IsDir() {
		return fmt.Errorf("Git repository not found here: %s", projectPath)
	}

	return nil
}