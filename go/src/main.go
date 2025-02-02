package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/nem0z/room-chat/src/grpc_server/.server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nem0z/room-chat/src/crypto"
	"github.com/nem0z/room-chat/src/server"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func sendPostRequestGRPC(url string, msg *pb.Message) error {
	// 1. Establish a gRPC connection
	// In a real application, you would likely have a more robust connection management
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials())) // Use insecure credentials for development/testing. For production, use TLS.
	if err != nil {
		return fmt.Errorf("did not connect: %w", err)
	}
	defer conn.Close()

	// 2. Create a gRPC client
	c := pb.NewChatServiceClient(conn)

	// 3. Call the gRPC method
	ctx := context.Background() // You might want to add timeouts or other context options
	resp, err := c.PostMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("could not post message: %w", err)
	}

	// 4. Handle the response
	if !resp.Success {
		// Convert the enum to a string for a more informative error
		errString := resp.ErrMessage
		if resp.ErrCode != pb.ErrorCode_UNKNOWN {
			errString = fmt.Sprintf("%v: %v", resp.ErrCode, resp.ErrMessage)
		}

		return fmt.Errorf("post message failed: %s", errString)
	}

	return nil
}

func main() {
	go server.Start(nil)

	pKey, err := crypto.GenPKey()
	Handle(err)

	data := "Salut"
	witness, err := crypto.GenWitness(pKey, []byte(data))
	Handle(err)

	msg := &pb.Message{
		PubKey:  crypto.PubKeyToBytes(pKey.PublicKey),
		Witness: witness,
		Tag:     "country",
		Data:    data,
	}

	grpcServerAddress := "localhost:50051"
	err = sendPostRequestGRPC(grpcServerAddress, msg)
	Handle(err)

	fmt.Println("gRPC request sent successfully")

	select {}
}
