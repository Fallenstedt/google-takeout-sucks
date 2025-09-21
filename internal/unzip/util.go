package unzip

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)



func UnzipFile(absolutePathOfFile string, cfg *Config) error {
	filename := filepath.Base(absolutePathOfFile)
	collisionSafeDir := strings.TrimSuffix(filename, filepath.Ext(filename))
	dst := filepath.Join(*cfg.OutDir, collisionSafeDir)

	if *cfg.DryRun {
		return nil
	}

	err := extract(absolutePathOfFile, dst)
	if err != nil {
		return fmt.Errorf("unable to extract files from zip %s: %w", absolutePathOfFile, err)
	}
	return nil
}

func GetZipFilesFromSourceDir(sourceDir *string, filepaths *[]string) error {

	return filepath.WalkDir(*sourceDir, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return fmt.Errorf("error walking directory %s: %w ", path, err)
		}

		if d.IsDir() {
			return nil
		}

		if !d.IsDir() {
			fileinfo, err := d.Info()
			if err != nil {
				return fmt.Errorf("error getting info of file %s: %w", path, err)
			}

			extension := filepath.Ext(fileinfo.Name())

			if extension != ".zip" {
				return fmt.Errorf("%s is not a zip file, move out of the directory", fileinfo.Name())
			}

			*filepaths = append(*filepaths, path)

			return nil

		}

		return nil

	})
}


// Given a source filename and a destination path, extract the ZIP archive
func extract(zipFilename, destPath string) error {
	// Extract the ZIP file and don't filter out any files
	return filterExtract(zipFilename, destPath, func(_ string) bool {
		return true
	})
}

// Given a source filename and a destination path, extract the ZIP archive.
// The filter function can be used to avoid extracting some filenames;
// when filterFunc returns true, the file is extracted.
func filterExtract(zipFilename, destPath string, filterFunc func(string) bool) error {

	// Open the source filename for reading
	zipReader, err := zip.OpenReader(zipFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// For each file in the archive
	for _, archiveReader := range zipReader.File {

		// Open the file in the archive
		archiveFile, err := archiveReader.Open()
		if err != nil {
			return err
		}
		defer archiveFile.Close()

		// Prepare to write the file
		finalPath := filepath.Join(destPath, archiveReader.Name)

		// Check if the file to extract is just a directory
		if archiveReader.FileInfo().IsDir() {
			err = os.MkdirAll(finalPath, 0755)
			if err != nil {
				return err
			}
			// Continue to the next file in the archive
			continue
		}

		if !filterFunc(finalPath) {
			// Skip this file
			continue
		}

		// Create all needed directories
		if os.MkdirAll(filepath.Dir(finalPath), 0755) != nil {
			return err
		}

		// Prepare to write the destination file
		destinationFile, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, archiveReader.Mode())
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// Write the destination file
		if _, err = io.Copy(destinationFile, archiveFile); err != nil {
			return err
		}
	}

	return nil
}
