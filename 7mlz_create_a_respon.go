package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Notification struct {
	Title string `json:"title"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/notify", handleNotification)

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", message)

		// Send message to all connected clients
		users := getUsers()
		for _, user := range users {
			err = user.WriteMessage(messageType, message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users := getUsers()
	for _, user := range users {
		err = user.WriteJSON(notification)
		if err != nil {
			fmt.Println(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func getUsers() []*websocket.Conn {
	// TODO: implement user storage and retrieval
	return []*websocket.Conn{}
}