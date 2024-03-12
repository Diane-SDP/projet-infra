package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

// défini la taille des msg en octets
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
var gameClients = make(map[string]map[*websocket.Conn]bool)
var gameBroadcast = make(map[string]chan []byte)
var gameAddClient = make(chan *websocket.Conn)
var gameRemoveClient = make(chan *websocket.Conn)

var listGame []string

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/BombGame", bombHandler)
	http.HandleFunc("/game/", gameHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/notfound", notfoundHandler)
	go handleMessages()
	http.ListenAndServe(":8080", nil)

	// Initialiser la gestion des messages pour chaque jeu créé
	for _, code := range listGame {
		gameClients[code] = make(map[*websocket.Conn]bool)
		gameBroadcast[code] = make(chan []byte)
		go handleGameMessages(code)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("home.html")
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

func handleGameMessages(code string) {
	for {
		select {
		case message := <-gameBroadcast[code]:
			clients := gameClients[code]
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error sending message to client:", err)
					client.Close()
					delete(clients, client)
				}
			}
		case client := <-gameAddClient:
			clients := gameClients[code]
			clients[client] = true
		case client := <-gameRemoveClient:
			clients := gameClients[code]
			delete(clients, client)
		}
	}
}

func notfoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("notfound.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func bombHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("bomba.html")
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

func createHandler(w http.ResponseWriter, r *http.Request) {
	println("gigapute")
	code := CodeGene()
	println("turbopute")
	listGame = append(listGame, code)
	println("pute")
	http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	code := parts[len(parts)-1]
	if code == "" {
		code = r.FormValue("code")
		http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
	}
	if !Contains(listGame, code) {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
	}

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

	// Récupérer le code de jeu depuis l'URL
	parts := strings.Split(r.URL.Path, "/")
	code := parts[len(parts)-1]

	// Ajouter la connexion WebSocket au jeu correspondant
	gameAddClient <- conn
	defer func() { gameRemoveClient <- conn }()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Diffuser le message uniquement au jeu correspondant
		gameBroadcast[code] <- message
	}
}

func handleMessages() {
	for {
		select {
		case message := <-broadcast:
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, message)
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

func CodeGene() string {
	alphabet := "azertyuiopqsdfghjklmwxcvbn"
	var code string
	var fini = false
	for !fini {
		for i := 0; i < 5; i++ {
			code += string(alphabet[rand.Intn(26)])
		}
		println("truc")
		if !Contains(listGame, code) {
			fini = true
			println("fini")
		}
	}
	return code
}

func Contains(liste []string, code string) bool {
	for i := 0; i < len(liste); i++ {
		if liste[i] == code {
			return true
		}
	}
	return false
}
