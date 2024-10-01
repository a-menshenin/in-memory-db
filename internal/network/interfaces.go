package network

type Handler interface {
	Handle(requestStr string) (string, error)
}
