package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SetupStatus reports whether the instance still needs first-run setup.
// It is intentionally unauthenticated so the frontend wizard can decide whether
// to show itself before any admin account exists.
func GetSetupStatus(c *gin.Context) {
	var userCount int64
	if err := config.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check setup status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"setup_required": userCount == 0,
		"user_count":     userCount,
	})
}

// CompleteSetup creates the first admin account (and optionally persists
// instance SMTP settings). It fails if any user already exists, so it can only
// run once as a genuine first-run setup.
func CompleteSetup(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Username  string `json:"username" binding:"required,min=3,max=50"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		SMTP      *struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"smtp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userCount int64
	if err := config.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing users"})
		return
	}
	if userCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Setup has already been completed; an admin account already exists"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashed),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
		IsAdmin:      true,
		IsVerified:   true,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	persistSetupEncryptionSalt(config.DB, user.ID, req.Password)

	// Persist optional SMTP settings as system settings if provided.
	if req.SMTP != nil {
		settings := map[string]string{
			"SMTP_HOST":     req.SMTP.Host,
			"SMTP_PORT":     req.SMTP.Port,
			"SMTP_USERNAME": req.SMTP.Username,
			"SMTP_PASSWORD": req.SMTP.Password,
		}
		for k, v := range settings {
			if v == "" {
				continue
			}
			val, _ := json.Marshal(v)
			setting := models.SystemSetting{Key: k, Value: string(val)}
			if err := config.DB.Save(&setting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save SMTP settings"})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Setup completed. Admin account created.",
		"user":    user.ToResponse(),
	})
}

func persistSetupEncryptionSalt(db *gorm.DB, userID uuid.UUID, password string) {
	km, err := utils.NewKeyManager(password)
	if err != nil {
		return
	}
	salt := base64.StdEncoding.EncodeToString(km.GetSalt())
	db.Model(&models.User{}).Where("id = ?", userID).Update("encryption_salt", salt)
}
