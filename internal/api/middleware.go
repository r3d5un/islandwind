package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"strings"
)

// BasicAuthConfig contains the username and password used in the basic authentication
// for the HTTP server.
type BasicAuthConfig struct {
	// Username is the admin username.
	Username string `json:"username"`
	// Password is the admin password.
	//
	// Field is safe for logging as the [BasicAuthConfig] contains a custom [BasicAuthConfig.LogValue] method.
	Password string `json:"password"`
}

func (c BasicAuthConfig) LogValue() slog.Value {
	return slog.GroupValue(slog.String("username", c.Username))
}

func BasicAuthMiddleware(next http.Handler, cfg BasicAuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(cfg.Username))
			expectedPasswordHash := sha256.Sum256([]byte(cfg.Password))
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authentication", `Basic realm="restricted", charset="UTF-8"`)
		UnauthorizedResponse(w, r)
	})
}

// CORSMiddleware returns middleware enabeling Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler, allow []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allow, ", "))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
