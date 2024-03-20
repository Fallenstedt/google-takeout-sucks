package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	logger "github.com/Fallenstedt/google-photo-organizer/internal"
)

var infoLog = logger.New("walk info")
var errorLog = logger.New("walk error")

type config struct {
	ext  string // extension to filter out
	list bool   // list files
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension to look for")
	flag.Parse()

	c := config{
		ext:  *ext,
		list: *list,
	}

	if err := run(*root, c); err != nil {
		fmt.Fprintln(errorLog.Writer(), err)
		os.Exit(1)
	}
}

func run(root string, cfg config) error {
	return filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, cfg.ext, info) {
				return nil
			}

			if cfg.list {
				return listFile(path, infoLog.Writer())
			}

			return listFile(path, infoLog.Writer())
		},
	)
}
