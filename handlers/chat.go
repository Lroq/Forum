package handlers

import (
	"Forum/database"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var HubInstance = Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (c *Client) readPump() {
	defer func() {
		HubInstance.unregister <- c
		c.conn.Close()
		log.Println("Client disconnected:", c.conn.RemoteAddr())
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println("Received message:", string(message))
		HubInstance.broadcast <- message
	}
}

// Méthode pour écrire les messages au client
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write error:", err)
			return
		}
	}
	// Si on sort de la boucle, le canal est fermé, donc on envoie un message de fermeture.
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client registered:", client.conn.RemoteAddr())
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client unregistered:", client.conn.RemoteAddr())
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					log.Println("Failed to send message to client, closing connection:", client.conn.RemoteAddr())
				}
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.IsNew {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Error(w, "Failed to get username from session", http.StatusInternalServerError)
		return
	}

	userUUID, ok := session.Values["user_uuid"].(string)
	if !ok {
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}

	role, err := database.GetUserRoleByUUID(userUUID)
	if err != nil {
		log.Printf("Failed to get user role: %v", err)
		http.Error(w, "Failed to validate user role", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	
	if role == "GUEST" {
		message := websocket.FormatCloseMessage(4001, "Unauthorized")
		conn.WriteMessage(websocket.CloseMessage, message)
		conn.Close()
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client
	go client.writePump()
	go client.readPump()
	log.Println("Client connected:", conn.RemoteAddr())
}



