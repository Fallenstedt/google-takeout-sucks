/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"
	"sync"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
	"github.com/Fallenstedt/google-photo-organizer/internal/unzip"
	"github.com/spf13/cobra"
)

var infoLog = logger.New("unzip info")
var errorLog = logger.New("unzip error")

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "Unzips zip files",
	Long: `This will unzip zip files found in a directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unzip called")
		dryRun, err := cmd.Flags().GetBool("dryRun");
		if (err != nil) {
			errorLog.Printf("Error finding dryRun flag: %e", err);
			return
		}

		source, err := cmd.Flags().GetString("source")
		if err != nil || source == "" {
			downloadErrorLog.Println("Error: --source is required")
			return
		}

		outDir, err := cmd.Flags().GetString("outDir")
		if err != nil || outDir == "" {
			downloadErrorLog.Println("Error: --outDir is required")
			return
		}

		cfg := &unzip.Config{
			SourceDir: &source,
			OutDir:      &outDir,
			DryRun: &dryRun,
		}

		infoLog.Println("Unzipping files")
		

		unzipFiles(cfg)
	},
}

func unzipFiles(cfg *unzip.Config) {
	var filepaths []string
	err := unzip.GetZipFilesFromSourceDir(cfg.SourceDir, &filepaths)
	if err != nil {
		errorLog.Fatalf("Unable to get zip files: %v", err)
	}

	infoLog.Printf("Found %d zip files\n", len(filepaths))
	for _, name := range filepaths {
		infoLog.Println(name)
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
				infoLog.Printf("unzipping file %s\n", zipFile)
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


}

func init() {
	rootCmd.AddCommand(unzipCmd)

	unzipCmd.Flags().String("source", ".", "Absolute path to directory containing zip files")
	unzipCmd.Flags().Bool("dryRun", false, "Dry run the operation, default is false")
	unzipCmd.Flags().String("outDir", ".", "The absolute path for downloaded files")

	unzipCmd.MarkFlagRequired("source")
	unzipCmd.MarkFlagRequired("outDir")
}
