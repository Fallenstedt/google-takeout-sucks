package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
)



type config struct {
	sourceDir *string
	dryRun      *bool
	outDir      *string
}

var infoLog = logger.New("unzip info")
var errorLog = logger.New("unzip error")


func main() {
	sourceDir := flag.String("source", ".", "The absolute path for the source directory containing zip files")
	dryRun := flag.Bool("dryRun", true, "Performs a dry run")
	outDir := flag.String("out", ".", "The absolute path for downloaded files")

	flag.Parse()

	cfg := &config{
		sourceDir: sourceDir,
		dryRun: dryRun,
		outDir: outDir,
	}

	run(cfg)
}

func run(cfg *config) {

	var filepaths []string
	err := getZipFilesFromSourceDir(cfg.sourceDir, &filepaths)
	if err != nil {
		errorLog.Fatalf("Unable to get zip files: %v", err)
	}

	if *cfg.dryRun {
		fmt.Printf("Found %d zip files", len(filepaths))
		for _, name := range filepaths {
			fmt.Println(name)
		}
		return 
	}
}


func getZipFilesFromSourceDir(sourceDir *string, filepaths *[]string) error {

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
				infoLog.Printf("ignoring %s because it is not a zip file\n", fileinfo.Name())
				return nil
			}

			*filepaths = append(*filepaths, path)

			return nil

		}

		return nil
		
	})
}