package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
)

type config struct {
	sourceDir *string
	dryRun    *bool
	outDir    *string
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
		dryRun:    dryRun,
		outDir:    outDir,
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
		infoLog.Printf("Found %d zip files\n", len(filepaths))
		for _, name := range filepaths {
			infoLog.Println(name)
		}
	}

	processCh := make(chan string)
	resCh := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// Loop through all zip files and store them in process channel
	go func() {
		defer close(processCh)
		for _, zipfile := range filepaths {
			processCh <- zipfile
		}
	}()

	// Create 1 worker for every CPU to process a zip file
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for zipFile := range processCh {
				err := unzipFile(zipFile, cfg)
				if err != nil {
					errCh <- err
				} else {
					resCh <- zipFile
				}
			}
		}()
	}

	// Wait for operation to complete
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <- errCh:
			errorLog.Println(err)
		case data := <-resCh:
			infoLog.Printf("extracted files from %s\n", data)
		case <-doneCh:
			fmt.Println("done")
			os.Exit(0)
		}
	}
}

func unzipFile(absolutePathOfFile string, cfg *config) error {
	infoLog.Printf("unzipping file %s\n", absolutePathOfFile)
	filename := filepath.Base(absolutePathOfFile)
	collisionSafeDir :=  strings.TrimSuffix(filename, filepath.Ext(filename)) 
	dst := filepath.Join(*cfg.outDir, collisionSafeDir)

	infoLog.Printf("extracting %s to %s\n", absolutePathOfFile, dst)

	if *cfg.dryRun {
		return nil
	}

	err := extract(absolutePathOfFile, dst)
	if err != nil {
		return fmt.Errorf("unable to extract files from zip %s: %w", absolutePathOfFile, err)
	}
	return nil
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
