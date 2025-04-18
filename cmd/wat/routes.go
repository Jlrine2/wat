package main

import (
	"net/http"
	"net/url"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /status", app.StatusHandler)

	// Media
	router.Handle("GET /media", app.RequireAuthMiddleware(
		http.HandlerFunc(app.MediaListHandler),
	))
	router.Handle("POST /media", app.RequireAdminMiddleware(
		http.HandlerFunc(app.MediaUploadHandler)))
	router.Handle("DELETE /media", app.RequireAdminMiddleware(
		http.HandlerFunc(app.MediaDeleteHandler)))
	router.Handle("/media/", app.RequireAuthMiddleware(
		http.StripPrefix("/media/",
			http.HandlerFunc(app.MediaHandler),
		)))

	// Watch Parties
	router.Handle("POST /watch-parties", app.RequireAdminMiddleware(
		http.HandlerFunc(app.CreateWatchPartyHandler)))
	router.Handle("GET /watch-parties", app.RequireAuthMiddleware(
		http.HandlerFunc(app.ListWatchPartyHandler)))
	router.Handle("DELETE /watch-parties", app.RequireAdminMiddleware(
		http.HandlerFunc(app.DeleteWatchPartyHandler)))

	// Client
	router.Handle("/", http.HandlerFunc(app.ClientHandler))

	// Auth
	router.HandleFunc("GET /auth/discord/login", app.DiscordLoginHandler)
	redirectUri, err := url.Parse(app.config.DiscordOauth.RedirectUri)
	if err != nil {
		panic("Invalid discord callback url")
	}
	router.HandleFunc("GET "+redirectUri.Path, app.DiscordCallbackHandler)
	router.Handle("GET /auth/me", app.RequireAuthMiddleware(
		http.HandlerFunc(app.GetAuthDetailsHandler)))

	// Websocket
	router.Handle("/ws", app.RequireAuthMiddleware(http.HandlerFunc(app.SyncHandler)))

	handler := app.LoggerMiddleware(router)
	return handler
}
