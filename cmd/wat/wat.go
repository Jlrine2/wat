package main

import "wat/pkg/server"

func main() {
	s := server.NewWatServer(8080, "./media")
	err := s.Run()
	if err != nil {
		panic(err)
	}
}
