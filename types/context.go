package types

import (
	"context"

	"golang.org/x/oauth2"
)

type AxonContext struct {
	Context            context.Context
	Settings           Settings           `json:"settings"`
	Oauth              oauth2.Config      `json:"oauth"`
	SessionId          string
}
type AxonContextKey string
