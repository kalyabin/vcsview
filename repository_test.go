package vcsview

import (
	"os"
	"testing"
)

func TestNewRepository(t *testing.T) {
	git := MakeGitMock(t)

	cases := []struct{
		projectPath string
		vcs Vcs
		wantFound bool
	}{
		{gitRepositoryPath, git, true},
		{hgRepositoryPath, git, false},
		{noRepositoryPath, git, false},
	}

	for key, testCase := range cases {
		repo, err := NewRepository(testCase.projectPath, testCase.vcs)

		if testCase.wantFound && (Repository{}) == repo {
			t.Errorf("[%d] NewRepository(%s, %v) got empty repository, want repository exists", key, testCase.projectPath, testCase.vcs)
			continue
		}

		if testCase.wantFound && err != nil {
			t.Errorf("[%d] NewRepository(%s, %v) got error, want nil", key, testCase.projectPath, testCase.vcs)
			continue
		}

		if !testCase.wantFound && (Repository{}) != repo {
			t.Errorf("[%d] NewRepository(%s, %v) got: %v, want empty repository", key, testCase.projectPath, testCase.vcs, repo)
			continue
		}

		if !testCase.wantFound && err == nil {
			t.Errorf("[%d] NewRepository(%s, %v) got no errors, want error", key, testCase.projectPath, testCase.vcs)
		}

		if testCase.wantFound {
			// check repository object
			if repo.ProjectPath != testCase.projectPath {
				t.Errorf("[%d] Repository.ProjectPath = %v, want: %v", key, repo.ProjectPath, testCase.projectPath)
			}

			repoPath := testCase.projectPath+string(os.PathSeparator)+testCase.vcs.RepositoryPathname()
			if repo.RepositoryPath != repoPath {
				t.Errorf("[%d] Repository.RepositoryPath = %v, want: %v", key, repo.RepositoryPath, repoPath)
			}
		}
	}
}
