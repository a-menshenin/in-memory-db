package compute

import (
	"fmt"
	"umemory/internal/storage"

	"go.uber.org/zap"
)

type Handler interface {
	Handle(p Parser) string
}

type ComputeHandler struct{
	storage storage.Storage
	requestParser Parser
	logger *zap.Logger
}

func NewComputeHandler(
	storage storage.Storage,
	requestParser Parser,
	logger *zap.Logger,
) *ComputeHandler {
	return &ComputeHandler{
		storage: storage,
		requestParser: requestParser,
		logger: logger,
	}
}

func (c *ComputeHandler) Handle(requestStr string) {
	command, args, err := c.requestParser.ParseArgs(requestStr)
	if err != nil {
		c.logger.Error("requestParser.ParseArgs error", zap.Error(err))
		fmt.Printf("Ошибка при парсинге аргументов: %s", err.Error())
		
		return
	}

	switch command {
	case GetCmd:
		v, found := c.storage.Get(args[0])
		if !found {
			c.logger.Error("storage.Get error: value not found")
			fmt.Printf("Значение по ключу %s не найдено", args[0])

			return
		}

		fmt.Printf("Найдено значение: %s", v)

		return
	case SetCmd:
		c.storage.Set(args[0], args[1])

		fmt.Printf("Значение %s сохранено", args[1])

		return
	case DeleteCmd:
		c.storage.Delete(args[0])

		fmt.Printf("Значение %s удалено", args[0])
	}
}
