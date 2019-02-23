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
	name string

	// Relative file path
	path string

	// True if file is directory
	isDir bool

	// True if file is existent
	isExists bool

	// File size
	size int64

	// File access mode
	mode os.FileMode
}

// Get file name (without file path)
func (f File) Name() string {
	return f.name
}

// Get file relative file path (without file name)
func (f File) Path() string {
	return f.path
}

// Returns relative file path with file name
func (f File) Pathname() string {
	if f.path == "" || f.path == pathSeparator {
		return f.name
	}

	return strings.TrimLeft(f.path+pathSeparator+f.name, pathSeparator)
}

// Returns true if file is directory
func (f File) IsDir() bool {
	return f.isDir
}

// Returns true if file is exists at the time
func (f File) IsExists() bool {
	return f.isExists
}

// Returns file bytes size
func (f File) Size() int64 {
	return f.size
}

// Returns permissions for the file
func (f File) Mode() os.FileMode {
	return f.mode
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