package logger

import (
	"log"
	"os"
)

func SetupLogger(filepath string) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file.")
	}

	log.SetOutput(file)
}
