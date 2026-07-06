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
		expectErr bool
	}{
		{
			name:     "Empty path",
			dataDir:  "/data",
			path:     "",
			expected: "/data",
			expectErr: false,
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
		{
			name:     "Data dir with trailing slash is normalized",
			dataDir:  "/data/",
			path:     "/sub",
			expected: "/data/sub",
		},
		{
			name:      "Relative traversal is rejected",
			dataDir:   "/data",
			path:      "../../../etc/passwd",
			expectErr: true,
		},
		{
			name:      "Deep relative traversal is rejected",
			dataDir:   "/data",
			path:      "../../../../../../etc/cron.d/evil",
			expectErr: true,
		},
		{
			name:      "Mixed traversal escaping the data dir is rejected",
			dataDir:   "/data",
			path:      "sub/../../etc/passwd",
			expectErr: true,
		},
 	}

	for _, test := range tests {
		pathConverter := NewPathConverter(test.dataDir)
		actual, err := pathConverter.GetAbsolutePath(test.path)

		if test.expectErr && err == nil {
			t.Errorf("Expected an error for %s (%s, %s)", test.name, test.dataDir, test.path)
		}

		if !test.expectErr && err != nil {
			t.Errorf("Got unexpected error for %s (%s, %s): %s", test.name, test.dataDir, test.path, err.Error())
		}

		if !test.expectErr && actual != test.expected {
			t.Errorf("Invalid absolute path for %s, expected %s, got %s", test.name, test.expected, actual)
		}
	}
}
