package routes

import (
	"database/sql"
	"net/http"
	"time"
	"tpt-titan-simple/backend/models"
	"tpt-titan-simple/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var calendarService *services.CalendarService

// InitCalendarService initializes the calendar service (called from main)
func InitCalendarService(db *sql.DB) {
	calendarService = services.NewCalendarService(db)
}

// GetCalendars returns all calendars for the authenticated user
func GetCalendars(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	calendars, err := calendarService.GetCalendars(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve calendars"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"calendars": calendars})
}

// GetCalendar returns a specific calendar
func GetCalendar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	calendarID := c.Param("id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
		return
	}

	id, err := uuid.Parse(calendarID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calendar ID"})
		return
	}

	calendar, err := calendarService.GetCalendar(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve calendar"})
		return
	}

	if calendar == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"calendar": calendar})
}

// CreateCalendar creates a new calendar
func CreateCalendar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calendar, err := calendarService.CreateCalendar(userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create calendar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"calendar": calendar})
}

// UpdateCalendar updates an existing calendar
func UpdateCalendar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	calendarID := c.Param("id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
		return
	}

	id, err := uuid.Parse(calendarID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calendar ID"})
		return
	}

	var req models.CalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calendar, err := calendarService.UpdateCalendar(userID.(uuid.UUID), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update calendar"})
		return
	}

	if calendar == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"calendar": calendar})
}

// DeleteCalendar deletes a calendar
func DeleteCalendar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	calendarID := c.Param("id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
		return
	}

	id, err := uuid.Parse(calendarID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calendar ID"})
		return
	}

	err = calendarService.DeleteCalendar(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete calendar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calendar deleted successfully"})
}

// GetEvents returns events within a date range
func GetEvents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end date parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use RFC3339"})
		return
	}

	events, err := calendarService.GetEvents(userID.(uuid.UUID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetEvent returns a specific event
func GetEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	id, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := calendarService.GetEvent(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// CreateEvent creates a new event
func CreateEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := calendarService.CreateEvent(userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}

// UpdateEvent updates an existing event
func UpdateEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	id, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req models.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := calendarService.UpdateEvent(userID.(uuid.UUID), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// DeleteEvent deletes an event
func DeleteEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	id, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = calendarService.DeleteEvent(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// GetCalendarEvents returns all events for a specific calendar
func GetCalendarEvents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	calendarID := c.Param("id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
		return
	}

	id, err := uuid.Parse(calendarID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calendar ID"})
		return
	}

	// Parse query parameters for date range
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")

	startDate := time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	endDate := time.Now().AddDate(0, 1, 0)    // Default to 1 month from now

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use RFC3339"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use RFC3339"})
			return
		}
	}

	// Get all events and filter by calendar
	allEvents, err := calendarService.GetEvents(userID.(uuid.UUID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	// Filter events by calendar
	var calendarEvents []models.EventResponse
	for _, event := range allEvents {
		if event.CalendarID == id {
			calendarEvents = append(calendarEvents, event)
		}
	}

	c.JSON(http.StatusOK, gin.H{"events": calendarEvents})
}
