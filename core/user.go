package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	axon_coredb "github.com/stephensanwo/axon-lib/coredb"
	github "github.com/stephensanwo/axon-lib/github"
	axon_types "github.com/stephensanwo/axon-lib/types"
	"golang.org/x/oauth2"
)

type User struct {
	AwsSession *aws_session.Session
}

func (u * User)CreateUser(a *axon_types.AxonContext, token *oauth2.Token) (*axon_types.User, error) {
	ctx := context.Background()
	
	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(u.AwsSession)
	if err != nil {
		return nil, err
	}
	
	// Get Authenticated User
	github_client := github.GetGithubClient(ctx, token.AccessToken)
	github_user, _, err := github.GetAuthenticatedUser(ctx, github_client)
	if err != nil {
		return nil, err
	}

	// Create User Object
	var user axon_types.User

	// Query the DynamoDB table for the user using the email from the Auth Client Response
	email := *github_user.Email
	result, err := db.QueryDatabase(axon_types.AXON_TABLE, fmt.Sprintf("USER#%s",email), &email)
	
	if err != nil {
		return nil, errors.New("could not authenticate user - " + err.Error())
	}

	if len(result.Item) > 0 { 
		// If the user exists, return the user
		err := dynamodbattribute.UnmarshalMap(result.Item, &user)
		return &user, err
	} else {
		// If the user does not exist, create a new user
		user = axon_types.User{
			UserId:    uuid.New().String(), // Using hex representation of ObjectID for DynamoDB
			Email:     *github_user.Email,
			UserName:  *github_user.Login,
			FirstName: strings.Split(*github_user.Name, " ")[0],
			LastName:  strings.Split(*github_user.Name, " ")[1],
			Avatar:    *github_user.AvatarURL,
		}

		err = db.MutateDatabase(axon_types.AXON_TABLE, fmt.Sprintf("USER#%s", user.Email), user.Email, &user)

		if err != nil {
			return nil, err
		}

		return &user, nil
	}
}

func (u *User) GetAuthenticatedUserData(a *axon_types.AxonContext) (axon_types.Session, error) {
	userSession := axon_types.Session{}

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(u.AwsSession)
	if err != nil {
		return userSession, errors.New("Error fetching user session" + err.Error())
	}

	// Find user session in the cache
	result, err := db.QueryDatabase(axon_types.AXON_USER_SESSION_TABLE, fmt.Sprintf("SESSION#%s", a.SessionId), &a.SessionId)

	if err != nil {
		return userSession, errors.New("Error fetching user session" + err.Error())
	}

	// Unmarshal the DynamoDB item into a Session struct
	dynamodbattribute.UnmarshalMap(result.Item, &userSession)

	return userSession, err

}
