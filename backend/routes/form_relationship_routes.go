package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// CreateRelationship creates a new relationship between forms
func CreateRelationship(c *gin.Context) {
	var rel services.Relationship
	if err := c.ShouldBindJSON(&rel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	// Validate relationship
	errors := relationshipSvc.ValidateRelationship(&rel)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	err := relationshipSvc.CreateRelationship(&rel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rel)
}

// GetFormRelationships gets relationships for a form
func GetFormRelationships(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	relationships, err := relationshipSvc.GetRelationshipsByForm(formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"relationships": relationships})
}

// CreateLookupField creates a lookup field
func CreateLookupField(c *gin.Context) {
	var lookup services.LookupField
	if err := c.ShouldBindJSON(&lookup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	err := relationshipSvc.CreateLookupField(&lookup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lookup)
}

// GetLookupData gets data for a lookup field
func GetLookupData(c *gin.Context) {
	lookupFieldIDStr := c.Param("lookupFieldId")
	lookupFieldID, err := uuid.Parse(lookupFieldIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lookup field ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	data, err := relationshipSvc.GetLookupFieldData(lookupFieldID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetFormHierarchy gets the hierarchy of related forms
func GetFormHierarchy(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	hierarchy, err := relationshipSvc.GetFormHierarchy(formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hierarchy)
}

// GetRelatedData gets related data for a record
func GetRelatedData(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	recordIDStr := c.Param("recordId")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	relatedData, err := relationshipSvc.GetRelatedData(formID, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, relatedData)
}
