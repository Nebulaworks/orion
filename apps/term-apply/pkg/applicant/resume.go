package applicant

import (
	"fmt"

	"github.com/nebulaworks/orion/apps/term-apply/pkg/s3file"
)

var checkS3KeyExists = s3file.S3keyExists

type resumeWatcher struct {
	bucket       string
	resumePrefix string
}

func newResumeWatcher(bucket, resumePrefix string) (*resumeWatcher, error) {
	return &resumeWatcher{
		bucket:       bucket,
		resumePrefix: resumePrefix,
	}, nil
}

func (r *resumeWatcher) isUploaded(userID string) bool {
	key := fmt.Sprintf("%s/%s-resume.pdf", r.resumePrefix, userID)
	return checkS3KeyExists(r.bucket, key)
}
