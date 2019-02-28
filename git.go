package vcsview

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	gitLogFormat = "%H%n%P%n%an%n%ae%n%ad%n%s"
	gitLogDateLayout = "Mon Jan 2 15:04:05 2006 -0700"
)

// CLI wrapper for GIT
type Git struct {
	Cli
}

// add specific params to command
func (g *Git) createCommand(dir string, params ...string) *exec.Cmd {
	return g.Cli.command(dir, append([]string{"--no-pager"}, params...)...)
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

	var (
		result string
		done = make(chan interface{}, 1)
	)

	cmd := g.createCommand(".", "--version")
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		for s.Scan() {
			result += s.Text() + "\n"
		}

		done <- struct{}{}
	})

	e := g.executor(cmd, reader)

	err := e.Run()

	<- done

	close(done)

	return versionPattern.FindString(result), err
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
	var (
		result string
		done = make(chan interface{}, 1)
	)

	cmd := g.createCommand(projectPath, "status", "--short")
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		for s.Scan() {
			result += s.Text()+"\n"
		}

		done <- struct{}{}
	})

	e := g.executor(cmd, reader)

	err := e.Run()

	return result, err
}

// Fetch repository branches asynchronously
// ProjectPath is the absolute path to project with Git repository
func (g Git) ReadBranches(projectPath string, result chan Branch) *Executor {
	// pattern to read branches line by line
	p := regexp.MustCompile(`^\*?[\s+|\t]+(?P<id>[^\s]+)[\s+|\t]+(?P<head>[a-fA-F0-9]+)[\s+|\t]+(?P<message>.*)$`)

	cmd := g.createCommand(projectPath, "branch", "-a", "-v")
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		for s.Scan() {
			line := s.Bytes()

			if !p.Match(line) {
				continue
			}

			matches := p.FindSubmatch(line)

			isCurrent := string(line[:1]) == "*"
			id := string(matches[1])
			head := string(matches[2])

			result <- Branch{id, head, isCurrent}
		}
	})

	return g.executor(cmd, reader)
}

// Wrapper for read commits from command line stdout
// Commits will going by such lines:
// 313604a7f4ecd265e56102fa2e22de35726f4687 <--- Commit sha256
// 1e16e4aeeef941bd037ed5f70e9d2abcf459ca2e 313604a7f4ecd265e56102fa2ee2de3572df4687 <--- Parents commit sha256
// Max Kalyabin <--- Author name
// maksim@kalyabin.ru <--- Author email
// Wed Feb 27 14:51:45 2019 +0300 <--- Commit date and time
// read git commit <--- Commit message
func (g *Git) readCommitsPipe(s *bufio.Scanner, result chan Commit) {
	data := make([]string, 6)
	key := 0

	for s.Scan() {
		str := s.Text()

		data[key] = str
		key++

		if key == 6 {
			time, _ := time.Parse(gitLogDateLayout, data[4])

			commit := Commit{
				id: data[0],
				parents: strings.Split(data[1], " "),
				author: Contributor{
					name: data[2],
					email: data[3],
				},
				date: time,
				message: data[5],
			}

			result <- commit

			runtime.Gosched()

			data = make([]string, 6)
			key = 0
		}
	}
}

// Fetch repository commit by identifier asynchronously
// ProjectPath is the absolute path to project with Git repository
// CommitId is the sha256 commit identifier (or short copy)
func (g Git) ReadCommit(projectPath string, commitId string, result chan Commit) *Executor {
	cmd := g.createCommand(projectPath, "show", "--quiet", commitId, `--pretty=format:`+gitLogFormat)
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		g.readCommitsPipe(s, result)
	})

	return g.executor(cmd, reader)
}

// Read commits history
// projectPath should contains absolute path to project with Git repository
// path should contains relative path of file for history
// If need provide whole repository history, path should be empty
// Branch should contain branch identifier if need get specified branch results
func (g Git) ReadHistory(projectPath string, path string, branch string, offset int, limit int, result chan Commit) *Executor {
	if branch == "" {
		branch = "*"
	} else {
		branch = "*"+branch+"*"
	}

	args := append(
		make([]string, 0, 6),
		"log",
		`--format=`+gitLogFormat,
		`-n`,
		fmt.Sprintf("%d", limit),
		fmt.Sprintf("--skip=%d", offset),
		"--branches="+branch)

	if path != "" {
		args = append(args, "--", path)
	}

	cmd := g.createCommand(projectPath, args...)
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		g.readCommitsPipe(s, result)
	})

	return g.executor(cmd, reader)
}