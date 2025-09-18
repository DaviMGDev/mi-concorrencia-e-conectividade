package config

import (
	"log"
	"os"
)

var (
	ServerAddress string

	logFilePath = "client.log"
	logFile     *os.File
)

func Initialize() {
	ServerAddress = "localhost:8080"
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
}

func Close() {
}
