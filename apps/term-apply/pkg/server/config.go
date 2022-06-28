package server

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	host            string
	port            int
	csvTmpFile      string
	resumeTmpDir    string
	s3Bucket        string
	s3CsvPrefix     string
	s3ResumePrefix  string
	ssmHostKeyParam string
	ssmHostKeyPath  string
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

	s3Bucket, ok := os.LookupEnv("TA_BUCKET")
	if !ok {
		s3Bucket = ""
	}
	log.Printf("TA_BUCKET set to '%s'", s3Bucket)

	s3CsvPrefix, ok := os.LookupEnv("TA_CSV_PREFIX")
	if !ok {
		s3CsvPrefix = "/term-apply/dev/data"
	}
	log.Printf("TA_CSV_PREFIX set to '%s'", s3CsvPrefix)

	s3ResumePrefix, ok := os.LookupEnv("TA_RESUME_PREFIX")
	if !ok {
		s3ResumePrefix = "/term-apply/dev/resumes"
	}
	log.Printf("TA_RESUME_PREFIX set to '%s'", s3ResumePrefix)

	ssmHostKeyParam, ok := os.LookupEnv("TA_SSM_HOST_KEY_PARAM")
	if !ok {
		ssmHostKeyParam = ""
	}

	ssmHostKeyPath, ok := os.LookupEnv("TA_SSM_HOST_KEY_PATH")
	if !ok {
		ssmHostKeyPath = ".ssh/term_info_ed25519"
	}

	return Config{
		host:            host,
		port:            port,
		csvTmpFile:      csvTmpFile,
		resumeTmpDir:    resumeTmpDir,
		s3Bucket:        s3Bucket,
		s3CsvPrefix:     s3CsvPrefix,
		s3ResumePrefix:  s3ResumePrefix,
		ssmHostKeyParam: ssmHostKeyParam,
		ssmHostKeyPath:  ssmHostKeyPath,
	}
}
