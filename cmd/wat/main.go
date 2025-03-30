package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"wat/internal/config"
	"wat/internal/database"
)

type application struct {
	config  *config.Config
	logger  slog.Logger
	db      database.DatabaseController
	version string
	syncHub *SyncHub
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "/etc/wat/config.yaml", "Path to config file for Wat")
	flag.Parse()
	fmt.Println(configPath)
	config, err := config.LoadConfig(configPath)
	if err != nil {
		panic("Unable to parse config file")
	}
	app := &application{
		config:  config,
		logger:  *slog.New(slog.NewTextHandler(os.Stdout, nil)),
		db:      database.NewMemoryDatabase(),
		version: "v0.0.0",
		syncHub: newHub(),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 15 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	go app.syncHub.Run()
	app.logger.Info("Listening for traffic", "address", fmt.Sprintf("http://%s:%s", app.config.Server.Host, app.config.Server.Port))
	err = srv.ListenAndServe()
	app.logger.Error(err.Error())
	os.Exit(1)

}
