package vcsview

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFile_Pathname(t *testing.T) {
	cases := []struct{
		path string
		name string
		pathname string
	}{
		{"", "testing.txt", "testing.txt"},
		{pathSeparator, "testing.txt", "testing.txt"},
		{"subdir", "testing.txt", "subdir"+pathSeparator+"testing.txt"},
		{"/subdir", "testing.txt", "subdir"+pathSeparator+"testing.txt"},
	}

	for key, testCase := range cases {
		f := File{}
		f.Name = testCase.name
		f.Path = testCase.path

		if pathname := f.Pathname(); pathname != testCase.pathname {
			t.Fatalf("[%d] File.Pathname() for %v want: %s, got: %s", key, f, testCase.pathname, pathname)
		}
	}
}

func TestNewFileFromProjectList(t *testing.T) {
	testingPath, _ := filepath.Abs(gitRepositoryPath)

	cases := []struct{
		relativePath string
		name string
		isDir bool
	}{
		{"", "testing.txt", false},
		{"", "second_testing.txt", false},
		{"", "testpath", true},
		{"testpath", "empty.txt", false},
	}

	for key, testCase := range cases {
		absPath := testingPath+pathSeparator+testCase.relativePath+pathSeparator+testCase.name
		stat, err := os.Stat(absPath)
		if err != nil {
			t.Fatalf("[%d] os.Stat(%s) got error: %v", key, absPath, err)
			continue
		}

		file := NewFileFromProjectList(stat, testCase.relativePath)

		if file.Path != testCase.relativePath {
			t.Errorf("[%d] file.Path for %s got: %s, want: %s", key, absPath, file.Path, testCase.relativePath)
		}

		if file.Name != testCase.name {
			t.Errorf("[%d] file.Name for %s got: %s, want: %s", key, absPath, file.Name, testCase.name)
		}

		if file.IsDir != testCase.isDir {
			t.Errorf("[%d] file.IsDir for %s got: %v, want: %v", key, absPath, file.IsDir, testCase.isDir)
		}

		if !file.IsExists {
			t.Errorf("[%d] file.IsExists for %s got: %v, want: %v", key, absPath, file.IsExists, true)
		}

		if file.Size != stat.Size() {
			t.Errorf("[%d] file.Size for %s got: %v, want: %v", key, absPath, file.Size, stat.Size())
		}

		if file.Mode != stat.Mode() {
			t.Errorf("[%d] file.Mode for %s got: %v, want: %v", key, absPath, file.Mode, stat.Mode())
		}
	}
}
