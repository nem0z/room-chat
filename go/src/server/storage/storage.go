package storage

import pb "github.com/nem0z/room-chat/src/grpc_server/.server"

type Storage interface {
	WriteOne(msg *pb.Message) error
	ReadAll(tag string) ([]*pb.Message, error)
}
