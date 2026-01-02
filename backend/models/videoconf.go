package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Meeting represents a video conference meeting
type Meeting struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	Description     *string    `json:"description,omitempty" db:"description"`
	HostID          uuid.UUID  `json:"host_id" db:"host_id"`
	RoomID          string     `json:"room_id" db:"room_id"` // Unique meeting room identifier
	MeetingType     string     `json:"meeting_type" db:"meeting_type"` // instant, scheduled, recurring
	StartTime       *time.Time `json:"start_time,omitempty" db:"start_time"`
	EndTime         *time.Time `json:"end_time,omitempty" db:"end_time"`
	TimeZone        string     `json:"time_zone" db:"time_zone"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsRecording     bool       `json:"is_recording" db:"is_recording"`
	RecordingPath   *string    `json:"recording_path,omitempty" db:"recording_path"`
	MaxParticipants int        `json:"max_participants" db:"max_participants"`
	RequireAuth     bool       `json:"require_auth" db:"require_auth"`
	PasswordHash    *string    `json:"-" db:"password_hash"` // For password-protected meetings
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// Participant represents a user in a meeting
type Participant struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	MeetingID  uuid.UUID  `json:"meeting_id" db:"meeting_id"`
	UserID     *uuid.UUID `json:"user_id,omitempty" db:"user_id"` // Null for guests
	Email      string     `json:"email" db:"email"`
	Name       string     `json:"name" db:"name"`
	Role       string     `json:"role" db:"role"` // host, cohost, participant, observer
	JoinedAt   time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt     *time.Time `json:"left_at,omitempty" db:"left_at"`
	IsMuted    bool       `json:"is_muted" db:"is_muted"`
	IsVideoOn  bool       `json:"is_video_on" db:"is_video_on"`
	IsScreenShare bool    `json:"is_screen_share" db:"is_screen_share"`
	IPAddress  string     `json:"ip_address" db:"ip_address"`
	UserAgent  string     `json:"user_agent" db:"user_agent"`
}

// WebRTCConnection represents a WebRTC peer connection
type WebRTCConnection struct {
	ID                uuid.UUID `json:"id" db:"id"`
	MeetingID         uuid.UUID `json:"meeting_id" db:"meeting_id"`
	ParticipantID     uuid.UUID `json:"participant_id" db:"participant_id"`
	ConnectionType    string    `json:"connection_type" db:"connection_type"` // offer, answer, ice
	SDP               string    `json:"sdp" db:"sdp"`                         // Session Description Protocol
	ICECandidates     []string  `json:"ice_candidates,omitempty" db:"ice_candidates"` // ICE candidates array
	TargetParticipant *uuid.UUID `json:"target_participant,omitempty" db:"target_participant"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	ExpiresAt         time.Time `json:"expires_at" db:"expires_at"`
}

// ChatMessage represents a chat message in a meeting
type MeetingChatMessage struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	MeetingID    uuid.UUID  `json:"meeting_id" db:"meeting_id"`
	ParticipantID uuid.UUID `json:"participant_id" db:"participant_id"`
	Message      string     `json:"message" db:"message"`
	MessageType  string     `json:"message_type" db:"message_type"` // text, system, file
	FileURL      *string    `json:"file_url,omitempty" db:"file_url"`
	FileName     *string    `json:"file_name,omitempty" db:"file_name"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// MeetingRecording represents a recorded meeting
type MeetingRecording struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	MeetingID  uuid.UUID  `json:"meeting_id" db:"meeting_id"`
	FilePath   string     `json:"file_path" db:"file_path"`
	FileSize   int64      `json:"file_size" db:"file_size"`
	Duration   int        `json:"duration" db:"duration"` // seconds
	Format     string     `json:"format" db:"format"`     // mp4, webm, etc.
	Quality    string     `json:"quality" db:"quality"`   // 720p, 1080p, etc.
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// ScreenShare represents an active screen share session
type ScreenShare struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	MeetingID      uuid.UUID  `json:"meeting_id" db:"meeting_id"`
	ParticipantID  uuid.UUID  `json:"participant_id" db:"participant_id"`
	StreamID       string     `json:"stream_id" db:"stream_id"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty" db:"ended_at"`
}

