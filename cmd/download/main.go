package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"
	"sync"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
	"google.golang.org/api/drive/v3"
)

type config struct {
	directoryId *string
	dryRun      *bool
	outDir      *string
}

var errorLog = logger.New("download error")
var infoLog = logger.New("download info")

func main() {

	directoryId := flag.String("directoryId", "", "The ID directory of your Google Takeout Folder")
	dryRun := flag.Bool("dryRun", true, "Performs a dry run")
	outDir := flag.String("outDir", ".", "The absolute path for downloaded files")

	flag.Parse()

	cfg := &config{
		directoryId: directoryId,
		dryRun:      dryRun,
		outDir:      outDir,
	}

	ctx := context.Background()
	srv, err := newGoogleDriveService(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	var r []*drive.File
	err = fetchFiles(&r, "", srv, cfg)
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	printFiles(r, func(s string) {
		infoLog.Println(s)
	})

	processCh := make(chan *drive.File)
	resCh := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// Fill the process channel
	// So each one will be processed when worker is available
	go func() {
		defer close(processCh)
		for _, driveFile := range r {
			processCh <- driveFile
		}
	}()

	// Create workers to download files
	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go downloadWorker(w, processCh, errCh, resCh, srv, cfg, &wg)
	}

	// Wait for all files to be downloaded and saved to disk
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			errorLog.Println(err)
		case data := <-resCh:
			infoLog.Printf("file has been saved: %s\n", data)
		case <-doneCh:
			infoLog.Println("done")
			os.Exit(0)
		}
	}
}
