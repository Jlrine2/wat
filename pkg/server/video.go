package server

import (
	"fmt"
	"html/template"
	"net/http"
)

type VideoTemplateParameters struct {
	VideoName string
	Scheme    string
	Host      string
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http" // Default to http if header is not present
	}
	data := &VideoTemplateParameters{
		VideoName: r.URL.Query().Get("videoName"),
		Scheme:    scheme,
		Host:      r.Host,
	}
	tmpl := template.Must(template.ParseFiles("web/video-player.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Error rendering template for Video Page %s", err)
	}
}
