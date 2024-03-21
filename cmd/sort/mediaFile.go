package main

import (
	"fmt"
	"path/filepath"
)

type MediaFile struct {
	Path string
	Filename string
}

func NewMediaFile(path string) MediaFile {
	filename := filepath.Base(path)

	return MediaFile{
		Filename: filename,
		Path: path,
	}
}

func (m *MediaFile) String() string {
	return fmt.Sprintf("filename: %s, path: %s", m.Filename, m.Path)
}