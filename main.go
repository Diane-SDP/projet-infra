package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Map pour stocker les connexions clients
var broadcast = make(chan string)            // Channel de diffusion des messages

var upgrader = websocket.Upgrader{} // Upgrader pour la mise à niveau de la connexion HTTP à une connexion WebSocket

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler) // Nouvelle route pour gérer les connexions WebSocket

	go handleMessages() // Démarrer la goroutine pour gérer les messages diffusés aux clients

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Mise à niveau de la connexion HTTP à une connexion WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	// Ajouter la nouvelle connexion client à la map
	clients[ws] = true

	// Attendre les messages du client (ceci peut être omis si vous n'avez pas besoin de recevoir des messages du client)
	for {
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		// Récupérer le message de la chaîne de diffusion
		msg := <-broadcast

		// Envoyer le message à tous les clients connectés
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				delete(clients, client)
				client.Close()
			}
		}
	}
}
