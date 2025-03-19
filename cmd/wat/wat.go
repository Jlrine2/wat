package main

import (
	"flag"
	"wat/pkg/server"
)

func main() {
	var port uint
	var mediaPath string

	flag.UintVar(&port, "p", 8080, "server port")

	const mediaUsage = "Path to media files to serve, either a flat or nested directory of .mp4 files"
	flag.StringVar(&mediaPath, "m", "./media", mediaUsage)

	flag.Parse()
	s := server.NewWatServer(uint16(port), mediaPath)
	err := s.Run()
	if err != nil {
		panic(err)
	}
}
