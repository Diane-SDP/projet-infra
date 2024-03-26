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

type Joueur struct {
	Id     int
	Pseudo string
	Client map[*websocket.Conn]bool
}

type Room struct {
	Code       string
	LesJoueurs []Joueur
	BroadCast  chan []byte
	Couleur    string
}

var LesRooms []Room

var (
	buttonColor  = "green"
	clients      = make(map[*websocket.Conn]bool)
	broadcast    = make(chan []byte)
	addClient    = make(chan *websocket.Conn)
	removeClient = make(chan *websocket.Conn)
)

var pseudo string
var listGame []string

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/game/", gameHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/notfound", notfoundHandler)
	go handleMessages()
	http.ListenAndServe(":8080", nil)
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

func createHandler(w http.ResponseWriter, r *http.Request) {
	code := CodeGene()
	listGame = append(listGame, code)
	pseudo = r.FormValue("pseudo")
	var joueur Joueur
	joueur.Pseudo = pseudo
	var room Room
	room.LesJoueurs = append(room.LesJoueurs, joueur)
	room.BroadCast = make(chan []byte)
	room.Couleur = buttonColor
	room.Code = code
	LesRooms = append(LesRooms, room)
	http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	code := parts[len(parts)-1]
	if code == "" {
		code = r.FormValue("code")
		pseudo = r.FormValue("pseudo")
		if code == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
	}
	if !Contains(listGame, code) {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
	}
	type Data struct {
		Color  string
		Pseudo string
	}
	var data Data
	data.Pseudo = pseudo
	data.Color = buttonColor

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
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
	log.Println("kefrioguergeirug ouije pute")
	clients[conn] = true
	addClient <- conn
	for {
		log.Println("allez ça part")
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			removeClient <- conn
			return
		}

		buttonColor = string(message)
		print(buttonColor)
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
