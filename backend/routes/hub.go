package routes

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"tpt-titan/backend/models"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

// Hub manages WebSocket client connections for real-time chat
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan models.WebSocketMessage
	mu         sync.RWMutex
}

// Client is a single WebSocket connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
	roomID *uuid.UUID
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// NewHub creates a new WebSocket hub
func NewHub(chatSvc *services.ChatService) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan models.WebSocketMessage, 256),
	}
}

// Run starts the hub event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				if message.RoomID != nil && client.roomID != nil && *client.roomID != *message.RoomID {
					continue
				}
				if message.RoomID != nil && client.roomID == nil {
					continue
				}
				data, err := json.Marshal(message)
				if err != nil {
					continue
				}
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage sends a message to all connected clients
func (h *Hub) BroadcastMessage(message models.WebSocketMessage) {
	h.broadcast <- message
}

// BroadcastToUser sends a message to a specific user's connections
func (h *Hub) BroadcastToUser(userID uuid.UUID, message models.WebSocketMessage) {
	message.UserID = userID
	h.broadcast <- message
}

// HandleWebSocket upgrades an HTTP request to a WebSocket connection
func (h *Hub) HandleWebSocket(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (cl *Client) readPump() {
	defer func() {
		cl.hub.unregister <- cl
		cl.conn.Close()
	}()

	cl.conn.SetReadLimit(1 << 20)
	cl.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	cl.conn.SetPongHandler(func(string) error {
		cl.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := cl.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (cl *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		cl.conn.Close()
	}()

	for {
		select {
		case message, ok := <-cl.send:
			cl.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				cl.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := cl.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			cl.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := cl.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
