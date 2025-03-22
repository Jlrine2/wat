package main

import (
	"net/http"
	"net/url"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /status", app.StatusHandler)
	router.HandleFunc("GET /media", app.MediaListHandler)
	router.HandleFunc("POST /media", app.MediaUploadHandler)
	router.HandleFunc("DELETE /media", app.MediaDeleteHandler)
	router.Handle("/media/", app.RequireAuthMiddleware(
		http.StripPrefix("/media/",
			http.HandlerFunc(app.MediaHandler),
		)))
	router.Handle("/", http.HandlerFunc(app.ClientHandler))

	router.HandleFunc("GET /auth/discord/login", app.DiscordLoginHandler)
	redirectUri, err := url.Parse(app.config.DiscordOauth.RedirectUri)
	if err != nil {
		panic("Invalid discord callback url")
	}
	router.HandleFunc("GET "+redirectUri.Path, app.DiscordCallbackHandler)
	router.HandleFunc("GET /auth/me", app.GetAuthDetailsHandler)

	router.Handle("/ws", app.RequireAuthMiddleware(http.HandlerFunc(app.SyncHandler)))

	return router
}
