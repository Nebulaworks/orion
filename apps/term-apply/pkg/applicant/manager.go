package applicant

import (
	"log"
	"reflect"
	"sync"
)

type ApplicantManager struct {
	mu            *sync.RWMutex          // wraps applicant slice access
	writeChan     chan applicationPacket // serializes application uploads
	resumes       *resumeWatcher
	bucket        string
	dynamodbTable string
	dynamodbIndex string
}

type applicationPacket struct {
	app           application
	prevEmail     string
	updateInPlace bool
}

func NewApplicantManager(uploadDir, bucket, resumePrefix, dynamodbTable, dynamodbIndex string) (*ApplicantManager, error) {

	writeChan := make(chan applicationPacket)

	resumes, err := newResumeWatcher(uploadDir, bucket, resumePrefix)
	if err != nil {
		return nil, err
	}

	am := &ApplicantManager{
		mu:            &sync.RWMutex{},
		writeChan:     writeChan,
		resumes:       resumes,
		bucket:        bucket,
		dynamodbTable: dynamodbTable,
		dynamodbIndex: dynamodbIndex,
	}

	go am.writeDynamoItem(writeChan)

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

	app, err := GetApplication(github, a.dynamodbTable, a.dynamodbIndex)
	if _, ok := err.(*emptyResultError); ok {
		log.Printf("Creating new application for applicant %s with (%s, %s, %s)", github, name, email, roleStr)
		a.mu.Unlock()
		a.writeChan <- applicationPacket{app: newApplication, updateInPlace: false}
		return nil
	} else if err != nil {
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
			a.writeChan <- applicationPacket{app: newApplication, updateInPlace: false}
		} else {
			newApplication.appliedDate = app.appliedDate
			if reflect.DeepEqual(newApplication, app) {
				log.Printf(
					"Found open application for applicant %s with identical fields, no changes with (%s, %s, %s)",
					github,
					name,
					email,
					roleStr,
				)
				a.mu.Unlock()
			} else {
				log.Printf(
					"Found open application for applicant %s with updated fields, updating in place with (%s, %s, %s)",
					github,
					name,
					email,
					roleStr,
				)
				a.mu.Unlock()
				a.writeChan <- applicationPacket{app: newApplication, prevEmail: app.email, updateInPlace: true}
			}
		}
		return nil
	} else {
		log.Printf("Creating new application for applicant %s with (%s, %s, %s)", github, name, email, roleStr)
		a.mu.Unlock()
		a.writeChan <- applicationPacket{app: newApplication, updateInPlace: false}
		return nil
	}
}

func (a *ApplicantManager) writeDynamoItem(writeChan chan applicationPacket) {
	for {
		packet := <-writeChan

		if packet.updateInPlace {
			if packet.prevEmail == packet.app.email {
				log.Printf("Updating dynamodb record in %s for %s", a.dynamodbTable, packet.app.github)
				if err := PutApplication(packet.app, a.dynamodbTable); err != nil {
					log.Printf("Error uploading application for %s in %s: %v", packet.app.github, a.dynamodbTable, err)
				} else {
					log.Printf("Succesful write")
				}
			} else {
				log.Printf("Deleting and recreating dynamodb record in %s for %s", a.dynamodbTable, packet.app.github)
				if err := UpdateApplication(packet.app, packet.prevEmail, a.dynamodbTable); err != nil {
					log.Printf("Error uploading application for %s in %s: %v", packet.app.github, a.dynamodbTable, err)
				} else {
					log.Printf("Succesful write")
				}
			}
		} else {
			log.Printf("Writing new record to dynamodb in %s for %s", a.dynamodbTable, packet.app.github)
			if err := PutApplication(packet.app, a.dynamodbTable); err != nil {
				log.Printf("Error uploading application for %s in %s: %v", packet.app.github, a.dynamodbTable, err)
			} else {
				log.Printf("Succesful write")
			}
		}
	}
}

func (a *ApplicantManager) HasResume(github string) bool {
	return a.resumes.isUploaded(github)
}
