package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// TODO create STRUCT
var (
	upgrader  = websocket.Upgrader{}
	clients   = make(map[*websocket.Conn]bool) // connected clients
	broadcast = make(chan string)              // broadcast channel
)

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	clients[conn] = true
	broadcast <- strconv.Itoa(len(clients))

	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			broadcast <- strconv.Itoa(len(clients))
			break
		}
		broadcast <- strconv.Itoa(len(clients))
	}

}
func main() {
	go handleMessages()
	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	if err := http.ListenAndServe(":4567", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
