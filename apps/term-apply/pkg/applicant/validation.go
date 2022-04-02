package applicant

import (
	"fmt"
	"regexp"
	"strings"
)

func checkForInputErrors(name, email, role string) error {
	var errs []string
	if !isValidName(name) {
		errs = append(errs, "name")
	}
	if !isValidEmail(email) {
		errs = append(errs, "e-mail")
	}
	if !isValidRole(role) {
		errs = append(errs, "role")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d invalid inputs %v", len(errs), strings.Join(errs, ","))
	}
	return nil
}

func isValidName(name string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_` {|}~-]+$")
	return re.MatchString(name)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}

func isValidRole(role string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9 ]+$")
	return re.MatchString(role)
}
