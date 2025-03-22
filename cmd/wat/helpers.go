package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any, status int, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func getHostandProto(r *http.Request) string {
	// Try X-Forwarded headers first
	host := r.Header.Get("X-Forwarded-Host")
	proto := r.Header.Get("X-Forwarded-Proto")

	// Fall back to request host if X-Forwarded-Host not present
	if host == "" {
		host = r.Host
	}

	// Fall back to request scheme if X-Forwarded-Proto not present
	if proto == "" {
		if r.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	return proto + "://" + host
}
