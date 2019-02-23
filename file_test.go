package vcsview

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFile_Name(t *testing.T) {
	cases := []string{
		"",
		"testing.txt",
		"empty.txt",
	}

	for key, testCase := range cases {
		f := File{}
		f.name = testCase
		if name := f.Name(); name != testCase {
			t.Errorf("[%d] File.Name() = %v, want: %v", key, name, testCase)
		}
	}
}

func TestFile_Path(t *testing.T) {
	cases := []string{
		"",
		"testpath",
		"testpath1/testpath2",
	}

	for key, testCase := range cases {
		f := File{}
		f.path = testCase
		if path := f.Path(); path != testCase {
			t.Errorf("[%d] File.Path() = %v, want: %v", key, path, testCase)
		}
	}
}

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
		f.name = testCase.name
		f.path = testCase.path

		if pathname := f.Pathname(); pathname != testCase.pathname {
			t.Fatalf("[%d] File.Pathname() for %v want: %s, got: %s", key, f, testCase.pathname, pathname)
		}
	}
}

func TestFile_IsDir(t *testing.T) {
	cases := []bool{
		true,
		false,
	}

	for key, testCase := range cases {
		f := File{}
		f.isDir = testCase
		if isDir := f.IsDir(); isDir != testCase {
			t.Errorf("[%d] File.IsDir() = %v, want: %v", key, isDir, testCase)
		}
	}
}

func TestFile_IsExists(t *testing.T) {
	cases := []bool{
		true,
		false,
	}

	for key, testCase := range cases {
		f := File{}
		f.isExists = testCase
		if isExists := f.IsExists(); isExists != testCase {
			t.Errorf("[%d] File.IsExists() = %v, want: %v", key, isExists, testCase)
		}
	}
}

func TestFile_Size(t *testing.T) {
	cases := []int64{
		0,
		10,
		200,
		1003,
	}

	for key, testCase := range cases {
		f := File{}
		f.size = testCase
		if size := f.Size(); size != testCase {
			t.Errorf("[%d] File.Size() = %v, want: %v", key, size, testCase)
		}
	}
}

func TestFile_Mode(t *testing.T) {
	cases := []struct{
		strRepresentation string
		permissions uint32
	}{
		{"----------", 0},
		{"-rwx------", 4<<6 | 2<<6 | 1<<6},
		{"-rwxrwx---", 4<<6 | 2<<6 | 1<<6 | 4<<3 | 2<<3 | 1<<3},
		{"-rwxrwxrwx", 4<<6 | 2<<6 | 1<<6 | 4<<3 | 2<<3 | 1<<3 | 4<<0 | 2<<0 | 1<<0},
	}

	for key, testCase := range cases {
		f := File{}
		f.mode = os.FileMode(testCase.permissions)

		if permissions := f.Mode().String(); permissions != testCase.strRepresentation {
			t.Errorf("[%d] f.Mode() = %v, want: %v", key, permissions, testCase.strRepresentation)
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

		if file.path != testCase.relativePath {
			t.Errorf("[%d] file.path for %s got: %s, want: %s", key, absPath, file.path, testCase.relativePath)
		}

		if file.name != testCase.name {
			t.Errorf("[%d] file.name for %s got: %s, want: %s", key, absPath, file.name, testCase.name)
		}

		if file.isDir != testCase.isDir {
			t.Errorf("[%d] file.isDir for %s got: %v, want: %v", key, absPath, file.isDir, testCase.isDir)
		}

		if !file.isExists {
			t.Errorf("[%d] file.isExists for %s got: %v, want: %v", key, absPath, file.isExists, true)
		}

		if file.size != stat.Size() {
			t.Errorf("[%d] file.size for %s got: %v, want: %v", key, absPath, file.size, stat.Size())
		}

		if file.mode != stat.Mode() {
			t.Errorf("[%d] file.mode for %s got: %v, want: %v", key, absPath, file.mode, stat.Mode())
		}
	}
}
