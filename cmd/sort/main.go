package main

import (
	"flag"
	"fmt"
	"os"
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
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	dryRun := flag.Bool("dryRun", true, "Performs a dry run")
	flag.Parse()

	c := config{
		dryRun: *dryRun,
	}

	if err := run(*root, c); err != nil {
		fmt.Fprintln(errorLog.Writer(), err)
		os.Exit(1)
	}
}

func run(root string, cfg config) error {
	jsonFiles, err := walkRootForFiles(root, cfg, ".json", false) // Get everything that is JSON
	if err != nil {
		return fmt.Errorf("unable to get all json files: %v", err)
	}
	mediaFilePaths, err := walkRootForFiles(root, cfg, ".json", true) // Get everything that is not JSON
	if err != nil {
		return fmt.Errorf("unable to get all media files: %v", err)
	}

	// Build media file map
	mediaFileMap := make(map[string]MediaFile)
	for _, mediaFilePath := range *mediaFilePaths {
		mf := NewMediaFile(mediaFilePath)

		_, exists := mediaFileMap[mf.Filename] 
		if exists {
			collisionLog.Println(mf)
		} else {
			mediaFileMap[mf.Filename] = mf
		}
	}
	
	processCh := make(chan GooglePhotoJsonFile)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	wg := sync.WaitGroup{}

	go func ()  {
		defer close(processCh)

		for _, jsonFilePath := range *jsonFiles {
			var d GooglePhotoJsonFile
			d.Path = &jsonFilePath
			err := d.GetPayload()
	
			if err != nil {
				errCh <- err
			}
	
			processCh <- d
		}
	}()
	
	// Create workers to process JSON files
	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)

		go func ()  {
			defer wg.Done()

			for jsonPayload := range processCh {
				filename := jsonPayload.Data.Title
				
				mf, ok := mediaFileMap[filename]
				if !ok {
					missingLog.Printf("cannot find file %s", *jsonPayload.Path)
				} else {
					infoLog.Printf("found media file %s", mf)
				}

			}
		}()
	}

	go func ()  {
		wg.Wait()
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


