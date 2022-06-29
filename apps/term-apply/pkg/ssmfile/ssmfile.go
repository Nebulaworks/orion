package ssmfile

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetParamFromSSM(paramName, path string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	ssm_client := ssm.New(sess)

	decrypt := true
	log.Printf("Getting parameter %s from ssm", paramName)
	output, err := ssm_client.GetParameter(&ssm.GetParameterInput{
		Name:           &paramName,
		WithDecryption: &decrypt,
	})
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("File %s does not exist... creating", path)
		dir, _ := filepath.Split(path)
		os.MkdirAll(dir, 0700)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.WriteString(*output.Parameter.Value)
	if err != nil {
		return err
	}

	return nil
}
