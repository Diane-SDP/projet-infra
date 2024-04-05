package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

// défini la taille des msg en octets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Joueur struct {
	Uid    string
	Pseudo string
	Client *websocket.Conn
}

type Room struct {
	Code       string
	LesJoueurs []Joueur
	BroadCast  chan []byte
	Couleur    string
}

var LesRooms []Room
var AllPlayer []Joueur
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
	cookie, err := r.Cookie("uid")
	var uid string
	if err != nil {
		panic(err)
	} else {
		uid = cookie.Value
	}
	println("uid : ", uid)
	pseudo = r.FormValue("pseudo")

	var joueur Joueur
	joueur.Pseudo = pseudo
	joueur.Uid = uid
	joueur.Client = GetClientByUid(uid)
	println(joueur.Client)

	var room Room
	room.LesJoueurs = append(room.LesJoueurs, joueur)
	room.BroadCast = make(chan []byte)
	room.Couleur = buttonColor
	room.Code = code
	LesRooms = append(LesRooms, room)

	http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {

	// cookie, err := r.Cookie("uid")
	// if err != nil {
	// 	value := ""
	// } else {
	// 	value := cookie.Value
	// }
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
		Code   string
	}
	var data Data
	data.Pseudo = pseudo
	data.Color = buttonColor
	data.Code = code
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
		panic(err)
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	var uid string
	for i := 0; i < 16; i++ {
		uid += string(charset[seededRand.Intn(len(charset))])
	}
	println("envoie de l'uid...")
	if err := conn.WriteMessage(websocket.TextMessage, []byte(uid)); err != nil {
		println("erreur lors de l'envoie de l'uid")
		return
	}
	joueurcurrent := Joueur{
		Pseudo: "",
		Uid:    uid,
		Client: conn,
	}
	println("on créer l'uid : ", uid)
	var present = false
	for _, joueur := range AllPlayer {
		if joueur.Client == conn || joueur.Uid == uid {
			present = true
		}
	}
	if !present {
		println("on ajoute l'uid : ",joueurcurrent.Uid)
		AllPlayer = append(AllPlayer, joueurcurrent)
	}

	defer conn.Close()
	clients[conn] = true
	addClient <- conn

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			removeClient <- conn
			return
		}

		buttonColor = string(message)
		broadcast <- message
	}
}

func handleMessages() {
	for {
		select {
		case message := <-broadcast:
			for _, room := range LesRooms {
				if room.Code == strings.Split(string(message), "|")[2] {
					println("room trouvé !")
					for _, client := range room.LesJoueurs {
						println("message envoyé a : ", client.Uid)
						println(client.Client)
						err := client.Client.WriteMessage(websocket.TextMessage, message)
						if err != nil {
							panic(err)
							client.Client.Close()
							delete(clients, client.Client)
						}
					}
				}
			}
			// for client := range clients {

			// 	//on recupere le cookie
			// 	//on verifie si la salle qu'on a eu dans le message correspond au cookie code
			// 	//si oui on envoie le message
			// 	//si non on fait un petit print pour verifier (hassoul)
			// 	println("le message : ",strings.Split(string(message),"|")[2])
			// 	for _,room := range LesRooms {
			// 		if room.Code == strings.Split(string(message),"|")[2]{

			// 		}
			// 	}

			// 	err := client.WriteMessage(websocket.TextMessage, message)
			// 	if err != nil {
			// 		log.Println("Error sending message to client:", err)
			// 		client.Close()
			// 		delete(clients, client)
			// 	}
			// }
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
		if !Contains(listGame, code) {
			fini = true
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

func GetClientByUid(uid string) *websocket.Conn {
	println("nombre de gens co :", len(AllPlayer))
	for _, joueur := range AllPlayer {
		println("l'uid selec : ", uid, " et l'uid du joueur : ", joueur.Uid)
		if joueur.Uid == uid {
			return joueur.Client
		}
	}
	var truc *websocket.Conn
	println("getclientbyuid marche pas")
	return truc
}
