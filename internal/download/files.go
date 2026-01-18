package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/api/drive/v3"
)

type DownloadFileInput struct {
	destination string
	srv         *drive.Service
	driveFile   *drive.File
	infoLog     *log.Logger
}

func FetchFiles(res *[]*drive.File, token string, srv *drive.Service, cfg *Config) error {
	q := fmt.Sprintf("'%s' in parents", *cfg.DirectoryId)

	r, err := srv.Files.List().Q(q).PageSize(50).Fields("nextPageToken, files(id, name, size)").PageToken(token).Do()

	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		*res = append(*res, r.Files...)
	}

	if len(r.NextPageToken) > 0 {
		return FetchFiles(res, r.NextPageToken, srv, cfg)
	}

	return nil
}

func DownloadWorker(
	id int,
	processCh <-chan *drive.File,
	errCh chan<- error,
	resCh chan<- string,
	srv *drive.Service,
	cfg *Config,
	infoLog *log.Logger,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for driveFile := range processCh {

		if *cfg.DryRun {
			resCh <- driveFile.Name
			return
		}

		err := downloadFileToDisk(DownloadFileInput{
			srv:         srv,
			destination: *cfg.OutDir,
			driveFile:   driveFile,
			infoLog:     infoLog,
		})

		if err != nil {
			errCh <- fmt.Errorf("cannot download file %s - %s: %w", driveFile.Name, driveFile.Id, err)
		} else {
			resCh <- driveFile.Name
		}

	}
}

func downloadFileToDisk(input DownloadFileInput) error {
	// Log which file we're working on
	input.infoLog.Printf("Processing file: %s", input.driveFile.Name)

	// Check if file already exists and matches the remote file size
	absdst := filepath.Join(input.destination, input.driveFile.Name)
	if fileInfo, err := os.Stat(absdst); err == nil {
		// File exists, check if size matches
		input.infoLog.Printf("File found locally: %s (local size: %d bytes, remote size: %d bytes)",
			input.driveFile.Name, fileInfo.Size(), input.driveFile.Size)

		if fileInfo.Size() == input.driveFile.Size {
			// File already exists with same size, skip download
			input.infoLog.Printf("Skipping file (already complete): %s", input.driveFile.Name)
			return nil
		}
		// File exists but size doesn't match, will re-download
		input.infoLog.Printf("File size mismatch, re-downloading: %s", input.driveFile.Name)
	} else {
		// File doesn't exist
		input.infoLog.Printf("New file, downloading: %s (size: %d bytes)", input.driveFile.Name, input.driveFile.Size)
	}

	// Create the file
	out, err := os.Create(absdst)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", input.destination, err)
	}
	defer out.Close()

	// Log start of download
	input.infoLog.Printf("Starting download: %s", input.driveFile.Name)

	// Fetch the file
	resp, err := input.srv.Files.Get(input.driveFile.Id).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status fetching file: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to copy response to destination: %w", err)
	}

	// Log completion
	input.infoLog.Printf("Download completed: %s", input.driveFile.Name)
	return nil
}
