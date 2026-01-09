package cmd

import (
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/Fallenstedt/google-takeout-sucks/internal/unzip"
	"github.com/spf13/cobra"
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "Unzips zip files",
	Long:  `This will unzip zip files found in a directory`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, err := cmd.Flags().GetBool("dryRun")
		if err != nil {
			log.Printf("Error finding dryRun flag: %e", err)
			return
		}

		source, err := cmd.Flags().GetString("source")
		if err != nil || source == "" {
			log.Println("Error: --source is required")
			return
		}

		outDir, err := cmd.Flags().GetString("outDir")
		if err != nil || outDir == "" {
			log.Println("Error: --outDir is required")
			return
		}

		cfg := &unzip.Config{
			SourceDir: &source,
			OutDir:    &outDir,
			DryRun:    &dryRun,
		}

		log.Println("Unzipping files")

		unzipFiles(cfg)
	},
}

func unzipFiles(cfg *unzip.Config) {
	var filepaths []string
	err := unzip.GetZipFilesFromSourceDir(cfg.SourceDir, &filepaths)
	if err != nil {
		log.Fatalf("Unable to get zip files: %v", err)
	}

	log.Printf("Found %d zip files\n", len(filepaths))
	for _, name := range filepaths {
		log.Println(name)
	}

	processCh := make(chan string)
	resCh := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}
	go func() {
		defer close(processCh)
		for _, zipfile := range filepaths {
			processCh <- zipfile
		}
	}()

	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for zipFile := range processCh {
				log.Printf("unzipping file %s\n", zipFile)
				err := unzip.UnzipFile(zipFile, cfg)
				if err != nil {
					errCh <- err
				} else {
					resCh <- zipFile
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			log.Println(err)
		case data := <-resCh:
			log.Printf("file has been saved: %s\n", data)
		case <-doneCh:
			log.Println("done")
			os.Exit(0)
		}
	}
}

func init() {
	rootCmd.AddCommand(unzipCmd)

	unzipCmd.Flags().String("source", ".", "Absolute path to directory containing zip files")
	unzipCmd.Flags().Bool("dryRun", false, "Dry run the operation, default is false")
	unzipCmd.Flags().String("outDir", ".", "The absolute path for downloaded files")

	unzipCmd.MarkFlagRequired("source")
	unzipCmd.MarkFlagRequired("outDir")
}
