package main

import (
	"fmt"
	"testing"
)

func TestParseAuthHeader(t *testing.T) {
	var tests = []struct {
		name        string
		header      string
		expected    string
		expectedErr error
	}{
		{
			name:        "Bearer token",
			header:      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectedErr: nil,
		},
		{
			name:        "Lacks token",
			header:      "Bearer",
			expected:    "",
			expectedErr: fmt.Errorf("Invalid header"),
		},
		{
			name:        "Lacks Bearer",
			header:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expected:    "",
			expectedErr: fmt.Errorf("Invalid header"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseAuthHeader(test.header)
			if actual != test.expected && err != test.expectedErr {
				t.Fail()
			}
		})
	}
}

func TestGenerateAuthToken(t *testing.T) {
	var tests = []struct {
		name   string
		secret string
	}{
		{
			name:   "Token with secret key",
			secret: "my_jwt_secret",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encodedToken, err := generateAuthToken([]byte(test.secret))
			if err != nil {
				t.Fail()
			}

			parsedToken, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(test.secret), nil
			})

			if err != nil || !parsedToken.Valid {
				t.Fail()
			}
		})
	}
}
