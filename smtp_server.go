package main

import (
	"container/list"
	"io"
	"log"
	"net"
	"os"
)

func StartSMTPServer() {
	// Implementation of SMTP server would go here
	logger := log.New(os.Stdout, "smtp-server: ", log.LstdFlags)
	logger.Println("SMTP Server started")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	logger.Println("Listening on :8080 for SMTP connections...")

	var connList list.List

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println("Error accepting connection:", err)
			continue
		}

		logger.Printf("New SMTP connection from %s\n", conn.RemoteAddr())
		connList.PushBack(conn)

		go handleSMTPConnection(conn, logger)
	}
}

func handleSMTPConnection(conn net.Conn, logger *log.Logger) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}

		// Simple SMTP command handling
		command := string(buf[:n])
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		// Respond to HELO command
		if command[:4] == "HELO" {
			response := "250 Hello\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Wait For MAIL FROM command

		n, err = conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}

		command = string(buf[:n])
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		if command[:9] == "MAIL FROM" {
			response := "250 OK\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Wait for RCPT TO command

		n, err = conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}

		command = string(buf[:n])
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		if command[:7] == "RCPT TO" {
			response := "250 OK\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Wait for DATA command

		n, err = conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}

		command = string(buf[:n])
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		if command[:4] == "DATA" {
			response := "354 End data with <CR><LF>.<CR><LF>\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Read email data until a single dot on a line

		var data string
		for {
			n, err = conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
				}
				return
			}

			line := string(buf[:n])
			if line == ".\r\n" {
				break
			}
			data += line
		}

		logger.Printf("Received email data from %s:\n%s", conn.RemoteAddr(), data)

		response := "250 OK: Queued\r\n"
		if _, err := conn.Write([]byte(response)); err != nil {
			logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
			return
		}

		logger.Printf("Email from %s processed successfully", conn.RemoteAddr())

		// Wait For QUIT command

		n, err = conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}

		command = string(buf[:n])
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		if command[:4] == "QUIT" {
			response := "221 Bye\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
			break
		}

	}
}
