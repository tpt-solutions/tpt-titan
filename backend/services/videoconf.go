package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
	"tpt-titan/backend/models"
)

// VideoConfService handles video conferencing operations
type VideoConfService struct {
	db *sql.DB
}

// NewVideoConfService creates a new video conference service
func NewVideoConfService(db *sql.DB) *VideoConfService {
	return &VideoConfService{db: db}
}

// CreateMeeting creates a new video conference meeting
func (s *VideoConfService) CreateMeeting(userID uuid.UUID, req models.MeetingRequest) (*models.MeetingResponse, error) {
	// Validate request
	meeting := models.Meeting{
		ID:              uuid.New(),
		Title:           req.Title,
		Description:     req.Description,
		HostID:          userID,
		RoomID:          models.GenerateRoomID(),
		MeetingType:     req.MeetingType,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		TimeZone:        req.TimeZone,
		IsActive:        false,
		MaxParticipants: req.MaxParticipants,
		RequireAuth:     req.RequireAuth,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := meeting.Validate(); err != nil {
		return nil, fmt.Errorf("invalid meeting data: %w", err)
	}

	// Encrypt password if provided
	if req.Password != nil && *req.Password != "" {
		// In production, hash the password
		meeting.PasswordHash = req.Password // Simplified
	}

	query := `
		INSERT INTO meetings (id, title, description, host_id, room_id, meeting_type, start_time, end_time, time_zone, is_active, max_participants, require_auth, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := s.db.Exec(query,
		meeting.ID, meeting.Title, meeting.Description, meeting.HostID, meeting.RoomID,
		meeting.MeetingType, meeting.StartTime, meeting.EndTime, meeting.TimeZone,
		meeting.IsActive, meeting.MaxParticipants, meeting.RequireAuth, meeting.PasswordHash,
		meeting.CreatedAt, meeting.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create meeting: %w", err)
	}

	// Add participants if specified
	if len(req.ParticipantEmails) > 0 {
		if err := s.addMeetingParticipants(meeting.ID, req.ParticipantEmails); err != nil {
			log.Printf("Failed to add participants: %v", err)
			// Don't fail the meeting creation
		}
	}

	// Get host name
	var hostName string
	err = s.db.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&hostName)
	if err != nil {
		hostName = "Unknown"
	}

	response := meeting.ToResponse(hostName, 0)
	return &response, nil
}

// GetMeetings retrieves meetings for a user
func (s *VideoConfService) GetMeetings(userID uuid.UUID) ([]models.MeetingResponse, error) {
	query := `
		SELECT m.id, m.title, m.description, m.host_id, m.room_id, m.meeting_type,
			   m.start_time, m.end_time, m.time_zone, m.is_active, m.is_recording,
			   m.max_participants, m.require_auth, m.created_at,
			   u.username,
			   COUNT(mp.id) as participant_count
		FROM meetings m
		JOIN users u ON m.host_id = u.id
		LEFT JOIN meeting_participants mp ON m.id = mp.meeting_id
		WHERE m.host_id = $1 OR mp.email IN (SELECT email FROM users WHERE id = $1)
		GROUP BY m.id, m.title, m.description, m.host_id, m.room_id, m.meeting_type,
				 m.start_time, m.end_time, m.time_zone, m.is_active, m.is_recording,
				 m.max_participants, m.require_auth, m.created_at, u.username
		ORDER BY m.start_time DESC NULLS LAST, m.created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings: %w", err)
	}
	defer rows.Close()

	var meetings []models.MeetingResponse
	for rows.Next() {
		var meeting models.Meeting
		var hostName string
		var participantCount int

		err := rows.Scan(
			&meeting.ID, &meeting.Title, &meeting.Description, &meeting.HostID, &meeting.RoomID,
			&meeting.MeetingType, &meeting.StartTime, &meeting.EndTime, &meeting.TimeZone,
			&meeting.IsActive, &meeting.IsRecording, &meeting.MaxParticipants, &meeting.RequireAuth,
			&meeting.CreatedAt, &hostName, &participantCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}

		response := meeting.ToResponse(hostName, participantCount)
		meetings = append(meetings, response)
	}

	return meetings, nil
}

