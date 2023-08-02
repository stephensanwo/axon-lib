package github

import (
	"context"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func GetAuthenticatedUser(ctx context.Context, client *github.Client) (*github.User, *github.Response, error) {
	// Empty string gets the authenticated user
	user, response, error := client.Users.Get(ctx, "")

	return user, response, error
}
