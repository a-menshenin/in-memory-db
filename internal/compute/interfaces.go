package compute

type Handler interface {
	Handle(p Parser) string
}

type Storage interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type Parser interface {
	ParseArgs(s string) (string, []string, error)
}
