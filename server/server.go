package server

import (
	"html/template"
	"net/http"
	"strconv"
)

type User struct {
	Username      string `json:"username"`
	ShortUsername string `json:"shortUsername"`
}

type Room struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Message struct {
	Sender  User   `json:"sender"`
	Message string `json:"message"`
}

type RoomPageData struct {
	DefaultRoom string `json:"defaultRoom"`
	Rooms       []Room `json:"rooms"`
}

type ChatPageData struct {
	Room    Room   `json:"room"`
	User    User   `json:"user"`
	Members []User `json:"members"`
}

var rooms = []Room{
	{Id: 1, Name: "Random Room 1", Description: "Random chat room 1 for everyone"},
	{Id: 2, Name: "Random Room 2", Description: "Random chat room 2 for everyone"},
	{Id: 3, Name: "Random Room 3", Description: "Random chat room 3 for everyone"},
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

	data := ChatPageData{
		User:    user,
		Room:    getRoomById(selectedRoomId),
		Members: GetMembers(selectedRoomId),
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
