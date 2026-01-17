package utils

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

func ReadLine(conn net.Conn, logger *log.Logger, reader *bufio.Reader) (string, error) {
	var buf []byte

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		defer conn.SetReadDeadline(time.Time{})

		chunk, err := reader.ReadSlice('\n')
		buf = append(buf, chunk...)
		logger.Printf("Read %d bytes: %s", len(chunk), string(chunk))
		if len(buf) > 512 {
			return "", io.ErrUnexpectedEOF
		}

		if err == nil {
			if len(buf) < 2 || buf[len(buf)-2] != '\r' {
				return "", io.ErrUnexpectedEOF
			}
			return string(buf[:len(buf)-2]), nil
		}

		if err == io.EOF {
			if len(buf) == 0 {
				return "", io.EOF
			}
			return string(buf), io.EOF
		}

		if err != bufio.ErrBufferFull {
			logger.Printf("Read error: %v", err)
			continue
		}

		return "", err
	}
}
