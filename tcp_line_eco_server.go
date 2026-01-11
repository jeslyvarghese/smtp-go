package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down Line Echo Server gracefully")
		listener.Close()
	}()

	var wg sync.WaitGroup

outer:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		select {
		case <-ctx.Done():
			break outer
		default:
			wg.Add(1)
			go handleEchoConnection(conn, &wg)

		}
	}

	wg.Wait()
	fmt.Println("Line Echo Server shut down")
}

func handleEchoConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
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