// MeetingSettings represents meeting configuration
type MeetingSettings struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	UserID                 uuid.UUID `json:"user_id" db:"user_id"`
	EnableWaitingRoom      bool      `json:"enable_waiting_room" db:"enable_waiting_room"`
	AllowGuests            bool      `json:"allow_guests" db:"allow_guests"`
	EnableRecording        bool      `json:"enable_recording" db:"enable_recording"`
	EnableChat             bool      `json:"enable_chat" db:"enable_chat"`
	MuteOnJoin             bool      `json:"mute_on_join" db:"mute_on_join"`
	VideoOnJoin            bool      `json:"video_on_join" db:"video_on_join"`
	MaxParticipants        int       `json:"max_participants" db:"max_participants"`
	MeetingTimeout         int       `json:"meeting_timeout" db:"meeting_timeout"` // minutes
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// MeetingInvite represents an invitation to a meeting
type MeetingInvite struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MeetingID uuid.UUID `json:"meeting_id" db:"meeting_id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Token     string    `json:"token" db:"token"`     // Unique invite token
	Status    string    `json:"status" db:"status"`   // pending, accepted, declined
	SentAt    time.Time `json:"sent_at" db:"sent_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty" db:"responded_at"`
}

// WebSocket Message Types for Video Conferencing
type MeetingMessageType string

const (
	MeetingMsgParticipantJoin    MeetingMessageType = "participant_join"
	MeetingMsgParticipantLeave   MeetingMessageType = "participant_leave"
	MeetingMsgWebRTCOffer        MeetingMessageType = "webrtc_offer"
	MeetingMsgWebRTCAnswer       MeetingMessageType = "webrtc_answer"
	MeetingMsgWebRTCIce          MeetingMessageType = "webrtc_ice"
	MeetingMsgChatMessage        MeetingMessageType = "chat_message"
	MeetingMsgScreenShareStart   MeetingMessageType = "screen_share_start"
	MeetingMsgScreenShareEnd     MeetingMessageType = "screen_share_end"
	MeetingMsgMuteToggle         MeetingMessageType = "mute_toggle"
	MeetingMsgVideoToggle        MeetingMessageType = "video_toggle"
	MeetingMsgRecordingStart     MeetingMessageType = "recording_start"
	MeetingMsgRecordingStop      MeetingMessageType = "recording_stop"
	MeetingMsgMeetingEnd         MeetingMessageType = "meeting_end"
)

// MeetingWebSocketMessage represents real-time messages in meetings
type MeetingWebSocketMessage struct {
	Type         MeetingMessageType `json:"type"`
	MeetingID    uuid.UUID          `json:"meeting_id"`
	ParticipantID uuid.UUID         `json:"participant_id"`
	UserID       *uuid.UUID         `json:"user_id,omitempty"`
	Data         interface{}        `json:"data"`
	Timestamp    time.Time          `json:"timestamp"`
}

// MeetingRequest represents the request payload for creating meetings
type MeetingRequest struct {
	Title           string     `json:"title" binding:"required"`
	Description     *string    `json:"description,omitempty"`
	MeetingType     string     `json:"meeting_type" binding:"required"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	TimeZone        string     `json:"time_zone,omitempty"`
	MaxParticipants int        `json:"max_participants,omitempty"`
	RequireAuth     bool       `json:"require_auth"`
	Password        *string    `json:"password,omitempty"`
	ParticipantEmails []string `json:"participant_emails,omitempty"`
}

