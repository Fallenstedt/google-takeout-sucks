package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func New(name string) *log.Logger {
	logpath, err := filepath.Abs(fmt.Sprintf("tmp/%s.log", name))
	fmt.Println(logpath)
	if err != nil {
		log.Fatalf("Unable to create absolute path: %v", err)
	}

	file, err := os.Create(logpath)

	if err != nil {
		log.Fatalf("Unable to create log file: %v", err)
	}

	instance := log.New(file, "", log.LstdFlags|log.Lshortfile)
	instance.Println("LogFile : " + logpath)

	return instance
}
