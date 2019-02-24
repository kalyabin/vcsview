package vcsview

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

type Git struct {
	Cli
}

// add specific params to execution
func (g *Git) execute(dir string, out io.Writer, params ...string) error {
	params = append([]string{"--no-pager"}, params...)
	return g.Cli.execute(dir, out, params...)
}

// Returns repository settings pathname
// like .git, .hg, etc.
func (g Git) RepositoryPathname() string {
	return ".git"
}

// Check Git version
// returns error if git command not found, or it hasn't version arguments
func (g Git) Version() (string, error) {
	versionPattern := regexp.MustCompile(`([\d]+\.?([\d]+)?\.([\d]+)?)`)

	buf := new(bytes.Buffer)

	if err := g.execute(".", buf, "--version"); err != nil {
		return "", err
	}

	return versionPattern.FindString(buf.String()), nil
}

// Check project repository
// projectPath is absolute path to project path
// Returns error if repository not found at provided projectPath
// Returns nil if repository found
func (g Git) CheckRepository(projectPath string) error {
	repoPath := projectPath+pathSeparator+g.RepositoryPathname()

	stats, err := os.Stat(repoPath)

	if err != nil {
		return err
	}

	if !stats.IsDir() {
		return fmt.Errorf("Git repository not found here: %s", projectPath)
	}

	return nil
}

// Check the repository status
// Throws an error if repository doesnt exists at the path
func (g Git) StatusRepository(projectPath string) (string, error) {
	buf := new(bytes.Buffer)

	if err := g.execute(projectPath, buf, "status", "--short"); err != nil {
		return "", err
	}

	return buf.String(), nil
}