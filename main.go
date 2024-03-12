package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	buttonColor  = "green"
	clients      = make(map[*websocket.Conn]bool)
	broadcast    = make(chan []byte)
	addClient    = make(chan *websocket.Conn)
	removeClient = make(chan *websocket.Conn)
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	go handleMessages()
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, buttonColor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	addClient <- conn

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			removeClient <- conn
			return
		}

		buttonColor = string(message)
		log.Println("Message reçu:", message, "type de message : ", messageType)
		log.Println(buttonColor)
		broadcast <- message
	}
}

func handleMessages() {
	for {
		select {
		case message := <-broadcast:
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				println("message broadcast ", message)
				if err != nil {
					log.Println("Error sending message to client:", err)
					client.Close()
					delete(clients, client)
				}
			}
		case client := <-addClient:
			clients[client] = true
		case client := <-removeClient:
			delete(clients, client)
		}
	}
}
