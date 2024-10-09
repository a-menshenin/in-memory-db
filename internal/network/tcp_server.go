package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
	"umemory/internal"

	"go.uber.org/zap"
)

type TCPServer struct {
	listener  net.Listener

	idleTimeout    time.Duration
	bufferSize     int
	activeConnections chan struct{}

	logger *zap.Logger
}

func NewTCPServer(config internal.Config, logger *zap.Logger) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	listener, err := net.Listen("tcp", config.Network.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &TCPServer{
		listener: listener,
		logger:   logger,
	}

	if config.Network.MaxMessageSize > 0 {
		server.bufferSize = config.Network.MaxMessageSize
	}

	server.activeConnections = make(chan struct{})
	if config.Network.MaxConnections > 0 {
		server.activeConnections = make(chan struct{}, config.Network.MaxConnections)
	}

	if config.Network.IdleTimeout != 0 {
		server.idleTimeout = config.Network.IdleTimeout
	}

	return server, nil
}

func (s *TCPServer) Handle(ctx context.Context, handler Handler) {
	defer s.listener.Close()

	select {
	case <-ctx.Done():
		break
	default:
		for {
			connection, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("Listener accept error", zap.Error(err))
				continue
			}

			s.activeConnections <- struct{}{}
			go func(conn net.Conn) {
				defer func() {
					if err := conn.Close(); err != nil {
						s.logger.Error("Connection close error", zap.Error(err))
					}
					<-s.activeConnections
				}()
				s.handleConnection(ctx, conn, handler)
			}(connection)
		}
	}
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn, handler Handler) {
	buffer := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			if err := connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Error("Set read deadline for connection error", zap.Error(err))
				s.writeResponse(connection, "Internal error")
				break
			}
			if err := connection.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("Set write deadline for connection error", zap.Error(err))
				s.writeResponse(connection, "Internal error")
				break
			}
		}

		readBytesCount, err := connection.Read(buffer)
		if err != nil && err != io.EOF {
			s.logger.Error(
				"Read data from connection error",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			s.writeResponse(connection, "Read data from connection error")
			break
		}
		if readBytesCount >= s.bufferSize {
			s.logger.Error("Read data error: small buffer size", zap.Int("buffer_size", s.bufferSize), zap.Int("readBytesCount", readBytesCount))
			s.writeResponse(connection, "Read data error: small buffer size")
			break
		}
		if readBytesCount == 0 {
			continue
		}

		response, err := handler.Handle(string(buffer[:readBytesCount]))
		if err != nil {
			s.logger.Error("handler.Handle error", zap.Error(err))
			s.writeResponse(connection, err.Error())
			break
		}
		s.writeResponse(connection, response)
	}
}

func (s *TCPServer) writeResponse(connection net.Conn, response string) {
	if _, err := connection.Write([]byte(response)); err != nil {
		s.logger.Error(
			"Write data to connection error",
			zap.String("address", connection.RemoteAddr().String()),
			zap.Error(err),
		)
	}
}
