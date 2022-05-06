package applicant

import (
	"fmt"
	"log"
	"sync"

	"github.com/nebulaworks/orion/apps/term-apply/pkg/s3file"
)

const RESUME_ROOT_KEY = "/term-apply/resumes/"

type resumeWatcher struct {
	uploadDir string
	uploaded  []string
	mu        *sync.Mutex
}

func newResumeWatcher(uploadDir string) (*resumeWatcher, error) {
	return &resumeWatcher{
		uploadDir: uploadDir,
		uploaded:  []string{},
		mu:        &sync.Mutex{},
	}, nil
}

func (r *resumeWatcher) isUploaded(userID string) bool {
	for _, u := range r.uploaded {
		if userID == u {
			return true
		}
	}
	if s3file.S3keyExists(fmt.Sprintf("%s/%s-resume.pdf", RESUME_ROOT_KEY, userID)) {
		log.Printf("Didn't find %s in uploaded list, but did find resume file in s3, adding...", userID)
		r.mu.Lock()
		defer r.mu.Unlock()

		r.uploaded = append(r.uploaded, userID)
		return true
	}
	return false
}
