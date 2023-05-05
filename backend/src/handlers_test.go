package main

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
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
			dir:      "/test",
			expected: true,
		},
		{
			name:     "File in another dir",
			file:     "/test2/myfile.txt",
			dir:      "/test",
			expected: false,
		},
		{
			name:     "File in a nested dir",
			file:     "/test/inner/myfile.txt",
			dir:      "/test",
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

func TestFileMatchesSearch(t *testing.T) {
	var tests = []struct {
		name     string
		file     string
		dir      string
		search   string
		expected bool
	}{
		{
			name:     "File in dir and matches search",
			file:     "/test/myfile.txt",
			dir:      "/test",
			search:   "file",
			expected: true,
		},
		{
			name:     "File in dir and doesn't match search",
			file:     "/test/myfile.txt",
			dir:      "/test",
			search:   "nope",
			expected: false,
		},
		{
			name:     "File not in dir and doesn't match search",
			file:     "/other/myfile.txt",
			dir:      "/test",
			search:   "nope",
			expected: false,
		},
		{
			name:     "File not in dir but matches search",
			file:     "/other/myfile.txt",
			dir:      "/test",
			search:   "file",
			expected: false,
		},
		{
			name:     "File in sub-dir and matches search",
			file:     "/test/sub/myfile.txt",
			dir:      "/test",
			search:   "file",
			expected: true,
		},
		{
			name:     "File in sub-dir and doesn't match search",
			file:     "/test/sub/myfile.txt",
			dir:      "/test",
			search:   "nope",
			expected: false,
		},
		{
			name:     "File in close-name dir and matches search",
			file:     "/testdir/sub/myfile.txt",
			dir:      "/test",
			search:   "file",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if fileMatchesSearch(test.file, test.dir, test.search) != test.expected {
				t.Fail()
			}
		})
	}
}
