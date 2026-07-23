package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tpt-titan/backend/models"
)

// ChatService handles chat-related business logic
type ChatService struct {
	db *sql.DB
}

// NewChatService creates a new chat service
func NewChatService(db *sql.DB) *ChatService {
	return &ChatService{db: db}
}

// GetChatRooms retrieves all chat rooms for a user
func (s *ChatService) GetChatRooms(userID uuid.UUID) ([]models.ChatRoomResponse, error) {
	query := `
		SELECT cr.id, cr.name, cr.description, cr.room_type, cr.owner_id, cr.is_private, cr.created_at, cr.updated_at,
			   COUNT(DISTINCT cp.user_id) as participant_count,
			   COUNT(DISTINCT cm.id) as message_count,
			   MAX(cm.created_at) as last_message_time
		FROM chat_rooms cr
		JOIN chat_participants cp ON cr.id = cp.room_id
		LEFT JOIN chat_messages cm ON cr.id = cm.room_id
		WHERE cp.user_id = $1
		GROUP BY cr.id, cr.name, cr.description, cr.room_type, cr.owner_id, cr.is_private, cr.created_at, cr.updated_at
		ORDER BY COALESCE(MAX(cm.created_at), cr.created_at) DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat rooms: %w", err)
	}
	defer rows.Close()

	var rooms []models.ChatRoomResponse
	for rows.Next() {
		var room models.ChatRoom
		var participantCount, messageCount int
		var lastMessageTime *time.Time

		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.RoomType, &room.OwnerID,
			&room.IsPrivate, &room.CreatedAt, &room.UpdatedAt,
			&participantCount, &messageCount, &lastMessageTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat room: %w", err)
		}

		response := room.ToResponse()

		// Get participants
		participants, err := s.GetRoomParticipants(room.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get participants for room %s: %w", room.ID, err)
		}
		response.Participants = participants

		// Get last message if exists
		if messageCount > 0 {
			lastMessage, err := s.getLastMessage(room.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get last message for room %s: %w", room.ID, err)
			}
			response.LastMessage = lastMessage
		}

		// Calculate unread count
		unreadCount, err := s.getUnreadCount(userID, room.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get unread count for room %s: %w", room.ID, err)
		}
		response.UnreadCount = unreadCount

		rooms = append(rooms, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat rooms: %w", err)
	}

	return rooms, nil
}

// GetChatRoom retrieves a single chat room by ID
func (s *ChatService) GetChatRoom(userID, roomID uuid.UUID) (*models.ChatRoomResponse, error) {
	// First check if user is a participant
	var isParticipant bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM chat_participants WHERE room_id = $1 AND user_id = $2)
	`, roomID, userID).Scan(&isParticipant)
	if err != nil {
		return nil, fmt.Errorf("failed to check participation: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this room")
	}

	query := `
		SELECT id, name, description, room_type, owner_id, is_private, created_at, updated_at
		FROM chat_rooms WHERE id = $1
	`

	var room models.ChatRoom
	err = s.db.QueryRow(query, roomID).Scan(
		&room.ID, &room.Name, &room.Description, &room.RoomType, &room.OwnerID,
		&room.IsPrivate, &room.CreatedAt, &room.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chat room: %w", err)
	}

	response := room.ToResponse()

	// Get participants
	participants, err := s.GetRoomParticipants(room.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	response.Participants = participants

	// Get last message
	lastMessage, err := s.getLastMessage(room.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}
	response.LastMessage = lastMessage

	// Calculate unread count
	unreadCount, err := s.getUnreadCount(userID, room.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread count: %w", err)
	}
	response.UnreadCount = unreadCount

	return &response, nil
}

