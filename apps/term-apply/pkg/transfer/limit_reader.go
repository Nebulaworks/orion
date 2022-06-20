package transfer

/*
limitReader is used to set a max file size on the inbound scp
This code was taken from the charmbracelet/wish project at
https://github.com/charmbracelet/wish/blob/30c2da8825ab7b0eaf195a080f6574c067413e34/scp/limit_reader.go
It looks like they eventually plan to make this usable through
their scp package directly, but since it's currently implemented
as an unexported type, we've had to replicate it here for now.
*/

import (
	"fmt"
	"io"
	"sync"
)

func newLimitReader(r io.Reader, limit int) io.Reader {
	return &limitReader{
		r:    r,
		left: limit,
	}
}

type limitReader struct {
	r io.Reader

	lock sync.Mutex
	left int
}

func (r *limitReader) Read(b []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.left <= 0 {
		return 0, fmt.Errorf("Uploaded file too large")
	}
	if len(b) > r.left {
		b = b[0:r.left]
	}
	n, err := r.r.Read(b)
	r.left -= n
	return n, err
}
