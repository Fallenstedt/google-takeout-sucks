package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			f := filterOut(tc.file, tc.ext, info)

			if f != tc.expected {
				t.Errorf("Expected %t, got %t instead\n", tc.expected, f)
			}
		})
	}
}
