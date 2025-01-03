package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/controllers"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

// Try to find the access token in the request.
// Returns error if the user is not authorized.
// If `nil` is returned instead of an error, it is safe to display protected content.
func ParseToken(r *http.Request) (*database.User, *database.Session, error) {
	tokenCookie, err := r.Cookie(sessions.TokenCookieName)

	// Anonymous visitor.
	if err != nil {
		return nil, nil, err
	}

	user, session, err := database.GetUserByToken(tokenCookie.Value)
	if err != nil {
		slog.Warn("session token is invalid", "err", err)
		return nil, nil, err
	}

	return &user, &session, nil
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if SSO is enabled.
		remoteUser := r.Header.Get("Remote-User")
		if config.Cfg.HeaderAuth && remoteUser != "" {
			_, err := database.CreateUser(remoteUser, nil)
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			count, err := database.CountUsers()
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			// If we have only one user, assign him all groups without owner.
			// We originally did not store remote users in the database,
			// so migrations cannot automatically assign groups to a user.
			if count == 1 {
				err := database.TransferGroups(nil, &remoteUser)
				if err != nil {
					controllers.ShowError(w, http.StatusInternalServerError, err)
					return
				}
			}

			r = sessions.AddUsernameToContext(r, remoteUser)
			next.ServeHTTP(w, r)
			return
		}

		user, sessionObj, err := ParseToken(r)

		// Most likely the user is not authorized.
		if err != nil {
			count, err := database.CountUsers()
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			if count < 1 {
				http.Redirect(w, r, "/registration", http.StatusFound)
			} else {
				http.Redirect(w, r, "/signin", http.StatusFound)
			}

			return
		}

		// Create a new token if the old one is about to expire.
		if sessionObj.CreatedAt.Add(database.TokenLifetime / 2).Before(time.Now()) {
			err := database.DeleteSession(sessionObj.Token)
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			newSession, err := database.CreateSession(user.Username)
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			http.SetCookie(w, sessions.SessionToCookie(newSession))
		}

		r = sessions.AddUsernameToContext(r, user.Username)
		next.ServeHTTP(w, r)
	})
}
