package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var ErrYearNotFound = errors.New("could not find year for file")

type MediaFile struct {
	Path     string
	Filename string
}

func NewMediaFile(path string) MediaFile {
	filename := filepath.Base(path)

	return MediaFile{
		Filename: filename,
		Path:     path,
	}
}

func (m *MediaFile) Year() (string, error) {

	splitlist := strings.Split(m.Path, "/")

	for _, d := range splitlist {
		if strings.Contains(d, "Photos from") {

			yearSlice := strings.Split(d, " ")
			year := yearSlice[len(yearSlice)-1]

			return year, nil
		}
	}

	return "", ErrYearNotFound
}

func (m *MediaFile) String() string {
	return fmt.Sprintf("filename: %s, path: %s", m.Filename, m.Path)
}
