package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
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

	server.bufferSize = config.Network.MaxMessageSize
	server.activeConnections = make(chan struct{}, config.Network.MaxConnections)

	return server, nil
}

func (s *TCPServer) Handle(ctx context.Context, handler Handler) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		
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
			go s.handleConnection(ctx, connection, handler)
		}
	}()

	<-ctx.Done()
	s.listener.Close()

	wg.Wait()
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn, handler Handler) {
	defer func() {
		if err := connection.Close(); err != nil {
			s.logger.Error("connection close error", zap.Error(err))
		}
		<- s.activeConnections
	}()

	buffer := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			if err := connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Error("Set read deadline for connection error", zap.Error(err))
				break
			}
			if err := connection.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("Set write deadline for connection error", zap.Error(err))
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
			break
		}
		if readBytesCount >= s.bufferSize {
			s.logger.Error("Read data error: small buffer size", zap.Int("buffer_size", s.bufferSize))
			break
		}

		response, err := handler.Handle(string(buffer[:readBytesCount]))
		if _, err := connection.Write([]byte(response)); err != nil {
			s.logger.Error(
				"Write data to connection error",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		}
	}
}
