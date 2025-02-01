package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nem0z/room-chat/src/crypto"
	"github.com/nem0z/room-chat/src/server/storage"
)

type handler func(w http.ResponseWriter, r *http.Request)

func GetMessages(store storage.Storage) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")

		messages, err := store.ReadAll(tag)
		if err != nil {
			http.Error(w, "Error reading messages", http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "Error marshalling messages", http.StatusInternalServerError)
			return
		}

		w.Write(json)
	}
}

func PostMessage(store storage.Storage) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
			log.Println("Failed to decode request body:", err)
			return
		}

		var msg storage.Message
		err = json.Unmarshal(body, &msg)
		if err != nil {
			http.Error(w, "Failed to unmarshall body to message: "+err.Error(), http.StatusBadRequest)
			log.Println("Failed to unmarshall body message:", err)
			return
		}

		if !crypto.VerifyWitness(msg.SenderPubKey, msg.Witness, []byte(msg.Data)) {
			http.Error(w, "Witness is not valid", http.StatusBadRequest)
			log.Println("Witness is not valid")
			return
		}

		err = store.WriteOne(storage.Message(msg))
		if err != nil {
			http.Error(w, "Failed to write message to storage: "+err.Error(), http.StatusInternalServerError)
			log.Println("Failed to write message to storage:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
