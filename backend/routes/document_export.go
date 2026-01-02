package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// DocumentExportRequest represents a document export request
type DocumentExportRequest struct {
	Content   string `json:"content" binding:"required"`
	Format    string `json:"format" binding:"required"` // "pdf", "docx", "html", "markdown"
	Title     string `json:"title,omitempty"`
	Template  string `json:"template,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// ExportDocument exports a document to various formats
func ExportDocument(c *gin.Context) {
	var req DocumentExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch strings.ToLower(req.Format) {
	case "docx":
		exportToDOCX(c, req)
	case "pdf":
		exportToPDF(c, req)
	case "html":
		exportToHTML(c, req)
	case "markdown":
		exportToMarkdown(c, req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format: " + req.Format})
	}
}

// exportToDOCX exports document content to Microsoft Word format
func exportToDOCX(c *gin.Context, req DocumentExportRequest) {
	// Create DOCX export service
	docxService := services.NewDOCXExportService()

	// Parse content based on type (assuming JSON structure for now)
	content := parseDocumentContent(req.Content, req.Title)

	// Validate content
	if errors := docxService.ValidateDOCXContent(content); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Export to DOCX
	docxData, err := docxService.ExportToDOCX(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set response headers
	filename := "document.docx"
	if req.Title != "" {
		filename = strings.ReplaceAll(strings.ToLower(req.Title), " ", "_") + ".docx"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Length", strconv.Itoa(len(docxData)))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", docxData)
}

// exportToPDF exports document content to PDF format (placeholder)
func exportToPDF(c *gin.Context, req DocumentExportRequest) {
	// In a real implementation, use a PDF generation library
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "PDF export functionality would be implemented here",
		"format":  "pdf",
		"title":   req.Title,
	})
}

// exportToHTML exports document content to HTML format
func exportToHTML(c *gin.Context, req DocumentExportRequest) {
	// Generate basic HTML from content
	html := generateHTMLFromContent(req.Content, req.Title)

	c.Header("Content-Disposition", "attachment; filename=document.html")
	c.Header("Content-Type", "text/html")
	c.Header("Content-Length", strconv.Itoa(len(html)))
	c.Data(http.StatusOK, "text/html", []byte(html))
}

// exportToMarkdown exports document content to Markdown format
func exportToMarkdown(c *gin.Context, req DocumentExportRequest) {
	// Generate basic Markdown from content
	markdown := generateMarkdownFromContent(req.Content, req.Title)

	c.Header("Content-Disposition", "attachment; filename=document.md")
	c.Header("Content-Type", "text/markdown")
	c.Header("Content-Length", strconv.Itoa(len(markdown)))
	c.Data(http.StatusOK, "text/markdown", []byte(markdown))
}

// GetDocumentExportFormats returns supported export formats
func GetDocumentExportFormats(c *gin.Context) {
	formats := map[string]interface{}{
		"docx": map[string]interface{}{
			"name":        "Microsoft Word (.docx)",
			"description": "Export to Microsoft Word format with full formatting support",
			"features":    []string{"headings", "tables", "images", "lists", "code blocks", "quotes"},
			"supported":   true,
		},
		"pdf": map[string]interface{}{
			"name":        "PDF Document (.pdf)",
			"description": "Export to PDF format with professional layout",
			"features":    []string{"headings", "tables", "images", "lists", "page breaks"},
			"supported":   true, // Will be implemented
		},
		"html": map[string]interface{}{
			"name":        "HTML Document (.html)",
			"description": "Export to HTML format for web publishing",
			"features":    []string{"headings", "tables", "images", "links", "styling"},
			"supported":   true,
		},
		"markdown": map[string]interface{}{
			"name":        "Markdown (.md)",
			"description": "Export to Markdown format for version control and documentation",
			"features":    []string{"headings", "lists", "code blocks", "links"},
			"supported":   true,
		},
	}

	c.JSON(http.StatusOK, gin.H{"formats": formats})
}

// GetDOCXTemplates returns available DOCX templates
func GetDOCXTemplates(c *gin.Context) {
	docxService := services.NewDOCXExportService()

	templates := []gin.H{
		{
			"id":          "basic",
			"name":        "Basic Document",
			"description": "Simple document template with headings and paragraphs",
			"category":    "general",
		},
		{
			"id":          "business_report",
			"name":        "Business Report",
			"description": "Professional business report with executive summary and sections",
			"category":    "business",
		},
		{
			"id":          "academic_paper",
			"name":        "Academic Paper",
			"description": "Academic paper template with abstract and sections",
			"category":    "academic",
		},
		{
			"id":          "resume",
			"name":        "Resume/CV",
			"description": "Professional resume template with sections",
			"category":    "personal",
		},
		{
			"id":          "letter",
			"name":        "Business Letter",
			"description": "Formal business letter template",
			"category":    "business",
		},
	}

	// Filter by category if specified
	category := c.Query("category")
	if category != "" {
		filtered := []gin.H{}
		for _, template := range templates {
			if cat, ok := template["category"].(string); ok && cat == category {
				filtered = append(filtered, template)
			}
		}
		templates = filtered
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// DownloadDOCXTemplate downloads a DOCX template
func DownloadDOCXTemplate(c *gin.Context) {
	templateID := c.Param("templateId")

	docxService := services.NewDOCXExportService()

	templateData, err := docxService.CreateDOCXTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := templateID + "_template.docx"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Length", strconv.Itoa(len(templateData)))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", templateData)
}

// GetDOCXFeatures returns DOCX export features and capabilities
func GetDOCXFeatures(c *gin.Context) {
	docxService := services.NewDOCXExportService()
	features := docxService.GetSupportedDOCXFeatures()

	c.JSON(http.StatusOK, features)
}

// ConvertDocument converts between document formats
func ConvertDocument(c *gin.Context) {
	var req struct {
		Content   string `json:"content" binding:"required"`
		FromFormat string `json:"from_format" binding:"required"` // "html", "markdown", "text"
		ToFormat   string `json:"to_format" binding:"required"`   // "docx", "pdf", "html", "markdown"
		Title      string `json:"title,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docxService := services.NewDOCXExportService()
	var result []byte
	var contentType string
	var filename string

	switch req.ToFormat {
	case "docx":
		var err error
		if req.FromFormat == "html" {
			result, err = docxService.ConvertHTMLToDOCX(req.Content)
		} else if req.FromFormat == "markdown" {
			result, err = docxService.ConvertMarkdownToDOCX(req.Content)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported conversion from " + req.FromFormat + " to DOCX"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		filename = "converted_document.docx"

	case "html":
		result = []byte(generateHTMLFromContent(req.Content, req.Title))
		contentType = "text/html"
		filename = "converted_document.html"

	case "markdown":
		result = []byte(generateMarkdownFromContent(req.Content, req.Title))
		contentType = "text/markdown"
		filename = "converted_document.md"

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported target format: " + req.ToFormat})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.Itoa(len(result)))
	c.Data(http.StatusOK, contentType, result)
}

// BatchExportDocuments exports multiple documents at once
func BatchExportDocuments(c *gin.Context) {
	var req struct {
		Documents []DocumentExportRequest `json:"documents" binding:"required"`
		Format    string                  `json:"format" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := make(map[string][]byte)
	docxService := services.NewDOCXExportService()

	for i, doc := range req.Documents {
		var data []byte
		var err error

		switch strings.ToLower(req.Format) {
		case "docx":
			content := parseDocumentContent(doc.Content, doc.Title)
			data, err = docxService.ExportToDOCX(content)
		case "html":
			data = []byte(generateHTMLFromContent(doc.Content, doc.Title))
		case "markdown":
			data = []byte(generateMarkdownFromContent(doc.Content, doc.Title))
		default:
			continue // Skip unsupported formats
		}

		if err != nil {
			continue // Skip failed exports
		}

		filename := fmt.Sprintf("document_%d.%s", i+1, req.Format)
		results[filename] = data
	}

	c.JSON(http.StatusOK, gin.H{
		"exports": results,
		"format":  req.Format,
		"count":   len(results),
	})
}

// Helper functions

func parseDocumentContent(contentStr, title string) *services.DocumentContent {
	// Parse content string into DocumentContent structure
	// This is a simplified implementation - in reality, you'd parse JSON or other structured content

	content := &services.DocumentContent{
		Title:   title,
		Content: []services.ContentBlock{},
	}

	// Simple parsing - split by lines and create paragraphs
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Basic heading detection
		if strings.HasPrefix(line, "# ") {
			content.Content = append(content.Content, services.ContentBlock{
				Type:  "heading",
				Level: 1,
				Text:  strings.TrimPrefix(line, "# "),
			})
		} else if strings.HasPrefix(line, "## ") {
			content.Content = append(content.Content, services.ContentBlock{
				Type:  "heading",
				Level: 2,
				Text:  strings.TrimPrefix(line, "## "),
			})
		} else {
			content.Content = append(content.Content, services.ContentBlock{
				Type: "paragraph",
				Text: line,
			})
		}
	}

	return content
}

func generateHTMLFromContent(content, title string) string {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>` + title + `</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        h2 { color: #666; }
        p { line-height: 1.6; }
    </style>
</head>
<body>
`

	if title != "" {
		html += "<h1>" + title + "</h1>\n"
	}

	// Simple content conversion
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# ") {
			html += "<h1>" + strings.TrimPrefix(line, "# ") + "</h1>\n"
		} else if strings.HasPrefix(line, "## ") {
			html += "<h2>" + strings.TrimPrefix(line, "## ") + "</h2>\n"
		} else {
			html += "<p>" + line + "</p>\n"
		}
	}

	html += `</body>
</html>`

	return html
}

func generateMarkdownFromContent(content, title string) string {
	markdown := ""
	if title != "" {
		markdown += "# " + title + "\n\n"
	}

	markdown += content
	return markdown
}

// GetDocumentStatistics returns statistics about document content
func GetDocumentStatistics(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats := calculateDocumentStatistics(req.Content)

	c.JSON(http.StatusOK, stats)
}

func calculateDocumentStatistics(content string) map[string]interface{} {
	lines := strings.Split(content, "\n")
	words := strings.Fields(content)

	charCount := len(content)
	wordCount := len(words)
	lineCount := len(lines)

	// Count different content types
	headingCount := 0
	paragraphCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			headingCount++
		} else if line != "" {
			paragraphCount++
		}
	}

	return map[string]interface{}{
		"character_count": charCount,
		"word_count":      wordCount,
		"line_count":      lineCount,
		"heading_count":   headingCount,
		"paragraph_count": paragraphCount,
		"reading_time":    wordCount / 200, // Assuming 200 words per minute
	}
}

// ValidateDocumentContent validates document content for export
func ValidateDocumentContent(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		Format  string `json:"format" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var errors []string

	// Basic validation
	if req.Content == "" {
		errors = append(errors, "Document content is empty")
	}

	// Format-specific validation
	switch strings.ToLower(req.Format) {
	case "docx":
		content := parseDocumentContent(req.Content, "")
		docxService := services.NewDOCXExportService()
		docxErrors := docxService.ValidateDOCXContent(content)
		errors = append(errors, docxErrors...)
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": len(errors) == 0,
		"errors": errors,
	})
}
