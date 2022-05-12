package server

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	host         string
	port         int
	csvTmpFile   string
	resumeTmpDir string
}

func NewConfig() Config {
	host, ok := os.LookupEnv("TA_HOST")
	if !ok {
		host = "0.0.0.0"
	}
	log.Printf("TA_HOST set to '%s'", host)

	var port int
	portStr, ok := os.LookupEnv("TA_PORT")
	port, err := strconv.Atoi(portStr)
	if !ok || err != nil {
		port = 23234
	}
	log.Printf("TA_PORT set to '%d'", port)

	resumeTmpDir, ok := os.LookupEnv("TA_UPLOAD_DIR")
	if !ok {
		resumeTmpDir = "./uploads"
	}
	log.Printf("TA_UPLOAD_DIR set to '%s'", resumeTmpDir)

	csvTmpFile, ok := os.LookupEnv("TA_DATAFILE")
	if !ok {
		csvTmpFile = "applicants.csv"
	}
	log.Printf("TA_DATAFILE set to '%s'", csvTmpFile)

	return Config{
		host:         host,
		port:         port,
		csvTmpFile:   csvTmpFile,
		resumeTmpDir: resumeTmpDir,
	}
}
