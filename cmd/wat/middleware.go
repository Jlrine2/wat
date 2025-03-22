package main

import (
	"net/http"
	"wat/internal/auth"
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
		accessToken, err := auth.GetSession(cookie.Value, app.db)
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
		_, err = auth.GetDiscordAuthDetails(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
