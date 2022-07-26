package applicant

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
func GetApplication(user, table string) (application, error) {
	var app application

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

	app = applicationFromItem(result.Items[len(result.Items)-1])

	return app, nil
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
		TableName: &table,
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
