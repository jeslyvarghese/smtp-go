package utils

import (
	"log"
	"net"
	"time"
)

func WriteLine(conn net.Conn, logger *log.Logger, message string) error {
	conn.SetWriteDeadline(time.Now().Add(5 * time.Minute))
	defer conn.SetWriteDeadline(time.Time{})

	n, err := conn.Write([]byte(message))
	if err != nil {
		logger.Printf("Write error: %v", err)
		return err
	}
	logger.Printf("Wrote %d bytes: %s", n, message)
	return nil
}
