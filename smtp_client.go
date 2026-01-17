package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"smtp-go.quadratic.xyz/internal/utils"
)

func StartSMTPClient() {
	logger := log.New(os.Stdout, "smtp-server: ", log.LstdFlags)

	logger.Println("SMTP Client started")

	config := ClientConfig{
		Address:       "localhost:8080",
		Timeout:       5 * time.Second,
		RetryAttempts: 3,
		RetryDelay:    2 * time.Second,
	}

	conn, err := net.DialTimeout("tcp", config.Address, config.Timeout)
	if err != nil {
		logger.Printf("Failed to connect after %d attempts: %v\n", config.RetryAttempts, err)
		return
	}
	defer conn.Close()

	logger.Printf("Connected to %s\n", config.Address)

	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set write deadline: %v\n", err)
		return
	}

	// Send HELO/EHLO command
	message := "HELO localhost\r\n"
	if _, err := conn.Write([]byte(message)); err != nil {
		logger.Printf("Failed to send data: %v\n", err)
		return
	}

	reader := bufio.NewReader(conn)

	response, err := utils.ReadLine(conn, logger, reader)
	if err != nil {
		logger.Printf("Failed to read data: %v\n", err)
		return
	}

	logger.Printf("Received response: %s\n", response)
	if response[:3] != "250" {
		logger.Printf("Unexpected response from server: %s\n", response)
		return
	}

	if response[3] == '-' {
		// Read multiline response
		for {
			line, err := utils.ReadLine(conn, logger, reader)
			if err != nil {
				logger.Printf("Failed to read data: %v\n", err)
				return
			}
			logger.Printf("Received response: %s\n", line)
			if line[3] != '-' {
				break
			}
		}
	}

	// Send MAIL FROM command

	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set write deadline: %v\n", err)
		return
	}

	mailFrom := "MAIL FROM:<someone@localhost>\r\n"
	if _, err := conn.Write([]byte(mailFrom)); err != nil {
		logger.Printf("Failed to send MAIL FROM command: %v\n", err)
		return
	}

	// Read response

	response, err = utils.ReadLine(conn, logger, reader)
	if err != nil {
		logger.Printf("Failed to read data: %v\n", err)
		return
	}

	logger.Printf("Received response: %s\n", response)
	if response[:3] != "250" {
		logger.Printf("Unexpected response from server: %s\n", response)
		return
	}

	// SEND RCPT TO command
	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set write deadline: %v\n", err)
		return
	}

	rcptTo := "RCPT TO:<recipient@localhost>\r\n"
	if _, err := conn.Write([]byte(rcptTo)); err != nil {
		logger.Printf("Failed to send RCPT TO command: %v\n", err)
		return
	}

	// Read response
	response, err = utils.ReadLine(conn, logger, reader)
	if err != nil {
		logger.Printf("Failed to read data: %v\n", err)
		return
	}
	logger.Printf("Received response: %s\n", response)
	if response[:3] != "250" {
		logger.Printf("Unexpected response from server: %s\n", response)
		return
	}

	// Send DATA command
	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set write deadline: %v\n", err)
		return
	}

	dataCmd := "DATA\r\n"
	if _, err := conn.Write([]byte(dataCmd)); err != nil {
		logger.Printf("Failed to send DATA command: %v\n", err)
		return
	}

	if err := conn.SetReadDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set read deadline: %v\n", err)
		return
	}

	// Read response
	response, err = utils.ReadLine(conn, logger, reader)
	if err != nil {
		logger.Printf("Failed to read data: %v\n", err)
		return
	}

	logger.Printf("Received response: %s\n", response)
	if response[:3] != "354" {
		logger.Printf("Unexpected response from server: %s\n", response)
		return
	}

	// Send email content
	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set write deadline: %v\n", err)
		return
	}

	emailContent := "Subject: Test Email\r\n\r\nThis is a test email.\r\n.\r\n"
	if _, err := conn.Write([]byte(emailContent)); err != nil {
		logger.Printf("Failed to send email content: %v\n", err)
		return
	}

	// Read response
	response, err = utils.ReadLine(conn, logger, reader)
	if err != nil {
		logger.Printf("Failed to read data: %v\n", err)
		return
	}

	logger.Printf("Received response: %s\n", response)
	if response[:3] != "250" {
		logger.Printf("Unexpected response from server: %s\n", response)
		return
	}

	// Send QUIT command
	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		logger.Printf("Failed to set write deadline: %v\n", err)
		return
	}

	quitCmd := "QUIT\r\n"
	if _, err := conn.Write([]byte(quitCmd)); err != nil {
		logger.Printf("Failed to send QUIT command: %v\n", err)
		return
	}

	// Read response
	response, err = utils.ReadLine(conn, logger, reader)
	if err != nil {
		logger.Printf("Failed to read data: %v\n", err)
		return
	}

	logger.Printf("Received response: %s\n", response)
	if response[:3] != "221" {
		logger.Printf("Unexpected response from server: %s\n", response)
		return
	}

	logger.Println("SMTP Client finished")
}
