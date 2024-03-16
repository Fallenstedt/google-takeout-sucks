package main

import (
	"fmt"
	"io"
	"os"

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

	r, err :=  srv.Files.List().Q(q).PageSize(50).Fields("nextPageToken, files(id, name)").PageToken(token).Do()
			
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


func downloadFile(srv *drive.Service, driveFile *drive.File) (error) {
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
	err = os.WriteFile(driveFile.Name, body, 0777)
	if err != nil {
		return err
	}

	return nil
}