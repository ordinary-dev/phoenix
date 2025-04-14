package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database/entities"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

// Try to find the access token in the request.
// Returns error if the user is not authorized.
// If `nil` is returned instead of an error, it is safe to display protected content.
func (m *Middleware) ParseToken(r *http.Request) (*entities.User, *entities.Session, error) {
	tokenCookie, err := r.Cookie(sessions.TokenCookieName)

	// Anonymous visitor.
	if err != nil {
		return nil, nil, err
	}

	user, session, err := m.db.GetUserByToken(tokenCookie.Value)
	if err != nil {
		slog.Warn("session token is invalid", "err", err)
		return nil, nil, err
	}

	return &user, &session, nil
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if SSO is enabled.
		remoteUser := r.Header.Get("Remote-User")
		if config.Cfg.HeaderAuth && remoteUser != "" {
			_, err := m.db.CreateUser(remoteUser, nil)
			if err != nil {
				m.ctrl.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			count, err := m.db.CountUsers()
			if err != nil {
				m.ctrl.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			// If we have only one user, assign him all groups without owner.
			// We originally did not store remote users in the database,
			// so migrations cannot automatically assign groups to a user.
			if count == 1 {
				err := m.db.TransferGroups(nil, &remoteUser)
				if err != nil {
					m.ctrl.ShowError(w, http.StatusInternalServerError, err)
					return
				}
			}

			r = sessions.AddUsernameToContext(r, remoteUser)
			next.ServeHTTP(w, r)
			return
		}

		user, sessionObj, err := m.ParseToken(r)

		// Most likely the user is not authorized.
		if err != nil {
			count, err := m.db.CountUsers()
			if err != nil {
				m.ctrl.ShowError(w, http.StatusInternalServerError, err)
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
		if sessionObj.CreatedAt.Add(entities.TokenLifetime / 2).Before(time.Now()) {
			err := m.db.DeleteSession(sessionObj.Token)
			if err != nil {
				m.ctrl.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			newSession, err := m.db.CreateSession(user.Username)
			if err != nil {
				m.ctrl.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			http.SetCookie(w, sessions.SessionToCookie(newSession))
		}

		r = sessions.AddUsernameToContext(r, user.Username)
		next.ServeHTTP(w, r)
	})
}
