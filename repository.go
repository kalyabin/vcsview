package vcsview

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Structure which provides access control to some project repository
type Repository struct {
	// Command line interfaces for specified version control system
	Cmd Vcs

	// Project absolute path (not a path to config directory)
	ProjectPath string

	// Repository absolute path (path to config directory, for example, /path/to/project/.git)
	RepositoryPath string
}

// Get project files list
// If subDir is empty - returns root directory path list
// If subDir is out of projectPath - returns error
func (r Repository) FilesList(subDir string) ([]File, error) {
	var result []File

	path, err := filepath.Abs(r.ProjectPath+pathSeparator+subDir)
	if err != nil {
		return result, err
	}

	path = filepath.Clean(path)

	if len([]rune(path)) < len([]rune(r.ProjectPath)) || string([]rune(path)[:len([]rune(r.ProjectPath))]) != r.ProjectPath {
		return result, fmt.Errorf("Directory %s out of %s", path, r.ProjectPath)
	}

	// get relative file path
	relativePath, _ := filepath.Rel(r.ProjectPath, path)
	relativePath = strings.TrimLeft(relativePath, "."+pathSeparator)

	// seek files
	files, err := ioutil.ReadDir(string(path))
	if err != nil {
		return result, err
	}

	result = make([]File, 0)

	for _, i := range files {
		if i.Name() == r.Cmd.RepositoryPathname() {
			// list doesn't need to provide repository path
			continue
		}

		file := File{i.Name(), string(relativePath), i.IsDir()}
		result = append(result, file)
	}

	return result, nil
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

	repoPath := projectPath+pathSeparator+vcs.RepositoryPathname()

	r = Repository{vcs, projectPath, repoPath}
	return r, nil
}