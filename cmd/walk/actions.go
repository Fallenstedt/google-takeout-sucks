package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func filterOut(path string, ext string, info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}

	filepathExtension := strings.ToLower(filepath.Ext(path))


	if ext != "" && filepathExtension != strings.ToLower(ext) {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
