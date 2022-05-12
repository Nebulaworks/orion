package applicant

import (
	"fmt"
	"log"
	"sync"

	"github.com/nebulaworks/orion/apps/term-apply/pkg/s3file"
)

type resumeWatcher struct {
	uploadDir    string
	uploaded     []string
	mu           *sync.Mutex
	bucket       string
	resumePrefix string
}

func newResumeWatcher(uploadDir, bucket, resumePrefix string) (*resumeWatcher, error) {
	return &resumeWatcher{
		uploadDir:    uploadDir,
		uploaded:     []string{},
		mu:           &sync.Mutex{},
		bucket:       bucket,
		resumePrefix: resumePrefix,
	}, nil
}

func (r *resumeWatcher) isUploaded(userID string) bool {
	for _, u := range r.uploaded {
		if userID == u {
			return true
		}
	}
	key := fmt.Sprintf("%s/%s-resume.pdf", r.resumePrefix, userID)
	if s3file.S3keyExists(r.bucket, key) {
		log.Printf("Didn't find %s in uploaded list, but did find resume file in s3, adding...", userID)
		r.mu.Lock()
		defer r.mu.Unlock()

		r.uploaded = append(r.uploaded, userID)
		return true
	}
	return false
}
