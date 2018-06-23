package main

import (
	"io"
	"log"
	"os"
)

func SetupLogger() {
	if config.FileLogEnabled {
		file, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	}
}
