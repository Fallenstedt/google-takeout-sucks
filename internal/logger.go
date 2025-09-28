package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func New(name string) *log.Logger {
	// Determine log directory inside the user's home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine user home directory: %v", err)
	}

	logDir := filepath.Join(home, ".google_takeout_sucks", "logs")
	// Create directory with secure permissions for the user.
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		log.Fatalf("Unable to create log directory %s: %v", logDir, err)
	}

	logpath := filepath.Join(logDir, fmt.Sprintf("%s.log", name))

	// Open (or create) the log file with secure permissions.
	file, err := os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Unable to create/open log file: %v", err)
	}

	// Inform the user where logs are stored.
	absDir, _ := filepath.Abs(logDir)
	fmt.Printf("Logs will be written to: %s\n", absDir)

	instance := log.New(file, "", log.LstdFlags|log.Lshortfile)
	instance.Println("LogFile : " + logpath)

	return instance
}
