package main

import (
	"context"
	"net/http"
	"wat/internal/auth"
	"wat/internal/models"
)

func (app *application) RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.config.Server.AuthEnabled {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("watAuth")
		if err != nil {
			app.logger.Info("Unable to get watAuth cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		session, err := auth.GetSession(cookie.Value, app.db)
		if err != nil {
			app.logger.Info("Unable to get session details from database", "error", err.Error())
			err = writeJSON(
				w,
				&map[string]bool{"authenticated": false},
				http.StatusUnauthorized,
				nil,
			)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

func (app *application) RequireAdminMiddleware(next http.Handler) http.Handler {
	return app.RequireAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.config.Server.AuthEnabled {
			next.ServeHTTP(w, r)
			return
		}
		session := r.Context().Value("session")
		if session == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s := session.(*models.Session)
		for _, adminUserId := range app.config.DiscordOauth.AdminUserIds {
			if s.User.User.ID == adminUserId {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusForbidden)
	}))
}

func (app *application) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Request received", "method", r.Method, "url", r.URL.String(), "client", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
