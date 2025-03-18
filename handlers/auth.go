package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"server/prisma/db"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (h *Handler) HandleSessionUser(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "session_name") // ✅ Ensure correct session name
	if err != nil {
		http.Error(w, "Session not found (error getting session)", http.StatusUnauthorized)
		fmt.Println("Error fetching session:", err)
		return
	}

	fmt.Println("Session Values:", session.Values) // ✅ Debugging session values

	user, exists := session.Values["user"]
	if !exists {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		fmt.Println("Session exists, but no user found!")
		return
	}

	fmt.Println("User found in session:", user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Failed to authenticate user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Try to find user in DB
	dbUser, err := h.db.User.FindFirst(
		db.User.GoogleID.Equals(user.UserID),
	).Exec(r.Context())
	// If user does not exist, create a new one
	if err != nil {
		var avatarURL *string
		if user.AvatarURL != "" {
			avatarURL = &user.AvatarURL
		}

		dbUser, err = h.db.User.CreateOne(
			db.User.Email.Set(user.Email),
			db.User.Name.Set(user.Name),
			db.User.GoogleID.Set(user.UserID),
			db.User.AvatarURL.SetOptional(avatarURL),
		).Exec(r.Context())
		if err != nil {
			http.Error(w, "Could not save user", http.StatusInternalServerError)
			return
		}
	}

	// Extract avatar URL properly
	var avatarURL string
	if avatar, ok := dbUser.AvatarURL(); ok { // Call the function to get the value
		avatarURL = avatar
	}

	// Convert `dbUser` to `goth.User` for session storage
	gothUser := goth.User{
		UserID:    dbUser.GoogleID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		AvatarURL: avatarURL, // Use extracted value
	}

	// Store the session
	err = h.auth.StoreUserSession(w, r, gothUser)
	if err != nil {
		http.Error(w, "Failed to store session", http.StatusInternalServerError)
		return
	}

	// ✅ Send a response first to confirm session is saved
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Session stored successfully",
		"redirect": "http://localhost:5173/dashboard?name=" + url.QueryEscape(user.Name) +
			"&email=" + url.QueryEscape(user.Email) +
			"&avatar=" + url.QueryEscape(avatarURL),
	})
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	h.auth.RemoveUserSession(w, r)

	http.Redirect(w, r, "http://localhost:5173/", http.StatusSeeOther)
}
