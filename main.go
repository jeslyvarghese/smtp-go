package main

import (
	"flag"
	"fmt"
)

func main() {
	server := flag.Bool("server", false, "Run as a server, otherwise starts a client")
	flag.Parse()

	if *server {
		fmt.Println("Starting server")
		StartServer()
	} else {
		fmt.Println("Starting client")
		StartClient()
	}
}
