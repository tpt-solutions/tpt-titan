package routes

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Form represents a basic form structure returned to the client
type Form struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Fields      []Field   `json:"fields"`
	Status      string    `json:"status"`
	Responses   int       `json:"responses"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Field represents a form field returned to the client
type Field struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Label       string    `json:"label"`
	Placeholder string    `json:"placeholder"`
	Required    bool      `json:"required"`
	Options     []string  `json:"options"`
	Order       int       `json:"order"`
}

// FormSchema is the serializable internal representation stored encrypted
type formSchema struct {
	Fields []Field `json:"fields"`
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return uuid.Nil, false
	}
	userID, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return uuid.Nil, false
	}
	return userID, true
}

func encryptJSON(userID uuid.UUID, value interface{}) ([]byte, []byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}
	key := utils.DeriveUserDocumentKey(userID)
	km, err := utils.DeriveKeyFromPassword(key, salt)
	if err != nil {
		return nil, nil, err
	}
	plain, err := json.Marshal(value)
	if err != nil {
		return nil, nil, err
	}
	cipher, err := km.Encrypt(plain)
	if err != nil {
		return nil, nil, err
	}
	return cipher, salt, nil
}

func decryptJSON(userID uuid.UUID, cipher, salt []byte, out interface{}) error {
	key := utils.DeriveUserDocumentKey(userID)
	km, err := utils.DeriveKeyFromPassword(key, salt)
	if err != nil {
		return err
	}
	plain, err := km.Decrypt(cipher)
	if err != nil {
		return err
	}
	return json.Unmarshal(plain, out)
}

// GetForms returns all forms for the authenticated user
func GetForms(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var stored []models.EncryptedForm
	if err := config.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&stored).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve forms"})
		return
	}

	forms := make([]Form, 0, len(stored))
	for _, s := range stored {
		f, err := toForm(s)
		if err != nil {
			continue
		}
		forms = append(forms, f)
	}

	c.JSON(http.StatusOK, gin.H{"forms": forms})
}

// GetForm returns a specific form
func GetForm(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var stored models.EncryptedForm
	if err := config.DB.Where("id = ? AND user_id = ?", formID, userID).First(&stored).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	form, err := toForm(stored)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode form"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// CreateForm creates a new form
func CreateForm(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var formData struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Fields      []Field `json:"fields"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "draft"
	schema := formSchema{Fields: normalizeFields(formData.Fields)}

	schemaData, salt, err := encryptJSON(userID, schema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt form"})
		return
	}

	stored := models.EncryptedForm{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           formData.Name,
		Description:    formData.Description,
		EncryptedSchema: schemaData,
		Salt:           salt,
		Algorithm:      "AES-256-GCM",
		ResponseCount:  0,
		Status:         status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := config.DB.Create(&stored).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form"})
		return
	}

	form := Form{
		ID:          stored.ID,
		UserID:      userID,
		Name:        stored.Name,
		Description: stored.Description,
		Fields:      schema.Fields,
		Status:      stored.Status,
		Responses:   0,
		CreatedAt:   stored.CreatedAt,
		UpdatedAt:   stored.UpdatedAt,
	}

	c.JSON(http.StatusCreated, form)
}

// UpdateForm updates an existing form
func UpdateForm(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var formData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Fields      []Field `json:"fields"`
		Status      string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stored models.EncryptedForm
	if err := config.DB.Where("id = ? AND user_id = ?", formID, userID).First(&stored).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	if formData.Name != "" {
		stored.Name = formData.Name
	}
	if formData.Description != "" {
		stored.Description = formData.Description
	}
	if formData.Status != "" {
		stored.Status = formData.Status
	}
	if formData.Fields != nil {
		schema := formSchema{Fields: normalizeFields(formData.Fields)}
		schemaData, salt, err := encryptJSON(userID, schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt form"})
			return
		}
		stored.EncryptedSchema = schemaData
		stored.Salt = salt
	}
	stored.UpdatedAt = time.Now()

	if err := config.DB.Save(&stored).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update form"})
		return
	}

	form, err := toForm(stored)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode form"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// DeleteForm deletes a form and its responses
func DeleteForm(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var stored models.EncryptedForm
	if err := config.DB.Where("id = ? AND user_id = ?", formID, userID).First(&stored).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	if err := config.DB.Where("form_id = ?", formID).Delete(&models.EncryptedFormResponse{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form responses"})
		return
	}
	if err := config.DB.Delete(&stored).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form deleted successfully"})
}

// GetFormResponses gets responses for a form
func GetFormResponses(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var stored models.EncryptedForm
	if err := config.DB.Where("id = ? AND user_id = ?", formID, userID).First(&stored).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	var responses []models.EncryptedFormResponse
	if err := config.DB.Where("form_id = ?", formID).Order("submitted_at DESC").Find(&responses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve responses"})
		return
	}

	out := make([]gin.H, 0, len(responses))
	for _, r := range responses {
		var data map[string]interface{}
		if err := decryptJSON(userID, r.EncryptedData, r.Salt, &data); err != nil {
			data = map[string]interface{}{}
		}
		out = append(out, gin.H{
			"id":           r.ID,
			"form_id":      r.FormID,
			"responses":    data,
			"submitted_at": r.SubmittedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"responses": out})
}

// SubmitFormResponse submits a response to a form
func SubmitFormResponse(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var responseData map[string]interface{}
	if err := c.ShouldBindJSON(&responseData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stored models.EncryptedForm
	if err := config.DB.Where("id = ? AND user_id = ?", formID, userID).First(&stored).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	data, salt, err := encryptJSON(userID, responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt response"})
		return
	}

	now := time.Now()
	response := models.EncryptedFormResponse{
		ID:             uuid.New(),
		FormID:         formID,
		UserID:         userID,
		EncryptedData:  data,
		Salt:           salt,
		Algorithm:      "AES-256-GCM",
		SubmittedAt:    now,
	}

	if err := config.DB.Create(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit response"})
		return
	}

	if err := config.DB.Model(&stored).Update("response_count", stored.ResponseCount+1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update response count"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           response.ID,
		"form_id":      response.FormID,
		"responses":    responseData,
		"submitted_at": response.SubmittedAt,
	})
}

// --- helpers ---

// toForm decodes a stored encrypted form into the client-facing Form struct
func toForm(s models.EncryptedForm) (Form, error) {
	var schema formSchema
	if err := decryptJSON(s.UserID, s.EncryptedSchema, s.Salt, &schema); err != nil {
		schema = formSchema{Fields: []Field{}}
	}
	return Form{
		ID:          s.ID,
		UserID:      s.UserID,
		Name:        s.Name,
		Description: s.Description,
		Fields:      schema.Fields,
		Status:      s.Status,
		Responses:   s.ResponseCount,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}, nil
}

// normalizeFields assigns stable IDs and ordering to incoming fields
func normalizeFields(fields []Field) []Field {
	out := make([]Field, 0, len(fields))
	for i, f := range fields {
		if f.ID == uuid.Nil {
			f.ID = uuid.New()
		}
		f.Order = i + 1
		if f.Options == nil {
			f.Options = []string{}
		}
		out = append(out, f)
	}
	return out
}
