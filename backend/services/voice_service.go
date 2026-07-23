package services

import (
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VoiceService handles voice notes and annotations
type VoiceService struct {
	db *gorm.DB
}

// NewVoiceService creates a new voice service instance
func NewVoiceService() *VoiceService {
	return &VoiceService{
		db: config.DB,
	}
}

// GetVoiceNotes retrieves voice notes for a user with filtering options
func (s *VoiceService) GetVoiceNotes(userID uuid.UUID, limit, offset int, favoritesOnly bool, tag string) ([]models.VoiceNote, error) {
	var notes []models.VoiceNote
	query := s.db.Where("user_id = ?", userID)

	if favoritesOnly {
		query = query.Where("is_favorite = ?", true)
	}

	if tag != "" {
		// Search for tag in JSON array
		query = query.Where("tags @> ?", `["`+tag+`"]`)
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notes).Error
	return notes, err
}

// GetVoiceNote retrieves a specific voice note
func (s *VoiceService) GetVoiceNote(userID, noteID uuid.UUID) (*models.VoiceNote, error) {
	var note models.VoiceNote
	err := s.db.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// CreateVoiceNote creates a new voice note
func (s *VoiceService) CreateVoiceNote(note *models.VoiceNote) error {
	return s.db.Create(note).Error
}

// UpdateVoiceNote updates an existing voice note
func (s *VoiceService) UpdateVoiceNote(userID, noteID uuid.UUID, title string, tags []string, isFavorite, isPublic bool) error {
	updates := map[string]interface{}{
		"title":       title,
		"tags":        tags,
		"is_favorite": isFavorite,
		"is_public":   isPublic,
	}
	return s.db.Model(&models.VoiceNote{}).Where("id = ? AND user_id = ?", noteID, userID).Updates(updates).Error
}

// DeleteVoiceNote deletes a voice note
func (s *VoiceService) DeleteVoiceNote(userID, noteID uuid.UUID) error {
	return s.db.Where("id = ? AND user_id = ?", noteID, userID).Delete(&models.VoiceNote{}).Error
}

// GetVoiceAnnotations retrieves voice annotations for specific content
func (s *VoiceService) GetVoiceAnnotations(userID uuid.UUID, contentType string, contentID uuid.UUID) ([]models.VoiceAnnotation, error) {
	var annotations []models.VoiceAnnotation
	query := s.db.Where("user_id = ? AND content_type = ? AND content_id = ?", userID, contentType, contentID)

	// Include public annotations from other users
	query = query.Or("is_public = ? AND content_type = ? AND content_id = ?", true, contentType, contentID)

	err := query.Order("created_at DESC").Find(&annotations).Error
	return annotations, err
}

// GetVoiceAnnotation retrieves a specific voice annotation
func (s *VoiceService) GetVoiceAnnotation(userID, annotationID uuid.UUID) (*models.VoiceAnnotation, error) {
	var annotation models.VoiceAnnotation
	err := s.db.Where("(user_id = ? OR is_public = ?) AND id = ?", userID, true, annotationID).First(&annotation).Error
	if err != nil {
		return nil, err
	}
	return &annotation, nil
}

// CreateVoiceAnnotation creates a new voice annotation
func (s *VoiceService) CreateVoiceAnnotation(annotation *models.VoiceAnnotation) error {
	return s.db.Create(annotation).Error
}

// DeleteVoiceAnnotation deletes a voice annotation
func (s *VoiceService) DeleteVoiceAnnotation(userID, annotationID uuid.UUID) error {
	return s.db.Where("id = ? AND user_id = ?", annotationID, userID).Delete(&models.VoiceAnnotation{}).Error
}

// SearchVoiceNotes searches voice notes by content or title
func (s *VoiceService) SearchVoiceNotes(userID uuid.UUID, query string, limit, offset int) ([]models.VoiceNote, error) {
	var notes []models.VoiceNote
	searchQuery := "%" + query + "%"
	err := s.db.Where("user_id = ? AND (title ILIKE ? OR content ILIKE ?)", userID, searchQuery, searchQuery).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&notes).Error
	return notes, err
}

// GetVoiceNoteTags returns all unique tags for a user's voice notes
func (s *VoiceService) GetVoiceNoteTags(userID uuid.UUID) ([]string, error) {
	var tags []string
	err := s.db.Raw(`
		SELECT DISTINCT jsonb_array_elements_text(tags) as tag
		FROM voice_notes
		WHERE user_id = ?
		ORDER BY tag
	`, userID).Pluck("tag", &tags).Error
	return tags, err
}

// GetVoiceNoteStats returns statistics about voice notes
func (s *VoiceService) GetVoiceNoteStats(userID uuid.UUID) (map[string]interface{}, error) {
	var stats struct {
		TotalNotes    int64
		TotalDuration int64
		FavoriteCount int64
		PublicCount   int64
		TotalSize     int64
	}

	// Get basic counts
	s.db.Model(&models.VoiceNote{}).Where("user_id = ?", userID).Count(&stats.TotalNotes)
	s.db.Model(&models.VoiceNote{}).Where("user_id = ? AND is_favorite = ?", userID, true).Count(&stats.FavoriteCount)
	s.db.Model(&models.VoiceNote{}).Where("user_id = ? AND is_public = ?", userID, true).Count(&stats.PublicCount)

	// Get sum of durations and file sizes
	s.db.Model(&models.VoiceNote{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&stats.TotalDuration)
	s.db.Model(&models.VoiceNote{}).Where("user_id = ?", userID).Select("COALESCE(SUM(file_size), 0)").Scan(&stats.TotalSize)

	result := map[string]interface{}{
		"total_notes":    stats.TotalNotes,
		"total_duration": stats.TotalDuration,
		"favorite_count": stats.FavoriteCount,
		"public_count":   stats.PublicCount,
		"total_size":     stats.TotalSize,
	}

	return result, nil
}
