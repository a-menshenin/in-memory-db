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

func (c *ComputeHandler) Handle(requestStr string) (string, error) {
	command, args, err := c.requestParser.ParseArgs(requestStr)
	if err != nil {
		c.logger.Error("requestParser.ParseArgs error", zap.Error(err))
		fmt.Printf("Arguments parse error: %s", err.Error())
		
		return "", fmt.Errorf("Arguments parse error: %s", err.Error())
	}

	switch command {
	case GetCmd:
		v, found := c.storage.Get(args[0])
		if !found {
			c.logger.Error("storage.Get error: value not found")
			fmt.Printf("Value by key=%s not found\n", args[0])

			return "value not found", fmt.Errorf("Value by key %s not found", args[0])
		}

		fmt.Printf("Value found: %s\n", v)

		return v, nil
	case SetCmd:
		c.storage.Set(args[0], args[1])

		fmt.Printf("Value %s saved\n", args[1])

		return "saved", nil
	case DeleteCmd:
		c.storage.Delete(args[0])

		fmt.Printf("Value %s deleted\n", args[0])

		return "deleted", nil
	default:
		return "Unknown command", nil
	}
}
