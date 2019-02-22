package vcsview

import (
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
		f := File{testCase.name, testCase.path, false}

		if pathname := f.Pathname(); pathname != testCase.pathname {
			t.Fatalf("[%d] File.Pathname() for %v want: %s, got: %s", key, f, testCase.pathname, pathname)
		}
	}
}
