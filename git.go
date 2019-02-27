package vcsview

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	gitLogFormat = "%H%n%P%n%an%n%ae%n%ad%n%s%n"
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
		defer close(result)

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

// Fetch repository commit by identifier asynchronously
// ProjectPath is the absolute path to project with Git repository
// CommitId is the sha256 commit identifier (or short copy)
func (g Git) ReadCommit(projectPath string, commitId string, result chan Commit) *Executor {
	cmd := g.createCommand(projectPath, "show", "--quiet", commitId, `--pretty=format:`+gitLogFormat)
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		defer close(result)

		data := make([]string, 6)
		key := 0

		for s.Scan() {
			data[key] = s.Text()
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

				key = 0
				data = make([]string, 6)
			}
		}
	})

	return g.executor(cmd, reader)
}