// CreateChatRoom creates a new chat room
func (s *ChatService) CreateChatRoom(userID uuid.UUID, req models.ChatRoomRequest) (*models.ChatRoomResponse, error) {
	// Validate request
	room := models.ChatRoom{
		Name:        &req.Name,
		Description: req.Description,
		RoomType:    req.RoomType,
		OwnerID:     userID,
		IsPrivate:   req.IsPrivate,
	}

	if err := room.Validate(); err != nil {
		return nil, fmt.Errorf("invalid room data: %w", err)
	}

	roomID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO chat_rooms (id, name, description, room_type, owner_id, is_private, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.Exec(query,
		roomID, req.Name, req.Description, req.RoomType, userID, req.IsPrivate, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat room: %w", err)
	}

	// Add participants including the creator
	allUserIDs := append([]uuid.UUID{userID}, req.UserIDs...)
	if err := s.addRoomParticipants(roomID, allUserIDs, userID); err != nil {
		return nil, fmt.Errorf("failed to add participants: %w", err)
	}

	// Return the created room
	return s.GetChatRoom(userID, roomID)
}

// AddRoomParticipants adds users to a chat room
func (s *ChatService) AddRoomParticipants(userID, roomID uuid.UUID, newUserIDs []uuid.UUID) error {
	// Check if user has permission (owner or admin)
	var userRole string
	err := s.db.QueryRow(`
		SELECT role FROM chat_participants WHERE room_id = $1 AND user_id = $2
	`, roomID, userID).Scan(&userRole)
	if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	}

	if userRole != "owner" && userRole != "admin" {
		return fmt.Errorf("insufficient permissions to add participants")
	}

	return s.addRoomParticipants(roomID, newUserIDs, userID)
}

