package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"umemory/internal"
	"umemory/internal/compute"
	"umemory/internal/network"
	"umemory/internal/storage"

	"github.com/joho/godotenv"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func init() {
    if err := godotenv.Load(); err != nil {
		fmt.Println("No env file found by path")
		os.Exit(1)
    }
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	
	cfg, err := internal.GetConfig()
	if err != nil {
		fmt.Println("Get config error: " + err.Error())

		return
	}
	logger, err := internal.CreateLogger(cfg)
	if err != nil {
		fmt.Println("Create logger error")

		return
	}
	defer logger.Sync()

	server, err := network.NewTCPServer(cfg, logger)
	if err != nil {
		logger.Error("Create tcp server error", zap.Error(err))
		fmt.Println("Create tcp server error")

		return
	}

	storage := storage.NewInMemoryStorage()
	requestParser := compute.NewRequestParser()
	handler := compute.NewComputeHandler(storage, requestParser, logger)

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		server.Handle(groupCtx, handler)

		return nil
	})

	if group.Wait() != nil {
		logger.Error("Server wait error", zap.Error(err))
		fmt.Println("Server wait error")

		return
	}
}
