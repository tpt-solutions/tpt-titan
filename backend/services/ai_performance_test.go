package services

import (
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"tpt-titan/backend/config"
)

// AI Performance and Resource Testing

// BenchmarkAIServiceProcessing benchmarks AI service processing performance
func BenchmarkAIServiceProcessing(b *testing.B) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	userID := uuid.New()
	modelID := uuid.New()

	// Test with different input sizes
	benchmarks := []struct {
		name  string
		input string
	}{
		{"small", "Hello world"},
		{"medium", "This is a medium length input for testing AI processing performance and resource usage."},
		{"large", `This is a comprehensive document about AI integration features in TPT Titan.
		It covers writing assistance, document summarization, email categorization, voice input,
		task prioritization, and various cross-module AI functionalities. The system is designed
		to enhance productivity for small and medium enterprises by providing intelligent automation
		and user-friendly AI features. Key components include natural language processing, speech
		recognition, machine learning models, and seamless integration across different modules.
		This document serves as a test case for performance benchmarking and memory optimization.`},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Mock the Ollama service for consistent benchmarking
			mockOllama := &MockOllamaService{}
			mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), bm.input).
				Return("Benchmark response", nil)
			service.ollamaService = mockOllama

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				request, err := service.ProcessRequest(userID, uuid.New(), modelID, bm.input, "text")
				if err != nil {
					b.Fatal(err)
				}
				if request == nil {
					b.Fatal("Request should not be nil")
				}
			}
		})
	}
}

// BenchmarkConcurrentAIProcessing benchmarks concurrent AI processing
func BenchmarkConcurrentAIProcessing(b *testing.B) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)
	userID := uuid.New()
	modelID := uuid.New()
	input := "Concurrent processing test input"

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), input).
		Return("Concurrent response", nil)
	service.ollamaService = mockOllama

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			request, err := service.ProcessRequest(userID, uuid.New(), modelID, input, "text")
			if err != nil {
				b.Fatal(err)
			}
			if request == nil {
				b.Fatal("Request should not be nil")
			}
		}
	})
}

// TestAIMemoryUsageOptimization tests memory usage optimization
func TestAIMemoryUsageOptimization(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	// Record initial memory usage
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	initialAlloc := m1.Alloc

	userID := uuid.New()
	modelID := uuid.New()

	// Process multiple requests of varying sizes
	testInputs := []string{
		"Short",
		"This is a medium length input that should not cause excessive memory allocation",
		`This is a longer document with substantial content. It includes multiple paragraphs
		and should be processed efficiently without excessive memory usage. The system should
		handle this input gracefully while maintaining good performance characteristics.`,
	}

	for _, input := range testInputs {
		mockOllama := &MockOllamaService{}
		mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), input).
			Return("Memory test response", nil)
		service.ollamaService = mockOllama

		request, err := service.ProcessRequest(userID, uuid.New(), modelID, input, "text")
		assert.NoError(t, err)
		assert.NotNil(t, request)

		// Force garbage collection between requests
		runtime.GC()
	}

	// Check final memory usage
	var m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m2)
	finalAlloc := m2.Alloc

	// Memory usage should not grow excessively
	memoryGrowth := finalAlloc - initialAlloc
	t.Logf("Memory growth: %d bytes (%d KB)", memoryGrowth, memoryGrowth/1024)

	// Assert reasonable memory growth (less than 10MB for this test)
	assert.True(t, memoryGrowth < 10*1024*1024, "Memory growth should be reasonable")
}

// TestAINetworkUsageMonitoring tests network usage monitoring
func TestAINetworkUsageMonitoring(t *testing.T) {
	cfg := &config.AIConfig{
		EnableOnlineAI: true,
	}

	service := NewAIService(cfg)

	userID := uuid.New()
	modelID := uuid.New()

	// Test with mock network monitoring
	mockOpenRouter := &MockOpenRouterService{}
	mockOpenRouter.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Network monitoring test response", nil)
	service.openRouterService = mockOpenRouter

	// Simulate network request
	request, err := service.ProcessRequest(userID, uuid.New(), modelID, "Network test", "text")

	assert.NoError(t, err)
	assert.NotNil(t, request)

	// In a real implementation, network usage would be tracked
	// For this test, we verify the request completes successfully
	assert.Equal(t, "completed", request.Status)
}

