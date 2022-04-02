package transfer

// Implements a purposly non-functional version of CopyToClientHandler
// https://pkg.go.dev/github.com/charmbracelet/wish/scp#CopyToClientHandler
// We don't allow clients to copy from the server

import (
	"fmt"
	"io/fs"

	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
)

type nilCopyHandler struct{}

func NewNilCopyHandler() *nilCopyHandler {
	return &nilCopyHandler{}
}

func (n *nilCopyHandler) Glob(_ ssh.Session, s string) ([]string, error) {
	return []string{s}, nil
}

func (n *nilCopyHandler) WalkDir(_ ssh.Session, path string, fn fs.WalkDirFunc) error {
	return nil
}

func (n *nilCopyHandler) NewDirEntry(_ ssh.Session, name string) (*scp.DirEntry, error) {
	return nil, fmt.Errorf("invalid request")
}

func (n *nilCopyHandler) NewFileEntry(_ ssh.Session, name string) (*scp.FileEntry, func() error, error) {
	return nil, nil, fmt.Errorf("invalid request")
}
