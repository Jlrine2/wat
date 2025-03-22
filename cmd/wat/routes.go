package main

import (
	"net/http"
	"net/url"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/status", app.StatusHandler)

	router.Handle("/media/", app.RequireAuthMiddleware(
		http.StripPrefix("/media/",
			http.HandlerFunc(app.MediaHandler),
		)))
	router.Handle("/", http.HandlerFunc(app.ClientHandler))

	router.HandleFunc("/auth/discord/login", app.DiscordLoginHandler)
	redirectUri, err := url.Parse(app.config.DiscordOauth.RedirectUri)
	if err != nil {
		panic("Invalid discord callback url")
	}
	router.HandleFunc(redirectUri.Path, app.DiscordCallbackHandler)
	router.HandleFunc("/auth/me", app.GetAuthDetailsHandler)

	router.Handle("/ws", app.RequireAuthMiddleware(http.HandlerFunc(app.SyncHandler)))

	return router
}
