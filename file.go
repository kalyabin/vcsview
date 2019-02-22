package vcsview

import (
	"os"
	"strings"
)

const (
	pathSeparator = string(os.PathSeparator)
)

type FileStatus string

const(
	FileAdded FileStatus = "A"
	FileCopied FileStatus = "C"
	FileDeleted FileStatus = "D"
	FileModified FileStatus = "M"
	FileRenamed FileStatus = "R"
	FileTyped FileStatus = "T"
	FileUnmerged FileStatus = "U"
	FileUnknownStatus FileStatus = "X"
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
	IsExists bool

	// File size
	Size int64

	// File mode
	Mode os.FileMode
}

// Create new file for project repository list
// In this case file should exist on the disk
// relativePath is relative path, where file located
func NewFileFromProjectList(i os.FileInfo, relativePath string) File {
	if relativePath == "." {
		relativePath = ""
	}

	f := File{i.Name(), relativePath, i.IsDir(), true, i.Size(), i.Mode()}
	return f
}

// Returns relative file path with file name
func (f File) Pathname() string {
	if f.Path == "" || f.Path == pathSeparator {
		return f.Name
	}

	return strings.TrimLeft(f.Path+pathSeparator+f.Name, pathSeparator)
}