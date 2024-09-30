package compute

import (
	"fmt"

	"go.uber.org/zap"
)

type ComputeHandler struct{
	storage Storage
	requestParser Parser
	logger *zap.Logger
}

func NewComputeHandler(
	storage Storage,
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
		fmt.Printf("Arguments parse error: %s", err.Error())
		
		return
	}

	switch command {
	case GetCmd:
		v, found := c.storage.Get(args[0])
		if !found {
			c.logger.Error("storage.Get error: value not found")
			fmt.Printf("Value by key=%s not found", args[0])

			return
		}

		fmt.Printf("Value found: %s", v)

		return
	case SetCmd:
		c.storage.Set(args[0], args[1])

		fmt.Printf("Value %s saved", args[1])

		return
	case DeleteCmd:
		c.storage.Delete(args[0])

		fmt.Printf("Value %s deleted", args[0])
	}
}
