package applicant

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/nebulaworks/orion/apps/term-apply/pkg/dynamofile"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/s3file"
)

type ApplicantManager struct {
	applicants  []applicant
	mu          *sync.RWMutex  // wraps applicant slice access
	writeChan   chan applicant // serializes csv file writes
	resumes     *resumeWatcher
	bucket      string
	csvKey      string
	dynamoTable string
}

func NewApplicantManager(path, uploadDir, bucket, resumePrefix, csvPrefix, dynamodbTable string) (*ApplicantManager, error) {
	if err := openOrCreateFile(path); err != nil {
		return nil, err
	}
	writeChan := make(chan applicant)

	resumes, err := newResumeWatcher(uploadDir, bucket, resumePrefix)
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(path)
	csvKey := fmt.Sprintf("%s/%s", csvPrefix, filename)
	am := &ApplicantManager{
		applicants:  []applicant{},
		mu:          &sync.RWMutex{},
		writeChan:   writeChan,
		resumes:     resumes,
		bucket:      bucket,
		csvKey:      csvKey,
		dynamoTable: dynamodbTable,
	}

	if err := am.readDataFile(path); err != nil {
		return nil, err
	}

	go am.writeDynamoItem(path, writeChan)

	return am, nil
}

func (a *ApplicantManager) AddApplicant(userID, name, email string, role int) error {
	roleStr := stringRole(role)
	newApplicant, err := newApplicant(userID, name, email, roleStr)
	if err != nil {
		log.Printf(
			"New applicant %s error (%v) with (%s, %s, %s)",
			userID,
			err,
			name,
			email,
			roleStr,
		)
		return err
	}

	a.mu.Lock()
	for i, applicant := range a.applicants {
		if applicant.userID == userID {
			if reflect.DeepEqual(newApplicant, applicant) {
				log.Printf(
					"Found existing applicant %s with identical fields, no changes with (%s, %s, %s)",
					userID,
					name,
					email,
					roleStr,
				)
				a.mu.Unlock()
			} else {
				log.Printf(
					"Found existing applicant %s with updated fields, updating in place with (%s, %s, %s)",
					userID,
					name,
					email,
					roleStr,
				)
				a.applicants[i] = newApplicant
				a.mu.Unlock()
				a.writeChan <- newApplicant
			}
			return nil
		}
	}

	log.Printf("Adding new applicant %s with (%s, %s, %s)", userID, name, email, roleStr)
	a.applicants = append(a.applicants, newApplicant)
	a.mu.Unlock()
	a.writeChan <- newApplicant
	return nil
}

func (a *ApplicantManager) readDataFile(filename string) error {
	s3file.CopyFromS3(a.bucket, a.csvKey, filename)
	records, err := readData(filename)
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.applicants = []applicant{}
	for _, record := range records {
		if len(record) != 4 {
			return fmt.Errorf("invalid record %v, fix csv file", record)
		}
		a.applicants = append(a.applicants,
			applicant{
				userID: record[0],
				name:   record[1],
				email:  record[2],
				role:   record[3],
			})
	}
	return nil
}

func (a *ApplicantManager) writeDynamoItem(filename string, writeChan chan applicant) error {
	log.Printf("Writing datafile")
	for {
		var application dynamofile.Application

		applicant := <-writeChan
		a.mu.RLock()

		application = dynamofile.NewApplication(
			"123456791", // TEMP APPLICATION TIME
			applicant.userID,
			applicant.name,
			applicant.email,
			applicant.role,
		)
		a.mu.RUnlock()

		if err := dynamofile.UploadApplication(application, a.dynamoTable); err != nil {
			log.Printf("Error writing to dynamodb %s, %s, %v", application.Github, a.dynamoTable, err)
		}
	}
}

func (a *ApplicantManager) HasResume(userID string) bool {
	return a.resumes.isUploaded(userID)
}
