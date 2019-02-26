package vcsview

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	gitRepoRealPath, _ = filepath.Abs(gitRepositoryPath)
	hgRepoRealPath, _ = filepath.Abs(hgRepositoryPath)
	noRepoRealPath, _ = filepath.Abs(noRepositoryPath)
)

func TestRepository_Cmd(t *testing.T) {
	r := Repository{}

	if cmd := r.Cmd(); cmd != nil {
		t.Errorf("Repository.Cmd() = %v, want nil", cmd)
	}

	r.cmd = Git{}
	cmd := r.Cmd()
	_, ok := cmd.(Git)

	if !ok {
		t.Errorf("Repository.Cmd() = %v, want Git{}", cmd)
	}
}

func TestNewRepository(t *testing.T) {
	git := MakeGitMock(t)

	cases := []struct{
		projectPath string
		vcs Vcs
		wantFound bool
	}{
		{gitRepositoryPath, git, true},
		{gitRepoRealPath, git, true},
		{"."+pathSeparator+gitRepositoryPath, git, true},
		{hgRepositoryPath, git, false},
		{hgRepoRealPath, git, false},
		{noRepositoryPath, git, false},
		{noRepoRealPath, git, false},
	}

	for key, testCase := range cases {
		repo, err := NewRepository(testCase.projectPath, testCase.vcs)

		if testCase.wantFound && err != nil {
			t.Errorf("[%d] NewRepository(%s, %v) got error: %v, want nil", key, testCase.projectPath, err, testCase.vcs)
			continue
		}

		if !testCase.wantFound && err == nil {
			t.Errorf("[%d] NewRepository(%s, %v) got no errors, want error", key, testCase.projectPath, testCase.vcs)
		}

		if testCase.wantFound && (Repository{}) == repo {
			t.Errorf("[%d] NewRepository(%s, %v) got empty repository, want repository exists", key, testCase.projectPath, testCase.vcs)
			continue
		}

		if !testCase.wantFound && (Repository{}) != repo {
			t.Errorf("[%d] NewRepository(%s, %v) got: %v, want empty repository", key, testCase.projectPath, testCase.vcs, repo)
			continue
		}

		if testCase.wantFound {
			// check repository object
			realPath, _ := filepath.Abs(filepath.Clean(testCase.projectPath))

			if p := repo.ProjectPath(); p != realPath {
				t.Errorf("[%d] Repository.ProjectPath() = %v, want: %v", key, p, realPath)
			}

			repoPath, _ := filepath.Abs(filepath.Clean(testCase.projectPath+pathSeparator+testCase.vcs.RepositoryPathname()))
			if p := repo.RepositoryPath(); p != repoPath {
				t.Errorf("[%d] Repository.RepositoryPath = %v, want: %v", key, p, repoPath)
			}
		}
	}
}

func TestReposotyr_AbsPath(t *testing.T) {
	gitRepoRealPath, _ := filepath.Abs(gitRepositoryPath)

	cases := []struct{
		projectPath string
		subDir string
		expectedAbsPath string
		gotError bool
	}{
		{gitRepoRealPath, "testpath/empty.txt", gitRepoRealPath+pathSeparator+"testpath/empty.txt", false},
		{gitRepoRealPath, "testpath/../testpath", gitRepoRealPath+pathSeparator+"testpath", false},
		{gitRepoRealPath, "testpath/../../git", gitRepoRealPath, false},
		{gitRepoRealPath, "testpath/../../../", gitRepoRealPath, true},
	}

	for key, testCase := range cases {
		r := Repository{}
		r.projectPath = testCase.projectPath

		result, err := r.AbsPath(testCase.subDir)
		if err != nil && !testCase.gotError {
			t.Errorf("[%d] Repository.AbsPath(%s) got error: %v, want no errors", key, testCase.subDir, err)
			continue
		}

		if err == nil && testCase.gotError {
			t.Errorf("[%d] Repository.AbsPath(%s) got no errors, want errors", key, testCase.subDir)
			continue
		}

		if !testCase.gotError && result != testCase.expectedAbsPath {
			t.Errorf("[%d] Repository.AbsPath(%s) = %v, want: %v", key, testCase.subDir, result, testCase.expectedAbsPath)
			continue
		}
	}
}

