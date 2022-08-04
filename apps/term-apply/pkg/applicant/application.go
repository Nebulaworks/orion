package applicant

import (
	"strconv"
	"time"
)

type application struct {
	appliedDate string
	github      string
	name        string
	email       string
	roleApplied string
	offerGiven  bool
	rejected    bool
}

func NewApplication(github, name, email, roleApplied string) (application, error) {
	if err := checkForInputErrors(name, email, roleApplied); err != nil {
		return application{}, err
	}

	appliedDate := strconv.FormatInt(time.Now().Unix(), 10)
	return application{
		appliedDate: appliedDate,
		github:      github,
		name:        name,
		email:       email,
		roleApplied: roleApplied,
		offerGiven:  false,
		rejected:    false,
	}, nil
}
