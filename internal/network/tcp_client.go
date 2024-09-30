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
	conn               net.Conn
	maxMessageSize     int
	idleTimeout        *time.Duration
	connectionDeadline *time.Time
	logger             *zap.Logger
}

type TCPClientConfig struct {
	Address            *string
	IdleTimeout        *time.Duration
	ConnectionDeadline *time.Time
	MaxMessageSize     *int
}

func NewTCPClient(cfg TCPClientConfig, conn net.Conn, logger *zap.Logger) (*TCPClient, error) {
	client := &TCPClient{
		conn:               conn,
		maxMessageSize:     *cfg.MaxMessageSize,
		idleTimeout:        cfg.IdleTimeout,
		connectionDeadline: cfg.ConnectionDeadline,
		logger:             logger,
	}

	return client, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	err := c.setConnectionDeadline()
	if err != nil {
		c.logger.Error("TCPClient Send: setIdleTimeout error", zap.Error(err))

		return nil, errors.New("Client internal error")
	}

	if _, err = c.conn.Write(request); err != nil {
		c.logger.Error("TCPClient Send: connection.Write request error", zap.Error(err))

		return nil, errors.New("Client send data error")
	}

	response := make([]byte, c.maxMessageSize)
	count, err := c.conn.Read(response)
	if err != nil && err != io.EOF {
		c.logger.Error("TCPClient Send: connection.Read response error", zap.Error(err))

		return nil, errors.New("Client read data error")
	}
	if count >= c.maxMessageSize {
		c.logger.Error("TCPClient Send: count >= c.maxMessageSize error", zap.Int("count", count), zap.Int("maxMessageSize", c.maxMessageSize))

		return nil, errors.New("Small buffer size")
	}

	return response[:count], nil
}

func (c *TCPClient) setConnectionDeadline() error {
	var deadline time.Time
	if c.connectionDeadline != nil {
		deadline = *c.connectionDeadline
	} else if c.idleTimeout != nil {
		deadline = time.Now().Add(*c.idleTimeout)
	}
	if err := c.conn.SetDeadline(deadline); err != nil {
		return fmt.Errorf("Connection set deadline error: %w", err)
	}

	return nil
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
