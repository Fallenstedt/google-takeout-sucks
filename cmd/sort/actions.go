package main

import (
	"fmt"
	"io"
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
			if  filepathExtension != strings.ToLower(ext) {
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

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func walkRootForFiles(root string, cfg config, ext string, exclude bool) (*[]string, error) {
	files := new([]string)
	err := filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("unexpected error walking %w", err)
			}

			if filterOut(path, ext, info, exclude) {
				return nil
			}

			if cfg.dryRun {
				return listFile(path, infoLog.Writer())
			}

			*files = append(*files, path)
			return nil
		},
	)
	return files, err
}

