package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"tpt-titan/backend/models"
	"github.com/ollama/ollama/api"
)

// AIModelInfo represents information about an AI model
type AIModelInfo struct {
	Name         string   `json:"name"`
	Size         string   `json:"size"`
	ModifiedAt   string   `json:"modified_at"`
	Digest       string   `json:"digest"`
	Capabilities []string `json:"capabilities"`
}

// OllamaService handles Ollama API interactions
type OllamaService struct {
	baseURL string
	client  *http.Client
}

// NewOllamaService creates a new Ollama service
func NewOllamaService(host, port string) *OllamaService {
	return &OllamaService{
		baseURL: fmt.Sprintf("http://%s:%s", host, port),
		client: &http.Client{
			Timeout: 5 * time.Minute, // Long timeout for AI processing
		},
	}
}

// ListModels gets available models from Ollama
func (s *OllamaService) ListModels() ([]AIModelInfo, error) {
	resp, err := s.client.Get(s.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API error: %s", resp.Status)
	}

	var result struct {
		Models []struct {
			Name       string `json:"name"`
			Size       int64  `json:"size"`
			ModifiedAt string `json:"modified_at"`
			Digest     string `json:"digest"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	var models []AIModelInfo
	for _, model := range result.Models {
		models = append(models, AIModelInfo{
			Name:         model.Name,
			Size:         formatBytes(model.Size),
			ModifiedAt:   model.ModifiedAt,
			Digest:       model.Digest,
			Capabilities: inferCapabilities(model.Name),
		})
	}

	return models, nil
}

// PullModel downloads a model from Ollama
func (s *OllamaService) PullModel(modelName string) error {
	req := map[string]string{"name": modelName}
	jsonData, _ := json.Marshal(req)

	resp, err := s.client.Post(s.baseURL+"/api/pull", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to pull model %s: %s", modelName, string(body))
	}

	return nil
}

// GenerateResponse generates a response using Ollama
func (s *OllamaService) GenerateResponse(modelName, prompt string) (string, error) {
	stream := false
	req := api.GenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: &stream,
	}

	jsonData, _ := json.Marshal(req)

	resp, err := s.client.Post(s.baseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("generation failed: %s", string(body))
	}

	var result api.GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}

// AnalyzeDocument processes a document using vision models for OCR and analysis
func (s *OllamaService) AnalyzeDocument(modelName string, imageData []byte, analysisType string) (*models.DocumentAnalysisResult, error) {
	// Convert image data to base64 for Ollama API
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Create analysis prompt based on type
	prompt := s.createAnalysisPrompt(analysisType)

	// Create the request with image
	req := api.GenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		Images: []api.ImageData{{Data: base64Data}},
		Stream: &[]bool{false}[0],
		Options: map[string]interface{}{
			"temperature": 0.1, // Lower temperature for more consistent analysis
		},
	}

	jsonData, _ := json.Marshal(req)

	resp, err := s.client.Post(s.baseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to analyze document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("document analysis failed: %s", string(body))
	}

	var result api.GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode analysis response: %w", err)
	}

	// Parse the AI response into structured data
	return s.parseAnalysisResponse(result.Response, analysisType)
}

// createAnalysisPrompt creates the appropriate prompt based on analysis type
func (s *OllamaService) createAnalysisPrompt(analysisType string) string {
	switch analysisType {
	case "ocr":
		return `Extract all text from this image. Return only the text content, preserving formatting and structure as much as possible. If there are tables, represent them in a structured format.`
	case "invoice":
		return `Analyze this invoice image and extract the following fields as JSON:
- vendor_name: Company or person who issued the invoice
- invoice_number: Invoice number/ID
- date: Invoice date
- total_amount: Total amount due
- currency: Currency code (e.g., USD, EUR)
- line_items: Array of items with description, quantity, unit_price, line_total
Return only valid JSON.`
	case "receipt":
		return `Analyze this receipt image and extract the following fields as JSON:
- merchant_name: Business/store name
- date: Transaction date
- total_amount: Total paid
- currency: Currency code
- items: Array of purchased items (if visible)
- tax_amount: Tax amount (if separate)
Return only valid JSON.`
	case "business_card":
		return `Extract contact information from this business card image as JSON:
- name: Person's full name
- title: Job title/position
- company: Company/organization name
- phone: Phone number(s)
- email: Email address(es)
- address: Physical address
- website: Website URL
Return only valid JSON.`
	case "contract":
		return `Analyze this contract document and extract key information as JSON:
- parties: Array of parties involved (names/companies)
- effective_date: Contract start date
- expiration_date: Contract end date (if applicable)
- key_terms: Array of important clauses/terms
- total_value: Contract value (if mentioned)
- payment_terms: Payment schedule/terms
Return only valid JSON.`
	default:
		return `Analyze this document image and extract all visible text and structured information. Return the results in a clear, organized format.`
	}
}

// parseAnalysisResponse parses the AI response into structured data
func (s *OllamaService) parseAnalysisResponse(response string, analysisType string) (*models.DocumentAnalysisResult, error) {
	result := &models.DocumentAnalysisResult{
		TextContent: response,
		Confidence:  0.8, // Default confidence - could be improved with model confidence scores
		Pages:       1,
		Language:    "en",
	}

	// Try to parse structured data based on analysis type
	switch analysisType {
	case "invoice", "receipt", "business_card", "contract":
		if fields, err := s.parseStructuredFields(response, analysisType); err == nil {
			result.Fields = fields
		}
	}

	// Try to extract tables if response contains tabular data
	if tables, err := s.extractTables(response); err == nil && len(tables) > 0 {
		result.Tables = tables
	}

	return result, nil
}

// parseStructuredFields attempts to parse JSON-structured fields from AI response
func (s *OllamaService) parseStructuredFields(response string, analysisType string) ([]models.ExtractedField, error) {
	// Find JSON in the response (AI might add extra text around it)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")
	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var fields []models.ExtractedField
	for key, value := range data {
		field := models.ExtractedField{
			Name:       key,
			Confidence: 0.8, // Default confidence
		}

		// Convert value to string
		switch v := value.(type) {
		case string:
			field.Value = v
		case float64:
			field.Value = fmt.Sprintf("%.2f", v)
		case int:
			field.Value = fmt.Sprintf("%d", v)
		case []interface{}:
			// Handle arrays (like line items)
			if jsonBytes, err := json.Marshal(v); err == nil {
				field.Value = string(jsonBytes)
			}
		default:
			if jsonBytes, err := json.Marshal(v); err == nil {
				field.Value = string(jsonBytes)
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// extractTables attempts to extract table data from the response
func (s *OllamaService) extractTables(response string) ([]models.ExtractedTable, error) {
	// Simple table detection - look for patterns that suggest tabular data
	lines := strings.Split(response, "\n")
	var tables []models.ExtractedTable

	// Look for lines that might be table headers (multiple columns separated by | or tabs)
	for i, line := range lines {
		if strings.Contains(line, "|") || strings.Contains(line, "\t") {
			// Potential table found
			table := models.ExtractedTable{
				PageIndex:  0,
				Confidence: 0.7,
			}

			// Parse headers
			headers := s.parseTableRow(line)
			if len(headers) > 1 {
				table.Headers = headers

				// Try to parse following rows as data
				var rows [][]string
				for j := i + 1; j < len(lines) && j < i+10; j++ { // Limit to 10 rows
					rowLine := strings.TrimSpace(lines[j])
					if rowLine == "" {
						break
					}
					if row := s.parseTableRow(rowLine); len(row) == len(headers) {
						rows = append(rows, row)
					} else {
						break // Row doesn't match header count
					}
				}

				if len(rows) > 0 {
					table.Rows = rows
					tables = append(tables, table)
				}
			}
		}
	}

	return tables, nil
}

// parseTableRow parses a single table row
func (s *OllamaService) parseTableRow(line string) []string {
	// Try pipe-separated first
	if strings.Contains(line, "|") {
		parts := strings.Split(line, "|")
		var cleaned []string
		for _, part := range parts {
			cleaned = append(cleaned, strings.TrimSpace(part))
		}
		return cleaned
	}

	// Try tab-separated
	if strings.Contains(line, "\t") {
		return strings.Split(line, "\t")
	}

	// Try comma-separated as fallback
	if strings.Contains(line, ",") {
		return strings.Split(line, ",")
	}

	return []string{strings.TrimSpace(line)}
}

// Helper functions

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func inferCapabilities(modelName string) []string {
	name := strings.ToLower(modelName)
	capabilities := []string{}

	// Infer capabilities based on model name
	if strings.Contains(name, "vision") || strings.Contains(name, "llava") {
		capabilities = append(capabilities, "vision", "ocr")
	}

	if strings.Contains(name, "coder") || strings.Contains(name, "code") {
		capabilities = append(capabilities, "coding", "analysis")
	}

	if strings.Contains(name, "llama") || strings.Contains(name, "mistral") {
		capabilities = append(capabilities, "writing", "analysis", "general")
	}

	if strings.Contains(name, "phi") || strings.Contains(name, "gemma") {
		capabilities = append(capabilities, "writing", "tasks", "forms")
	}

	// Default capabilities if none inferred
	if len(capabilities) == 0 {
		capabilities = []string{"general"}
	}

	return capabilities
}