// GetMeeting retrieves a specific meeting
func (s *VideoConfService) GetMeeting(userID, meetingID uuid.UUID) (*models.MeetingResponse, error) {
	query := `
		SELECT m.id, m.title, m.description, m.host_id, m.room_id, m.meeting_type,
			   m.start_time, m.end_time, m.time_zone, m.is_active, m.is_recording,
			   m.max_participants, m.require_auth, m.created_at,
			   u.username,
			   COUNT(mp.id) as participant_count
		FROM meetings m
		JOIN users u ON m.host_id = u.id
		LEFT JOIN meeting_participants mp ON m.id = mp.meeting_id
		WHERE m.id = $1 AND (m.host_id = $2 OR mp.email IN (SELECT email FROM users WHERE id = $2))
		GROUP BY m.id, m.title, m.description, m.host_id, m.room_id, m.meeting_type,
				 m.start_time, m.end_time, m.time_zone, m.is_active, m.is_recording,
				 m.max_participants, m.require_auth, m.created_at, u.username
	`

	var meeting models.Meeting
	var hostName string
	var participantCount int

	err := s.db.QueryRow(query, meetingID, userID).Scan(
		&meeting.ID, &meeting.Title, &meeting.Description, &meeting.HostID, &meeting.RoomID,
		&meeting.MeetingType, &meeting.StartTime, &meeting.EndTime, &meeting.TimeZone,
		&meeting.IsActive, &meeting.IsRecording, &meeting.MaxParticipants, &meeting.RequireAuth,
		&meeting.CreatedAt, &hostName, &participantCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get meeting: %w", err)
	}

	response := meeting.ToResponse(hostName, participantCount)
	return &response, nil
}

