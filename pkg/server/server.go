package server

import (
	"fmt"
	"net/http"
)

type WatServer struct {
	Port uint16
	mux  *http.ServeMux
}

func NewWatServer(port uint16, mediaPath string) *WatServer {
	fileHandler := http.FileServer(http.Dir(mediaPath))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			fmt.Printf("Error writing response: %v", err)
		}
	})
	mux.Handle("/media/", http.StripPrefix("/media/", fileHandler))
	return &WatServer{
		Port: port,
		mux:  mux,
	}
}

func (s *WatServer) Run() error {
	fmt.Printf("Listening on port %d\n", s.Port)
	addr := fmt.Sprintf(":%d", s.Port)
	err := http.ListenAndServe(addr, s.mux)
	if err != nil {
		return err
	}
	return nil
}
