package handlers

import (
	"net/http"
	"server/prisma/db"
	"server/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Failed to authenticate user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dbUser, err := h.db.User.FindFirst(
		db.User.GoogleID.Equals(user.UserID),
	).Exec(r.Context())
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

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     expirationTime.Unix(),
	})

	tokenString, err := claims.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	err = utils.StoreUserSession(w, r, dbUser.ID, tokenString)
	if err != nil {
		http.Error(w, "Failed to store session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	utils.RemoveUserSession(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
