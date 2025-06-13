package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Conn        *websocket.Conn
	SessionCode string
}

var sessionClients = make(map[string]map[*Client]bool)
var mu sync.Mutex

func HandleWS(c *gin.Context) {
	sessionCode := c.Param("code")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		Conn:        conn,
		SessionCode: sessionCode,
	}

	// Register client
	mu.Lock()
	if sessionClients[sessionCode] == nil {
		sessionClients[sessionCode] = make(map[*Client]bool)
	}
	sessionClients[sessionCode][client] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(sessionClients[sessionCode], client)
		mu.Unlock()
		client.Conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// No read handling (one-way push)
	}
}

func BroadcastToSession(sessionCode string, event string, payload any) {
	mu.Lock()
	clients := sessionClients[sessionCode]
	mu.Unlock()

	for client := range clients {
		msg := gin.H{"event": event, "data": payload}
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("WebSocket send error:", err)
			client.Conn.Close()
			mu.Lock()
			delete(sessionClients[sessionCode], client)
			mu.Unlock()
		}
	}
}
