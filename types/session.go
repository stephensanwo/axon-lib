package types

const (
	AUTH_SESSION string = "axon_auth_session"
)

type Session struct {
	SessionId   string
	SessionData UserCache
}
