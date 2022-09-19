package applicant

import (
	"log"
	"reflect"
	"sync"
)

type writeState int64

const (
	newApp      writeState = 0
	updateApp   writeState = 1
	recreateApp writeState = 2
)

type ApplicantManager struct {
	mu            *sync.RWMutex          // wraps applicant slice access
	writeChan     chan applicationPacket // serializes application uploads
	resumes       *resumeWatcher
	dynamodbTable string
	dynamodbIndex string
}

type applicationPacket struct {
	app        application
	prevEmail  string
	writeState writeState
}

func NewApplicantManager(bucket, resumePrefix, dynamodbTable, dynamodbIndex string) (*ApplicantManager, error) {

	writeChan := make(chan applicationPacket)

	resumes, err := newResumeWatcher(bucket, resumePrefix)
	if err != nil {
		return nil, err
	}

	am := &ApplicantManager{
		mu:            &sync.RWMutex{},
		writeChan:     writeChan,
		resumes:       resumes,
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

	// No application exists: new applicant
	if _, ok := err.(*emptyResultError); ok {
		log.Printf("Creating new application for applicant %s with (%s, %s, %s)", github, name, email, roleStr)
		a.mu.Unlock()
		a.writeChan <- applicationPacket{app: newApplication, writeState: newApp}
		return nil
	} else if err != nil {
		return err
	}

	// Closed application exists: returning applicant
	if app.rejected || app.offerGiven {
		log.Printf(
			"Found closed application for applicant %s, creating new application (%s, %s, %s)",
			github,
			name,
			email,
			roleStr,
		)
		a.mu.Unlock()
		a.writeChan <- applicationPacket{app: newApplication, writeState: newApp}

		return nil
	}

	// Keep original applied date for open applications
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

		return nil
	}

	// Updated application with unchanged email
	if newApplication.email == app.email {
		log.Printf(
			"Found open application for applicant %s with updated fields, updating in place with (%s, %s, %s)",
			github,
			name,
			email,
			roleStr,
		)
		a.mu.Unlock()
		a.writeChan <- applicationPacket{app: newApplication, writeState: updateApp}

		return nil
	}

	// Updated application with modified email (recreate necessary)
	log.Printf(
		"Found open application for applicant %s with updated fields, updating in place with (%s, %s, %s)",
		github,
		name,
		email,
		roleStr,
	)
	a.mu.Unlock()
	a.writeChan <- applicationPacket{app: newApplication, prevEmail: app.email, writeState: recreateApp}

	return nil
}

func (a *ApplicantManager) writeDynamoItem(writeChan chan applicationPacket) {
	for {
		packet := <-writeChan

		switch packet.writeState {
		case newApp:
			{
				log.Printf("Writing new record to dynamodb in %s for %s", a.dynamodbTable, packet.app.github)
				if err := PutApplication(packet.app, a.dynamodbTable); err != nil {
					log.Printf("Error uploading application for %s in %s: %v", packet.app.github, a.dynamodbTable, err)
				} else {
					log.Printf("Succesful write")
				}
			}
		case updateApp:
			{
				log.Printf("Updating dynamodb record in %s for %s", a.dynamodbTable, packet.app.github)
				if err := UpdateApplication(packet.app, a.dynamodbTable); err != nil {
					log.Printf("Error uploading application for %s in %s: %v", packet.app.github, a.dynamodbTable, err)
				} else {
					log.Printf("Succesful write")
				}
			}
		case recreateApp:
			{
				log.Printf("Deleting and recreating dynamodb record in %s for %s", a.dynamodbTable, packet.app.github)
				if err := RecreateApplication(packet.app, packet.prevEmail, a.dynamodbTable); err != nil {
					log.Printf("Error uploading application for %s in %s: %v", packet.app.github, a.dynamodbTable, err)
				} else {
					log.Printf("Succesful write")
				}
			}
		default:
			{
				log.Printf("Invalid write state")
			}
		}
	}
}

func (a *ApplicantManager) HasResume(github string) bool {
	return a.resumes.isUploaded(github)
}