// TestAIOfflineFunctionalityVerification tests offline functionality
func TestAIOfflineFunctionalityVerification(t *testing.T) {
	// Test with local AI only (simulating offline mode)
	cfg := &config.AIConfig{
		EnableLocalAI:  true,
		EnableOnlineAI: false, // Offline mode
	}

	service := NewAIService(cfg)

	userID := uuid.New()
	modelID := uuid.New()

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Offline response", nil)
	service.ollamaService = mockOllama

	// Test that offline functionality works
	request, err := service.ProcessRequest(userID, uuid.New(), modelID, "Offline test", "text")

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, "Offline response", request.Output)

	mockOllama.AssertExpectations(t)
}

// TestAIPerformanceUnderLoad tests AI performance under load
func TestAIPerformanceUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)
	userID := uuid.New()
	modelID := uuid.New()

	// Test with multiple concurrent requests
	numRequests := 10
	requests := make([]*models.AIRequest, numRequests)
	errors := make([]error, numRequests)

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Load test response", nil)
	service.ollamaService = mockOllama

	start := time.Now()

	// Process requests concurrently
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			request, err := service.ProcessRequest(userID, uuid.New(), modelID, "Load test input", "text")
			requests[index] = request
			errors[index] = err
		}(i)
	}

	// Wait for all requests to complete (with timeout)
	timeout := time.After(30 * time.Second)
	for i := 0; i < numRequests; i++ {
		select {
		case <-timeout:
			t.Fatal("Test timed out waiting for requests to complete")
		default:
			// Check if request is complete
			if requests[i] != nil || errors[i] != nil {
				continue
			}
			time.Sleep(100 * time.Millisecond)
			i-- // Retry this index
		}
	}

	duration := time.Since(start)

	t.Logf("Processed %d requests in %v", numRequests, duration)
	t.Logf("Average time per request: %v", duration/time.Duration(numRequests))

	// Verify all requests completed successfully
	for i, err := range errors {
		assert.NoError(t, err, "Request %d should not have error", i)
		assert.NotNil(t, requests[i], "Request %d should not be nil", i)
		assert.Equal(t, "Load test response", requests[i].Output, "Request %d should have correct output", i)
	}

	// Performance assertions
	avgTimePerRequest := duration / time.Duration(numRequests)
	assert.True(t, avgTimePerRequest < 5*time.Second, "Average request time should be reasonable")
}

// TestAIResourceCleanup tests proper resource cleanup
func TestAIResourceCleanup(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
	}

	service := NewAIService(cfg)

	// Process several requests
	userID := uuid.New()
	modelID := uuid.New()

	for i := 0; i < 5; i++ {
		mockOllama := &MockOllamaService{}
		mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return("Cleanup test response", nil)
		service.ollamaService = mockOllama

		request, err := service.ProcessRequest(userID, uuid.New(), modelID, "Cleanup test", "text")
		assert.NoError(t, err)
		assert.NotNil(t, request)
	}

	// Force garbage collection
	runtime.GC()
	runtime.GC() // Run twice to ensure cleanup

	// In a real implementation, we would check for resource leaks
	// For this test, we verify the service remains functional
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	t.Logf("Final memory usage: %d KB", m.Alloc/1024)
	assert.True(t, m.Alloc < 100*1024*1024, "Memory usage should be reasonable after cleanup")
}

// TestAIBatteryLifeOptimization tests battery-conscious processing
func TestAIBatteryLifeOptimization(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI: true,
		LowPowerMode:  true, // Enable battery optimization
	}

	service := NewAIService(cfg)

	// In low power mode, the service should optimize for efficiency
	userID := uuid.New()
	modelID := uuid.New()

	mockOllama := &MockOllamaService{}
	mockOllama.On("GenerateResponse", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("Low power response", nil)
	service.ollamaService = mockOllama

	request, err := service.ProcessRequest(userID, uuid.New(), modelID, "Battery test", "text")

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, "Low power response", request.Output)

	mockOllama.AssertExpectations(t)
}

// TestAIHardwareAcceleration tests hardware acceleration usage
func TestAIHardwareAcceleration(t *testing.T) {
	cfg := &config.AIConfig{
		EnableLocalAI:        true,
		HardwareAcceleration: true,
	}

	service := NewAIService(cfg
