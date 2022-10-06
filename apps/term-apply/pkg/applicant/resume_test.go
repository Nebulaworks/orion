package applicant

import "testing"

func swapS3KeyExists(value func(string, string) bool) {
	checkS3KeyExists = value
}

func shouldExist(string, string) bool {
	return true
}

func shouldNotExist(string, string) bool {
	return false
}

func TestIsUploaded(t *testing.T) {
	defer swapS3KeyExists(checkS3KeyExists)

	watcher, _ := newResumeWatcher("notarealbucket", "fakeprefix")

	swapS3KeyExists(shouldExist)
	if !watcher.isUploaded("nothing") {
		t.Fatalf("It should show as uploaded")
	}

	swapS3KeyExists(shouldNotExist)
	if watcher.isUploaded("nothing") {
		t.Fatalf("It should not show as uploaded")
	}
}
