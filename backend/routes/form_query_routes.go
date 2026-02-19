package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// BuildSQL builds SQL from visual query elements
func BuildSQL(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	sql, err := queryBuilder.BuildSQLFromElements(elements)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sql": sql,
		"elements": elements,
	})
}

// ExecuteVisualQuery executes a visual query and returns results
func ExecuteVisualQuery(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	result, err := queryBuilder.ExecuteQuery(elements)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAvailableTables returns tables available for query building
func GetAvailableTables(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	tables, err := queryBuilder.GetAvailableTables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}

// ValidateVisualQuery validates a visual query
func ValidateVisualQuery(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	errors := queryBuilder.ValidateQuery(elements)
	c.JSON(http.StatusOK, gin.H{
		"valid": len(errors) == 0,
		"errors": errors,
	})
}

// GetQuerySuggestions provides suggestions for query building
func GetQuerySuggestions(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	suggestions := queryBuilder.GetQuerySuggestions(elements)
	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

// SaveQueryTemplate saves a visual query as a template
func SaveQueryTemplate(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name        string                   `json:"name" binding:"required"`
		Description string                   `json:"description"`
		Elements    []services.QueryElement `json:"elements" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	err := queryBuilder.SaveQuery(userID.String(), req.Name, req.Description, req.Elements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Query template saved successfully"})
}

// GetQueryTemplates returns saved query templates
func GetQueryTemplates(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	templates, err := queryBuilder.LoadSavedQueries(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}