// LeaveChatRoom removes a user from a chat room
func (s *ChatService) LeaveChatRoom(userID, roomID uuid.UUID) error {
	query := `DELETE FROM chat_participants WHERE room_id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, roomID, userID)
	if err != nil {
		return fmt.Errorf("failed to leave chat room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user is not a participant in this room")
	}

	// If room becomes empty, we could delete it, but for now we'll leave it
	return nil
}

// GetMessages retrieves messages for a chat room with pagination
func (s *ChatService) GetMessages(userID, roomID uuid.UUID, limit, offset int) ([]models.ChatMessageResponse, error) {
	// Check if user is a participant
	var isParticipant bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM chat_participants WHERE room_id = $1 AND user_id = $2)
	`, roomID, userID).Scan(&isParticipant)
	if err != nil {
		return nil, fmt.Errorf("failed to check participation: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this room")
	}

	query := `
		SELECT m.id, m.room_id, m.sender_id, m.content, m.message_type, m.file_url, m.file_name,
			   m.file_size, m.reply_to_id, m.edited_at, m.created_at,
			   u.username, u.email
		FROM chat_messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
	`

	args := []interface{}{roomID}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []models.ChatMessageResponse
	for rows.Next() {
		var message models.ChatMessage
		var senderName, senderEmail string

		err := rows.Scan(
			&message.ID, &message.RoomID, &message.SenderID, &message.Content, &message.MessageType,
			&message.FileURL, &message.FileName, &message.FileSize, &message.ReplyToID,
			&message.EditedAt, &message.CreatedAt, &senderName, &senderEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		response := message.ToResponse(senderName, senderEmail)

		// Get reactions for this message
		reactions, err := s.getMessageReactions(message.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get reactions for message %s: %w", message.ID, err)
		}
		response.Reactions = reactions

		messages = append(messages, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// SendMessage sends a new message to a chat room
func (s *ChatService) SendMessage(userID, roomID uuid.UUID, req models.ChatMessageRequest) (*models.ChatMessageResponse, error) {
	// Check if user is a participant
	var isParticipant bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM chat_participants WHERE room_id = $1 AND user_id = $2)
	`, roomID, userID).Scan(&isParticipant)
	if err != nil {
		return nil, fmt.Errorf("failed to check participation: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this room")
	}

	messageID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO chat_messages (id, room_id, sender_id, content, message_type, reply_to_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	messageType := req.MessageType
	if messageType == "" {
		messageType = "text"
	}

	_, err = s.db.Exec(query,
		messageID, roomID, userID, req.Content, messageType, req.ReplyToID, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Update room's updated_at timestamp
	_, err = s.db.Exec(`UPDATE chat_rooms SET updated_at = $1 WHERE id = $2`, now, roomID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update room timestamp: %v\n", err)
	}

	// Return the sent message
	return s.getMessageWithDetails(messageID)
}

// MarkRoomAsRead marks a room as read for a user
func (s *ChatService) MarkRoomAsRead(userID, roomID uuid.UUID) error {
	query := `UPDATE chat_participants SET last_read_at = $1 WHERE room_id = $2 AND user_id = $3`

	_, err := s.db.Exec(query, time.Now(), roomID, userID)
	return err
}

// AddReaction adds a reaction to a message
func (s *ChatService) AddReaction(userID, messageID uuid.UUID, reaction string) error {
	// Check if user can react to this message (must be in the same room)
	var roomID uuid.UUID
	err := s.db.QueryRow(`
		SELECT cm.room_id FROM chat_messages cm
		JOIN chat_participants cp ON cm.room_id = cp.room_id
		WHERE cm.id = $1 AND cp.user_id = $2
	`, messageID, userID).Scan(&roomID)
	if err != nil {
		return fmt.Errorf("message not found or user cannot react to it: %w", err)
	}

	reactionID := uuid.New()

	query := `
		INSERT INTO message_reactions (id, message_id, user_id, reaction, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (message_id, user_id, reaction) DO NOTHING
	`

	_, err = s.db.Exec(query, reactionID, messageID, userID, reaction, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

// RemoveReaction removes a reaction from a message
func (s *ChatService) RemoveReaction(userID, messageID uuid.UUID, reaction string) error {
	query := `DELETE FROM message_reactions WHERE message_id = $1 AND user_id = $2 AND reaction = $3`

	result, err := s.db.Exec(query, messageID, userID, reaction)
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reaction not found")
	}

	return nil
}

// Helper functions

func (s *ChatService) addRoomParticipants(roomID uuid.UUID, userIDs []uuid.UUID, firstUserID uuid.UUID) error {
	for i, userID := range userIDs {
		participantID := uuid.New()
		role := "member"
		if i == 0 && userID == firstUserID {
			role = "owner"
		}

		query := `
			INSERT INTO chat_participants (id, room_id, user_id, role, joined_at)
			VALUES ($1, $2, $3, $4, $5)
		`

		_, err := s.db.Exec(query, participantID, roomID, userID, role, time.Now())
		if err != nil {
			return fmt.Errorf("failed to add participant %s: %w", userID, err)
		}
	}
	return nil
}

func (s *ChatService) GetRoomParticipants(roomID uuid.UUID) ([]models.ChatParticipantResponse, error) {
	query := `
		SELECT cp.id, cp.room_id, cp.user_id, cp.role, cp.joined_at, cp.last_read_at,
			   u.username, u.email
		FROM chat_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.room_id = $1
		ORDER BY cp.joined_at ASC
	`

	rows, err := s.db.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	var participants []models.ChatParticipantResponse
	for rows.Next() {
		var participant models.ChatParticipant
		var userName, userEmail string

		err := rows.Scan(
			&participant.ID, &participant.RoomID, &participant.UserID, &participant.Role,
			&participant.JoinedAt, &participant.LastReadAt, &userName, &userEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}

		response := participant.ToResponse(userName, userEmail)
		participants = append(participants, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating participants: %w", err)
	}

	return participants, nil
}

func (s *ChatService) getLastMessage(roomID uuid.UUID) (*models.ChatMessageResponse, error) {
	query := `
		SELECT m.id, m.room_id, m.sender_id, m.content, m.message_type, m.file_url, m.file_name,
			   m.file_size, m.reply_to_id, m.edited_at, m.created_at,
			   u.username, u.email
		FROM chat_messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT 1
	`

	var message models.ChatMessage
	var senderName, senderEmail string

	err := s.db.QueryRow(query, roomID).Scan(
		&message.ID, &message.RoomID, &message.SenderID, &message.Content, &message.MessageType,
		&message.FileURL, &message.FileName, &message.FileSize, &message.ReplyToID,
		&message.EditedAt, &message.CreatedAt, &senderName, &senderEmail,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}

	response := message.ToResponse(senderName, senderEmail)
	return &response, nil
}

func (s *ChatService) getUnreadCount(userID, roomID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM chat_messages m
		WHERE m.room_id = $1 AND m.created_at > (
			SELECT COALESCE(cp.last_read_at, '1970-01-01'::timestamp)
			FROM chat_participants cp
			WHERE cp.room_id = $1 AND cp.user_id = $2
		)
	`

	var count int
	err := s.db.QueryRow(query, roomID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (s *ChatService) getMessageWithDetails(messageID uuid.UUID) (*models.ChatMessageResponse, error) {
	query := `
		SELECT m.id, m.room_id, m.sender_id, m.content, m.message_type, m.file_url, m.file_name,
			   m.file_size, m.reply_to_id, m.edited_at, m.created_at,
			   u.username, u.email
		FROM chat_messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = $1
	`

	var message models.ChatMessage
	var senderName, senderEmail string

	err := s.db.QueryRow(query, messageID).Scan(
		&message.ID, &message.RoomID, &message.SenderID, &message.Content, &message.MessageType,
		&message.FileURL, &message.FileName, &message.FileSize, &message.ReplyToID,
		&message.EditedAt, &message.CreatedAt, &senderName, &senderEmail,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	response := message.ToResponse(senderName, senderEmail)

	// Get reactions
	reactions, err := s.getMessageReactions(message.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}
	response.Reactions = reactions

	return &response, nil
}

func (s *ChatService) getMessageReactions(messageID uuid.UUID) ([]models.MessageReactionResponse, error) {
	query := `
		SELECT mr.id, mr.message_id, mr.user_id, mr.reaction, mr.created_at, u.username
		FROM message_reactions mr
		JOIN users u ON mr.user_id = u.id
		WHERE mr.message_id = $1
		ORDER BY mr.created_at ASC
	`

	rows, err := s.db.Query(query, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reactions: %w", err)
	}
	defer rows.Close()

	var reactions []models.MessageReactionResponse
	for rows.Next() {
		var reaction models.MessageReaction
		var userName string

		err := rows.Scan(
			&reaction.ID, &reaction.MessageID, &reaction.UserID, &reaction.Reaction,
			&reaction.CreatedAt, &userName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reaction: %w", err)
		}

		response := reaction.ToResponse(userName)
		reactions = append(reactions, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reactions: %w", err)
	}

	return reactions, nil
}

// CreateDirectMessage creates or finds a direct message room between two users
func (s *ChatService) CreateDirectMessage(userID, otherUserID uuid.UUID) (*models.ChatRoomResponse, error) {
	// Check if a direct message room already exists between these users
	query := `
		SELECT cr.id FROM chat_rooms cr
		WHERE cr.room_type = 'direct' AND cr.owner_id = $1
		AND EXISTS (
			SELECT 1 FROM chat_participants cp1
			JOIN chat_participants cp2 ON cp1.room_id = cp2.room_id
			WHERE cp1.room_id = cr.id AND cp1.user_id = $1 AND cp2.user_id = $2
		)
		LIMIT 1
	`

	var existingRoomID uuid.UUID
	err := s.db.QueryRow(query, userID, otherUserID).Scan(&existingRoomID)
	if err == nil {
		// Room exists, return it
		return s.GetChatRoom(userID, existingRoomID)
	}

	// Room doesn't exist, create a new one
	req := models.ChatRoomRequest{
		RoomType:  "direct",
		UserIDs:   []uuid.UUID{otherUserID},
		IsPrivate: true,
	}

	return s.CreateChatRoom(userID, req)
}
