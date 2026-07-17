package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

// HandwritingRecognitionService provides handwriting recognition for mathematical expressions
type HandwritingRecognitionService struct {
	apiKey     string
	apiURL     string
	modelName  string
	timeout    time.Duration
}

// HandwritingStroke represents a single stroke in handwriting
type HandwritingStroke struct {
	X []float64 `json:"x"` // X coordinates
	Y []float64 `json:"y"` // Y coordinates
	T []int64   `json:"t,omitempty"` // Timestamps (optional)
}

// RecognitionRequest represents a handwriting recognition request
type RecognitionRequest struct {
	Strokes      []HandwritingStroke `json:"strokes"`
	Width        int                 `json:"width,omitempty"`
	Height       int                 `json:"height,omitempty"`
	Language     string              `json:"language,omitempty"` // "en", "es", "fr", etc.
	MathMode     bool                `json:"math_mode"`          // True for mathematical expressions
	PreContext   string              `json:"pre_context,omitempty"`   // Previous text context
	PostContext  string              `json:"post_context,omitempty"`  // Following text context
}

// RecognitionResult represents the result of handwriting recognition
type RecognitionResult struct {
	Text         string              `json:"text"`                   // Recognized text
	Confidence   float64             `json:"confidence"`            // Recognition confidence (0-1)
	Alternatives []RecognitionAlternative `json:"alternatives,omitempty"` // Alternative interpretations
	MathExpr     *MathExpression     `json:"math_expr,omitempty"`   // Parsed mathematical expression
	Error        string              `json:"error,omitempty"`
	ProcessingTime int64             `json:"processing_time_ms"`    // Processing time in milliseconds
}

// RecognitionAlternative represents an alternative recognition result
type RecognitionAlternative struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

// MathExpression represents a parsed mathematical expression
type MathExpression struct {
	LaTeX       string                 `json:"latex"`                 // LaTeX representation
	MathML      string                 `json:"mathml"`                // MathML representation
	SVG         string                 `json:"svg,omitempty"`         // SVG representation
	Text        string                 `json:"text"`                  // Plain text representation
	Variables   []string               `json:"variables"`             // Variables used
	Functions   []string               `json:"functions"`             // Functions used
	Complexity  int                    `json:"complexity"`            // Expression complexity score
	Validated   bool                   `json:"validated"`             // Whether expression is mathematically valid
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EquationTemplate represents a saved equation template
type EquationTemplate struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Category    string              `json:"category"`    // "algebra", "calculus", "geometry", "physics", etc.
	LaTeX       string              `json:"latex"`
	MathML      string              `json:"mathml"`
	SVG         string              `json:"svg"`
	Variables   []string            `json:"variables"`
	Tags        []string            `json:"tags"`
	IsPublic    bool                `json:"is_public"`
	CreatedBy   uuid.UUID           `json:"created_by"`
	UsageCount  int                 `json:"usage_count"`
	CreatedAt   time.Time           `json:"created_at"`
}

// NewHandwritingRecognitionService creates a new handwriting recognition service
func NewHandwritingRecognitionService(apiKey, apiURL, modelName string) *HandwritingRecognitionService {
	return &HandwritingRecognitionService{
		apiKey:    apiKey,
		apiURL:    apiURL,
		modelName: modelName,
		timeout:   30 * time.Second, // 30 second timeout
	}
}

// RecognizeHandwriting recognizes handwriting from strokes
func (hrs *HandwritingRecognitionService) RecognizeHandwriting(strokes []HandwritingStroke, width, height int) (*RecognitionResult, error) {
	startTime := time.Now()

	// Create recognition request
	request := RecognitionRequest{
		Strokes:  strokes,
		Width:    width,
		Height:   height,
		MathMode: true, // Enable math mode for equations
		Language: "en",
	}

	// Call external recognition API (simulated - would call actual ML service)
	result, err := hrs.callRecognitionAPI(request)
	if err != nil {
		return nil, fmt.Errorf("recognition API call failed: %w", err)
	}

	result.ProcessingTime = time.Since(startTime).Milliseconds()

	// Post-process mathematical expressions
	if result.MathExpr == nil && result.Text != "" {
		mathExpr, err := hrs.ParseMathematicalExpression(result.Text)
		if err == nil {
			result.MathExpr = mathExpr
		}
	}

	return result, nil
}

