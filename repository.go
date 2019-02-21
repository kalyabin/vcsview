package vcsview

import "os"

// Structure which provides access control to some project repository
type Repository struct {
	// Command line interfaces for specified version control system
	Cmd Vcs

	// Project absolute path (not a path to config directory)
	ProjectPath string

	// Repository absolute path (path to config directory, for example, /path/to/project/.git)
	RepositoryPath string
}

// Create new repository object for the project path and provided version control system
// Returns error if repository not found at the path
// Returns repository object if repository found at the path
func NewRepository(projectPath string, vcs Vcs) (Repository, error) {
	var r Repository

	// repository not found at the path
	if err := vcs.CheckRepository(projectPath); err != nil {
		return r, err
	}

	r = Repository{vcs, projectPath, projectPath+string(os.PathSeparator)+vcs.RepositoryPathname()}
	return r, nil
}