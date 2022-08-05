package server

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	host            string
	port            int
	resumeTmpDir    string
	s3Bucket        string
	s3ResumePrefix  string
	dynamodbTable   string
	dynamodbIndex   string
	ssmHostKeyParam string
	hostKeyPath     string
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

	s3Bucket, ok := os.LookupEnv("TA_BUCKET")
	if !ok {
		s3Bucket = ""
	}
	log.Printf("TA_BUCKET set to '%s'", s3Bucket)

	s3ResumePrefix, ok := os.LookupEnv("TA_RESUME_PREFIX")
	if !ok {
		s3ResumePrefix = "/term-apply/dev/resumes"
	}
	log.Printf("TA_RESUME_PREFIX set to '%s'", s3ResumePrefix)

	dynamodbTable, ok := os.LookupEnv("TA_DYNAMODB_TABLE")
	if !ok {
		dynamodbTable = ""
	}
	log.Printf("TA_DYNAMODB_TABLE set to '%s'", dynamodbTable)

	dynamodbIndex, ok := os.LookupEnv("TA_DYNAMODB_GSI")
	if !ok {
		dynamodbIndex = ""
	}
	log.Printf("TA_DYNAMODB_GSI set to '%s'", dynamodbIndex)

	ssmHostKeyParam, ok := os.LookupEnv("TA_SSM_HOST_KEY_PARAM")
	if !ok {
		ssmHostKeyParam = ""
	}
	log.Printf("TA_SSM_HOST_KEY_PARAM set to '%s'", ssmHostKeyParam)
	hostKeyPath, ok := os.LookupEnv("TA_HOST_KEY_PATH")
	if !ok {
		hostKeyPath = ".ssh/term_info_ed25519"
	}
	log.Printf("TA_HOST_KEY_PATH set to '%s'", hostKeyPath)

	return Config{
		host:            host,
		port:            port,
		resumeTmpDir:    resumeTmpDir,
		s3Bucket:        s3Bucket,
		s3ResumePrefix:  s3ResumePrefix,
		dynamodbTable:   dynamodbTable,
		dynamodbIndex:   dynamodbIndex,
		ssmHostKeyParam: ssmHostKeyParam,
		hostKeyPath:     hostKeyPath,
	}
}