// RecognizeEquationFromImage recognizes equations from image data
func (hrs *HandwritingRecognitionService) RecognizeEquationFromImage(imageData []byte, format string) (*RecognitionResult, error) {
	// Convert image to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Create API request for image recognition
	request := map[string]interface{}{
		"image": base64Image,
		"format": format,
		"math_mode": true,
	}

	// Call image recognition API (simulated)
	_, _ = json.Marshal(request) // serialise for future use

	// Simulate API call
	result := &RecognitionResult{
		Text:       "E = mc^{2}",
		Confidence: 0.95,
		ProcessingTime: 150,
	}

	// Parse the recognized equation
	mathExpr, err := hrs.ParseMathematicalExpression(result.Text)
	if err == nil {
		result.MathExpr = mathExpr
	}

	return result, nil
}

// callRecognitionAPI calls the external recognition API
func (hrs *HandwritingRecognitionService) callRecognitionAPI(request RecognitionRequest) (*RecognitionResult, error) {
	// In a real implementation, this would make an HTTP call to a machine learning service
	// For now, simulate recognition based on stroke patterns

	_, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Simulate API call delay
	time.Sleep(100 * time.Millisecond)

	// Analyze strokes to determine likely mathematical expression
	recognizedText := hrs.analyzeStrokes(request.Strokes)

	result := &RecognitionResult{
		Text:       recognizedText,
		Confidence: hrs.calculateConfidence(request.Strokes, recognizedText),
		Alternatives: []RecognitionAlternative{
			{Text: recognizedText, Confidence: 0.9},
			{Text: hrs.generateAlternative(recognizedText), Confidence: 0.7},
		},
	}

	return result, nil
}

// analyzeStrokes analyzes stroke patterns to recognize mathematical symbols
func (hrs *HandwritingRecognitionService) analyzeStrokes(strokes []HandwritingStroke) string {
	if len(strokes) == 0 {
		return ""
	}

	// Simple stroke analysis (in real implementation, this would use ML models)
	totalStrokes := len(strokes)

	// Analyze stroke patterns to recognize common mathematical symbols
	if totalStrokes == 1 {
		// Single stroke - could be 1, l, i, etc.
		return hrs.analyzeSingleStroke(strokes[0])
	} else if totalStrokes == 2 {
		// Two strokes - could be =, +, x, etc.
		return hrs.analyzeTwoStrokes(strokes[0], strokes[1])
	} else if totalStrokes >= 3 {
		// Multiple strokes - likely fractions, integrals, etc.
		return hrs.analyzeMultipleStrokes(strokes)
	}

	return "x" // Default fallback
}

// analyzeSingleStroke analyzes a single stroke
func (hrs *HandwritingRecognitionService) analyzeSingleStroke(stroke HandwritingStroke) string {
	if len(stroke.X) < 2 {
		return "."
	}

	// Calculate stroke characteristics
	startX, startY := stroke.X[0], stroke.Y[0]
	endX, endY := stroke.X[len(stroke.X)-1], stroke.Y[len(stroke.Y)-1]

	deltaX := math.Abs(endX - startX)
	deltaY := math.Abs(endY - startY)

	// Vertical stroke (1, l, i)
	if deltaY > deltaX*2 {
		return "1"
	}

	// Horizontal stroke (-)
	if deltaX > deltaY*2 {
		return "-"
	}

	// Diagonal stroke (/)
	if deltaX > deltaY && endY < startY {
		return "/"
	}

	// Curved stroke (possible parentheses or other symbols)
	return "("
}

// analyzeTwoStrokes analyzes two strokes
func (hrs *HandwritingRecognitionService) analyzeTwoStrokes(stroke1, stroke2 HandwritingStroke) string {
	// Check if strokes are parallel (likely =)
	if hrs.areStrokesParallel(stroke1, stroke2) {
		return "="
	}

	// Check for + symbol (cross)
	if hrs.formsCross(stroke1, stroke2) {
		return "+"
	}

	// Check for x symbol (cross with different orientation)
	return "x"
}

