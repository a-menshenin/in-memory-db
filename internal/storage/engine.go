package storage

type InMemoryStorage struct{
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	value, found := s.data[key]

	return value, found
}

func (s *InMemoryStorage) Set(key string, value string) {
	s.data[key] = value
}

func (s *InMemoryStorage) Delete(key string) {
	delete(s.data, key)
}

