package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nem0z/room-chat/src/crypto"
	"github.com/nem0z/room-chat/src/server"
	"github.com/nem0z/room-chat/src/server/storage"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func sendPostRequest(url string, msg storage.Message) error {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with status: %d => %v", resp.StatusCode, resp)
	}

	return nil
}

func main() {
	go server.Start()

	pKey, err := crypto.GenPKey()
	Handle(err)

	data := "Salut"
	witness, err := crypto.GenWitness(pKey, []byte(data))
	Handle(err)

	msg := storage.Message{
		Alias:   crypto.GetAlias(pKey.PublicKey),
		PubKey:  crypto.PubKeyToBytes(pKey.PublicKey),
		Witness: witness,
		Tag:     "country",
		Data:    data,
	}

	err = sendPostRequest("http://localhost:8080/message", msg)
	Handle(err)

	select {}
}
