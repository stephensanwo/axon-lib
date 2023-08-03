package session

import (
	"fmt"

	aws_session "github.com/aws/aws-sdk-go/aws/session"
	axon_coredb "github.com/stephensanwo/axon-lib/coredb"
	axon_types "github.com/stephensanwo/axon-lib/types"

	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type SessionManager struct {
	CookieName string
	SessionId  string
	AwsSession *aws_session.Session
}

// type Session struct {
// 	SessionId   string
// 	SessionData axon_types.UserCache
// }

func (s SessionManager) CreateSession(w http.ResponseWriter, a *axon_types.AxonContext, sessionData *axon_types.Session) {
	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(s.AwsSession)
	
	if err != nil {
		log.Panicln("Error creating user session" + err.Error())
	}

	expiration := time.Now().Add(365 * 24 * 12 * time.Hour)
	cookie := http.Cookie{Name: s.CookieName, Value: s.SessionId, Path: "/", HttpOnly: true, Expires: expiration}
	http.SetCookie(w, &cookie)

	// Cache Session Data
	err = db.CacheData(axon_coredb.AXON_USER_SESSION_TABLE, fmt.Sprintf("SESSION#%s", s.SessionId), s.SessionId, sessionData, 12 * 60 * 60)

	if err != nil {
		log.Panicln("Error saving session in cache")
	}

}

func NewSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
