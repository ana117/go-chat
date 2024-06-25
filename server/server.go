package server

import (
	"html/template"
	"net/http"
	"strconv"
)

type User struct {
	Username      string
	ShortUsername string
}

type Room struct {
	Id          int
	Name        string
	Description string
}

type RoomPageData struct {
	DefaultRoom string
	Rooms       []Room
}

type ChatPageData struct {
	Room     Room
	User     User
	Members  []User
	Messages []ChatMessage
}

type ChatMessage struct {
	Sender  User
	Message string
}

var rooms = []Room{
	{Id: 1, Name: "Random Room 1", Description: "Random chat room 1 for everyone"},
	{Id: 2, Name: "Random Room 2", Description: "Random chat room 2 for everyone"},
	{Id: 3, Name: "Random Room 3", Description: "Random chat room 3 for everyone"},
}

var mockUsers = []User{
	{Username: "User 1", ShortUsername: "U1"},
	{Username: "User 2", ShortUsername: "U2"},
	{Username: "User 3", ShortUsername: "U3"},
}

var mockChatMessages = []ChatMessage{
	{Sender: mockUsers[0], Message: "Hello!"},
	{Sender: mockUsers[1], Message: "Hi!"},
	{Sender: mockUsers[2], Message: "Hey!"},
	{Sender: mockUsers[0], Message: "How's it going?"},
	{Sender: mockUsers[1], Message: "Great, thanks!"},
	{Sender: mockUsers[2], Message: "Not too bad."},
	{Sender: mockUsers[0], Message: "Any plans for the weekend?"},
	{Sender: mockUsers[1], Message: "I'm going hiking."},
	{Sender: mockUsers[2], Message: "I'll probably just relax at home."},
	{Sender: mockUsers[0], Message: "Sounds nice."},
	{Sender: mockUsers[1], Message: "Yeah, I'm looking forward to it."},
	{Sender: mockUsers[2], Message: "Me too."},
	{Sender: mockUsers[0], Message: "By the way, have you seen the latest movie?"},
	{Sender: mockUsers[1], Message: "No, not yet."},
	{Sender: mockUsers[2], Message: "I haven't either."},
	{Sender: mockUsers[0], Message: "We should go watch it together."},
	{Sender: mockUsers[1], Message: "That's a great idea!"},
	{Sender: mockUsers[2], Message: "I'm in!"},
	{Sender: mockUsers[0], Message: "Alright, let's plan it."},
	{Sender: mockUsers[1], Message: "Sure, let's do it."},
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		selectedRoom := r.FormValue("room")

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "room",
			Value: selectedRoom,
		})
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	data := RoomPageData{
		DefaultRoom: rooms[0].Name,
		Rooms:       rooms,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, data)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
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

	user := User{Username: username, ShortUsername: username[:1]}
	members := []User{user, mockUsers[0], mockUsers[1], mockUsers[2]}

	chatMessages := make([]ChatMessage, len(mockChatMessages))
	copy(chatMessages, mockChatMessages)
	newMessage := ChatMessage{Sender: user, Message: "Hello, everyone!"}
	chatMessages = append(chatMessages, newMessage)

	data := ChatPageData{
		User:     user,
		Room:     getRoomById(selectedRoomId),
		Members:  members,
		Messages: chatMessages,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/chat.html"))
	tmpl.Execute(w, data)
}

func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "room",
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getCookieValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func getRoomById(id int) Room {
	for _, room := range rooms {
		if room.Id == id {
			return room
		}
	}

	return Room{}
}
