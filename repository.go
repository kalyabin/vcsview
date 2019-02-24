package vcsview

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Structure which provides access control to some project repository
type Repository struct {
	// Command line interface for specified version control system
	cmd Vcs

	// Project absolute path (not a path to config directory)
	projectPath string
}

// Repository absolute path (path to config directory, for example, /path/to/project/.git)
func (r Repository) RepositoryPath() string {
	return r.projectPath+pathSeparator+r.cmd.RepositoryPathname()
}

// Returns project path (not a path to repository)
func (r Repository) ProjectPath() string {
	return r.projectPath
}

// Get command line interface
func (r *Repository) Cmd() Vcs {
	return r.cmd
}

// Get absolute path of subdir path
// Returns error if file path out of project path
func (r Repository) AbsPath(subDir string) (string, error) {
	path, err := filepath.Abs(r.projectPath +pathSeparator+subDir)
	if err != nil {
		return "", err
	}

	path = filepath.Clean(path)

	if rPath, pPath := []rune(path), []rune(r.projectPath); len(rPath) < len(pPath) || string(rPath[:len(pPath)]) != r.projectPath {
		return "", fmt.Errorf("Directory %s out of %s", path, r.projectPath)
	}

	return path, nil
}

// Get project files list
// If subDir is empty - returns root directory path list
// If subDir is out of projectPath - returns error
func (r Repository) FilesList(subDir string) ([]File, error) {
	var result []File

	path, err := r.AbsPath(subDir)
	if err != nil {
		return result, err
	}

	// get relative file path
	relativePath, _ := filepath.Rel(r.projectPath, path)
	relativePath = strings.TrimPrefix(relativePath, "."+pathSeparator)

	// seek files
	files, err := ioutil.ReadDir(string(path))
	if err != nil {
		return result, err
	}

	result = make([]File, 0)

	for _, i := range files {
		if i.Name() == r.cmd.RepositoryPathname() && subDir == "" {
			// list doesn't need to provide repository path
			continue
		}

		result = append(result, NewFileFromProjectList(i, relativePath))
	}

	return result, nil
}

// Check the repository
// Repository exists and well works if the vcs doesnt throw an error while fetch repository status
func (r Repository) Check() (err error) {
	_, err = r.cmd.StatusRepository(r.projectPath)
	return
}

// Get repository branches
func (r Repository) GetBranches() ([]Branch, error) {
	return r.cmd.GetBranches(r.projectPath)
}

// Create new repository object for the project path and provided version control system
// Returns error if repository not found at the path
// Returns repository object if repository found at the path
func NewRepository(projectPath string, vcs Vcs) (Repository, error) {
	var r Repository

	projectPath, err := filepath.Abs(projectPath)

	if err != nil {
		return r, err
	}

	projectPath = filepath.Clean(projectPath)

	// repository not found at the path
	if err = vcs.CheckRepository(projectPath); err != nil {
		return r, err
	}

	r = Repository{vcs, projectPath}
	return r, nil
}