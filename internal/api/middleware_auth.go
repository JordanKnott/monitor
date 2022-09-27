package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/db"
	"github.com/jordanknott/monitor/internal/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// AuthenticationMiddleware is a middleware that requires a valid JWT token to be passed via the Authorization header
type AuthenticationMiddleware struct {
	Data db.Data
}

func GetAuthorization(r *http.Request) (string, bool) {
	token := r.Header.Get("Authorization")
	log.WithField("authorization", token).Info("getting authorization token")
	if token != "" {
		return token, true
	}
	return "", false
}
func resolveToken(r *http.Request) (string, error) {
	c, err := r.Cookie("authToken")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

// Middleware returns the middleware handler
func (m *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Auth with app token
		if tokenRaw, ok := GetAuthorization(r); ok {
			token, err := uuid.Parse(tokenRaw)
			if err != nil {
				log.WithError(err).Error("error while parsing authorization token")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			t, err := m.Data.GetAppToken(ctx, token)
			if err != nil {
				log.WithError(err).Error("while getting app token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, utils.UserIDKey, t.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Auth with cookie
		if token, err := resolveToken(r); err == nil {
			userID, err := m.Data.GetUserIDForToken(r.Context(), token)
			if err != nil {
				log.WithError(err).Error("error while fetching authToken")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			ctx = context.WithValue(ctx, utils.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusBadRequest)
				logrus.Warn("no cookie")
				return
			}
			log.WithError(err).Error("error while fetching authToken cookie")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (api *MonitorApi) GetUserID(r *http.Request) (string, bool) {
	if val, ok := r.Context().Value(utils.UserIDKey).(string); ok {
		return val, true
	}
	return "", false
}
