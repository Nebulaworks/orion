package s3file

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CopyFromS3(bucket, key, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create %s", filename)
	}
	defer file.Close()

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			return fmt.Errorf(err.Error())
		}
	}
	log.Printf("Downloaded %d bytes; %s to %s", numBytes, key, filename)
	return nil
}

func CopyToS3(bucket, filename, key string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	content, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:    content,
		Bucket:  aws.String(bucket),
		Key:     aws.String(key),
		Tagging: aws.String("owner=term-apply"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			return fmt.Errorf(err.Error())
		}
	}
	log.Printf("S3 success %v", result)
	return nil
}

func S3keyExists(bucket, key string) bool {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Printf("%v", err)
		return false
	}
	svc := s3.New(sess)
	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				log.Printf("NotFound")
				return false
			default:
				log.Printf("%v", err)
				return false
			}
		}
		log.Printf("%v", err)
		return false
	}
	return true
}
