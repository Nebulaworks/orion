package applicant

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type resumeWatcher struct {
	uploadDir string
	uploaded  []string
	mu        *sync.Mutex
}

func newResumeWatcher(uploadDir string) (*resumeWatcher, error) {
	existing, err := getAllExisting(uploadDir)
	if err != nil {
		return &resumeWatcher{}, err
	}
	return &resumeWatcher{
		uploadDir: uploadDir,
		uploaded:  existing,
		mu:        &sync.Mutex{},
	}, nil
}

func (r *resumeWatcher) isUploaded(userID string) bool {
	for _, u := range r.uploaded {
		if userID == u {
			return true
		}
	}
	if checkForFile(fmt.Sprintf("%s/%s-resume.pdf", r.uploadDir, userID)) {
		log.Printf("Didn't find %s in uploaded list, but did find resume file, adding...", userID)
		r.mu.Lock()
		defer r.mu.Unlock()

		r.uploaded = append(r.uploaded, userID)
		return true
	}
	return false
}

func checkForFile(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getAllExisting(uploadDir string) ([]string, error) {
	log.Printf("Importing existing resumes from %s", uploadDir)
	var results []string
	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		return results, err
	}
	for _, file := range files {
		log.Printf("found %s", file.Name())
		var userID string
		if !file.IsDir() {
			userID = strings.Split(file.Name(), "-resume.pdf")[0]
		}
		if len(userID) > 0 {
			results = append(results, userID)
		}
	}
	return results, nil
}
