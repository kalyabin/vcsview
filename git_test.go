package vcsview

import (
	"strings"
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

		e := g.ReadBranches(testCase, result)

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()

			loop: for {
				select {
				case <-e.ctx.Done():
					close(result)
					break loop
				case <- result:
					gotBranches++
				}
			}
		}()

		go func() {
			defer wg.Done()

			gotError = e.Run() != nil
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

	e := g.ReadBranches(projectPath, result)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		loop: for {
			select {
			case <- e.ctx.Done():
				close(result)
				break loop
			case branch := <- result:
				branches = append(branches, branch)
			}
		}
	}()

	var err error
	go func() {
		defer wg.Done()

		err = e.Run()
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

func TestGit_ReadCommitFail(t *testing.T) {
	g := MakeGitMock(t)

	cases := []struct{
		repoPath string
		commitId string
	}{
		{hgRepositoryPath, "xxx"},
		{gitRepositoryPath, "xxx"},
		{hgRepositoryPath, gitReadCommitTestCase.commitId},
	}

	for key, testCase := range cases {
		var (
			result chan Commit
			gotError bool
			gotCommit int
		)

		result = make(chan Commit)

		e := g.ReadCommit(testCase.repoPath, testCase.commitId, result)

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()

			loop: for {
				select {
				case <-e.ctx.Done():
					close(result)
					break loop
				case <-result:
					gotCommit++
				}
			}
		}()

		go func() {
			defer wg.Done()

			gotError = e.Run() != nil
		}()

		wg.Wait()

		if !gotError {
			t.Errorf("[%d] Git.ReadCommit(%s, %s, ...) has no errors, want error", key, testCase.repoPath, testCase.commitId)
		}

		if gotCommit > 0 {
			t.Errorf("[%d] Git.ReadCommit(%s, %s, ...) got %v commits, want: 0", key, testCase.repoPath, testCase.commitId, gotCommit)
		}
	}
}

func TestGit_ReadCommitOk(t *testing.T) {
	g := MakeGitMock(t)

	testCase := gitReadCommitTestCase

	var err error

	commit := make([]Commit, 0)
	result := make(chan Commit)

	e := g.ReadCommit(testCase.repoPath, testCase.commitId, result)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		loop: for {
			select {
			case <-e.ctx.Done():
				close(result)
				break loop
			case c := <- result:
				commit = append(commit, c)
			}
		}
	}()

	go func() {
		defer wg.Done()

		err = e.Run()
	}()

	wg.Wait()

	if err != nil {
		t.Fatalf("Git.ReadCommit(%s, %s, ...) got error: %v, want no errors", testCase.repoPath, testCase.commitId, err)
	}

	if len(commit) != 1 {
		t.Fatalf("Git.ReadCommit(%s, %s, ...) got %d commits, want: 1", testCase.repoPath, testCase.commitId, len(commit))
	} else {
		e := testCase.commit
		c := commit[0]

		if id, eId := c.Id(), e.Id(); id != eId {
			t.Fatalf("Git.ReadCommit(%s, %s, ...) commitId = %s, want: %s", testCase.repoPath, testCase.commitId, id, eId)
		}

		if date, eDate := c.Date().Format(gitLogDateLayout), e.Date().Format(gitLogDateLayout); date != eDate {
			t.Fatalf("Git.ReadCommit(%s, %s, ...) date = %s, want: %s", testCase.repoPath, testCase.commitId, date, eDate)
		}

		if author, eAuthor := c.Author().String(), e.Author().String(); author != eAuthor {
			t.Fatalf("Git.ReadCommit(%s, %s, ...) author = %s, want: %s", testCase.repoPath, testCase.commitId, author, eAuthor)
		}

		if parents, eParents := strings.Join(c.Parents(), " "), strings.Join(e.parents, " "); parents != eParents {
			t.Fatalf("Git.ReadCommit(%s, %s, ...) parents = %s, want: %s", testCase.repoPath, testCase.commitId, parents, eParents)
		}

		if message, eMessage := c.Message(), e.Message(); message != eMessage {
			t.Fatalf("Git.ReadCommit(%s, %s, ...) message = %s, want: %s", testCase.repoPath, testCase.commitId, message, eMessage)
		}
	}
}

