package applicant

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

type ApplicantManager struct {
	applicants []applicant
	mu         *sync.RWMutex    // wraps applicant slice access
	writeChan  chan []applicant // serializes csv file writes
	resumes    *resumeWatcher
}

func NewApplicantManager(filename, uploadDir string) (*ApplicantManager, error) {
	if err := openOrCreateFile(filename); err != nil {
		return nil, err
	}
	writeChan := make(chan []applicant)

	resumes, err := newResumeWatcher(uploadDir)
	if err != nil {
		return nil, err
	}

	am := &ApplicantManager{
		applicants: []applicant{},
		mu:         &sync.RWMutex{},
		writeChan:  writeChan,
		resumes:    resumes,
	}

	if err := am.readDataFile(filename); err != nil {
		return nil, err
	}

	go am.writeDataFile(filename, writeChan)

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
				a.writeChan <- a.applicants
			}
			return nil
		}
	}

	log.Printf("Adding new applicant %s with (%s, %s, %s)", userID, name, email, roleStr)
	a.applicants = append(a.applicants, newApplicant)
	a.mu.Unlock()
	a.writeChan <- a.applicants
	return nil
}

func (a *ApplicantManager) readDataFile(filename string) error {
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

func (a *ApplicantManager) writeDataFile(filename string, writeChan chan []applicant) error {
	for {
		updatedApplicants := <-writeChan
		a.mu.RLock()
		records := [][]string{}
		for _, applicant := range updatedApplicants {
			records = append(records, []string{
				applicant.userID,
				applicant.name,
				applicant.email,
				applicant.role,
			})
		}
		a.mu.RUnlock()

		if err := writeData(filename, records); err != nil {
			log.Printf("error writing file %s, %v", filename, err)
			return err
		}
	}
}

func (a *ApplicantManager) HasResume(userID string) bool {
	return a.resumes.isUploaded(userID)
}
