package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	cmd "smtp-go.quadratic.xyz/internal/smtp/client"
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

	// Send HELO/EHLO command
	if err := cmd.Helo(conn, logger); err != nil {
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

	if err := cmd.MailFrom(conn, logger, "someone@localhost"); err != nil {
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
	if err := cmd.RcptTo(conn, logger, "recipient@localhost"); err != nil {
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
	if err := cmd.Data(conn, logger); err != nil {
		logger.Printf("Failed to send DATA command: %v\n", err)
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
	emailContent := "Subject: Test Email\r\n\r\nThis is a test email sent from the SMTP client."
	if err := cmd.SendEmailContent(conn, logger, emailContent); err != nil {
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
	if err := cmd.Quit(conn, logger); err != nil {
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
