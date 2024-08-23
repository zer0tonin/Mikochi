package browser

import (
	"testing"
)

func TestGetAbsolutePath(t *testing.T) {
	tests := []struct {
		name     string
		dataDir  string
		path     string
		expected string
	}{
		{
			name:     "Empty path",
			dataDir:  "/data",
			path:     "",
			expected: "/data",
		},
		{
			name:     "Two sub dirs",
			dataDir:  "/data",
			path:     "/sub/nested",
			expected: "/data/sub/nested",
		},
		{
			name:     "Unclean path",
			dataDir:  "/data",
			path:     "/sub/nested/..",
			expected: "/data/sub",
		},
	}

	for _, test := range tests {
		pathConverter := NewPathConverter(test.dataDir)
		actual := pathConverter.GetAbsolutePath(test.path)

		if actual != test.expected {
			t.Errorf("Invalid absolute path, expected %s, got %s", test.expected, actual)
		}
	}
}
