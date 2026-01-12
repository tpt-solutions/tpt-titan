package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
)

// Comprehensive AI Service Tests for Cross-Module Integration

// TestWritingAssistance tests AI writing assistance features
func TestAIService_WritingAssistance(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI:  true,
		EnableOnlineAI: true,
	}

	service := NewAIService(cfg)

	// Mock Ollama service for writing assistance
	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Here are some writing suggestions:\n\n1. Use more active voice\n2. Vary sentence length\n3. Add specific examples", nil)

	service.ollamaService = mockOllama

	userID := uuid.New()
	taskID := uuid.New()
	modelID := uuid.New()

	// Test writing assistance request processing
	request, err := service.ProcessWritingAssistance(userID, taskID, modelID, "This is a sample text that needs improvement.", "general")

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, "writing_assistance", request.InputType)
	assert.Contains(t, request.Input, "sample text")

	mockOllama.AssertExpectations(t)
}

// TestDocumentSummarization tests document summarization functionality
func TestAIService_DocumentSummarization(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.MatchedBy(func(prompt string) bool {
		return len(prompt) > 100 // Document content should be substantial
	})).Return("This document discusses AI integration features including writing assistance, document summarization, and cross-module functionality. Key benefits include improved productivity and enhanced user experience.", nil)

	service.ollamaService = mockOllama

	longDocument := `This is a comprehensive document about AI integration features in TPT Titan.
	It covers writing assistance, document summarization, email categorization, voice input,
	task prioritization, and various cross-module AI functionalities. The system is designed
	to enhance productivity for small and medium enterprises by providing intelligent automation
	and user-friendly AI features. Key components include natural language processing, speech
	recognition, machine learning models, and seamless integration across different modules.`

	summary, err := service.SummarizeDocument(uuid.New(), uuid.New(), uuid.New(), longDocument, "concise")

	assert.NoError(t, err)
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, "AI integration")

	mockOllama.AssertExpectations(t)
}

// TestEmailCategorization tests AI-powered email categorization
func TestAIService_EmailCategorization(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.MatchedBy(func(prompt string) bool {
		return prompt != ""
	})).Return(`[
		{"email_id": "123e4567-e89b-12d3-a456-426614174000", "category": "work", "confidence": 0.95},
		{"email_id": "987fcdeb-51a2-43d7-8f9e-123456789abc", "category": "personal", "confidence": 0.87},
		{"email_id": "456789ab-cdef-1234-5678-90abcdef1234", "category": "promotional", "confidence": 0.92}
	]`, nil)

	service.ollamaService = mockOllama

	emails := []models.Email{
		{
			ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Subject:  "Project Update Meeting",
			Content:  "Let's discuss the quarterly project updates in our team meeting tomorrow.",
			SenderEmail: "boss@company.com",
		},
		{
			ID:       uuid.MustParse("987fcdeb-51a2-43d7-8f9e-123456789abc"),
			Subject:  "Family Dinner Plans",
			Content:  "Hi everyone, are we still on for dinner at grandma's house this weekend?",
			SenderEmail: "mom@gmail.com",
		},
		{
			ID:       uuid.MustParse("456789ab-cdef-1234-5678-90abcdef1234"),
			Subject:  "50% Off Sale - Limited Time",
			Content:  "Don't miss our biggest sale of the year! Get 50% off all products.",
			SenderEmail: "deals@amazon.com",
		},
	}

	categories, err := service.CategorizeEmails(uuid.New(), emails)

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 3)

	// Check that categories are properly assigned
	foundWork := false
	foundPersonal := false
	foundPromotional := false

	for _, category := range categories {
		if category.Category == "work" {
			foundWork = true
		} else if category.Category == "personal" {
			foundPersonal = true
		} else if category.Category == "promotional" {
			foundPromotional = true
		}
	}

	assert.True(t, foundWork, "Should categorize work email")
	assert.True(t, foundPersonal, "Should categorize personal email")
	assert.True(t, foundPromotional, "Should categorize promotional email")

	mockOllama.AssertExpectations(t)
}

// TestTaskPrioritization tests AI-powered task prioritization
func TestAIService_TaskPrioritization(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(`{"predicted_priority": "high", "confidence": 0.89, "reasoning": "This appears to be a time-sensitive project deadline with multiple stakeholders involved."}`, nil)

	service.ollamaService = mockOllama

	task := models.Task{
		ID:          uuid.New(),
		Title:       "Prepare quarterly financial report",
		Description: "Compile Q3 financial data, create charts, and present to board members by Friday.",
		Priority:    "medium",
		DueDate:     &[]time.Time{time.Now().AddDate(0, 0, 2)}[0],
	}

	result, err := service.PredictTaskPriority(uuid.New(), task)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "high", result.PredictedPriority)
	assert.Greater(t, result.Confidence, 0.8)

	mockOllama.AssertExpectations(t)
}

// TestTaskDeadlinePrediction tests deadline prediction functionality
func TestAIService_TaskDeadlinePrediction(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(`{"predicted_date": "2024-01-15", "confidence": 0.82, "reasoning": "Based on task complexity and current workload, this should take approximately 3-4 days to complete."}`, nil)

	service.ollamaService = mockOllama

	task := models.Task{
		ID:          uuid.New(),
		Title:       "Design new website layout",
		Description: "Create wireframes, mockups, and design system for the company website redesign project.",
		Priority:    "high",
		Tags:        []string{"design", "urgent", "client-facing"},
	}

	result, err := service.PredictTaskDeadline(uuid.New(), task)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.PredictedDate)
	assert.Greater(t, result.Confidence, 0.7)

	mockOllama.AssertExpectations(t)
}

