package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type ClientConfig struct {
	Address       string
	Timeout       time.Duration
	RetryAttempts int
	RetryDelay    time.Duration
}

func StartClient() {
	config := ClientConfig{
		Address:       "localhost:8080",
		Timeout:       5 * time.Second,
		RetryAttempts: 3,
		RetryDelay:    2 * time.Second,
	}

	if err := sendMessage(config, "Hello server!"); err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		return
	}
}

func sendMessage(config ClientConfig, message string) error {
	var conn net.Conn
	var err error

	for attempt := 1; attempt <= config.RetryAttempts; attempt++ {
		conn, err = net.DialTimeout("tcp", config.Address, config.Timeout)
		if err == nil {
			break
		}

		if attempt < config.RetryAttempts {
			fmt.Printf("Connection attempt %d failed: %v. Retrying in %v...\n",
				attempt, err, config.RetryDelay)
			time.Sleep(config.RetryDelay)
			continue
		}

		return fmt.Errorf("failed to connect after %d attempts: %w", config.RetryAttempts, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", config.Address)

	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if _, err := io.WriteString(conn, message); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(config.Timeout)); err != nil {
		return fmt.Errorf("failed to set read deadline: %w", err)
	}

	response := make([]byte, len(message))
	n, err := io.ReadFull(conn, response)

	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Recieved from serer: %s\n", string(response[:n]))
	return nil
}
