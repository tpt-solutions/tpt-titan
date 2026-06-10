package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
)

func isAllowedOrigin(origin string) bool {
	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		allowed = "http://localhost:3000,http://localhost:5173"
	}
	for _, o := range strings.Split(allowed, ",") {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // same-origin or non-browser client
		}
		return isAllowedOrigin(origin)
	},
}

// Hub maintains the set of active clients and broadcasts messages to the
type Hub struct {
	// Registered clients, keyed by user ID
	clients map[uuid.UUID]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan models.WebSocketMessage

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Chat service for database operations
	chatService *services.ChatService
}

// Client represents a WebSocket client
type Client struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan models.WebSocketMessage
	Hub    *Hub
}

// NewHub creates a new hub
func NewHub(chatService *services.ChatService) *Hub {
	return &Hub{
		clients:     make(map[uuid.UUID]map[*Client]bool),
		broadcast:   make(chan models.WebSocketMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		chatService: chatService,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// addClient adds a client to the hub
func (h *Hub) addClient(client *Client) {
	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*Client]bool)
	}
	h.clients[client.UserID][client] = true
	log.Printf("Client %s connected for user %s", client.ID, client.UserID)
}

// removeClient removes a client from the hub
func (h *Hub) removeClient(client *Client) {
	if clients, ok := h.clients[client.UserID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)
			log.Printf("Client %s disconnected for user %s", client.ID, client.UserID)

			// Clean up empty user map
			if len(clients) == 0 {
				delete(h.clients, client.UserID)
			}
		}
	}
}

// broadcastMessage broadcasts a message to appropriate clients
func (h *Hub) broadcastMessage(message models.WebSocketMessage) {
	// Get participants of the room
	if message.RoomID != nil {
		participants, err := h.chatService.GetRoomParticipants(*message.RoomID)
		if err != nil {
			log.Printf("Failed to get room participants: %v", err)
			return
		}

		// Send to all participants except sender
		for _, participant := range participants {
			if participant.UserID != message.UserID {
				if clients, ok := h.clients[participant.UserID]; ok {
					for client := range clients {
						select {
						case client.Send <- message:
						default:
							h.removeClient(client)
						}
					}
				}
			}
		}
	}
}

// HandleWebSocket handles WebSocket connections
func (h *Hub) HandleWebSocket(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		ID:     uuid.New(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan models.WebSocketMessage, 256),
		Hub:    h,
	}

	client.Hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMessage models.WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		wsMessage.UserID = c.UserID
		wsMessage.Timestamp = time.Now()

		// Handle different message types
		switch wsMessage.Type {
		case models.WSMessageSent:
			// Message was already sent via REST API, just broadcast it
			c.Hub.broadcast <- wsMessage
		case models.WSReactionAdded, models.WSReactionRemoved:
			c.Hub.broadcast <- wsMessage
		case models.WSTypingStart, models.WSTypingStop:
			c.Hub.broadcast <- wsMessage
		case models.WSUserStatus:
			// Broadcast user status changes
			c.Hub.broadcast <- wsMessage
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Failed to write message: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// BroadcastMessage broadcasts a message to all connected clients of a room
func (h *Hub) BroadcastMessage(message models.WebSocketMessage) {
	h.broadcast <- message
}

// BroadcastToUser broadcasts a message to all connected clients of a specific user
func (h *Hub) BroadcastToUser(userID uuid.UUID, message models.WebSocketMessage) {
	message.UserID = userID
	message.Timestamp = time.Now()

	if clients, ok := h.clients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				h.removeClient(client)
			}
		}
	}
}
