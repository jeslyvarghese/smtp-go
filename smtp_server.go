package main

import (
	"bufio"
	"container/list"
	"io"
	"log"
	"net"
	"os"

	"smtp-go.quadratic.xyz/internal/utils"
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
	reader := bufio.NewReader(conn)
	for {
		// Simple SMTP command handling
		command, err := utils.ReadLine(conn, logger, reader)
		if err != nil {
			logger.Printf("Error reading command from %s: %v\n", conn.RemoteAddr(), err)
			return
		}
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		// Respond to HELO command
		if command[:4] == "HELO" {
			response := "250-us11-012mrc.dh.atmailcloud.com Hello 27-33-184-51.static.tpgi.com.au [27.33.184.51]" + "\r\n" + "250-SIZE 52428800" + "\r\n" + "250-8BITMIME" + "\r\n" + "250-PIPELINING" + "\r\n" + "250-AUTH LOGIN PLAIN" + "\r\n" + "250-CHUNKING" + "\r\n" + "250 HELO" + "\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Wait For MAIL FROM command
		command, err = utils.ReadLine(conn, logger, reader)
		if err != nil {
			logger.Printf("Error reading command from %s: %v\n", conn.RemoteAddr(), err)
			return
		}
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		if command[:9] == "MAIL FROM" {
			response := "250 OK\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Wait for RCPT TO command

		command, err = utils.ReadLine(conn, logger, reader)
		if err != nil {
			logger.Printf("Error reading command from %s: %v\n", conn.RemoteAddr(), err)
			return
		}
		logger.Printf("Received command from %s: %s", conn.RemoteAddr(), command)

		if command[:7] == "RCPT TO" {
			response := "250 OK\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		// Wait for DATA command

		command, err = utils.ReadLine(conn, logger, reader)
		if err != nil {
			logger.Printf("Error reading command from %s: %v\n", conn.RemoteAddr(), err)
			return
		}
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
			line, err := utils.ReadLine(conn, logger, reader)
			if err != nil {
				if err != io.EOF {
					logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
				}
				break
			}

			logger.Printf("Received data line from %s: %s", conn.RemoteAddr(), line)
			if line == "." {
				logger.Printf("End of data from %s", conn.RemoteAddr())
				break
			}
			data += line + "\r\n"
		}

		logger.Printf("Received email data from %s:\n%s", conn.RemoteAddr(), data)

		response := "250 OK: Queued\r\n"
		if _, err := conn.Write([]byte(response)); err != nil {
			logger.Printf("Error writing to %s: %v\n", conn.RemoteAddr(), err)
			return
		}

		logger.Printf("Email from %s processed successfully", conn.RemoteAddr())

		// Wait For QUIT command

		command, err = utils.ReadLine(conn, logger, reader)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}
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
