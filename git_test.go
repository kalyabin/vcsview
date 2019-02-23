package vcsview

import (
	"testing"
)

func MakeGitMockWithCmd(cmd string, t *testing.T) Git {
	g := Git{
		Cli{
			cmd: cmd,
			Debugger: DebugFunc(func(msg string) {
				t.Log(msg)
			}),
		},
	}
	return g
}

func MakeGitMock(t *testing.T) Git {
	return MakeGitMockWithCmd("git", t)
}

func TestGit_RepositoryPathname(t *testing.T) {
	g := Git{}

	want := ".git"

	if result := g.RepositoryPathname(); result != want {
		t.Fatalf("Unexpected git repository pathname, want: %v, got: %v", want, result)
	}
}

func TestGit_Version(t *testing.T) {
	g := MakeGitMock(t)

	version, err := g.Version();

	if err != nil {
		t.Errorf("Unexpected git version error: %v", err)
	}

	if version == "" {
		t.Errorf("Git version is empty string.")
	}

	g = MakeGitMockWithCmd("non_git", t)

	version, err = g.Version()

	if err == nil {
		t.Errorf("Expected git version error, none given")
	}

	if version != "" {
		t.Errorf("Expected empty git version, non given")
	}
}

func TestGit_CheckRepository(t *testing.T) {
	g := MakeGitMock(t)

	cases := []struct{
		projectPath string
		wantFound bool
	}{
		{gitRepositoryPath, true},
		{hgRepositoryPath, false},
	}

	for key, testCase := range cases {
		err := g.CheckRepository(testCase.projectPath)

		if testCase.wantFound && err != nil {
			t.Errorf("[%d] Git.CheckRepository(%s) = %v, want: nil", key, testCase.projectPath, err)
		} else if !testCase.wantFound && err == nil {
			t.Errorf("[%d] Git.CheckRepository(%s) = nil, want: error", key, testCase.projectPath)
		}
	}
}
