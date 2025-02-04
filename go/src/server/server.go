package server

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/nem0z/room-chat/src/crypto"
	pb "github.com/nem0z/room-chat/src/grpc/chat"
	"github.com/nem0z/room-chat/src/server/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	mu    sync.Mutex
	store storage.Storage
}

func (s *Server) PostMessage(ctx context.Context, msg *pb.Message) (*pb.PostMessageResp, error) {
	if !crypto.VerifyWitness(msg.PubKey, msg.Witness, []byte(msg.Data)) {
		return &pb.PostMessageResp{
			ErrCode:    pb.ErrorCode_INVALID_WITNESS,
			ErrMessage: "Invalid signature",
		}, nil
	}

	err := s.store.WriteOne(msg)
	if err != nil {
		return &pb.PostMessageResp{
			ErrCode:    pb.ErrorCode_INTERNAL_ERROR,
			ErrMessage: err.Error(),
		}, nil
	}

	return &pb.PostMessageResp{
		Success: true,
	}, nil
}

func (s *Server) GetMessagesByTag(ctx context.Context, req *pb.GetMessagesReq) (*pb.GetMessagesResp, error) {
	tag := req.GetTag()
	messages, err := s.store.ReadAll(tag)
	if err != nil {
		return &pb.GetMessagesResp{
			ErrCode:    pb.ErrorCode_INTERNAL_ERROR,
			ErrMessage: err.Error(),
		}, nil
	}

	return &pb.GetMessagesResp{
		Messages: messages,
		Success:  true,
	}, nil
}

func Start(store storage.Storage) {
	if store == nil {
		store = storage.NewMemStore()
	}

	serv := &Server{
		mu:    sync.Mutex{},
		store: store,
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, serv)
	reflection.Register(grpcServer)

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
