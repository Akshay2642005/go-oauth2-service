package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "session"
)

type SessionOptions struct {
	CookiesKey string
	MaxAge     int    // Session expiration time in seconds
	HttpOnly   bool   // Prevents JavaScript from accessing cookies
	Secure     bool   // Ensures cookies are sent only over HTTPS
	SameSite   string // "Strict", "Lax", or "None"
}

// NewCookieStore initializes and configures a Gorilla session store
func NewCookieStore(opts SessionOptions) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))

	// Configure session options
	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	// Corrected SameSite settings using http package
	switch opts.SameSite {
	case "Strict":
		store.Options.SameSite = http.SameSiteStrictMode
	case "Lax":
		store.Options.SameSite = http.SameSiteLaxMode
	case "None":
		store.Options.SameSite = http.SameSiteNoneMode
	default:
		store.Options.SameSite = http.SameSiteNoneMode // Default to Lax for security
	}

	return store
}
