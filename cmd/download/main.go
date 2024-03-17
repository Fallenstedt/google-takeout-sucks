package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type config struct {
	directoryId *string
	dryRun      *bool
	outDir      *string
}

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
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	var r []*drive.File
	err = fetchFiles(&r, "", srv, cfg)
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	printFiles(r)
	if *cfg.dryRun {
		return
	}

	processCh := make(chan *drive.File, 5) // Only process 5 files at a time
	resCh := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	//File the process channel
	// So each one will be processed when worker is available
	go func() {
		defer close(processCh)
		for _, driveFile := range r {
			wg.Add(1)
			processCh <- driveFile
		}
	}()

	// Create 3 workers to download files
	for w := 1; w <= 3; w++ {
		go downloadWorker(w, processCh, errCh, resCh, srv, cfg, &wg)
	}

	
	// Wait for all files to be downloaded and saved to disk
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	errorLog := logger.New("error")
	for {
		select {
		case err := <-errCh:
			errorLog.Println(err)
		case data := <-resCh:
			fmt.Printf("file has been saved: %s\n", data)
		case <-doneCh:
			fmt.Println("done")
			os.Exit(0)
		}
	}
}

func downloadWorker(
	id int, 
	processCh <-chan *drive.File, 
	errCh chan<- error, 
	resCh chan<- string, 
	srv *drive.Service, 
	cfg *config,
	wg *sync.WaitGroup) {
	for driveFile := range processCh {
	
		fmt.Printf("worker %d: downloading file: %s\n", id, driveFile.Name)
		err := downloadFile(srv, driveFile, cfg)
		if err != nil {
			fmt.Printf("error downloading file: %s\n", driveFile.Name)

			errCh <- fmt.Errorf("cannot download file %s - %s: %w", driveFile.Name, driveFile.Id, err)
		} else {
			resCh <- driveFile.Name
		}

		wg.Done()
	}
}