// MeetingResponse represents the response payload for meetings
type MeetingResponse struct {
	ID              uuid.UUID            `json:"id"`
	Title           string               `json:"title"`
	Description     *string              `json:"description,omitempty"`
	HostID          uuid.UUID            `json:"host_id"`
	HostName        string               `json:"host_name"`
	RoomID          string               `json:"room_id"`
	MeetingType     string               `json:"meeting_type"`
	StartTime       *time.Time           `json:"start_time,omitempty"`
	EndTime         *time.Time           `json:"end_time,omitempty"`
	TimeZone        string               `json:"time_zone"`
	IsActive        bool                 `json:"is_active"`
	IsRecording     bool                 `json:"is_recording"`
	MaxParticipants int                  `json:"max_participants"`
	RequireAuth     bool                 `json:"require_auth"`
	ParticipantCount int                 `json:"participant_count"`
	Participants    []ParticipantResponse `json:"participants,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
}

// ParticipantResponse represents the response payload for participants
type ParticipantResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Email      string     `json:"email"`
	Name       string     `json:"name"`
	Role       string     `json:"role"`
	JoinedAt   time.Time  `json:"joined_at"`
	LeftAt     *time.Time `json:"left_at,omitempty"`
	IsMuted    bool       `json:"is_muted"`
	IsVideoOn  bool       `json:"is_video_on"`
	IsScreenShare bool    `json:"is_screen_share"`
}

// MeetingStats represents meeting statistics
type MeetingStats struct {
	TotalMeetings    int           `json:"total_meetings"`
	ActiveMeetings   int           `json:"active_meetings"`
	TotalParticipants int          `json:"total_participants"`
	TotalDuration    time.Duration `json:"total_duration"`
	AverageParticipants float64    `json:"average_participants"`
}

// ToResponse converts a Meeting to MeetingResponse
func (m *Meeting) ToResponse(hostName string, participantCount int) MeetingResponse {
	return MeetingResponse{
		ID:              m.ID,
		Title:           m.Title,
		Description:     m.Description,
		HostID:          m.HostID,
		HostName:        hostName,
		RoomID:          m.RoomID,
		MeetingType:     m.MeetingType,
		StartTime:       m.StartTime,
		EndTime:         m.EndTime,
		TimeZone:        m.TimeZone,
		IsActive:        m.IsActive,
		IsRecording:     m.IsRecording,
		MaxParticipants: m.MaxParticipants,
		RequireAuth:     m.RequireAuth,
		ParticipantCount: participantCount,
		Participants:    []ParticipantResponse{}, // Populated separately
		CreatedAt:       m.CreatedAt,
	}
}

// ToResponse converts a Participant to ParticipantResponse
func (p *Participant) ToResponse() ParticipantResponse {
	return ParticipantResponse{
		ID:         p.ID,
		UserID:     p.UserID,
		Email:      p.Email,
		Name:       p.Name,
		Role:       p.Role,
		JoinedAt:   p.JoinedAt,
		LeftAt:     p.LeftAt,
		IsMuted:    p.IsMuted,
		IsVideoOn:  p.IsVideoOn,
		IsScreenShare: p.IsScreenShare,
	}
}

// GenerateRoomID generates a unique room ID for meetings
func GenerateRoomID() string {
	return uuid.New().String()[:8] // 8-character room ID
}

// Validate checks if the meeting has valid data
func (m *Meeting) Validate() error {
	if m.Title == "" {
		return fmt.Errorf("meeting title is required")
	}

	validTypes := []string{"instant", "scheduled", "recurring"}
	for _, mt := range validTypes {
		if m.MeetingType == mt {
			goto typeValid
		}
	}
	return fmt.Errorf("invalid meeting type: %s", m.MeetingType)

typeValid:
	if m.MaxParticipants < 1 || m.MaxParticipants > 1000 {
		return fmt.Errorf("max participants must be between 1 and 1000")
	}

	if m.StartTime != nil && m.EndTime != nil && m.EndTime.Before(*m.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}

// IsActive checks if the meeting is currently active
func (m *Meeting) IsCurrentlyActive() bool {
	now := time.Now()
	if m.StartTime == nil || m.EndTime == nil {
		return false
	}
	return now.After(*m.StartTime) && now.Before(*m.EndTime)
}

// CanJoin checks if a participant can join the meeting
func (m *Meeting) CanJoin(participantCount int) bool {
	return participantCount < m.MaxParticipants
}

// GetDuration returns the meeting duration in minutes
func (m *Meeting) GetDuration() int {
	if m.StartTime == nil || m.EndTime == nil {
		return 0
	}
	return int(m.EndTime.Sub(*m.StartTime).Minutes())
}
