package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func filterOut(path string, ext string, info os.FileInfo, exclude bool) bool {
	if info.IsDir() {
		return true
	}

	filepathExtension := strings.ToLower(filepath.Ext(path))

	if ext != "" {
		if !exclude {
			if filepathExtension != strings.ToLower(ext) {
				return true
			}
		}
		if exclude {
			if filepathExtension == strings.ToLower(ext) {
				return true
			}
		}

	}

	return false
}

func walkRootForFiles(root string, ext string, exclude bool) (*[]string, error) {
	files := new([]string)
	err := filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("unexpected error walking %w", err)
			}

			if filterOut(path, ext, info, exclude) {
				return nil
			}

			// Only get photos which are from Google photos "Photos from XXXX" albums.
			if strings.Contains(path, "Photos from") {
				*files = append(*files, path)
			}

			return nil
		},
	)
	return files, err
}
