package storage

type MemStore struct {
	messages map[string][]Message
}

func NewMemStore() *MemStore {
	return &MemStore{
		messages: make(map[string][]Message),
	}
}

func (store *MemStore) WriteOne(msg Message) error {
	store.messages[msg.Tag] = append(store.messages[msg.Tag], msg)
	return nil
}

func (store *MemStore) ReadAll(tag string) ([]Message, error) {
	messages := make([]Message, len(store.messages[tag]))
	copy(messages, store.messages[tag])
	return messages, nil
}
