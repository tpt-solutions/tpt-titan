package routes

import (
	"database/sql"
	"net/http"
	"strconv"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var chatService *services.ChatService
var wsHub *Hub

// WebSocketHub is a type alias for the Hub from the main package
type WebSocketHub = Hub

// InitChatService initializes the chat service (called from main)
func InitChatService(db *sql.DB) {
	chatService = services.NewChatService(db)
}

// GetChatService returns the chat service instance
func GetChatService() *services.ChatService {
	return chatService
}

// InitWebSocketHub initializes the WebSocket hub
func InitWebSocketHub(chatSvc *services.ChatService) *WebSocketHub {
	wsHub = NewHub(chatSvc)
	go wsHub.Run()
	return wsHub
}

// GetChatRooms returns all chat rooms for the authenticated user
func GetChatRooms(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	rooms, err := chatService.GetChatRooms(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// GetChatRoom returns a specific chat room
func GetChatRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	id, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	room, err := chatService.GetChatRoom(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat room"})
		return
	}

	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// CreateChatRoom creates a new chat room
func CreateChatRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.ChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := chatService.CreateChatRoom(userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat room"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"room": room})
}

// CreateDirectMessage creates or finds a direct message room
func CreateDirectMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	otherUserIDStr := c.Query("user_id")
	if otherUserIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Other user ID is required"})
		return
	}

	otherUserID, err := uuid.Parse(otherUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid other user ID"})
		return
	}

	room, err := chatService.CreateDirectMessage(userID.(uuid.UUID), otherUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create direct message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// AddRoomParticipants adds users to a chat room
func AddRoomParticipants(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	id, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var req struct {
		UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = chatService.AddRoomParticipants(userID.(uuid.UUID), id, req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participants added successfully"})
}

// LeaveChatRoom removes the user from a chat room
func LeaveChatRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	id, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	err = chatService.LeaveChatRoom(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave chat room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left chat room successfully"})
}

// GetMessages returns messages for a chat room
func GetMessages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	id, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Parse pagination parameters
	limit := 50 // default
	offset := 0 // default

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := chatService.GetMessages(userID.(uuid.UUID), id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage sends a new message to a chat room
func SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	id, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var req models.ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := chatService.SendMessage(userID.(uuid.UUID), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Broadcast the message via WebSocket
	if wsHub != nil {
		wsMessage := models.WebSocketMessage{
			Type:      models.WSMessageSent,
			RoomID:    &id,
			UserID:    userID.(uuid.UUID),
			Data:      message,
			Timestamp: message.CreatedAt,
		}
		wsHub.BroadcastMessage(wsMessage)
	}

	c.JSON(http.StatusCreated, gin.H{"message": message})
}

// MarkRoomAsRead marks a room as read for the user
func MarkRoomAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	id, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	err = chatService.MarkRoomAsRead(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark room as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room marked as read"})
}

// AddReaction adds a reaction to a message
func AddReaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID is required"})
		return
	}

	id, err := uuid.Parse(messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var req struct {
		Reaction string `json:"reaction" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = chatService.AddReaction(userID.(uuid.UUID), id, req.Reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction added successfully"})
}

// RemoveReaction removes a reaction from a message
func RemoveReaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID is required"})
		return
	}

	id, err := uuid.Parse(messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var req struct {
		Reaction string `json:"reaction" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = chatService.RemoveReaction(userID.(uuid.UUID), id, req.Reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed successfully"})
}
