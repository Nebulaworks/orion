package applicant

type applicant struct {
	userID string
	name   string
	email  string
	role   string
}

func newApplicant(userID, name, email, role string) (applicant, error) {
	err := checkForInputErrors(name, email, role)
	if err != nil {
		return applicant{}, err
	}
	return applicant{
		userID: userID,
		name:   name,
		email:  email,
		role:   role,
	}, nil
}
