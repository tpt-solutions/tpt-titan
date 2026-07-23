package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ChatRoom represents a chat room/conversation
type ChatRoom struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        *string   `json:"name,omitempty" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	RoomType    string    `json:"room_type" db:"room_type"` // direct, group, channel
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ChatParticipant represents a user in a chat room
type ChatParticipant struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	RoomID     uuid.UUID  `json:"room_id" db:"room_id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	Role       string     `json:"role" db:"role"` // owner, admin, member
	JoinedAt   time.Time  `json:"joined_at" db:"joined_at"`
	LastReadAt *time.Time `json:"last_read_at,omitempty" db:"last_read_at"`
}

// ChatMessage represents a message in a chat room
type ChatMessage struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	RoomID      uuid.UUID  `json:"room_id" db:"room_id"`
	SenderID    uuid.UUID  `json:"sender_id" db:"sender_id"`
	Content     string     `json:"content" db:"content"`
	MessageType string     `json:"message_type" db:"message_type"` // text, file, image, system
	FileURL     *string    `json:"file_url,omitempty" db:"file_url"`
	FileName    *string    `json:"file_name,omitempty" db:"file_name"`
	FileSize    *int64     `json:"file_size,omitempty" db:"file_size"`
	ReplyToID   *uuid.UUID `json:"reply_to_id,omitempty" db:"reply_to_id"`
	EditedAt    *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// MessageReaction represents a reaction to a message
type MessageReaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Reaction  string    `json:"reaction" db:"reaction"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserStatus represents a user's online status
type UserStatus struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Status       string    `json:"status" db:"status"` // online, away, busy, offline
	CustomStatus *string   `json:"custom_status,omitempty" db:"custom_status"`
	LastSeen     time.Time `json:"last_seen" db:"last_seen"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// VoiceNote represents an audio note recorded by a user
type VoiceNote struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"` // Auto-transcribed text
	AudioData   []byte    `json:"-" db:"audio_data"`    // Encrypted audio data
	AudioFormat string    `json:"audio_format" db:"audio_format"`
	Duration    int       `json:"duration" db:"duration"` // Duration in seconds
	FileSize    int64     `json:"file_size" db:"file_size"`
	Tags        []string  `json:"tags" db:"tags"` // JSON array of tags
	IsFavorite  bool      `json:"is_favorite" db:"is_favorite"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// VoiceAnnotation represents a voice annotation attached to content
type VoiceAnnotation struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	ContentType string    `json:"content_type" db:"content_type"` // document, task, email, calendar, contact
	ContentID   uuid.UUID `json:"content_id" db:"content_id"`     // ID of the content being annotated
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"` // Auto-transcribed text
	AudioData   []byte    `json:"-" db:"audio_data"`    // Encrypted audio data
	AudioFormat string    `json:"audio_format" db:"audio_format"`
	Duration    int       `json:"duration" db:"duration"` // Duration in seconds
	FileSize    int64     `json:"file_size" db:"file_size"`
	Position    *string   `json:"position,omitempty" db:"position"` // JSON position data for highlighting
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ChatRoomRequest represents the request payload for creating/updating chat rooms
type ChatRoomRequest struct {
	Name        string      `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	RoomType    string      `json:"room_type" binding:"required"`
	UserIDs     []uuid.UUID `json:"user_ids,omitempty"` // For adding participants
	IsPrivate   bool        `json:"is_private"`
}

// ChatMessageRequest represents the request payload for sending messages
type ChatMessageRequest struct {
	Content     string     `json:"content" binding:"required"`
	MessageType string     `json:"message_type,omitempty"`
	ReplyToID   *uuid.UUID `json:"reply_to_id,omitempty"`
}

// ChatRoomResponse represents the response payload for chat rooms
type ChatRoomResponse struct {
	ID           uuid.UUID                 `json:"id"`
	Name         *string                   `json:"name,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	RoomType     string                    `json:"room_type"`
	OwnerID      uuid.UUID                 `json:"owner_id"`
	IsPrivate    bool                      `json:"is_private"`
	Participants []ChatParticipantResponse `json:"participants,omitempty"`
	LastMessage  *ChatMessageResponse      `json:"last_message,omitempty"`
	UnreadCount  int                       `json:"unread_count"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// ChatParticipantResponse represents the response payload for chat participants
type ChatParticipantResponse struct {
	ID         uuid.UUID  `json:"id"`
	RoomID     uuid.UUID  `json:"room_id"`
	UserID     uuid.UUID  `json:"user_id"`
	UserName   string     `json:"user_name"`
	UserEmail  string     `json:"user_email"`
	Role       string     `json:"role"`
	JoinedAt   time.Time  `json:"joined_at"`
	LastReadAt *time.Time `json:"last_read_at,omitempty"`
}

// ChatMessageResponse represents the response payload for chat messages
type ChatMessageResponse struct {
	ID          uuid.UUID                 `json:"id"`
	RoomID      uuid.UUID                 `json:"room_id"`
	SenderID    uuid.UUID                 `json:"sender_id"`
	SenderName  string                    `json:"sender_name"`
	SenderEmail string                    `json:"sender_email"`
	Content     string                    `json:"content"`
	MessageType string                    `json:"message_type"`
	FileURL     *string                   `json:"file_url,omitempty"`
	FileName    *string                   `json:"file_name,omitempty"`
	FileSize    *int64                    `json:"file_size,omitempty"`
	ReplyToID   *uuid.UUID                `json:"reply_to_id,omitempty"`
	ReplyTo     *ChatMessageResponse      `json:"reply_to,omitempty"`
	Reactions   []MessageReactionResponse `json:"reactions,omitempty"`
	EditedAt    *time.Time                `json:"edited_at,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
}

