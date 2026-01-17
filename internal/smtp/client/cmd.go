package cmd

import (
	"log"
	"net"

	"smtp-go.quadratic.xyz/internal/utils"
)

func Ehlo(conn net.Conn, logger *log.Logger) error {
	if err := utils.WriteLine(conn, logger, "EHLO localhost\r\n"); err != nil {
		logger.Printf("Failed to send EHLO command: %v\n", err)
		return err
	}
	return nil
}

func Helo(conn net.Conn, logger *log.Logger) error {
	if err := utils.WriteLine(conn, logger, "HELO localhost\r\n"); err != nil {
		logger.Printf("Failed to send HELO command: %v\n", err)
		return err
	}
	return nil
}

func MailFrom(conn net.Conn, logger *log.Logger, from string) error {
	command := "MAIL FROM:<" + from + ">\r\n"
	if err := utils.WriteLine(conn, logger, command); err != nil {
		logger.Printf("Failed to send MAIL FROM command: %v\n", err)
		return err
	}
	return nil
}

func RcptTo(conn net.Conn, logger *log.Logger, to string) error {
	command := "RCPT TO:<" + to + ">\r\n"
	if err := utils.WriteLine(conn, logger, command); err != nil {
		logger.Printf("Failed to send RCPT TO command: %v\n", err)
		return err
	}
	return nil
}

func Data(conn net.Conn, logger *log.Logger) error {
	if err := utils.WriteLine(conn, logger, "DATA\r\n"); err != nil {
		logger.Printf("Failed to send DATA command: %v\n", err)
		return err
	}
	return nil
}

func SendEmailContent(conn net.Conn, logger *log.Logger, content string) error {
	if err := utils.WriteLine(conn, logger, content+"\r\n.\r\n"); err != nil {
		logger.Printf("Failed to send email content: %v\n", err)
		return err
	}
	return nil
}

func Quit(conn net.Conn, logger *log.Logger) error {
	if err := utils.WriteLine(conn, logger, "QUIT\r\n"); err != nil {
		logger.Printf("Failed to send QUIT command: %v\n", err)
		return err
	}
	return nil
}
