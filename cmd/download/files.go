package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"google.golang.org/api/drive/v3"
)

func printFiles(r []*drive.File) {
	fmt.Println("Files:")

	if len(r) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

func fetchFiles(res *[]*drive.File, token string, srv *drive.Service, cfg *config) error {
	q := fmt.Sprintf("'%s' in parents", *cfg.directoryId)

	r, err := srv.Files.List().Q(q).PageSize(50).Fields("nextPageToken, files(id, name)").PageToken(token).Do()

	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		*res = append(*res, r.Files...)
	}

	if len(r.NextPageToken) > 0 {
		return fetchFiles(res, r.NextPageToken, srv, cfg)
	}

	return nil
}

func downloadFile(srv *drive.Service, driveFile *drive.File, cfg *config) error {
	resp, err := srv.Files.Get(driveFile.Id).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// TODO absolute path to file
	// https://stackoverflow.com/questions/63218663/how-to-create-an-empty-file-in-golang-at-a-specified-path-lets-say-at-home-new

	location := fmt.Sprintf("%s/%s", *cfg.outDir, driveFile.Name)
	err = os.WriteFile(location, body, 0777)
	if err != nil {
		return err
	}

	return nil
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
