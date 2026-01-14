package main

import (
	"flag"
	"fmt"
)

func main() {
	server := flag.String("type", "tcp-server", "Run as a server, otherwise starts a client")
	flag.Parse()

	switch *server {
	case "tcp-server":
		fmt.Println("Starting TCP server")
		StartServer()
	case "tcp-client":
		fmt.Println("Starting TCP client")
		StartClient()
	case "line-echo":
		fmt.Println("Starting line echo server")
		StartLineEchoServer()
	case "smtp-server":
		fmt.Println("Starting SMTP server")
		StartSMTPServer()
	case "smtp-client":
		fmt.Println("Starting SMTP client")
		StartSMTPClient()
	default:
		fmt.Println("Invalid server type")
	}
}
