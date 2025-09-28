package unzip

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testdataPath(t *testing.T) string {
	// internal/unzip -> ../../testdata
	p := filepath.Join("..", "..", "testdata")
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("unable to resolve testdata path: %v", err)
	}
	return abs
}

func TestGetZipFilesFromSourceDir_FindsZips(t *testing.T) {
	td := testdataPath(t)

	var files []string
	if err := GetZipFilesFromSourceDir(&td, &files); err != nil {
		t.Fatalf("GetZipFilesFromSourceDir returned error: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("expected at least one zip file, got 0")
	}

	// ensure expected files are present
	found1, found2 := false, false
	for _, f := range files {
		base := filepath.Base(f)
		if base == "unzipTest1.zip" {
			found1 = true
		}
		if base == "unzipTest2.zip" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Fatalf("expected unzipTest1.zip and unzipTest2.zip in results, got: %v", files)
	}
}

func TestUnzipFile_ExtractsContents(t *testing.T) {
	td := testdataPath(t)

	zipPath := filepath.Join(td, "unzipTest1.zip")
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("test zip not found at %s: %v", zipPath, err)
	}

	outDir := t.TempDir()
	dryRun := false
	cfg := &Config{
		OutDir: &outDir,
		DryRun: &dryRun,
	}

	if err := UnzipFile(zipPath, cfg); err != nil {
		t.Fatalf("UnzipFile returned error: %v", err)
	}

	// destination dir is OutDir/<zipname without ext>
	filename := filepath.Base(zipPath)
	collisionSafeDir := strings.TrimSuffix(filename, filepath.Ext(filename))
	destDir := filepath.Join(outDir, collisionSafeDir)

	info, err := os.Stat(destDir)
	if err != nil {
		t.Fatalf("expected destination dir %s to exist, stat error: %v", destDir, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", destDir)
	}

	entries, err := os.ReadDir(destDir)
	if err != nil {
		t.Fatalf("reading dest dir failed: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected files inside %s, found none", destDir)
	}
}

func TestGetZipFilesFromSourceDir_NonZipReturnsError(t *testing.T) {
	tmp := t.TempDir()
	// create a non-zip file
	notzip := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(notzip, []byte("not a zip"), 0600); err != nil {
		t.Fatalf("unable to create test file: %v", err)
	}

	var files []string
	err := GetZipFilesFromSourceDir(&tmp, &files)
	if err == nil {
		t.Fatalf("expected error when non-zip file present, got nil; files: %v", files)
	}
}
