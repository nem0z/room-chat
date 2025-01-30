package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nem0z/room-chat/src/server/storage"
)

type handler func(w http.ResponseWriter, r *http.Request)

func GetMessages(store storage.Storage) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		messages, err := store.ReadAll()
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
		message, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = store.WriteOne(storage.Message(message))
		if err != nil {
			http.Error(w, "Failed to write message to storage: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
