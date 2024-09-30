package network

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
	"umemory/internal"
	"umemory/internal/network"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type TestHandler struct {}
func (h TestHandler) Handle(requestStr string) (string, error) {
	return "Response for " + requestStr, nil
}

func TestTCPServer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := internal.Config{
		Engine: struct{
			EngineType string `yaml:"engine_type"`
		}{
			EngineType: "in memory",
		},
		Network: struct{
			Address string `yaml:"address"`
			MaxConnections int `yaml:"max_connections",omitempty`
			MaxMessageSize int `yaml:"max_message_size",omitempty`
			IdleTimeout time.Duration `yaml:"idle_timeout",omitempty`
		}{
			Address: "localhost:22222",
			MaxConnections: 2,
			MaxMessageSize: 1024,
		},
	}
	server, err := network.NewTCPServer(cfg, zap.NewNop())
	if err != nil {
		fmt.Errorf("network.NewTCPServer error: %w", err)

		t.Fail()
	}

	go func() {
		server.Handle(ctx, TestHandler{})
	}()

	time.Sleep(100 * time.Millisecond)

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		connection, clientErr := net.Dial("tcp", cfg.Network.Address)
		if clientErr != nil {
			return fmt.Errorf("Client1 net.Dial error: %w", clientErr)
		}

		_, clientErr = connection.Write([]byte("client1"))
		if clientErr != nil {
			return fmt.Errorf("Client1 connection.Write error: %w", clientErr)
		}

		buffer := make([]byte, 1024)
		size, clientErr := connection.Read(buffer)
		if clientErr != nil {
			return fmt.Errorf("Client1 connection.Read error: %w", clientErr)
		}

		clientErr = connection.Close()
		if clientErr != nil {
			return fmt.Errorf("Client1 connection.Close error")
		}

		assert.Equal(t, "Response for client1", string(buffer[:size]))

		return nil
	})

	g.Go(func() error {
		connection, clientErr := net.Dial("tcp", cfg.Network.Address)
		if clientErr != nil {
			return fmt.Errorf("Client2 net.Dial error: %w", clientErr)
		}

		_, clientErr = connection.Write([]byte("client2"))
		if clientErr != nil {
			return fmt.Errorf("Client2 connection.Write error: %w", clientErr)
		}

		buffer := make([]byte, 1024)
		size, clientErr := connection.Read(buffer)
		if clientErr != nil {
			return fmt.Errorf("Client2 connection.Read error: %w", clientErr)
		}

		clientErr = connection.Close()
		if clientErr != nil {
			return fmt.Errorf("Client2 connection.Close error: %w", clientErr)
		}

		assert.Equal(t, "Response for client2", string(buffer[:size]))

		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())

		t.Fail()
	}
}
