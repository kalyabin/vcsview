package vcsview

import (
	"path/filepath"
	"strings"
	"testing"
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

	gitRepoRealPath, _ := filepath.Abs(gitRepositoryPath)
	hgRepoRealPath, _ := filepath.Abs(hgRepositoryPath)
	noRepoRealPath, _ := filepath.Abs(noRepositoryPath)

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

			if repo.ProjectPath != realPath {
				t.Errorf("[%d] Repository.ProjectPath = %v, want: %v", key, repo.ProjectPath, realPath)
			}

			repoPath, _ := filepath.Abs(filepath.Clean(testCase.projectPath+pathSeparator+testCase.vcs.RepositoryPathname()))
			if repo.RepositoryPath != repoPath {
				t.Errorf("[%d] Repository.RepositoryPath = %v, want: %v", key, repo.RepositoryPath, repoPath)
			}
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
