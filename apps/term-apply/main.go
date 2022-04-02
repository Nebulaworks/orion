package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nebulaworks/orion/apps/term-apply/pkg/server"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/version"
)

func parseEnv() (string, int, string, string) {
	host, ok := os.LookupEnv("TA_HOST")
	if !ok {
		log.Printf("TA_HOST not found, defaulting to '0.0.0.0'")
		host = "0.0.0.0"
	}
	var port int
	portStr, ok := os.LookupEnv("TA_PORT")
	port, err := strconv.Atoi(portStr)
	if !ok || err != nil {
		log.Printf("bad ports string %s or TA_PORT not found, defaulting to 23234", portStr)
		port = 23234
	}
	uploadDir, ok := os.LookupEnv("TA_UPLOAD_DIR")
	if !ok {
		log.Printf("TA_UPLOAD_DIR not found, defaulting to './uploads'")
		uploadDir = "./uploads"
	}
	dataFile, ok := os.LookupEnv("TA_DATAFILE")
	if !ok {
		log.Printf("TA_DATAFILE not found, defaulting to 'applicants.csv'")
		dataFile = "applicants.csv"
	}
	return host, port, uploadDir, dataFile
}

func main() {

	log.Print(version.BuildInfo())
	host, port, uploadDir, dataFile := parseEnv()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s, err := server.NewServer(host, uploadDir, dataFile, port)
	if err != nil {
		log.Printf("Cannot create server %v", err)
	}
	s.Start()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	s.Stop(ctx)

}
