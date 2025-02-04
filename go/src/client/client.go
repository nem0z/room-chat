package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/nem0z/room-chat/src/crypto"
	pb "github.com/nem0z/room-chat/src/grpc/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.ChatServiceClient
	key    *ecdsa.PrivateKey
}

func New(target string, key *ecdsa.PrivateKey) (*Client, error) {
	var err error

	if key == nil {
		key, err = crypto.GenPKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate key: %v", err)
		}
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %v", err)
	}

	client := pb.NewChatServiceClient(conn)

	return &Client{
		client: client,
		key:    key,
	}, nil
}

func (c *Client) PostMessage(tag string, data string) error {
	witness, err := crypto.GenWitness(c.key, []byte(data))
	if err != nil {
		return err
	}

	msg := &pb.Message{
		PubKey:  crypto.PubKeyToBytes(c.key.PublicKey),
		Witness: witness,
		Tag:     tag,
		Data:    data,
	}

	resp, err := c.client.PostMessage(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("request to post message failed: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to post message (%v) : %v", resp.ErrCode, resp.ErrMessage)
	}

	return nil
}

func (c *Client) GetMessagesByTag(tag string) ([]*pb.Message, error) {
	resp, err := c.client.GetMessagesByTag(context.Background(), &pb.GetMessagesReq{Tag: tag})
	if err != nil {
		return nil, fmt.Errorf("request to get messages by tags failed: %v", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get messages by tag (%v): %v", resp.ErrCode, resp.ErrMessage)
	}

	return resp.Messages, nil
}
