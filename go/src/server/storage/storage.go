package storage

type Message string

type Storage interface {
	WriteOne(msg Message) error
	ReadAll() ([]Message, error)
}
