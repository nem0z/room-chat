package storage

type Message struct {
	Tag  string `json:"tag"`
	Data string `json:"data"`
}

type Storage interface {
	WriteOne(msg Message) error
	ReadAll(tag string) ([]Message, error)
}
