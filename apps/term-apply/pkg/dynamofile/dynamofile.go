package dynamofile

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Application struct {
	Applied_date int
	Github       string
	Name         string
	Email        string
	Role_applied string
}

func applicationFromItem(item map[string]*dynamodb.AttributeValue) Application {
	var app Application

	applied_date, exists := item["applied_date"]
	if exists {
		i_64, _ := strconv.ParseInt(*(applied_date.N), 10, 32)
		app.Applied_date = int(i_64)
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
func GetApplication(user string, table string) (Application, error) {
	var app Application

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "sandbox",
	})
	if err != nil {
		return app, err
	}

	svc := dynamodb.New(sess)

	// query := "email = :candy@date.com"

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
func GetApplications(user string, table string) ([]Application, error) {
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

func NewApplication() {

}

func SetApplicationName() {

}

func SetApplicationEmail() {

}

func SetApplicationJob() {

}
