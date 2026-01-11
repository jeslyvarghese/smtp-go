package main

import (
	"bufio"
	"fmt"
	"net"
)

func StartLineEchoServer() {
	// Listen on TCP port 9090
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Line Echo Server listening on :9090")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleEchoConnection(conn)
	}
}

func handleEchoConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	fmt.Printf("Received: %s", message)

	_, err = writer.WriteString("ECHO: " + message)
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}
	writer.Flush()
}
