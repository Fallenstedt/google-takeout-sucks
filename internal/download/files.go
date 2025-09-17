package download

import (
	"fmt"
	"io"
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
}


func FetchFiles(res *[]*drive.File, token string, srv *drive.Service, cfg *Config) error {
	q := fmt.Sprintf("'%s' in parents", *cfg.DirectoryId)

	r, err := srv.Files.List().Q(q).PageSize(50).Fields("nextPageToken, files(id, name)").PageToken(token).Do()

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
		})

		if err != nil {
			errCh <- fmt.Errorf("cannot download file %s - %s: %w", driveFile.Name, driveFile.Id, err)
		} else {
			resCh <- driveFile.Name
		}

	}
}


func downloadFileToDisk(input DownloadFileInput) error {
	// Create the file
	absdst := filepath.Join(input.destination, input.driveFile.Name)
	out, err := os.Create(absdst)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", input.destination, err)
	}
	defer out.Close()

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
	return nil
}