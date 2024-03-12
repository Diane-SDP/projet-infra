package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var buttonColor = "green"

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Écoute des messages WebSocket
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Message reçu:", message, "type de message : ", messageType)
		buttonColor = string(message)
		// Vous pouvez ajouter ici la logique pour diffuser le message à tous les autres clients
		// ou effectuer toute autre action en fonction du message reçu
		// Par exemple, si le message est une couleur, vous pouvez le diffuser à tous les autres clients connectés.

		// Vous pouvez également envoyer une réponse au client si nécessaire.
		// Par exemple, si le client envoie une demande de confirmation, vous pouvez répondre avec une confirmation.
		err = conn.WriteMessage(messageType, []byte(buttonColor))
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}