// JoinMeeting allows a user to join a meeting
func (s *VideoConfService) JoinMeeting(meetingID uuid.UUID, email, name string) (*models.Participant, error) {
	// Check if meeting exists and is accessible
	var meeting models.Meeting
	err := s.db.QueryRow(`
		SELECT id, max_participants, require_auth, password_hash FROM meetings WHERE id = $1
	`, meetingID).Scan(
		&meeting.ID, &meeting.MaxParticipants, &meeting.RequireAuth, &meeting.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("meeting not found: %w", err)
	}

	// Check participant count
	var currentCount int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM meeting_participants WHERE meeting_id = $1 AND left_at IS NULL`, meetingID).Scan(&currentCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check participant count: %w", err)
	}

	if currentCount >= meeting.MaxParticipants {
		return nil, fmt.Errorf("meeting is full")
	}

	// Check if user is already a participant
	var existingParticipant models.Participant
	err = s.db.QueryRow(`
		SELECT id, joined_at, left_at FROM meeting_participants
		WHERE meeting_id = $1 AND email = $2
	`, meetingID, email).Scan(&existingParticipant.ID, &existingParticipant.JoinedAt, &existingParticipant.LeftAt)

	if err == nil && existingParticipant.LeftAt == nil {
		// User is already in the meeting
		return &existingParticipant, nil
	}

	// Create new participant
	participant := models.Participant{
		ID:         uuid.New(),
		MeetingID:  meetingID,
		Email:      email,
		Name:       name,
		Role:       "participant",
		JoinedAt:   time.Now(),
		IsMuted:    true,  // Default to muted
		IsVideoOn:  false, // Default to video off
		IsScreenShare: false,
		IPAddress:  "", // Would be populated from request
		UserAgent:  "", // Would be populated from request
	}

	query := `
		INSERT INTO meeting_participants (id, meeting_id, email, name, role, joined_at, is_muted, is_video_on, is_screen_share, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = s.db.Exec(query,
		participant.ID, participant.MeetingID, participant.Email, participant.Name,
		participant.Role, participant.JoinedAt, participant.IsMuted, participant.IsVideoOn,
		participant.IsScreenShare, participant.IPAddress, participant.UserAgent,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to join meeting: %w", err)
	}

	return &participant, nil
}

// LeaveMeeting allows a participant to leave a meeting
func (s *VideoConfService) LeaveMeeting(meetingID, participantID uuid.UUID) error {
	query := `UPDATE meeting_participants SET left_at = $1 WHERE id = $2 AND meeting_id = $3`

	result, err := s.db.Exec(query, time.Now(), participantID, meetingID)
	if err != nil {
		return fmt.Errorf("failed to leave meeting: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("participant not found in meeting")
	}

	return nil
}

// GetMeetingParticipants retrieves all participants in a meeting
func (s *VideoConfService) GetMeetingParticipants(meetingID uuid.UUID) ([]models.ParticipantResponse, error) {
	query := `
		SELECT mp.id, mp.email, mp.name, mp.role, mp.joined_at, mp.left_at,
			   mp.is_muted, mp.is_video_on, mp.is_screen_share
		FROM meeting_participants mp
		WHERE mp.meeting_id = $1 AND mp.left_at IS NULL
		ORDER BY mp.joined_at ASC
	`

	rows, err := s.db.Query(query, meetingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	var participants []models.ParticipantResponse
	for rows.Next() {
		var participant models.Participant

		err := rows.Scan(
			&participant.ID, &participant.Email, &participant.Name, &participant.Role,
			&participant.JoinedAt, &participant.LeftAt, &participant.IsMuted,
			&participant.IsVideoOn, &participant.IsScreenShare,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}

		response := participant.ToResponse()
		participants = append(participants, response)
	}

	return participants, nil
}

// StartWebRTCConnection creates a new WebRTC peer connection
func (s *VideoConfService) StartWebRTCConnection(meetingID, participantID uuid.UUID) (*webrtc.PeerConnection, error) {
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Store connection info in database
	connection := models.WebRTCConnection{
		ID:             uuid.New(),
		MeetingID:      meetingID,
		ParticipantID:  participantID,
		ConnectionType: "offer",
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	query := `
		INSERT INTO webrtc_connections (id, meeting_id, participant_id, connection_type, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = s.db.Exec(query,
		connection.ID, connection.MeetingID, connection.ParticipantID,
		connection.ConnectionType, connection.CreatedAt, connection.ExpiresAt,
	)

	if err != nil {
		peerConnection.Close()
		return nil, fmt.Errorf("failed to store connection: %w", err)
	}

	return peerConnection, nil
}

// HandleWebRTCSignal handles WebRTC signaling messages
func (s *VideoConfService) HandleWebRTCSignal(meetingID, participantID uuid.UUID, signalType string, signalData interface{}) error {
	// Store the signaling data
	data, err := json.Marshal(signalData)
	if err != nil {
		return fmt.Errorf("failed to marshal signal data: %w", err)
	}

	connection := models.WebRTCConnection{
		ID:             uuid.New(),
		MeetingID:      meetingID,
		ParticipantID:  participantID,
		ConnectionType: signalType,
		SDP:            string(data),
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(1 * time.Hour),
	}

	query := `
		INSERT INTO webrtc_connections (id, meeting_id, participant_id, connection_type, sdp, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = s.db.Exec(query,
		connection.ID, connection.MeetingID, connection.ParticipantID,
		connection.ConnectionType, connection.SDP, connection.CreatedAt, connection.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to store WebRTC signal: %w", err)
	}

	return nil
}

// SendMeetingMessage sends a chat message in a meeting
func (s *VideoConfService) SendMeetingMessage(meetingID, participantID uuid.UUID, message string) error {
	// Verify participant is in the meeting
	var isParticipant bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM meeting_participants WHERE meeting_id = $1 AND id = $2 AND left_at IS NULL)
	`, meetingID, participantID).Scan(&isParticipant)
	if err != nil || !isParticipant {
		return fmt.Errorf("participant not found in meeting")
	}

	meetingMessage := models.MeetingChatMessage{
		ID:           uuid.New(),
		MeetingID:    meetingID,
		ParticipantID: participantID,
		Message:      message,
		MessageType:  "text",
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO meeting_chat_messages (id, meeting_id, participant_id, message, message_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = s.db.Exec(query,
		meetingMessage.ID, meetingMessage.MeetingID, meetingMessage.ParticipantID,
		meetingMessage.Message, meetingMessage.MessageType, meetingMessage.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// GetMeetingMessages retrieves chat messages for a meeting
func (s *VideoConfService) GetMeetingMessages(meetingID uuid.UUID, limit, offset int) ([]models.MeetingChatMessage, error) {
	query := `
		SELECT mcm.id, mcm.meeting_id, mcm.participant_id, mcm.message, mcm.message_type,
			   mcm.file_url, mcm.file_name, mcm.created_at,
			   mp.name as participant_name
		FROM meeting_chat_messages mcm
		JOIN meeting_participants mp ON mcm.participant_id = mp.id
		WHERE mcm.meeting_id = $1
		ORDER BY mcm.created_at DESC
	`

	args := []interface{}{meetingID}

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

	var messages []models.MeetingChatMessage
	for rows.Next() {
		var message models.MeetingChatMessage
		var participantName string

		err := rows.Scan(
			&message.ID, &message.MeetingID, &message.ParticipantID, &message.Message,
			&message.MessageType, &message.FileURL, &message.FileName, &message.CreatedAt,
			&participantName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		messages = append(messages, message)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// StartRecording starts recording a meeting
func (s *VideoConfService) StartRecording(meetingID uuid.UUID) error {
	query := `UPDATE meetings SET is_recording = true WHERE id = $1`

	result, err := s.db.Exec(query, meetingID)
	if err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("meeting not found")
	}

	return nil
}

// StopRecording stops recording a meeting
func (s *VideoConfService) StopRecording(meetingID uuid.UUID) error {
	query := `UPDATE meetings SET is_recording = false WHERE id = $1`

	_, err := s.db.Exec(query, meetingID)
	if err != nil {
		return fmt.Errorf("failed to stop recording: %w", err)
	}

	return nil
}

// UpdateParticipantStatus updates a participant's status (mute, video, etc.)
func (s *VideoConfService) UpdateParticipantStatus(meetingID, participantID uuid.UUID, updates map[string]interface{}) error {
	query := `UPDATE meeting_participants SET `
	args := []interface{}{}
	argCount := 0

	if isMuted, ok := updates["is_muted"].(bool); ok {
		argCount++
		query += fmt.Sprintf("is_muted = $%d, ", argCount)
		args = append(args, isMuted)
	}

	if isVideoOn, ok := updates["is_video_on"].(bool); ok {
		argCount++
		query += fmt.Sprintf("is_video_on = $%d, ", argCount)
		args = append(args, isVideoOn)
	}

	if isScreenShare, ok := updates["is_screen_share"].(bool); ok {
		argCount++
		query += fmt.Sprintf("is_screen_share = $%d, ", argCount)
		args = append(args, isScreenShare)
	}

	if argCount == 0 {
		return fmt.Errorf("no valid updates provided")
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d AND meeting_id = $%d", argCount+1, argCount+2)
	args = append(args, participantID, meetingID)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update participant status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("participant not found in meeting")
	}

	return nil
}

// Helper functions

func (s *VideoConfService) addMeetingParticipants(meetingID uuid.UUID, emails []string) error {
	for _, email := range emails {
		invite := models.MeetingInvite{
			ID:        uuid.New(),
			MeetingID: meetingID,
			Email:     email,
			Token:     uuid.New().String(),
			Status:    "pending",
			SentAt:    time.Now(),
		}

		query := `
			INSERT INTO meeting_invites (id, meeting_id, email, token, status, sent_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err := s.db.Exec(query,
			invite.ID, invite.MeetingID, invite.Email, invite.Token,
			invite.Status, invite.SentAt,
		)

		if err != nil {
			return fmt.Errorf("failed to add invite for %s: %w", email, err)
		}
	}

	return nil
}
