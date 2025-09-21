package cmd

import (
	"os"
	"runtime"
	"sync"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
	download "github.com/Fallenstedt/google-photo-organizer/internal/download"
	"google.golang.org/api/drive/v3"

	"github.com/spf13/cobra"
)

var downloadErrorLog = logger.New("download error")
var downloadInfoLog = logger.New("download info")

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download your google takeout zip files from google drive",
	Run: func(cmd *cobra.Command, args []string) {

		dryRun, err := cmd.Flags().GetBool("dryRun")
		if err != nil {
			downloadErrorLog.Printf("Error finding dryRun flag: %e", err)
			return
		}

		directoryId, err := cmd.Flags().GetString("directoryId")
		if err != nil || directoryId == "" {
			downloadErrorLog.Println("Error: --directoryId is required")
			return
		}

		outDir, err := cmd.Flags().GetString("outDir")
		if err != nil || outDir == "" {
			downloadErrorLog.Println("Error: --outDir is required")
			return
		}

		cfg := &download.Config{
			DirectoryId: &directoryId,
			OutDir:      &outDir,
			DryRun:      &dryRun,
		}

		downloadInfoLog.Println("Downloading google takeout zip files from google drive...")

		downloadFiles(cmd, cfg)
	},
}

func downloadFiles(cmd *cobra.Command, cfg *download.Config) {

	srv, err := download.NewGoogleDriveService(cmd.Context())
	if err != nil {
		downloadErrorLog.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	var r []*drive.File
	err = download.FetchFiles(&r, "", srv, cfg)
	if err != nil {
		downloadErrorLog.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(r) == 0 {
		downloadErrorLog.Fatalf("No files found for downloading")
	}

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
		go download.DownloadWorker(w, processCh, errCh, resCh, srv, cfg, &wg)
	}

	// Wait for all files to be downloaded and saved to disk
	go func() {
		wg.Wait()
		close(doneCh)
		close(errCh)
		close(resCh)
	}()

	for {
		select {
		case err := <-errCh:
			downloadErrorLog.Println(err)
		case data := <-resCh:
			downloadInfoLog.Printf("file has been saved: %s\n", data)
		case <-doneCh:
			downloadInfoLog.Println("done")
			os.Exit(0)
		}
	}

}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().Bool("dryRun", false, "Dry run the operation, default is false")
	downloadCmd.Flags().String("directoryId", "", "ID of the Google Drive directory containing your takeout files (required)")
	downloadCmd.Flags().String("outDir", "", "Directory to save downloaded files (required)")

	downloadCmd.MarkFlagRequired("directoryId")
	downloadCmd.MarkFlagRequired("outDir")
}
