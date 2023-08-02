package types

import "net/http"

const (
	PrivateRoute = "private"
	PublicRoute  = "public"
)

type Route struct {
	Path    string                                                 `json:"path"`
	Auth    string                                                 `json:"auth"`
	Handler func(http.ResponseWriter, *http.Request, *AxonContext) `json:"handler"`
	Method  string                                                 `json:"method"`
}
