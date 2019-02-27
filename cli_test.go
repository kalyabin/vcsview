package vcsview

import (
	"strings"
	"testing"
)

func TestCli_Command(t *testing.T) {
	c := Cli{}
	c.cmd = "git"

	wantArgs := "git --version"
	wantDir := gitRepoRealPath

	cmd := c.command(gitRepoRealPath, "--version")

	if dir := cmd.Dir; dir != wantDir {
		t.Errorf("Cli.CreateCommand(%s, %s).Dir = %v, want: %v", gitRepoRealPath, "--version", dir, wantDir)
	}

	if args := strings.Join(cmd.Args, " "); args != wantArgs {
		t.Errorf("Cli.CreateCommand(%s, %s).Args = %v, want: %v", gitRepoRealPath, "--version", args, wantArgs)
	}
}
