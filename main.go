package main

import (
    "log"
    "math/rand"
    "net/http"
    "strings"
    "text/template"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Room struct {
    ID       string
    Clients  map[*websocket.Conn]bool
    Broadcast chan []byte
}

var Rooms map[string]*Room

func main() {
    Rooms = make(map[string]*Room)
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
    err = tmpl.Execute(w, nil)
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
    Rooms[code] = &Room{
        ID:       code,
        Clients:  make(map[*websocket.Conn]bool),
        Broadcast: make(chan []byte),
    }
    http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    code := parts[len(parts)-1]
    if code == "" {
        code = r.FormValue("code")
        http.Redirect(w, r, "/game/"+code, http.StatusSeeOther)
    }
    if _, ok := Rooms[code]; !ok {
        http.Redirect(w, r, "/notfound", http.StatusSeeOther)
    }

    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, code)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    roomID := r.URL.Query().Get("room")
    room, ok := Rooms[roomID]
    if !ok {
        http.Error(w, "Room not found", http.StatusNotFound)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    room.Clients[conn] = true

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            delete(room.Clients, conn)
            return
        }

        log.Println("Message reçu:", message, "type de message : ", messageType)
        room.Broadcast <- message
    }
}

func handleMessages() {
    for _, room := range Rooms {
        go func(room *Room) {
            for {
                select {
                case message := <-room.Broadcast:
                    for client := range room.Clients {
                        err := client.WriteMessage(websocket.TextMessage, message)
                        if err != nil {
                            log.Println("Error sending message to client:", err)
                            client.Close()
                            delete(room.Clients, client)
                        }
                    }
                }
            }
        }(room)
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
        if _, exists := Rooms[code]; !exists {
            fini = true
        }
    }
    return code
}
