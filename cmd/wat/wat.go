package main

import server2 "wat/pkg/server"

func main() {
	server := server2.NewWatServer(8080, "./media")
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