// analyzeMultipleStrokes analyzes multiple strokes
func (hrs *HandwritingRecognitionService) analyzeMultipleStrokes(strokes []HandwritingStroke) string {
	strokeCount := len(strokes)

	// Common mathematical expressions based on stroke count
	switch strokeCount {
	case 3:
		// Could be π (pi), or fractions like 1/2
		if hrs.hasHorizontalFractionBar(strokes) {
			return "1/2"
		}
		return "\\pi"
	case 4:
		// Could be integral symbol ∫
		return "\\int"
	case 5:
		// Could be summation Σ
		return "\\sum"
	default:
		// Complex expression - try to build it
		return hrs.buildComplexExpression(strokes)
	}
}

// Helper methods for stroke analysis
func (hrs *HandwritingRecognitionService) areStrokesParallel(s1, s2 HandwritingStroke) bool {
	// Simplified parallel check
	if len(s1.X) < 2 || len(s2.X) < 2 {
		return false
	}

	// Calculate slopes
	slope1 := hrs.calculateSlope(s1)
	slope2 := hrs.calculateSlope(s2)

	// Consider parallel if slopes are similar
	return math.Abs(slope1-slope2) < 0.2
}

func (hrs *HandwritingRecognitionService) formsCross(s1, s2 HandwritingStroke) bool {
	// Check if strokes intersect forming a cross
	return hrs.strokesIntersect(s1, s2)
}

func (hrs *HandwritingRecognitionService) hasHorizontalFractionBar(strokes []HandwritingStroke) bool {
	// Check if one stroke is horizontal and spans multiple other strokes
	for _, stroke := range strokes {
		if hrs.isHorizontalStroke(stroke) {
			// Check if other strokes are positioned above/below this line
			return true
		}
	}
	return false
}

func (hrs *HandwritingRecognitionService) calculateSlope(stroke HandwritingStroke) float64 {
	if len(stroke.X) < 2 {
		return 0
	}

	startIdx := 0
	endIdx := len(stroke.X) - 1

	deltaX := stroke.X[endIdx] - stroke.X[startIdx]
	deltaY := stroke.Y[endIdx] - stroke.Y[startIdx]

	if deltaX == 0 {
		return math.Inf(1) // Vertical line
	}

	return deltaY / deltaX
}

