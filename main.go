package main

import (
	"log"
	"net/http"

	"github.com/ana117/go-chat/server"
)

func main() {
	http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/chat", server.ChatHandler)
	http.HandleFunc("/leave", server.LeaveHandler)
	http.HandleFunc("/ws", server.WebSocketHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
