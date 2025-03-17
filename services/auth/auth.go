package auth

import (
	"fmt"
	"log"
	"net/http"
	"server/config"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthService struct{}

// NewAuthService initializes Gothic with session store and Google OAuth provider
func NewAuthService(store sessions.Store) *AuthService {
	gothic.Store = store

	goth.UseProviders(
		google.New(
			config.Envs.GoogleClientID,
			config.Envs.GoogleClientSecret, // Fixed incorrect duplicate usage of GoogleClientID
			buildCallbackURL("google"),
		),
	)

	return &AuthService{}
}

// GetSessionUser retrieves the authenticated user from session
func (s *AuthService) GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return goth.User{}, err
	}

	u, exists := session.Values["user"]
	if !exists || u == nil {
		return goth.User{}, fmt.Errorf("user is not authenticated")
	}

	// Ensure correct type assertion to goth.User
	user, ok := u.(goth.User)
	if !ok {
		return goth.User{}, fmt.Errorf("session contains invalid user data")
	}

	return user, nil
}

// StoreUserSession saves the authenticated user in session
func (s *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		log.Println("Failed to get session:", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return err
	}

	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		log.Println("Failed to save session:", err)
		http.Error(w, "Failed to store session", http.StatusInternalServerError)
		return err
	}

	return nil
}

// RemoveUserSession deletes the user session (logout)
func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		log.Println("Error retrieving session:", err)
		http.Error(w, "Session retrieval failed", http.StatusInternalServerError)
		return
	}

	// Clear user session data
	session.Values["user"] = nil
	session.Options.MaxAge = -1 // Expire the session immediately
	err = session.Save(r, w)
	if err != nil {
		log.Println("Failed to clear session:", err)
		http.Error(w, "Failed to remove session", http.StatusInternalServerError)
	}
}

// RequireAuth is a middleware that ensures user authentication
func RequireAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("User is authenticated: %s!", user.FirstName)

		handlerFunc(w, r)
	}
}

// buildCallbackURL constructs the OAuth callback URL dynamically
func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s:%s/auth/%s/callback", config.Envs.PublicHost, config.Envs.Port, provider)
}