func (hrs *HandwritingRecognitionService) isHorizontalStroke(stroke HandwritingStroke) bool {
	if len(stroke.Y) < 2 {
		return false
	}

	minY := stroke.Y[0]
	maxY := stroke.Y[0]

	for _, y := range stroke.Y {
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	// Consider horizontal if vertical variation is small
	return (maxY - minY) < 10 // pixels
}

func (hrs *HandwritingRecognitionService) strokesIntersect(s1, s2 HandwritingStroke) bool {
	b1 := strokeBounds(s1)
	b2 := strokeBounds(s2)
	if b1 == nil || b2 == nil {
		return false
	}
	// Axis-aligned bounding-box overlap test.
	return b1.minX <= b2.maxX && b1.maxX >= b2.minX &&
		b1.minY <= b2.maxY && b1.maxY >= b2.minY
}

type bounds struct {
	minX, minY, maxX, maxY float64
}

func strokeBounds(s HandwritingStroke) *bounds {
	if len(s.X) == 0 || len(s.Y) == 0 {
		return nil
	}
	b := bounds{minX: s.X[0], minY: s.Y[0], maxX: s.X[0], maxY: s.Y[0]}
	for i := range s.X {
		if s.X[i] < b.minX {
			b.minX = s.X[i]
		}
		if s.X[i] > b.maxX {
			b.maxX = s.X[i]
		}
	}
	for i := range s.Y {
		if s.Y[i] < b.minY {
			b.minY = s.Y[i]
		}
		if s.Y[i] > b.maxY {
			b.maxY = s.Y[i]
		}
	}
	return &b
}

func (hrs *HandwritingRecognitionService) calculateConfidence(strokes []HandwritingStroke, recognizedText string) float64 {
	// Calculate confidence based on stroke characteristics
	baseConfidence := 0.8

	// Adjust based on stroke count and complexity
	strokeCount := len(strokes)
	if strokeCount > 10 {
		baseConfidence -= 0.1 // Complex expressions harder to recognize
	}

	// Adjust based on text length vs stroke count
	if len(recognizedText) > strokeCount*2 {
		baseConfidence -= 0.1 // Too many characters for strokes
	}

	if baseConfidence < 0.1 {
		baseConfidence = 0.1
	}

	return baseConfidence
}

func (hrs *HandwritingRecognitionService) generateAlternative(originalText string) string {
	// Generate alternative interpretation
	switch originalText {
	case "=":
		return "\\neq"
	case "+":
		return "-"
	case "x":
		return "\\times"
	default:
		return originalText + "'"
	}
}

func (hrs *HandwritingRecognitionService) buildComplexExpression(strokes []HandwritingStroke) string {
	// Derive a structural description of the drawing from its geometry. A full
	// recognizer would map this to LaTeX via an ML model; here we summarize the
	// strokes (count, total length, bounding box) so the output reflects the
	// actual input rather than a hardcoded constant.
	if len(strokes) == 0 {
		return ""
	}

	totalPoints := 0
	totalLen := 0.0
	all := bounds{}
	first := true
	for _, s := range strokes {
		b := strokeBounds(s)
		if b == nil {
			continue
		}
		totalPoints += len(s.X)
		for i := 1; i < len(s.X); i++ {
			dx := s.X[i] - s.X[i-1]
			dy := s.Y[i] - s.Y[i-1]
			totalLen += math.Sqrt(dx*dx + dy*dy)
		}
		if first {
			all = *b
			first = false
		} else {
			if b.minX < all.minX {
				all.minX = b.minX
			}
			if b.minY < all.minY {
				all.minY = b.minY
			}
			if b.maxX > all.maxX {
				all.maxX = b.maxX
			}
			if b.maxY > all.maxY {
				all.maxY = b.maxY
			}
		}
	}

	return fmt.Sprintf("drawing{strokes:%d,points:%d,length:%.1f,w:%.1f,h:%.1f}",
		len(strokes), totalPoints, totalLen, all.maxX-all.minX, all.maxY-all.minY)
}

// ParseMathematicalExpression parses a recognized mathematical expression
func (hrs *HandwritingRecognitionService) ParseMathematicalExpression(text string) (*MathExpression, error) {
	expr := &MathExpression{
		Text:       text,
		Variables:  hrs.extractVariables(text),
		Functions:  hrs.extractFunctions(text),
		Complexity: hrs.calculateComplexity(text),
		Validated:  true, // Assume valid for now
		Metadata:   make(map[string]interface{}),
	}

	// Generate LaTeX representation
	expr.LaTeX = hrs.convertToLaTeX(text)

	// Generate MathML representation
	expr.MathML = hrs.convertToMathML(text)

	// Generate SVG (placeholder)
	expr.SVG = hrs.generateSVG(text)

	return expr, nil
}

// extractVariables extracts variable names from expression
func (hrs *HandwritingRecognitionService) extractVariables(text string) []string {
	variables := []string{}

	// Simple variable extraction (single letters)
	for _, char := range text {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			varName := string(char)
			found := false
			for _, existing := range variables {
				if existing == varName {
					found = true
					break
				}
			}
			if !found {
				variables = append(variables, varName)
			}
		}
	}

	return variables
}

// extractFunctions extracts function names from expression
func (hrs *HandwritingRecognitionService) extractFunctions(text string) []string {
	functions := []string{}

	// Common mathematical functions
	commonFunctions := []string{"sin", "cos", "tan", "log", "ln", "sqrt", "exp", "abs"}

	for _, fn := range commonFunctions {
		if strings.Contains(text, fn) {
			functions = append(functions, fn)
		}
	}

	return functions
}

// calculateComplexity calculates expression complexity score
func (hrs *HandwritingRecognitionService) calculateComplexity(text string) int {
	complexity := 0

	// Count operators
	operators := []string{"+", "-", "*", "/", "=", "^", "(", ")", "[", "]", "{", "}"}
	for _, op := range operators {
		complexity += strings.Count(text, op)
	}

	// Count functions
	functions := hrs.extractFunctions(text)
	complexity += len(functions) * 5 // Functions are more complex

	// Count variables
	variables := hrs.extractVariables(text)
	complexity += len(variables) * 2

	return complexity
}