func TestGit_ReadHistoryFail(t *testing.T) {
	g := MakeGitMock(t)

	cases := []struct{
		repoPath string
		path string
		branch string
		offset int
		limit int
	}{
		{noRepositoryPath, "", "", 0, 10},
	}

	for key, testCase := range cases {
		var (
			result chan Commit
			gotError bool
			gotCommit int
		)

		result = make(chan Commit)

		e := g.ReadHistory(testCase.repoPath, testCase.path, testCase.branch, testCase.offset, testCase.limit, result)

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()

			loop: for {
				select {
				case <- e.ctx.Done():
					close(result)
					break loop
				case <- result:
					gotCommit++
				}
			}
		}()

		go func() {
			defer wg.Done()

			gotError = e.Run() != nil
		}()

		wg.Wait()

		if !gotError {
			t.Errorf("[%d] Git.ReadHistory(%v) has no errors, want error", key, testCase)
		}

		if gotCommit > 0 {
			t.Errorf("[%d] Git.ReadHistory(%v) got %v commits, want: 0", key, testCase, gotCommit)
		}
	}
}

func TestGit_ReadHistoryOk(t *testing.T) {
	g := MakeGitMock(t)

	cases := []struct{
		repoPath string
		path string
		branch string
		offset int
		limit int
	}{
		{gitRepositoryPath, "", "", 0, 10},
		{gitRepositoryPath, "testpath", "", 0, 1},
		{gitRepositoryPath, "testpath", "", 0, 1},
		{gitRepositoryPath, "", "master", 0, 1},
		// @todo fix travis testing for other branches
		//{gitRepositoryPath, "", expectedGitBranches[0], 0, 1},
		//{gitRepositoryPath, "", expectedGitBranches[1], 0, 1},
		//{gitRepositoryPath, "", expectedGitBranches[2], 0, 1},
		//{gitRepositoryPath, "", expectedGitBranches[3], 0, 1},
		{gitRepositoryPath, "", "", 2, 2},
	}

	for key, testCase := range cases {
		var (
			gotError error
			gotCommit int
		)

		result := make(chan Commit, 1)
		e := g.ReadHistory(testCase.repoPath, testCase.path, testCase.branch, testCase.offset, testCase.limit, result)
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()

			loop: for {
				select {
				case <-e.ctx.Done():
					close(result)
					break loop
				case commit := <- result:
					gotCommit++

					if commit.Id() == "" {
						t.Fatalf("[%d] Git.ReadHistory(%v) commit has empty identifier", key, testCase)
					}
					if len(commit.Parents()) == 0 {
						t.Fatalf("[%d] Git.ReadHistory(%v) commit has empty parents", key, testCase)
					}
					if commit.Message() == "" {
						t.Fatalf("[%d] Git.ReadHistory(%v) commit has empty message", key, testCase)
					}
					if commit.Author().String() == "" {
						t.Fatalf("[%d] Git.ReadHistory(%v) commit has empty author", key, testCase)
					}
					if commit.Date().Unix() < 0 {
						t.Fatalf("[%d] Git.ReadHistory(%v) commit has empty date", key, testCase)
					}
				}
			}
		}()

		gotError = e.Run()

		wg.Wait()

		if gotError != nil {
			t.Fatalf("[%d] Git.ReadHistory(%v) has error: %v, want no errors", key, testCase, gotError)
		}

		if gotCommit != testCase.limit {
			t.Fatalf("[%d] Git.ReadHistory(%v) got %v commits, want: %v", key, testCase, gotCommit, testCase.limit)
		}
	}
}
