package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"go.uber.org/zap"
)

type TCPClient struct {
	conn     net.Conn
	maxMessageSize int
	logger *zap.Logger
}

type TCPClientConfig struct {
	Address        *string
	IdleTimeout    *time.Duration
	MaxMessageSize *int
}

func NewTCPClient(cfg TCPClientConfig, logger *zap.Logger) (*TCPClient, error) {
	conn, err := net.Dial("tcp", *cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("Connection create error: %w", err)
	}

	client := &TCPClient{
		conn:     conn,
		maxMessageSize: *cfg.MaxMessageSize,
		logger: logger,
	}

	if err := conn.SetDeadline(time.Now().Add(*cfg.IdleTimeout)); err != nil {
		return nil, fmt.Errorf("Connection set deadline error: %w", err)
	}

	return client, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if _, err := c.conn.Write(request); err != nil {
		c.logger.Error("TCPClient Send: connection.Write request error", zap.Error(err))

		return nil, err
	}

	response := make([]byte, c.maxMessageSize)
	count, err := c.conn.Read(response)
	if err != nil && err != io.EOF {
		c.logger.Error("TCPClient Send: connection.Read response error", zap.Error(err))

		return nil, err
	}
	if count >= c.maxMessageSize {
		c.logger.Error("TCPClient Send: count >= c.maxMessageSize error", zap.Int("count", count), zap.Int("maxMessageSize", c.maxMessageSize))

		return nil, errors.New("small buffer size")
	}

	return response[:count], nil
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
