package server

import (
	"net/http"

	"github.com/nem0z/room-chat/src/server/storage"
)

func Start() {
	store := storage.NewMemStore()
	router := http.NewServeMux()

	router.HandleFunc("GET /messages", GetMessages(store))

	router.HandleFunc("POST /message", PostMessage(store))

	http.ListenAndServe(":8080", router)
}