// convertToLaTeX converts mathematical expression to LaTeX
func (hrs *HandwritingRecognitionService) convertToLaTeX(text string) string {
	latex := text

	// Basic conversions
	conversions := map[string]string{
		"^":        "^{",
		"sqrt(":    "\\sqrt{",
		"sin(":     "\\sin{",
		"cos(":     "\\cos{",
		"tan(":     "\\tan{",
		"log(":     "\\log{",
		"ln(":      "\\ln{",
		"exp(":     "\\exp{",
		"pi":       "\\pi",
		"alpha":    "\\alpha",
		"beta":     "\\beta",
		"gamma":    "\\gamma",
		"delta":    "\\delta",
		"epsilon":  "\\epsilon",
		"theta":    "\\theta",
		"lambda":   "\\lambda",
		"mu":       "\\mu",
		"sigma":    "\\sigma",
		"omega":    "\\omega",
		"infty":    "\\infty",
		"sum":      "\\sum",
		"prod":     "\\prod",
		"int":      "\\int",
		"lim":      "\\lim",
		"frac":     "\\frac",
		"sqrt":     "\\sqrt",
	}

	for old, new := range conversions {
		latex = strings.ReplaceAll(latex, old, new)
	}

	// Handle fractions
	latex = hrs.convertFractionsToLaTeX(latex)

	// Handle superscripts
	latex = hrs.convertSuperscriptsToLaTeX(latex)

	// Handle subscripts
	latex = hrs.convertSubscriptsToLaTeX(latex)

	return latex
}

// convertFractionsToLaTeX converts fractions to LaTeX format
func (hrs *HandwritingRecognitionService) convertFractionsToLaTeX(text string) string {
	// Simple fraction conversion (e.g., a/b -> \frac{a}{b})
	if strings.Contains(text, "/") {
		parts := strings.Split(text, "/")
		if len(parts) == 2 {
			return fmt.Sprintf("\\frac{%s}{%s}", parts[0], parts[1])
		}
	}
	return text
}

// convertSuperscriptsToLaTeX converts superscripts to LaTeX format
func (hrs *HandwritingRecognitionService) convertSuperscriptsToLaTeX(text string) string {
	// Handle simple superscripts (e.g., x^2 -> x^{2})
	if strings.Contains(text, "^") {
		// Find the character after ^ and wrap it in braces
		parts := strings.Split(text, "^")
		if len(parts) == 2 {
			return fmt.Sprintf("%s^{%s}", parts[0], parts[1])
		}
	}
	return text
}

// convertSubscriptsToLaTeX converts subscripts to LaTeX format
func (hrs *HandwritingRecognitionService) convertSubscriptsToLaTeX(text string) string {
	// Handle simple subscripts (e.g., x_2 -> x_{2})
	if strings.Contains(text, "_") {
		parts := strings.Split(text, "_")
		if len(parts) == 2 {
			return fmt.Sprintf("%s_{%s}", parts[0], parts[1])
		}
	}
	return text
}

// convertToMathML converts mathematical expression to MathML
func (hrs *HandwritingRecognitionService) convertToMathML(text string) string {
	// Basic MathML structure
	mathml := fmt.Sprintf(`<math xmlns="http://www.w3.org/1998/Math/MathML">
  <mrow>
    <mi>%s</mi>
  </mrow>
</math>`, text)

	// For more complex expressions, this would need proper parsing
	// This is a simplified implementation
	return mathml
}

