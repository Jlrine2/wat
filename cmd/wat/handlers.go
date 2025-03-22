package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"wat/internal/auth"
	"wat/internal/models"
	"wat/internal/watchParties"
)

// status Handler

func (app *application) StatusHandler(w http.ResponseWriter, r *http.Request) {
	err := writeJSON(
		w,
		map[string]string{
			"status":  "up",
			"version": app.version,
		},
		http.StatusOK,
		nil,
	)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

// auth Handlers
func (app *application) DiscordLoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, auth.GetDiscordOauthUrl(app.config.DiscordOauth), http.StatusTemporaryRedirect)
}

func (app *application) DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken, err := auth.GetDiscordAccessToken(code, app.config.DiscordOauth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.logger.Info("Unable to get access token", "error", err.Error())
		return
	}
	sessionId, err := auth.CreateSession(accessToken, app.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.logger.Error("error creating session token", "error", err.Error())
	}
	isGuildMember, err := auth.GetDiscordGuildMembership(accessToken, app.config.DiscordOauth.GuildId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.logger.Error("error getting guild membership", "error", err.Error())
		return
	}
	if !isGuildMember {
		w.WriteHeader(http.StatusForbidden)
		app.logger.Info("User is not in expected guild")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "watAuth",
		Value:  sessionId,
		MaxAge: 60 * 60 * 12,
		Path:   "/",
	})
	http.Redirect(w, r, getHostandProto(r), http.StatusPermanentRedirect)
}

func (app *application) GetAuthDetailsHandler(w http.ResponseWriter, r *http.Request) {
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
	authDetails, err := auth.GetDiscordAuthDetails(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	authenticated := authDetails != nil

	if authenticated {
		err = writeJSON(w, &authDetails, http.StatusOK, nil)
	} else {
		err = writeJSON(
			w,
			&map[string]bool{"authenticated": false},
			http.StatusUnauthorized,
			nil,
		)
	}
	if err != nil {
		http.Error(w, "We are unable to process your request right now", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Media Handler
func (app *application) MediaHandler(w http.ResponseWriter, r *http.Request) {
	fileHandler := http.FileServer(http.Dir(app.config.Server.MediaLocation))
	fileHandler.ServeHTTP(w, r)
}

func (app *application) MediaListHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(app.config.Server.MediaLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileNames := []string{}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	err = writeJSON(w, &fileNames, http.StatusOK, nil)
	if err != nil {
		http.Error(w, "We are unable to process your request right now", http.StatusInternalServerError)
	}
}

func (app *application) MediaUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 10GB max file size
	r.ParseMultipartForm(10 << 30)

	file, handler, err := r.FormFile("file")
	if err != nil {
		app.logger.Error("Error retrieving file from form", "error", err.Error())
		http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create destination file
	dst, err := os.Create(filepath.Join(app.config.Server.MediaLocation, handler.Filename))
	if err != nil {
		app.logger.Error("Error creating destination file", "error", err.Error())
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		app.logger.Error("Error saving file", "error", err.Error())
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *application) MediaDeleteHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	err := os.Remove(filepath.Join(app.config.Server.MediaLocation, filename))
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Client Handler
func (app *application) ClientHandler(w http.ResponseWriter, r *http.Request) {
	if app.config.Server.ClientLocation[:4] == "http" {
		resp, err := http.Get(app.config.Server.ClientLocation + r.RequestURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return
		}
	} else {
		http.FileServer(http.Dir(app.config.Server.ClientLocation)).ServeHTTP(w, r)
	}
}

// Watch Party Handlers
func (app *application) CreateWatchPartyHandler(w http.ResponseWriter, r *http.Request) {

	watchParty := &models.WatchParty{}
	err := readJSON(r, watchParty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	watchPartyId, err := watchParties.CreateWatchParty(app.db, watchParty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, map[string]*models.WatchParty{watchPartyId: watchParty}, http.StatusCreated, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) ListWatchPartyHandler(w http.ResponseWriter, r *http.Request) {
	watchParties, err := watchParties.GetAllWatchParties(app.db)
	err = writeJSON(w, watchParties, http.StatusOK, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) DeleteWatchPartyHandler(w http.ResponseWriter, r *http.Request) {
	watchPartyId := r.URL.Query().Get("watchPartyId")
	err := watchParties.DeleteWatchParty(app.db, watchPartyId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
