package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func StartServer() {
	// Create a TCP listener on port 8080

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
	}
	defer listener.Close()

	fmt.Println("Listening on :8080..")

	// Setup a graceful shutdown

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	go func() {
		<-sigChan
		fmt.Println("\nShutting down gracefully")
		cancel()
		listener.Close()
	}()

	for {
		// Set a timeout for Accept to check shutdown signal
		listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))

		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				select {
				case <-ctx.Done():
					wg.Wait()
					return
				default:
					continue
				}
			}
			fmt.Println("Error accepting connection:", err)
			return
		}

		wg.Add(1)
		go handleConnection(ctx, conn, &wg)
	}
}

func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	// Set connection deadlines
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Echo recieved data back to client
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
				}
				return
			}

			if _, err := conn.Write(buf[:n]); err != nil {
				fmt.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}

			// Reset deadline after successful operation
			conn.SetDeadline(time.Now().Add(30 * time.Second))

		}

	}

}
