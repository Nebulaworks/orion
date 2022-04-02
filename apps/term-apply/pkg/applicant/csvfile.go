package applicant

import (
	"encoding/csv"
	"os"
)

func openOrCreateFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func readData(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func writeData(filename string, records [][]string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		return err
	}
	return nil
}
