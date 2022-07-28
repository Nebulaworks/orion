package csvfile

import (
	"encoding/csv"
	"log"
	"os"
)

type Applicant struct {
	Github string
	Name   string
	Email  string
	Role   string
	Date   string
}

func CreateApplicant(data []string) Applicant {
	var app Applicant
	for i, field := range data {
		switch i {
		case 0:
			app.Github = field
		case 1:
			app.Name = field
		case 2:
			app.Email = field
		case 3:
			app.Role = field
		case 4:
			app.Date = field
		}
	}

	return app
}

func ReadFile(path string) []Applicant {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := csv.NewReader(f)

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var apps []Applicant
	for _, app := range data {
		apps = append(apps, CreateApplicant(app))
	}

	return apps
}
