package transfer

// Implements CopyFromClientHanlder interface
// https://pkg.go.dev/github.com/charmbracelet/wish/scp#CopyFromClientHandler

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/wish/scp"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gliderlabs/ssh"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/s3file"
)

type copyFromClientHandler struct {
	root         string
	bucket       string
	resumePrefix string
}

func NewCopyFromClientHandler(root, bucket, resumePrefix string) *copyFromClientHandler {
	rootInfo, err := os.Stat(root)
	if os.IsNotExist(err) {
		log.Fatal(root + " doesn't exist")
	}
	if !rootInfo.IsDir() {
		log.Fatal(root + " is not a directory")
	}
	return &copyFromClientHandler{
		root:         filepath.Clean(root),
		bucket:       bucket,
		resumePrefix: resumePrefix,
	}
}

func (c *copyFromClientHandler) Mkdir(s ssh.Session, entry *scp.DirEntry) error {
	//username is more appropriate since a user could have multiple keys tied to github
	fin := s.User()
	if err := os.Mkdir(c.prefixed(fin+"/"+entry.Filepath), entry.Mode); err != nil {
		return fmt.Errorf("failed to create dir: %q: %w", entry.Filepath, err)
	}
	return c.chtimes(entry.Filepath, entry.Mtime, entry.Atime)
}

func (c *copyFromClientHandler) Write(s ssh.Session, entry *scp.FileEntry) (int64, error) {

	user := s.User()
	filename := fmt.Sprintf("%s-resume.pdf", user)
	fileKey := fmt.Sprintf("%s/%s", c.resumePrefix, filename)
	localFile := fmt.Sprintf("%s/%s", c.root, filename)

	// Check if resume has been uploaded
	_, err := os.Stat(c.prefixed(filename))

	if s3file.S3keyExists(c.bucket, fileKey) {
		log.Printf("Resume %s already exists: uploading replacement resume for %s.", filename, user)
	} else {
		log.Printf("Resume %s has not been uploaded: initial upload for %s.", filename, user)
	}

	// Write scp input to temp file for validity checking
	t, err := os.OpenFile(c.prefixed("temp"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, entry.Mode)
	if err != nil {
		return 0, fmt.Errorf("\nfailed to open file: %q: %w", entry.Filepath, err)
	}

	// check size constraint usin limit reader
	const BYTES_TEN_MEGABYTES = 10485760

	lr := newLimitReader(entry.Reader, BYTES_TEN_MEGABYTES)

	written, err := io.Copy(t, lr)
	if err != nil {
		log.Printf("error writing file %s, %v", filename, err)
		return 0, fmt.Errorf("\nProvided file is too large. Maximum size is 10MB\n%s", getLastResumeStatus(c.bucket, fileKey, user))
	}

	// validate contents of uploaded file
	tempFile := fmt.Sprintf("%s/%s", c.root, "temp")
	mtype, err := mimetype.DetectFile(tempFile)
	if err != nil {
		log.Printf("error checking pdf validity: %v", err)
		return 0, fmt.Errorf("error occured while pdf validity")
	}
	if !(mtype.String() == "application/pdf" || mtype.String() == "application/x-pdf") {
		log.Printf("Provided file failed PDF validity check")
		return 0, fmt.Errorf("\nProvided file failed PDF validity check\n%s", getLastResumeStatus(c.bucket, fileKey, user))
	}
	log.Printf("Provided file passed PDF vaildity check with type %s", mtype.String())

	// copy valid upload file to final file path
	t, _ = os.Open(tempFile)

	f, err := os.OpenFile(c.prefixed(filename), os.O_TRUNC|os.O_RDWR|os.O_CREATE, entry.Mode)
	if err != nil {
		return 0, fmt.Errorf("\nfailed to open file: %q: %w", entry.Filepath, err)
	}

	written, err = io.Copy(f, t)
	if err != nil {
		log.Printf("error writing file %s, %v", filename, err)
		return 0, fmt.Errorf("\nfailed to write file: %q", entry.Filepath)
	}

	// delete temporary file
	err = os.Remove(tempFile)
	if err != nil {
		log.Printf("failed to delete temp upload file")
	}

	// copy validated file to s3

	if err := s3file.CopyToS3(c.bucket, localFile, fileKey); err != nil {
		log.Printf("error writing to s3 %s, %s, %v", filename, fileKey, err)
		return 0, fmt.Errorf("\nfailed to write file: %q", entry.Filepath)
	}

	return written, c.chtimes(entry.Filepath, entry.Mtime, entry.Atime)
}

func getLastResumeStatus(bucket, fileKey string, user string) string {
	var sts string
	if s3file.S3keyExists(bucket, fileKey) {
		sts = fmt.Sprintf("Last valid upload by user %s on %s", user, s3file.S3keyLastModified(bucket, fileKey))
	} else {
		sts = fmt.Sprintf("No valid file has been uploaded by user %s", user)
	}
	return sts
}

func (c *copyFromClientHandler) chtimes(path string, mtime, atime int64) error {
	if mtime == 0 || atime == 0 {
		return nil
	}
	if err := os.Chtimes(
		c.prefixed(path),
		time.Unix(atime, 0),
		time.Unix(mtime, 0),
	); err != nil {
		return fmt.Errorf("\nfailed to chtimes: %q: %w", path, err)
	}
	return nil
}

func (c *copyFromClientHandler) prefixed(path string) string {
	path = filepath.Clean(path)
	if strings.HasPrefix(path, c.root) {
		return path
	}
	return filepath.Join(c.root, path)
}
