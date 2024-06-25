package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	MessageType string          `json:"messageType"`
	Message     json.RawMessage `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

var roomClients = make(map[int]map[*websocket.Conn]string)

func GetMembers(selectedRoomId int) []User {
	members := make([]User, 0)
	for _, client := range roomClients[selectedRoomId] {
		members = append(members, User{
			Username:      client,
			ShortUsername: client[:1],
		})
	}
	return members
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	username, err := getCookieValue(r, "username")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	selectedRoom, err := getCookieValue(r, "room")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	selectedRoomId, err := strconv.Atoi(selectedRoom)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	if roomClients[selectedRoomId] == nil {
		roomClients[selectedRoomId] = make(map[*websocket.Conn]string)
	}
	roomClients[selectedRoomId][conn] = username
	defer conn.Close()

	broadcastUserAction(selectedRoomId, username, "join")
	for {
		var wsMessage WebSocketMessage
		err := conn.ReadJSON(&wsMessage)
		if err != nil {
			delete(roomClients[selectedRoomId], conn)
			broadcastUserAction(selectedRoomId, username, "leave")
			break
		}

		switch wsMessage.MessageType {
		case "chat":
			broadcastMessage(selectedRoomId, wsMessage)
		}
	}
}

func broadcastMessage(selectedRoomId int, wsMessage WebSocketMessage) {
	for client := range roomClients[selectedRoomId] {
		client.WriteJSON(wsMessage)
	}
}

func broadcastUserAction(selectedRoomId int, username string, actionType string) {
	userBytes, _ := json.Marshal(username)

	for client := range roomClients[selectedRoomId] {
		client.WriteJSON(WebSocketMessage{
			MessageType: actionType,
			Message:     userBytes,
		})
	}
}
