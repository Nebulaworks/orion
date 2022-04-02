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
	"github.com/gliderlabs/ssh"
)

type copyFromClientHandler struct {
	root string
}

func NewCopyFromClientHandler(root string) *copyFromClientHandler {
	rootInfo, err := os.Stat(root)
	if os.IsNotExist(err) {
		log.Fatal(root + " doesn't exist")
	}
	if !rootInfo.IsDir() {
		log.Fatal(root + " is not a directory")
	}
	return &copyFromClientHandler{
		root: filepath.Clean(root),
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
	fin := s.User()
	f, err := os.OpenFile(c.prefixed(fin+"-"+"resume.pdf"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, entry.Mode)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %q: %w", entry.Filepath, err)
	}

	const BYTES_TEN_MEGABYTES = 10485760
	lr := newLimitReader(entry.Reader, BYTES_TEN_MEGABYTES)

	written, err := io.Copy(f, lr)
	if err != nil {
		return 0, fmt.Errorf("failed to write file: %q: %w", entry.Filepath, err)
	}
	return written, c.chtimes(entry.Filepath, entry.Mtime, entry.Atime)
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
		return fmt.Errorf("failed to chtimes: %q: %w", path, err)
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