func TestRepository_FilesList(t *testing.T) {
	git := MakeGitMock(t)

	cases := []struct{
		projectPath string
		subDir string
		expectedSubDir string
		vcs Vcs
		wantFound bool
	}{
		{gitRepositoryPath, "", "", git, true},
		{gitRepositoryPath, "testpath", "testpath", git, true},
		{gitRepositoryPath, "../git", "", git, true},
		{gitRepositoryPath, "..//git", "", git, true},
		{gitRepositoryPath, "..//..//testdata//git", "", git, true},
		{gitRepositoryPath, "../../testdata", "", git, false},
		{gitRepositoryPath, "../../.git", "", git, false},
	}

	for key, testCase := range cases {
		r, err := NewRepository(testCase.projectPath, testCase.vcs)

		if err != nil {
			t.Fatalf("Can't create repository for %s. Got error: %v", testCase.projectPath, err)
			break
		}

		files, err := r.FilesList(testCase.subDir)

		if testCase.wantFound && err != nil {
			t.Errorf("[%d] Repository.FilesList(%s) got error: %v, want no errors", key, testCase.subDir, err)
			continue
		}

		if !testCase.wantFound && err == nil {
			t.Errorf("[%d] Repository.FilesList(%s) want error, no errors got", key, testCase.subDir)
			continue
		}

		if testCase.wantFound && len(files) == 0 {
			t.Errorf("[%d] Repository.FilesList(%s) got empty files list", key, testCase.subDir)

			continue
		}

		path := filepath.Clean(testCase.projectPath+pathSeparator+testCase.subDir)
		path, _ = filepath.Abs(path)

		for _, f := range files {
			expectedPathname := strings.TrimLeft(testCase.expectedSubDir+pathSeparator+f.name, pathSeparator)

			if pathname := f.Pathname(); expectedPathname != pathname {
				t.Errorf("[%d] Unexpected pathname for file in %s. Want: %s, got: %s", key, path+pathSeparator+f.name, expectedPathname, pathname)
				continue
			}
		}
	}

}

func TestRepository_Check(t *testing.T) {
	git := MakeGitMock(t)

	cases := []struct{
		projectPath string
		vcs Vcs
		gotError bool
	}{
		{gitRepoRealPath, git, false},
		{hgRepoRealPath, git, true},
		{noRepoRealPath, git, true},
	}

	for key, testCase := range cases {
		r := Repository{}
		r.cmd = testCase.vcs
		r.projectPath = testCase.projectPath

		err := r.Check()

		if err != nil && !testCase.gotError {
			t.Errorf("[%d] Repository.Check() = %v, want no errors", key, err)
		}

		if err == nil && testCase.gotError {
			t.Errorf("[%d] Repository.Check() = nil, want errors", key)
		}
	}
}

func TestRepository_GetBranchesFail(t *testing.T) {
	git := MakeGitMock(t)

	cases := []struct{
		vcs Vcs
		path string
	}{
		{git, hgRepositoryPath},
		{git, noRepositoryPath},
	}

	for key, testCase := range cases {
		r := Repository{}
		r.projectPath = testCase.path
		r.cmd = testCase.vcs

		var (
			result chan Branch
			err chan error
			gotError bool
			gotBranches int
		)

		result = make(chan Branch)
		err = make(chan error)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			for range err {
				gotError = true
			}
		}()

		go func() {
			defer wg.Done()

			for range result {
				gotBranches++
			}
		}()

		r.GetBranches(result, err)

		wg.Wait()

		if !gotError {
			t.Errorf("[%d] Repository.GetBranches(...) for %s has no errors, want error", key, testCase)
		}

		if gotBranches > 0 {
			t.Errorf("[%d] Repository.GetBranches(...) for %s got %v branches, want: 0", key, testCase, gotBranches)
		}
	}
}

func TestRepository_GetBranchesOk(t *testing.T) {
	git := MakeGitMock(t)

	cases := []struct{
		vcs Vcs
		path string
		expectedBranches []string
	}{
		{git, gitRepositoryPath, expectedGitBranches},
	}

	for key, testCase := range cases {
		r := Repository{}
		r.projectPath = testCase.path
		r.cmd = testCase.vcs

		var gotError error = nil
		branches := make([]Branch, 0)
		result := make(chan Branch)
		err := make(chan error)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			for e := range err {
				gotError = e
			}
		}()

		go func() {
			defer wg.Done()

			for branch := range result {
				branches = append(branches, branch)
			}
		}()

		r.GetBranches(result, err)

		wg.Wait()

		if gotError != nil {
			t.Errorf("Repository.GetBranches(...) = %v, %v, want no errors", branches, gotError)
		}

		if len(branches) != len(testCase.expectedBranches) {
			t.Errorf("Repository.GetBranches(...) = %v, %v, want %d branches", branches, gotError, len(expectedGitBranches))
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
			for _, expectedBranch := range testCase.expectedBranches {
				if branch.Id() == expectedBranch {
					gotBranches++
				}
			}
		}

		if gotBranches != len(testCase.expectedBranches) {
			t.Errorf("[%d] Repository.GetBranches(...) doesnt contain expected branches: %v", key, testCase.expectedBranches)
		}

		if !gotCurrent {
			t.Errorf("[%d] Repository.GetBranches(...) doesnt contain current branch", key)
		}
	}
}
