package applicant

import (
	"log"
	"reflect"
	"sync"
)

type ApplicantManager struct {
	mu          *sync.RWMutex    // wraps applicant slice access
	writeChan   chan application // serializes csv file writes
	resumes     *resumeWatcher
	bucket      string
	dynamoTable string
}

func NewApplicantManager(path, uploadDir, bucket, resumePrefix, dynamodbTable string) (*ApplicantManager, error) {

	writeChan := make(chan application)

	resumes, err := newResumeWatcher(uploadDir, bucket, resumePrefix)
	if err != nil {
		return nil, err
	}

	am := &ApplicantManager{
		mu:          &sync.RWMutex{},
		writeChan:   writeChan,
		resumes:     resumes,
		bucket:      bucket,
		dynamoTable: dynamodbTable,
	}

	go am.writeDynamoItem(path, writeChan)

	return am, nil
}

func (a *ApplicantManager) AddApplicant(github, name, email string, roleApplied int) error {
	roleStr := stringRole(roleApplied)
	newApplication, err := NewApplication(github, name, email, roleStr)
	if err != nil {
		log.Printf(
			"New applicant %s error (%v) with (%s, %s, %s)",
			github,
			err,
			name,
			email,
			roleStr,
		)
		return err
	}

	a.mu.Lock()
	app, err := GetApplication(github, a.dynamoTable)
	if err != nil {
		return err
	}
	if app.github == github {
		if app.rejected || app.offerGiven {
			log.Printf(
				"Found closed application for applicant %s, creating new application (%s, %s, %s)",
				github,
				name,
				email,
				roleStr,
			)
			a.mu.Unlock()
			a.writeChan <- newApplication
		} else {
			newApplication.appliedDate = app.appliedDate
			if reflect.DeepEqual(newApplication, app) {
				log.Printf(
					"Found existing application for applicant %s with identical fields, no changes with (%s, %s, %s)",
					github,
					name,
					email,
					roleStr,
				)
				a.mu.Unlock()
			} else {
				log.Printf(
					"Found existing application for applicant %s with updated fields, updating in place with (%s, %s, %s)",
					github,
					name,
					email,
					roleStr,
				)
				a.mu.Unlock()
				a.writeChan <- newApplication
			}
		}
		return nil
	} else {
		log.Printf("Creating new application for applicant %s with (%s, %s, %s)", github, name, email, roleStr)
		a.mu.Unlock()
		a.writeChan <- newApplication
		return nil
	}
}

func (a *ApplicantManager) writeDynamoItem(filename string, writeChan chan application) error {
	for {
		application := <-writeChan

		if err := PutApplication(application, a.dynamoTable); err != nil {
			log.Printf("Error wrPutApplication %s, %s, %v", application.github, a.dynamoTable, err)
		} else {
			log.Printf("Writing to dynamodb %s, %s", application.github, a.dynamoTable)
		}
	}
}

func (a *ApplicantManager) HasResume(github string) bool {
	return a.resumes.isUploaded(github)
}