// generateSVG generates SVG representation of the mathematical expression
func (hrs *HandwritingRecognitionService) generateSVG(text string) string {
	// In a real implementation, this would use a math rendering library
	// For now, return a placeholder SVG
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="50">
  <text x="10" y="30" font-family="Times New Roman" font-size="16">%s</text>
</svg>`, text)

	return svg
}

// SaveEquationTemplate saves an equation as a template
func (hrs *HandwritingRecognitionService) SaveEquationTemplate(template *EquationTemplate) error {
	template.ID = uuid.New()
	template.CreatedAt = time.Now()
	template.UsageCount = 0

	// In a real implementation, save to database
	return nil
}

// GetEquationTemplates gets available equation templates
func (hrs *HandwritingRecognitionService) GetEquationTemplates(category string) ([]EquationTemplate, error) {
	// Return some predefined templates
	templates := []EquationTemplate{
		{
			ID:          uuid.New(),
			Name:        "Pythagorean Theorem",
			Description: "The square of the hypotenuse equals the sum of the squares of the other two sides",
			Category:    "geometry",
			LaTeX:       "a^{2} + b^{2} = c^{2}",
			Variables:   []string{"a", "b", "c"},
			Tags:        []string{"theorem", "triangle", "geometry"},
			IsPublic:    true,
			UsageCount:  150,
		},
		{
			ID:          uuid.New(),
			Name:        "Quadratic Formula",
			Description: "Solutions to the quadratic equation ax² + bx + c = 0",
			Category:    "algebra",
			LaTeX:       "x = \\frac{-b \\pm \\sqrt{b^{2} - 4ac}}{2a}",
			Variables:   []string{"a", "b", "c", "x"},
			Tags:        []string{"quadratic", "algebra", "equation"},
			IsPublic:    true,
			UsageCount:  200,
		},
		{
			ID:          uuid.New(),
			Name:        "Einstein's Mass-Energy Equivalence",
			Description: "Mass and energy are equivalent",
			Category:    "physics",
			LaTeX:       "E = mc^{2}",
			Variables:   []string{"E", "m", "c"},
			Tags:        []string{"physics", "relativity", "energy"},
			IsPublic:    true,
			UsageCount:  300,
		},
	}

	// Filter by category if specified
	if category != "" {
		filtered := []EquationTemplate{}
		for _, template := range templates {
			if template.Category == category {
				filtered = append(filtered, template)
			}
		}
		templates = filtered
	}

	return templates, nil
}

// GetEquationTemplateCategories returns available template categories
func (hrs *HandwritingRecognitionService) GetEquationTemplateCategories() []string {
	return []string{
		"algebra",
		"calculus",
		"geometry",
		"trigonometry",
		"physics",
		"chemistry",
		"statistics",
		"engineering",
	}
}

// SearchEquations searches for equations by text or LaTeX
func (hrs *HandwritingRecognitionService) SearchEquations(query string) ([]EquationTemplate, error) {
	// In a real implementation, search database
	templates, _ := hrs.GetEquationTemplates("")
	filtered := []EquationTemplate{}

	for _, template := range templates {
		if strings.Contains(strings.ToLower(template.Name), strings.ToLower(query)) ||
		   strings.Contains(strings.ToLower(template.Description), strings.ToLower(query)) ||
		   strings.Contains(template.LaTeX, query) {
			filtered = append(filtered, template)
		}
	}

	return filtered, nil
}

// ExportEquation exports an equation to different formats
func (hrs *HandwritingRecognitionService) ExportEquation(expr *MathExpression, format string) ([]byte, error) {
	switch format {
	case "latex":
		return []byte(expr.LaTeX), nil
	case "mathml":
		return []byte(expr.MathML), nil
	case "svg":
		return []byte(expr.SVG), nil
	case "png":
		// Render a real PNG of the expression text.
		data, _, err := RenderEquationImage(expr.LaTeX, "png", "medium", "")
		if err != nil {
			return nil, err
		}
		return data, nil
	case "pdf":
		// Render a real PDF of the expression text.
		data, _, err := RenderEquationImage(expr.LaTeX, "pdf", "medium", "")
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// ValidateMathematicalExpression validates if an expression is mathematically correct
func (hrs *HandwritingRecognitionService) ValidateMathematicalExpression(expr string) (bool, string) {
	// Basic validation (simplified)
	if expr == "" {
		return false, "Expression is empty"
	}

	// Check for balanced parentheses
	openCount := strings.Count(expr, "(")
	closeCount := strings.Count(expr, ")")
	if openCount != closeCount {
		return false, "Unbalanced parentheses"
	}

	// Check for basic mathematical syntax
	// This would need a full mathematical expression parser for complete validation

	return true, ""
}

// OptimizeExpression optimizes a mathematical expression
func (hrs *HandwritingRecognitionService) OptimizeExpression(expr string) (string, error) {
	// Basic optimizations (simplified)
	optimized := expr

	// Remove unnecessary spaces
	optimized = strings.ReplaceAll(optimized, " ", "")

	// Simplify basic expressions
	optimized = strings.ReplaceAll(optimized, "1*x", "x")
	optimized = strings.ReplaceAll(optimized, "x*1", "x")
	optimized = strings.ReplaceAll(optimized, "0+x", "x")
	optimized = strings.ReplaceAll(optimized, "x+0", "x")

	return optimized, nil
}
