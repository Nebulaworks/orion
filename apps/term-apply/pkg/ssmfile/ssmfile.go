package ssmfile

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetParamFromSSM(paramName, path string) error {
	log.Printf("Searching for Path")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir, _ := filepath.Split(path)
		os.MkdirAll(dir, 0700)
	}

	log.Printf("Creating File")
	file, err := os.Create(path)
	if err != nil {
		return err //fmt.Errorf("Cannot create ssh host key file: %s", path)
	}

	defer file.Close()

	log.Printf("Opening new session")
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	log.Printf("Creating SSM Client")
	ssm_client := ssm.New(sess)

	decrypt := true
	log.Printf("Getting param: %s", paramName)
	output, err := ssm_client.GetParameter(&ssm.GetParameterInput{
		Name:           &paramName,
		WithDecryption: &decrypt,
	})
	if err != nil {
		return err
	}

	log.Printf("Writing to file")
	_, err = file.WriteString(*output.Parameter.Value)
	if err != nil {
		return err
	}

	return nil
}
