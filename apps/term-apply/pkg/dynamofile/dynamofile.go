package dynamofile

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Application struct {
	Applied_date string
	Github       string
	Name         string
	Email        string
	Role_applied string
}

func NewApplication(applied_date, github, name, email, role_applied string) Application {
	return Application{
		Applied_date: applied_date,
		Github:       github,
		Name:         name,
		Email:        email,
		Role_applied: role_applied,
	}
}

func applicationFromItem(item map[string]*dynamodb.AttributeValue) Application {
	var app Application

	applied_date, exists := item["applied_date"]
	if exists {
		app.Applied_date = *applied_date.N
	}

	email, exists := item["email"]
	if exists {
		app.Email = *email.S
	}
	name, exists := item["name"]
	if exists {
		app.Name = *name.S
	}
	github, exists := item["github"]
	if exists {
		app.Github = *github.S
	}
	role_applied, exists := item["role_applied"]
	if exists {
		app.Role_applied = *role_applied.S
	}

	return app
}

// Returns the provided user's most recent application
func GetApplication(user, table string) (Application, error) {
	var app Application

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "sandbox",
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
func GetApplications(user, table string) ([]Application, error) {
	var app []Application

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "sandbox",
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

func UploadApplication(app Application, table string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "sandbox",
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	svc.PutItem(&dynamodb.PutItemInput{
		TableName: &table,
		Item: map[string]*dynamodb.AttributeValue{
			"applied_date": {
				N: &app.Applied_date,
			},
			"github": {
				S: &app.Github,
			},
			"name": {
				S: &app.Name,
			},
			"email": {
				S: &app.Email,
			},
			"role_applied": {
				S: &app.Role_applied,
			},
		},
	})

	return nil
}
