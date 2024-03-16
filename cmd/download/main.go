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
	dryRun *bool
}

func main() {

	directoryId := flag.String("directoryId", "", "The ID directory of your Google Takeout Folder")
	dryRun := flag.Bool("dryRun", true, "Performs a dry run")

	flag.Parse()

	cfg := &config{
		directoryId: directoryId,
		dryRun: dryRun,
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

	resCh := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// Download all the files
	for _, driveFile := range r {
		wg.Add(1)
		
		go func (driveFile *drive.File)  {
			defer wg.Done()

			fmt.Printf("downloading file: %s\n", driveFile.Name)
			err := downloadFile(srv, driveFile)
			if err != nil {
				errCh <- fmt.Errorf("cannot download file %s - %s: %w", driveFile.Name, driveFile.Id, err)
			} else {
				resCh <- driveFile.Name
			}

		}(driveFile)
	}

	// Wait for all files to be downloaded and saved to disk
	go func ()  {
		wg.Wait()
		close(doneCh)
	}()

	errorLog := logger.New("error")
	for {
		select {
		case err := <-errCh:
			errorLog.Println(err)
		case data := <- resCh:
			fmt.Printf("file has been saved: %s\n", data)
		case <-doneCh:
			fmt.Println("done")
			os.Exit(0)
		}
	}
}
