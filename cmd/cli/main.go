package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"umemory/internal"
	"umemory/internal/network"

	"go.uber.org/zap"
)

func main() {
	cfg, err := internal.GetConfig()
	if err != nil {
		fmt.Println("Get config error: " + err.Error())

		return
	}
	logger, err := internal.CreateLogger(cfg)
	if err != nil {
		fmt.Println("Create logger error: " + err.Error())

		return
	}
	defer logger.Sync()

	tcpCfg := network.TCPClientConfig{}
	tcpCfg.Address = flag.String("address", cfg.Network.Address, "Connection host:port")
	tcpCfg.IdleTimeout = flag.Duration("idle_timeout", cfg.Network.IdleTimeout, "Connection Idle timeout")
	tcpCfg.MaxMessageSize = flag.Int("max_message_size", cfg.Network.MaxMessageSize, "Connection Max message size")
	flag.Parse()

	bufferReader := bufio.NewReader(os.Stdin)

	conn, err := net.Dial("tcp", *tcpCfg.Address)
	if err != nil {
		logger.Error("Connection create error", zap.Error(err))
		fmt.Println("Connection create error")

		return
	}
	tcpClient, err := network.NewTCPClient(tcpCfg, conn, logger)
	if err != nil {
		logger.Error("Create tcp client error", zap.Error(err))
		fmt.Println("Create tcp client error")

		return
	}
	defer tcpClient.Close()

	fmt.Println("\nSave/Get/Delete value by key")
	fmt.Println(`key/value available symbols: [a-zA-Zа-яА-Я0-9!?,.;:\"\'\ *#-=_@+№%$^/\|[]]`)

	for {
		fmt.Println("\nCommands: set key value || get key || delete key")
		fmt.Print("Your command: ")

		request, err := bufferReader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			logger.Error("Arguments read error", zap.Error(err))
			fmt.Println("Arguments read error")

			return
		}

		response, err := tcpClient.Send(request)
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				logger.Fatal("Connection was closed", zap.Error(err))
			}

			logger.Error("Send client request error", zap.Error(err))
			fmt.Println("Send client request error")

			continue
		}

		fmt.Print("Response: ")
		fmt.Println(string(response))
	}
}
