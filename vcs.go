package vcsview

// Common interfaces for each one version control system like git, mercurial, etc
type Vcs interface {
	// check for VCS version
	Version() (string, error)

	// Returns pathname for repository config (.git, .hg, etc)
	RepositoryPathname() string

	// Check project repository
	// projectPath is absolute path to project path
	// Returns error if repository not found at provided projectPath
	// Returns nil if repository found
	CheckRepository(projectPath string) error

	// Fetch repository current status
	// Returns error if the repository doesn't exists at specified path
	StatusRepository(projectPath string) (string, error)

	// Fetch repository branches
	GetBranches(projectPath string) ([]Branch, error)
}
