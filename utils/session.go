package utils

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

// Session store (use a secure key in production)
var store = sessions.NewCookieStore([]byte("your-secret-key"))

// Register custom types for encoding in sessions
func init() {
	gob.Register(map[string]string{})
}

// StoreUserSession saves user info in a session
func StoreUserSession(w http.ResponseWriter, r *http.Request, userID string, token string) error {
	session, _ := store.Get(r, "user-session")
	session.Values["user_id"] = userID
	session.Values["token"] = token
	return session.Save(r, w)
}

// GetUserSession retrieves user session data
func GetUserSession(r *http.Request) (string, string, error) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return "", "", err
	}

	userID, ok1 := session.Values["user_id"].(string)
	token, ok2 := session.Values["token"].(string)

	if !ok1 || !ok2 {
		return "", "", err
	}

	return userID, token, nil
}

// RemoveUserSession clears the session
func RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Options.MaxAge = -1 // Expire session immediately
	session.Save(r, w)
}
