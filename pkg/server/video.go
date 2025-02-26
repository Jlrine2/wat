package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed video.html
var VideoTemplate string

type VideoTemplateParameters struct {
	VideoName string
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http" // Default to http if header is not present
	}
	data := &VideoTemplateParameters{
		VideoName: r.URL.Query().Get("videoName"),
	}
	tmpl := template.Must(template.New("video").Parse(VideoTemplate))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Error rendering template for Video Page %s", err)
	}
}
