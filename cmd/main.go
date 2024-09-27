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
		fmt.Println("Ошибка приложения: невозможно создать logger")
        os.Exit(1)
    }
	defer logger.Sync()

	storage := storage.NewInMemoryStorage()
	bufferReader := bufio.NewReader(os.Stdin)
	requestParser := compute.NewRequestParser()
	handler := compute.NewComputeHandler(storage, requestParser, logger)

	fmt.Println("\nЭто твоё хранилище. Сохрани / получи / удали данные по ключу")
	fmt.Println(`key и value могут состоять из символов [a-zA-Zа-яА-Я0-9!?,.;:\"\'\ *#-=_@+№%$^/\|[]]`)

	for {
		fmt.Println("\nКоманды: set key value || get key || delete key || exit")
		fmt.Print("Твоя команда: ")

		requestStr, err := bufferReader.ReadString('\n')
		if err != nil && err != io.EOF {
			logger.Error("bufferReader.ReadString error", zap.Error(err))
			fmt.Println("Ошибка при чтении аргументов")

			return
		}

		handler.Handle(requestStr)
	}
}
