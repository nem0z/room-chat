package storage

type MemStore struct {
	messages []Message
}

func NewMemStore() *MemStore {
	return &MemStore{
		messages: make([]Message, 0, 10),
	}
}

func (store *MemStore) WriteOne(msg Message) error {
	store.messages = append(store.messages, msg)
	return nil
}

func (store *MemStore) ReadAll() ([]Message, error) {
	messages := make([]Message, len(store.messages))
	copy(messages, store.messages)

	return messages, nil
}