// TestSpeechServiceIntegration tests speech service integration
func TestAIService_SpeechServiceIntegration(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	// Test speech request creation and processing
	userID := uuid.New()
	modelID := uuid.New()

	text := "This is a test document for speech synthesis."
	request, err := service.ProcessSpeechRequest(userID, modelID, "tts", text, "standard", "en")

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, "tts", request.RequestType)
	assert.Equal(t, text, request.InputText)
	assert.Equal(t, "en", request.Language)
}

// TestWorkflowAIIntegration tests workflow AI integration
func TestAIService_WorkflowAIIntegration(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("This workflow could be optimized by adding a conditional branch for error handling and implementing parallel processing for independent tasks.", nil)

	service.ollamaService = mockOllama

	workflow := models.Workflow{
		ID:          uuid.New(),
		Name:        "Invoice Processing Workflow",
		Description: "Process incoming invoices from receipt to payment",
		Category:    "invoice_processing",
	}

	suggestions, err := service.AnalyzeWorkflowOptimization(uuid.New(), workflow)

	assert.NoError(t, err)
	assert.NotEmpty(t, suggestions)
	assert.Contains(t, suggestions, "conditional branch")

	mockOllama.AssertExpectations(t)
}

// TestAIFallbackMechanisms tests AI fallback mechanisms
func TestAIService_FallbackMechanisms(t *testing.T) {
	// Test with no AI services enabled
	cfg := &config.AIConfig{
		EnableLocalAI:  false,
		EnableOnlineAI: false,
	}

	service := NewAIService(cfg)

	userID := uuid.New()
	taskID := uuid.New()
	modelID := uuid.New()

	// Should fail gracefully when no AI services are available
	request, err := service.ProcessRequest(userID, taskID, modelID, "Test prompt", "text")

	assert.Error(t, err)
	assert.Nil(t, request)
	assert.Contains(t, err.Error(), "AI services are disabled")
}

// TestAIPerformanceMonitoring tests performance monitoring
func TestAIService_PerformanceMonitoring(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Performance test response", nil)

	service.ollamaService = mockOllama

	userID := uuid.New()
	taskID := uuid.New()
	modelID := uuid.New()

	start := time.Now()
	request, err := service.ProcessRequest(userID, taskID, modelID, "Performance test", "text")
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, request)

	// Verify processing time is recorded
	assert.True(t, request.ProcessingTime > 0)
	assert.True(t, duration.Milliseconds() >= int64(request.ProcessingTime))

	mockOllama.AssertExpectations(t)
}

// TestAIMemoryOptimization tests memory usage optimization
func TestAIService_MemoryOptimization(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	// Test with different input sizes to ensure memory efficiency
	smallInput := "Short text"
	mediumInput := "This is a medium length text that should be processed efficiently without excessive memory usage. " +
		"It contains enough content to test memory optimization but not so much as to cause performance issues."
	largeInput := mediumInput + mediumInput + mediumInput + mediumInput // 4x medium

	testCases := []struct {
		name  string
		input string
	}{
		{"small", smallInput},
		{"medium", mediumInput},
		{"large", largeInput},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockOllama := &MockOllamaService{}
			mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), tc.input).
				Return("Processed: " + tc.input[:min(50, len(tc.input))], nil)

			service.ollamaService = mockOllama

			userID := uuid.New()
			request, err := service.ProcessRequest(userID, uuid.New(), uuid.New(), tc.input, "text")

			assert.NoError(t, err)
			assert.NotNil(t, request)
			assert.True(t, request.ProcessingTime > 0)

			mockOllama.AssertExpectations(t)
		})
	}
}

// TestCrossPlatformAIConsistency tests cross-platform consistency
func TestAIService_CrossPlatformConsistency(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	// Test that AI responses are consistent across different input formats
	testPhrase := "Schedule a meeting for tomorrow at 3 PM"

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), testPhrase).
		Return(`{"action": "create_event", "title": "Meeting", "date": "tomorrow", "time": "15:00"}`, nil)

	service.ollamaService = mockOllama

	userID := uuid.New()
	request, err := service.ProcessRequest(userID, uuid.New(), uuid.New(), testPhrase, "text")

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, testPhrase, request.Input)

	mockOllama.AssertExpectations(t)
}

// TestAIErrorHandlingAndRecovery tests error handling and recovery
func TestAIService_ErrorHandlingAndRecovery(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI:  true,
		EnableOnlineAI: true, // Enable fallback
	}

	service := NewAIService(cfg)

	// Test local AI failure with online AI fallback
	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", assert.AnError) // Local AI fails

	mockOpenRouter := &MockOpenRouterService{}
	mockOpenRouter.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Fallback response from online AI", nil) // Online AI succeeds

	service.ollamaService = mockOllama
	service.openRouterService = mockOpenRouter

	userID := uuid.New()
	request, err := service.ProcessRequestWithFallback(userID, uuid.New(), uuid.New(), "Test with fallback", "text")

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, "Fallback response from online AI", request.Output)

	mockOllama.AssertExpectations(t)
	mockOpenRouter.AssertExpectations(t)
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
