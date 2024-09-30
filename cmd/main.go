package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"umemory/internal/compute"
	"umemory/internal/storage"

	"go.uber.org/zap"
)

var (
	err error
	logger *zap.Logger
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"app.log"}
	cfg.ErrorOutputPaths = []string{"app.log"}
	logger, err = cfg.Build()
    if err != nil {
		logger.Error("Create logger error", zap.Error(err))
		fmt.Println("Create logger error")
		
        os.Exit(1)
    }
	defer logger.Sync()

	storage := storage.NewInMemoryStorage()
	bufferReader := bufio.NewReader(os.Stdin)
	requestParser := compute.NewRequestParser()
	handler := compute.NewComputeHandler(storage, requestParser, logger)

	fmt.Println("\nSave/Get/Delete value by key")
	fmt.Println(`key/value available symbols: [a-zA-Zа-яА-Я0-9!?,.;:\"\'\ *#-=_@+№%$^/\|[]]`)

	for {
		fmt.Println("\nCommands: set key value || get key || delete key || exit")
		fmt.Print("You: ")

		requestStr, err := bufferReader.ReadString('\n')
		if err != nil && err != io.EOF {
			logger.Error("bufferReader.ReadString error", zap.Error(err))
			fmt.Println("Arguments read error")

			return
		}

		handler.Handle(requestStr)
	}
}
