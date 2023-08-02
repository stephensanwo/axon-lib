package types

import (
	"context"

	oauthTypes "golang.org/x/oauth2"
)

type AxonContext struct {
	Context            context.Context
	Settings           Settings           `json:"settings"`
	Oauth              oauthTypes.Config      `json:"oauth"`
	SessionId          string
}
type AxonContextKey string
