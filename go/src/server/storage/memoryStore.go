package storage

import pb "github.com/nem0z/room-chat/src/grpc/chat"

type MemStore struct {
	messages map[string][]*pb.Message
}

func NewMemStore() *MemStore {
	return &MemStore{
		messages: make(map[string][]*pb.Message),
	}
}

func (store *MemStore) WriteOne(msg *pb.Message) error {
	store.messages[msg.Tag] = append(store.messages[msg.Tag], msg)
	return nil
}

func (store *MemStore) ReadAll(tag string) ([]*pb.Message, error) {
	messages := make([]*pb.Message, len(store.messages[tag]))
	copy(messages, store.messages[tag])
	return messages, nil
}
