package vcsview

// Represents VCS branch model
type Branch struct {
	// Branch identifier
	id string

	// Head commit identifier of the branch
	head string

	// That is the current branch
	isCurrent bool
}

// Get branch identifier
func (b Branch) Id() string {
	return b.id
}

// Get branch head commit identifier
func (b Branch) Head() string {
	return b.head
}

// Returns true if branch is current
func (b Branch) IsCurrent() bool {
	return b.isCurrent
}

