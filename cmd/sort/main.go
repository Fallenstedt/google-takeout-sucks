package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
)

var infoLog = logger.New("sort info")
var collisionLog = logger.New("sort collision")
var missingLog = logger.New("sort missing")
var errorLog = logger.New("sort error")

type config struct {
	dryRun bool   // perform dry run
	dist   string // the destination for all files
}

type MediaFileMap map[string]MediaFile

func main() {
	root := flag.String("root", ".", "Root directory to start")
	dist := flag.String("dist", ".", "Absolute path to destination directory")
	dryRun := flag.Bool("dryRun", true, "Performs a dry run")
	flag.Parse()

	c := config{
		dryRun: *dryRun,
		dist:   *dist,
	}

	if err := run(*root, c); err != nil {
		fmt.Fprintln(errorLog.Writer(), err)
		os.Exit(1)
	}
}

func run(root string, cfg config) error {
	mediaFilePaths, err := walkRootForFiles(root, ".json", true) // Get everything that is not JSON
	if err != nil {
		return fmt.Errorf("unable to get all media files: %v", err)
	}

	processCh := make(chan MediaFile)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	wg := sync.WaitGroup{}

	// Get unique years
	years := make(map[string]bool)
	for _, mediaFile := range *mediaFilePaths {
		mf := NewMediaFile(mediaFile)
		year, err := mf.Year()

		if err != nil {
			errCh <- err
		}
		years[year] = true
	}

	// Create year directories
	for k := range years {
		yearDir := fmt.Sprintf("%s/%s", cfg.dist, k)

		infoLog.Printf("creating directory %s", yearDir)

		err := os.MkdirAll(yearDir, 0700)
		if err != nil {
			errCh <- err
		}
	}

	// Fill process channel with all paths to media files
	go func() {

		defer close(processCh)
		for _, mediaFile := range *mediaFilePaths {
			mf := NewMediaFile(mediaFile)
			processCh <- mf
		}
	}()

	// Create workers to process media file paths

	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for mediaFile := range processCh {

				year, err := mediaFile.Year()
				if err != nil {
					errCh <- err
					return
				}

				dest := filepath.Join(cfg.dist, year, mediaFile.Filename)

				if _, err := os.Stat(dest); err == nil {
					collisionLog.Printf("already exists %s", dest)
				} else if errors.Is(err, os.ErrNotExist) {
					infoLog.Printf("Moving %s to %s", mediaFile.Filename, dest)

					if cfg.dryRun {
						f, err := os.Create(dest)
						if err != nil {
							errCh <- err
						}
						f.Close()
					} else {
						err := os.Rename(mediaFile.Path, dest)
						if err != nil {
							errCh <- err
						}
					}

				} else {
					errCh <- err
				}

			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			errorLog.Println(err)
		case <-doneCh:
			infoLog.Println("done")
			os.Exit(0)
		}
	}
}
