package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const sessionCookie = "workbench_session"
const sessionTTL = 30 * 24 * time.Hour

// hashPassword and checkPassword wrap bcrypt so the admin password is
// never stored or compared in plaintext.
func hashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func createSession(db *sql.DB) (string, error) {
	token, err := newSessionToken()
	if err != nil {
		return "", err
	}
	expires := time.Now().Add(sessionTTL).UTC().Format(timeLayout)
	_, err = db.Exec("INSERT INTO sessions (token, expires_at) VALUES (?, ?)", token, expires)
	return token, err
}

func sessionValid(db *sql.DB, token string) bool {
	if token == "" {
		return false
	}
	var expiresAt string
	err := db.QueryRow("SELECT expires_at FROM sessions WHERE token = ?", token).Scan(&expiresAt)
	if err != nil {
		return false
	}
	expires, err := time.Parse(timeLayout, expiresAt)
	if err != nil || time.Now().UTC().After(expires) {
		return false
	}
	return true
}

func deleteSession(db *sql.DB, token string) {
	db.Exec("DELETE FROM sessions WHERE token = ?", token)
}

// setSessionCookie marks the cookie Secure whenever the request arrived
// over TLS or behind a TLS-terminating proxy (as on Railway/Fly), but
// leaves it unset for plain-http local development.
func setSessionCookie(w http.ResponseWriter, r *http.Request, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isTLS(r),
		SameSite: http.SameSiteStrictMode,
	})
}

func isTLS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

// requireAuth protects admin routes: no valid session cookie means a
// redirect to the login page instead of the handler running at all.
func requireAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdmin(db, r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// isAdmin reports whether the request carries a valid session, without
// redirecting. Public pages use it to decide whether to show admin
// affordances (an "edit this post" link, a "+ new post" button, drafts
// mixed into the listing) inline, the way a logged-in author expects on
// an ordinary blog.
func isAdmin(db *sql.DB, r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return false
	}
	return sessionValid(db, cookie.Value)
}
