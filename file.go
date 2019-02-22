package vcsview

import (
	"os"
	"strings"
)

const (
	pathSeparator = string(os.PathSeparator)
)

// Project file with relative path
type File struct {
	// File name
	Name string

	// Relative file path
	Path string

	// True if file is directory
	IsDir bool

	// True if file is existent
}

// Returns relative file path with file name
func (f File) Pathname() string {
	if f.Path == "" || f.Path == pathSeparator {
		return f.Name
	}

	return strings.TrimLeft(f.Path+pathSeparator+f.Name, pathSeparator)
}