// MessageReactionResponse represents the response payload for message reactions
type MessageReactionResponse struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Reaction  string    `json:"reaction"`
	CreatedAt time.Time `json:"created_at"`
}

// WebSocket message types for real-time communication
type WSMessageType string

const (
	WSMessageSent        WSMessageType = "message_sent"
	WSMessageEdited      WSMessageType = "message_edited"
	WSMessageDeleted     WSMessageType = "message_deleted"
	WSReactionAdded      WSMessageType = "reaction_added"
	WSReactionRemoved    WSMessageType = "reaction_removed"
	WSUserJoined         WSMessageType = "user_joined"
	WSUserLeft           WSMessageType = "user_left"
	WSUserStatus         WSMessageType = "user_status"
	WSTypingStart        WSMessageType = "typing_start"
	WSTypingStop         WSMessageType = "typing_stop"
	WSDocumentProcessing WSMessageType = "document_processing"
	WSDocumentProcessed  WSMessageType = "document_processed"
	WSDocumentFailed     WSMessageType = "document_failed"
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      WSMessageType `json:"type"`
	RoomID    *uuid.UUID    `json:"room_id,omitempty"`
	UserID    uuid.UUID     `json:"user_id"`
	UserName  string        `json:"user_name"`
	Data      interface{}   `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
}

// ToResponse converts a ChatRoom to ChatRoomResponse
func (cr *ChatRoom) ToResponse() ChatRoomResponse {
	return ChatRoomResponse{
		ID:           cr.ID,
		Name:         cr.Name,
		Description:  cr.Description,
		RoomType:     cr.RoomType,
		OwnerID:      cr.OwnerID,
		IsPrivate:    cr.IsPrivate,
		Participants: []ChatParticipantResponse{}, // Populated by service
		LastMessage:  nil,                         // Populated by service
		UnreadCount:  0,                           // Populated by service
		CreatedAt:    cr.CreatedAt,
		UpdatedAt:    cr.UpdatedAt,
	}
}

// ToResponse converts a ChatMessage to ChatMessageResponse
func (cm *ChatMessage) ToResponse(senderName, senderEmail string) ChatMessageResponse {
	return ChatMessageResponse{
		ID:          cm.ID,
		RoomID:      cm.RoomID,
		SenderID:    cm.SenderID,
		SenderName:  senderName,
		SenderEmail: senderEmail,
		Content:     cm.Content,
		MessageType: cm.MessageType,
		FileURL:     cm.FileURL,
		FileName:    cm.FileName,
		FileSize:    cm.FileSize,
		ReplyToID:   cm.ReplyToID,
		ReplyTo:     nil,                         // Populated by service if needed
		Reactions:   []MessageReactionResponse{}, // Populated by service
		EditedAt:    cm.EditedAt,
		CreatedAt:   cm.CreatedAt,
	}
}

// ToResponse converts a MessageReaction to MessageReactionResponse
func (mr *MessageReaction) ToResponse(userName string) MessageReactionResponse {
	return MessageReactionResponse{
		ID:        mr.ID,
		MessageID: mr.MessageID,
		UserID:    mr.UserID,
		UserName:  userName,
		Reaction:  mr.Reaction,
		CreatedAt: mr.CreatedAt,
	}
}

// ToResponse converts a ChatParticipant to ChatParticipantResponse
func (cp *ChatParticipant) ToResponse(userName, userEmail string) ChatParticipantResponse {
	return ChatParticipantResponse{
		ID:         cp.ID,
		RoomID:     cp.RoomID,
		UserID:     cp.UserID,
		UserName:   userName,
		UserEmail:  userEmail,
		Role:       cp.Role,
		JoinedAt:   cp.JoinedAt,
		LastReadAt: cp.LastReadAt,
	}
}

// Validate checks if the chat room has valid data
func (cr *ChatRoom) Validate() error {
	if cr.RoomType != "direct" && cr.RoomType != "group" && cr.RoomType != "channel" {
		return fmt.Errorf("invalid room type: must be direct, group, or channel")
	}

	if cr.RoomType != "direct" && (cr.Name == nil || *cr.Name == "") {
		return fmt.Errorf("name is required for group and channel rooms")
	}

	return nil
}

// GetDisplayName returns a display name for the chat room
func (cr *ChatRoom) GetDisplayName() string {
	if cr.Name != nil && *cr.Name != "" {
		return *cr.Name
	}

	// For direct messages, this would be populated with participant names
	return "Direct Message"
}

// IsUserParticipant checks if a user is a participant in the room
func (cr *ChatRoom) IsUserParticipant(userID uuid.UUID, participants []ChatParticipant) bool {
	for _, p := range participants {
		if p.UserID == userID {
			return true
		}
	}
	return false
}
