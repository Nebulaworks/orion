package applicant

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/*

The relevant expected schema for the DynamoDB table is as follows:

> applied_date: unix time - number - Sort DDB Key <br>
> email: candy@date.com - string - Primary/Partition DDB Key <br>
> github: candydate100 - string - Secondary Global Index <br>
> name: Candy Date - string <br>
> role_applied: sr. software engineer - string <br>
> offer_given: bool <br>
> rejected: bool <br>

term-apply users have the ability to modify their email after submitting
an application. If this happens, the existing record will be deleted after
copying all other data

*/

type emptyResultError struct {
	user string
}

func newEmptyResult(user string) *emptyResultError {
	var err emptyResultError
	err.user = user
	return &err
}

func (err *emptyResultError) Error() string {
	return fmt.Sprintf("No uploads found for user %s", err.user)
}

func applicationFromItem(item map[string]*dynamodb.AttributeValue) application {
	var app application

	appliedDate, exists := item["applied_date"]
	if exists {
		app.appliedDate = *appliedDate.N
	}
	github, exists := item["github"]
	if exists {
		app.github = *github.S
	}
	email, exists := item["email"]
	if exists {
		app.email = *email.S
	}
	name, exists := item["name"]
	if exists {
		app.name = *name.S
	}
	roleApplied, exists := item["role_applied"]
	if exists {
		app.roleApplied = *roleApplied.S
	}
	offferGiven, exists := item["offer_given"]
	if exists {
		app.offerGiven = *offferGiven.BOOL
	}
	rejected, exists := item["rejected"]
	if exists {
		app.rejected = *rejected.BOOL
	}

	return app
}

// Returns the provided user's most recent application
func GetApplication(user, table, index string) (application, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return application{}, err
	}

	svc := dynamodb.New(sess)

	result, err := svc.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		IndexName:              aws.String(index),
		KeyConditionExpression: aws.String("github = :github"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":github": {S: aws.String(user)},
		},
	})
	if err != nil {
		return application{}, err
	}

	if len(result.Items) > 0 {
		app := applicationFromItem(result.Items[len(result.Items)-1])
		return app, nil
	} else {
		log.Printf("No applications found for %s", user)
		return application{}, newEmptyResult(user)
	}
}

// Returns all the provided user's applications
func GetApplications(user, table string) ([]application, error) {
	var app []application

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return app, err
	}

	svc := dynamodb.New(sess)

	result, err := svc.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("github = :github"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":github": {S: aws.String(user)},
		},
	})
	if err != nil {
		return app, err
	}

	for _, attr := range result.Items {
		app = append(app, applicationFromItem(attr))
	}

	return app, nil
}

func PutApplication(app application, table string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"applied_date": {
				N: &app.appliedDate,
			},
			"github": {
				S: &app.github,
			},
			"name": {
				S: &app.name,
			},
			"email": {
				S: &app.email,
			},
			"role_applied": {
				S: &app.roleApplied,
			},
		},
	})

	return nil
}

func UpdateApplication(app application, prevEmail, table string) error {
	// Create DynamoDB Session
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	// Query record to update to ensure all unchanged values are preserved
	result, err := svc.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("applied_date = :a and email = :e"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {N: aws.String(app.appliedDate)},
			":e": {S: aws.String(prevEmail)},
		},
	})
	if err != nil {
		return err
	} else if len(result.Items) == 0 {
		return fmt.Errorf("expected record not found: %s (applied at %s)", prevEmail, app.appliedDate)
	}

	record := result.Items[0]

	// Update values in record
	record["email"] = &dynamodb.AttributeValue{
		S: aws.String(app.email),
	}
	record["name"] = &dynamodb.AttributeValue{
		S: aws.String(app.name),
	}
	record["role_applied"] = &dynamodb.AttributeValue{
		S: aws.String(app.roleApplied),
	}

	items := []*dynamodb.TransactWriteItem{
		// Queue deleting original record
		&dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				TableName: aws.String(table),
				Key: map[string]*dynamodb.AttributeValue{
					"applied_date": {
						N: aws.String(app.appliedDate),
					},
					"email": {
						S: aws.String(prevEmail),
					},
				},
			},
		},

		// Queue recreating record with updated values
		&dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(table),
				Item:      record,
			},
		},
	}

	_, err = svc.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	})
	if err != nil {
		return err
	}

	return nil
}
