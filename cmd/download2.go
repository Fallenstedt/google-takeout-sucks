/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// download2Cmd represents the download2 command
var download2Cmd = &cobra.Command{
	Use:   "download",
	Short: "Download your google takeout zip files from google drive",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
	},
	
}

func init() {
	rootCmd.AddCommand(download2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// download2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// download2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
