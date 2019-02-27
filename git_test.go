package vcsview

import (
	"sync"
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

func TestGit_ReadBranchesFail(t *testing.T) {
	g := MakeGitMock(t)

	cases := []string{
		hgRepositoryPath,
		noRepositoryPath,
	}

	for key, testCase := range cases {
		var (
			result chan Branch
			gotError bool
			gotBranches int
		)

		result = make(chan Branch)

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			for range result {
				gotBranches++
			}
		}()

		go func() {
			defer wg.Done()

			gotError = g.ReadBranches(testCase, result).Start() != nil
		}()

		wg.Wait()

		if !gotError {
			t.Errorf("[%d] Git.ReadBranches(%s, ...) has no errors, want error", key, testCase)
		}

		if gotBranches > 0 {
			t.Errorf("[%d] Git.ReadBranches(%s, ...) got %v branches, want: 0", key, testCase, gotBranches)
		}
	}
}

func TestGit_ReadBranchesOk(t *testing.T) {
	g := MakeGitMock(t)

	branches := make([]Branch, 0)
	result := make(chan Branch)

	projectPath := gitRepositoryPath

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for branch := range result {
			branches = append(branches, branch)
		}
	}()

	var err error
	go func() {
		defer wg.Done()

		err = g.ReadBranches(projectPath, result).Start()
	}()

	wg.Wait()

	if err != nil {
		t.Errorf("Git.ReadBranches(%s) = %v, %v, want no errors", projectPath, branches, err)
	}

	if len(branches) != len(expectedGitBranches) {
		t.Errorf("Git.ReadBranches(%s) = %v, %v, want %d branches", projectPath, branches, err, len(expectedGitBranches))
	}

	gotBranches := 0
	gotCurrent := false

	for key, branch := range branches {
		if branch.Id() == "" {
			t.Errorf("Branch %d got empty identifier", key)
		}
		if branch.Head() == "" {
			t.Errorf("Branch %d got empty head commit", key)
		}
		if branch.IsCurrent() {
			gotCurrent = true
		}
		for _, expectedBranch := range expectedGitBranches {
			if branch.Id() == expectedBranch {
				gotBranches++
			}
		}
	}

	if gotBranches != len(expectedGitBranches) {
		t.Errorf("Git.ReadBranches(%s) doesnt contain expected branches: %v", projectPath, expectedGitBranches)
	}

	if !gotCurrent {
		t.Errorf("Git.ReadBranches(%s) doesnt contain current branch", projectPath)
	}
}