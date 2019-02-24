package vcsview

import "time"

// Represents a commit model
type Commit struct {
	// Commit identifier (sha256 for a git one and number for a mercurial one)
	id string

	// Commit date and time
	date time.Time

	// Commit author
	author Contributor

	// Commit message
	message string

	// Parent commits identifiers
	parents []string
}

// Get commit identifier
func (c Commit) Id() string {
	return c.id
}

// Get commit date time
func (c Commit) Date() time.Time {
	return c.date
}

// Get commit author
func (c Commit) Author() Contributor {
	return c.author
}

// Get commit message
func (c Commit) Message() string {
	return c.message
}

// Get commit parents identifiers
func (c Commit) Parents() []string {
	return c.parents
}
