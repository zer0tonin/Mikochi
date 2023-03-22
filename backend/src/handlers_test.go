package main

import (
	"testing"
)

func TestFileInDir(t *testing.T) {
	var tests = []struct {
		name     string
		file     string
		dir      string
		expected bool
	}{
		{
			name:     "File in dir",
			file:     "/test/myfile.txt",
			dir:      "/test/",
			expected: true,
		},
		{
			name:     "File in another dir",
			file:     "/test2/myfile.txt",
			dir:      "/test/",
			expected: false,
		},
		{
			name:     "File in a nested dir",
			file:     "/test/inner/myfile.txt",
			dir:      "/test/",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if fileInDir(test.file, test.dir) != test.expected {
				t.Fail()
			}
		})
	}
}
