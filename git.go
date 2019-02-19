package vcsview

import (
	"bytes"
	"regexp"
)

var (
	versionPattern = regexp.MustCompile(`([\d]+\.?([\d]+)?\.([\d]+)?)`)
)

type Git struct {
	Cli
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
