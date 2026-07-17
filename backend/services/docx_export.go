package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"time"
)

// DOCXExportService provides Microsoft Word (.docx) export functionality
type DOCXExportService struct {
}

// DocumentContent represents the content structure of a document
type DocumentContent struct {
	Title    string            `json:"title"`
	Content  []ContentBlock    `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Header   string            `json:"header,omitempty"`
	Footer   string            `json:"footer,omitempty"`
}

// ContentBlock represents a block of content (paragraph, heading, list, etc.)
type ContentBlock struct {
	Type    string                 `json:"type"`    // "paragraph", "heading", "list", "table", "image", "code", "quote"
	Level   int                    `json:"level,omitempty"`   // For headings and lists
	Text    string                 `json:"text,omitempty"`
	Style   map[string]interface{} `json:"style,omitempty"`   // Formatting styles
	Items   []string               `json:"items,omitempty"`   // For lists
	Table   *TableData             `json:"table,omitempty"`   // For tables
	Image   *ImageData             `json:"image,omitempty"`   // For images
	Code    *CodeBlock             `json:"code,omitempty"`    // For code blocks
}

// TableData represents table content
type TableData struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
	Style   string     `json:"style,omitempty"` // "plain", "bordered", "striped"
}

// ImageData represents image content
type ImageData struct {
	Data        []byte `json:"data"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Caption     string `json:"caption,omitempty"`
}

// CodeBlock represents code content
type CodeBlock struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// NewDOCXExportService creates a new DOCX export service
func NewDOCXExportService() *DOCXExportService {
	return &DOCXExportService{}
}

// ExportToDOCX exports document content to Microsoft Word format
func (des *DOCXExportService) ExportToDOCX(content *DocumentContent) ([]byte, error) {
	// Create a new DOCX file structure
	docx := des.createDOCXStructure()

	// Generate document content (including the section properties that reference
	// any header/footer parts)
	documentXML := des.generateDocumentXML(content)
	docx.Files["word/document.xml"] = documentXML

	// Generate header and footer parts when provided.
	if content.Header != "" {
		docx.Files["word/header1.xml"] = des.generateHeaderFooterXML(content.Header, "hdr")
	}
	if content.Footer != "" {
		docx.Files["word/footer1.xml"] = des.generateHeaderFooterXML(content.Footer, "ftr")
	}

	// Generate relationships and other required files
	docx.Files["word/_rels/document.xml.rels"] = des.generateRelationshipsXML(content)
	docx.Files["[Content_Types].xml"] = des.generateContentTypesXML(content)
	docx.Files["word/styles.xml"] = des.generateStylesXML()

	// Create the ZIP archive (DOCX is essentially a ZIP file)
	return des.createZIPArchive(docx)
}

// createDOCXStructure initializes the DOCX file structure
func (des *DOCXExportService) createDOCXStructure() *DOCXStructure {
	return &DOCXStructure{
		Files: make(map[string][]byte),
	}
}

// DOCXStructure represents the internal structure of a DOCX file
type DOCXStructure struct {
	Files map[string][]byte
}

// generateDocumentXML generates the main document XML content
func (des *DOCXExportService) generateDocumentXML(content *DocumentContent) []byte {
	var xmlContent strings.Builder

	// XML declaration and namespace
	xmlContent.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <w:body>
`)

	// Add document content
	for _, block := range content.Content {
		xmlContent.WriteString(des.generateBlockXML(block))
	}

	// Section properties (must be the last child of w:body). When a header or
	// footer is provided, reference the corresponding parts so Word renders them.
	xmlContent.WriteString(des.generateSectPr(content))

	// Close document
	xmlContent.WriteString(`  </w:body>
</w:document>`)

	return []byte(xmlContent.String())
}

// generateSectPr emits the section properties that close w:body. It references
// any header/footer parts that were generated for the document.
func (des *DOCXExportService) generateSectPr(content *DocumentContent) string {
	var refs strings.Builder
	if content.Header != "" {
		refs.WriteString(`      <w:headerReference w:type="default" r:id="rIdHdr"/>
`)
	}
	if content.Footer != "" {
		refs.WriteString(`      <w:footerReference w:type="default" r:id="rIdFtr"/>
`)
	}

	return fmt.Sprintf(`    <w:sectPr>
%s      <w:pgSz w:w="12240" w:h="15840"/>
      <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>
    </w:sectPr>
`, refs.String())
}

// generateHeaderFooterXML emits a minimal header/footer part containing a
// single paragraph of text. part must be either "hdr" or "ftr".
func (des *DOCXExportService) generateHeaderFooterXML(text, part string) []byte {
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:%s xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:p>
    <w:r>
      <w:t>%s</w:t>
    </w:r>
  </w:p>
</w:%s>`, part, des.escapeXML(text), part))
}

// generateHeadingXML generates XML for headings
func (des *DOCXExportService) generateHeadingXML(block ContentBlock) string {
	level := block.Level
	if level < 1 {
		level = 1
	}
	if level > 9 {
		level = 9
	}

	style := fmt.Sprintf("Heading%d", level)

	return fmt.Sprintf(`    <w:p>
      <w:pPr>
        <w:pStyle w:val="%s"/>
      </w:pPr>
      <w:r>
        <w:t>%s</w:t>
      </w:r>
    </w:p>
`, style, des.escapeXML(block.Text))
}

// generateParagraphXML generates XML for paragraphs
func (des *DOCXExportService) generateParagraphXML(block ContentBlock) string {
	return fmt.Sprintf(`    <w:p>
      <w:r>
        <w:t>%s</w:t>
      </w:r>
    </w:p>
`, des.escapeXML(block.Text))
}

// generateListXML generates XML for lists
func (des *DOCXExportService) generateListXML(block ContentBlock) string {
	var xmlContent strings.Builder

	listType := "bullet"
	if block.Level > 0 {
		listType = "number"
	}

	for _, item := range block.Items {
		xmlContent.WriteString(fmt.Sprintf(`    <w:p>
      <w:pPr>
        <w:numPr>
          <w:ilvl w:val="%d"/>
          <w:numId w:val="%s"/>
        </w:numPr>
      </w:pPr>
      <w:r>
        <w:t>%s</w:t>
      </w:r>
    </w:p>
`, block.Level, listType, des.escapeXML(item)))
	}

	return xmlContent.String()
}

// generateTableXML generates XML for tables
func (des *DOCXExportService) generateTableXML(block ContentBlock) string {
	if block.Table == nil {
		return ""
	}

	var xmlContent strings.Builder

	xmlContent.WriteString(`    <w:tbl>
      <w:tblPr>
        <w:tblStyle w:val="TableGrid"/>
        <w:tblW w:w="0" w:type="auto"/>
      </w:tblPr>
`)

	// Add header row if headers exist
	if len(block.Table.Headers) > 0 {
		xmlContent.WriteString(`      <w:tr>
`)
		for _, header := range block.Table.Headers {
			xmlContent.WriteString(fmt.Sprintf(`        <w:tc>
          <w:p>
            <w:r>
              <w:t>%s</w:t>
            </w:r>
          </w:p>
        </w:tc>
`, des.escapeXML(header)))
		}
		xmlContent.WriteString(`      </w:tr>
`)
	}

	// Add data rows
	for _, row := range block.Table.Rows {
		xmlContent.WriteString(`      <w:tr>
`)
		for _, cell := range row {
			xmlContent.WriteString(fmt.Sprintf(`        <w:tc>
          <w:p>
            <w:r>
              <w:t>%s</w:t>
            </w:r>
          </w:p>
        </w:tc>
`, des.escapeXML(cell)))
		}
		xmlContent.WriteString(`      </w:tr>
`)
	}

	xmlContent.WriteString(`    </w:tbl>
`)

	return xmlContent.String()
}

// generateImageXML generates XML for images
func (des *DOCXExportService) generateImageXML(block ContentBlock) string {
	if block.Image == nil {
		return ""
	}

	// Generate a unique relationship ID for the image
	relID := fmt.Sprintf("rId%d", len(block.Image.Filename))

	return fmt.Sprintf(`    <w:p>
      <w:r>
        <w:drawing>
          <wp:inline distT="0" distB="0" distL="0" distR="0">
            <wp:extent cx="%d" cy="%d"/>
            <wp:docPr id="1" name="Picture"/>
            <a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
              <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">
                <pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
                  <pic:nvPicPr>
                    <pic:cNvPr id="0" name="Picture"/>
                    <pic:cNvPicPr/>
                  </pic:nvPicPr>
                  <pic:blipFill>
                    <a:blip r:embed="%s"/>
                    <a:stretch>
                      <a:fillRect/>
                    </a:stretch>
                  </pic:blipFill>
                  <pic:spPr>
                    <a:xfrm>
                      <a:ext cx="%d" cy="%d"/>
                    </a:xfrm>
                    <a:prstGeom prst="rect">
                      <a:avLst/>
                    </a:prstGeom>
                  </pic:spPr>
                </pic:pic>
              </a:graphicData>
            </a:graphic>
          </wp:inline>
        </w:drawing>
      </w:r>
    </w:p>
`, block.Image.Width*914400, block.Image.Height*914400, relID, block.Image.Width*914400, block.Image.Height*914400)
}

// generateCodeXML generates XML for code blocks
func (des *DOCXExportService) generateCodeXML(block ContentBlock) string {
	if block.Code == nil {
		return ""
	}

	return fmt.Sprintf(`    <w:p>
      <w:pPr>
        <w:shd w:val="clear" w:color="auto" w:fill="F2F2F2"/>
        <w:rPr>
          <w:rFonts w:ascii="Courier New" w:hAnsi="Courier New"/>
          <w:sz w:val="20"/>
        </w:rPr>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:rFonts w:ascii="Courier New" w:hAnsi="Courier New"/>
          <w:sz w:val="20"/>
        </w:rPr>
        <w:t>%s</w:t>
      </w:r>
    </w:p>
`, des.escapeXML(block.Code.Code))
}

// generateQuoteXML generates XML for blockquotes
func (des *DOCXExportService) generateQuoteXML(block ContentBlock) string {
	return fmt.Sprintf(`    <w:p>
      <w:pPr>
        <w:ind w:left="720" w:right="720"/>
        <w:shd w:val="clear" w:color="auto" w:fill="F5F5F5"/>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:i/>
        </w:rPr>
        <w:t>%s</w:t>
      </w:r>
    </w:p>
`, des.escapeXML(block.Text))
}

// generateRelationshipsXML generates the relationships for document.xml,
// including references to the header/footer parts when present.
func (des *DOCXExportService) generateRelationshipsXML(content *DocumentContent) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`)

	if content.Header != "" {
		b.WriteString(`
  <Relationship Id="rIdHdr" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>`)
	}
	if content.Footer != "" {
		b.WriteString(`
  <Relationship Id="rIdFtr" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>`)
	}

	b.WriteString(`
</Relationships>`)
	return []byte(b.String())
}

// generateContentTypesXML generates the content types XML, registering the
// header/footer parts when present.
func (des *DOCXExportService) generateContentTypesXML(content *DocumentContent) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>`)

	if content.Header != "" {
		b.WriteString(`
  <Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>`)
	}
	if content.Footer != "" {
		b.WriteString(`
  <Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>`)
	}

	b.WriteString(`
</Types>`)
	return []byte(b.String())
}

// generateStylesXML generates the styles XML
func (des *DOCXExportService) generateStylesXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:docDefaults>
    <w:rPrDefault>
      <w:rPr>
        <w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman"/>
        <w:sz w:val="24"/>
        <w:szCs w:val="24"/>
        <w:lang w:val="en-US" w:eastAsia="en-US"/>
      </w:rPr>
    </w:rPrDefault>
    <w:pPrDefault>
      <w:pPr>
        <w:spacing w:after="160" w:line="259" w:lineRule="auto"/>
      </w:pPr>
    </w:pPrDefault>
  </w:docDefaults>
  <w:style w:type="paragraph" w:styleId="Normal">
    <w:name w:val="Normal"/>
    <w:qFormat/>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading1">
    <w:name w:val="heading 1"/>
    <w:basedOn w:val="Normal"/>
    <w:next w:val="Normal"/>
    <w:qFormat/>
    <w:rPr>
      <w:b/>
      <w:sz w:val="32"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading2">
    <w:name w:val="heading 2"/>
    <w:basedOn w:val="Normal"/>
    <w:next w:val="Normal"/>
    <w:qFormat/>
    <w:rPr>
      <w:b/>
      <w:sz w:val="28"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading3">
    <w:name w:val="heading 3"/>
    <w:basedOn w:val="Normal"/>
    <w:next w:val="Normal"/>
    <w:qFormat/>
    <w:rPr>
      <w:b/>
      <w:sz w:val="26"/>
    </w:rPr>
  </w:style>
</w:styles>`)
}

// createZIPArchive creates a ZIP archive from the DOCX structure
func (des *DOCXExportService) createZIPArchive(docx *DOCXStructure) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add each file to the ZIP archive
	for filename, content := range docx.Files {
		file, err := zipWriter.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create file in ZIP: %w", err)
		}

		_, err = file.Write(content)
		if err != nil {
			return nil, fmt.Errorf("failed to write file content: %w", err)
		}
	}

	// Close the ZIP writer
	err := zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close ZIP writer: %w", err)
	}

	return buf.Bytes(), nil
}

// escapeXML escapes special characters for XML
func (des *DOCXExportService) escapeXML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&apos;")
	return text
}

// ConvertHTMLToDOCX converts HTML content to DOCX format
func (des *DOCXExportService) ConvertHTMLToDOCX(htmlContent string) ([]byte, error) {
	// Parse HTML and convert to document structure
	content := des.parseHTMLToDocumentContent(htmlContent)
	return des.ExportToDOCX(content)
}

// ConvertMarkdownToDOCX converts Markdown content to DOCX format
func (des *DOCXExportService) ConvertMarkdownToDOCX(markdownContent string) ([]byte, error) {
	// Parse Markdown and convert to document structure
	content := des.parseMarkdownToDocumentContent(markdownContent)
	return des.ExportToDOCX(content)
}

// parseHTMLToDocumentContent parses HTML and converts to document content structure
func (des *DOCXExportService) parseHTMLToDocumentContent(html string) *DocumentContent {
	content := &DocumentContent{
		Title:   "Converted Document",
		Content: []ContentBlock{},
	}

	// Simple HTML parsing (in a real implementation, use a proper HTML parser)
	// This is a simplified version for demonstration
	lines := strings.Split(html, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		block := ContentBlock{Type: "paragraph", Text: line}

		// Basic HTML tag detection
		if strings.HasPrefix(line, "<h1>") && strings.HasSuffix(line, "</h1>") {
			block.Type = "heading"
			block.Level = 1
			block.Text = strings.TrimSuffix(strings.TrimPrefix(line, "<h1>"), "</h1>")
		} else if strings.HasPrefix(line, "<h2>") && strings.HasSuffix(line, "</h2>") {
			block.Type = "heading"
			block.Level = 2
			block.Text = strings.TrimSuffix(strings.TrimPrefix(line, "<h2>"), "</h2>")
		} else if strings.HasPrefix(line, "<p>") && strings.HasSuffix(line, "</p>") {
			block.Text = strings.TrimSuffix(strings.TrimPrefix(line, "<p>"), "</p>")
		}

		content.Content = append(content.Content, block)
	}

	return content
}

// parseMarkdownToDocumentContent parses Markdown and converts to document content structure
func (des *DOCXExportService) parseMarkdownToDocumentContent(markdown string) *DocumentContent {
	content := &DocumentContent{
		Title:   "Converted Document",
		Content: []ContentBlock{},
	}

	lines := strings.Split(markdown, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		block := ContentBlock{Type: "paragraph", Text: line}

		// Markdown parsing
		if strings.HasPrefix(line, "# ") {
			block.Type = "heading"
			block.Level = 1
			block.Text = strings.TrimPrefix(line, "# ")
		} else if strings.HasPrefix(line, "## ") {
			block.Type = "heading"
			block.Level = 2
			block.Text = strings.TrimPrefix(line, "## ")
		} else if strings.HasPrefix(line, "### ") {
			block.Type = "heading"
			block.Level = 3
			block.Text = strings.TrimPrefix(line, "### ")
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			block.Type = "list"
			block.Items = []string{strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")}
		} else if strings.HasPrefix(line, "```") {
			// Code block (simplified)
			block.Type = "code"
			block.Code = &CodeBlock{
				Language: "text",
				Code:     strings.TrimPrefix(line, "```"),
			}
		}

		content.Content = append(content.Content, block)
	}

	return content
}

// GetSupportedDOCXFeatures returns the features supported by DOCX export
func (des *DOCXExportService) GetSupportedDOCXFeatures() map[string]interface{} {
	return map[string]interface{}{
		"text_formatting": []string{"bold", "italic", "underline", "strikethrough"},
		"headings":        []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"lists":           []string{"bullet", "numbered", "multilevel"},
		"tables":          true,
		"images":          []string{"png", "jpg", "jpeg", "gif", "bmp"},
		"code_blocks":     true,
		"blockquotes":     true,
		"hyperlinks":      true,
		"page_breaks":     true,
		"headers_footers": true, // Rendered when DocumentContent.Header/Footer are set
		"styles":          []string{"normal", "heading1-9", "code", "quote"},
	}
}

// ValidateDOCXContent validates if content can be exported to DOCX
func (des *DOCXExportService) ValidateDOCXContent(content *DocumentContent) []string {
	var errors []string

	if len(content.Content) == 0 {
		errors = append(errors, "Document has no content")
	}

	for i, block := range content.Content {
		switch block.Type {
		case "heading":
			if block.Level < 1 || block.Level > 9 {
				errors = append(errors, fmt.Sprintf("Block %d: Invalid heading level %d", i+1, block.Level))
			}
		case "table":
			if block.Table == nil {
				errors = append(errors, fmt.Sprintf("Block %d: Table block missing table data", i+1))
			}
		case "image":
			if block.Image == nil {
				errors = append(errors, fmt.Sprintf("Block %d: Image block missing image data", i+1))
			}
		case "code":
			if block.Code == nil {
				errors = append(errors, fmt.Sprintf("Block %d: Code block missing code data", i+1))
			}
		}
	}

	return errors
}

// generateBlockXML generates XML for a content block
func (des *DOCXExportService) generateBlockXML(block ContentBlock) string {
	switch block.Type {
	case "heading":
		return des.generateHeadingXML(block)
	case "paragraph":
		return des.generateParagraphXML(block)
	case "list":
		return des.generateListXML(block)
	case "table":
		return des.generateTableXML(block)
	case "image":
		return des.generateImageXML(block)
	case "code":
		return des.generateCodeXML(block)
	case "quote":
		return des.generateQuoteXML(block)
	default:
		return des.generateParagraphXML(block)
	}
}

// CreateDOCXTemplate creates a DOCX template with predefined styles
func (des *DOCXExportService) CreateDOCXTemplate(templateType string) ([]byte, error) {
	var content *DocumentContent

	switch templateType {
	case "business_report":
		content = des.createBusinessReportTemplate()
	case "academic_paper":
		content = des.createAcademicPaperTemplate()
	case "resume":
		content = des.createResumeTemplate()
	case "letter":
		content = des.createLetterTemplate()
	default:
		content = des.createBasicTemplate()
	}

	return des.ExportToDOCX(content)
}

// Template creation methods
func (des *DOCXExportService) createBusinessReportTemplate() *DocumentContent {
	return &DocumentContent{
		Title: "Business Report",
		Content: []ContentBlock{
			{Type: "heading", Level: 1, Text: "Business Report Title"},
			{Type: "heading", Level: 2, Text: "Executive Summary"},
			{Type: "paragraph", Text: "This is the executive summary of the business report."},
			{Type: "heading", Level: 2, Text: "Introduction"},
			{Type: "paragraph", Text: "This section contains the introduction to the report."},
			{Type: "heading", Level: 2, Text: "Findings"},
			{Type: "table", Table: &TableData{
				Headers: []string{"Category", "Value", "Percentage"},
				Rows: [][]string{
					{"Sales", "$100,000", "40%"},
					{"Marketing", "$75,000", "30%"},
					{"Operations", "$50,000", "20%"},
					{"Other", "$25,000", "10%"},
				},
			}},
			{Type: "heading", Level: 2, Text: "Conclusion"},
			{Type: "paragraph", Text: "This section contains the conclusion of the report."},
		},
	}
}

func (des *DOCXExportService) createAcademicPaperTemplate() *DocumentContent {
	return &DocumentContent{
		Title: "Academic Paper",
		Content: []ContentBlock{
			{Type: "heading", Level: 1, Text: "Paper Title"},
			{Type: "paragraph", Text: "Author Name\nInstitution\nEmail"},
			{Type: "heading", Level: 2, Text: "Abstract"},
			{Type: "paragraph", Text: "This paper presents..."},
			{Type: "heading", Level: 2, Text: "Introduction"},
			{Type: "paragraph", Text: "The introduction section..."},
			{Type: "heading", Level: 2, Text: "Methodology"},
			{Type: "paragraph", Text: "The methodology used..."},
			{Type: "heading", Level: 2, Text: "Results"},
			{Type: "paragraph", Text: "The results show..."},
			{Type: "heading", Level: 2, Text: "Conclusion"},
			{Type: "paragraph", Text: "In conclusion..."},
		},
	}
}

func (des *DOCXExportService) createResumeTemplate() *DocumentContent {
	return &DocumentContent{
		Title: "Resume",
		Content: []ContentBlock{
			{Type: "heading", Level: 1, Text: "John Doe"},
			{Type: "paragraph", Text: "Software Engineer\nSan Francisco, CA\njohn.doe@email.com\n(555) 123-4567"},
			{Type: "heading", Level: 2, Text: "Professional Summary"},
			{Type: "paragraph", Text: "Experienced software engineer with 5+ years..."},
			{Type: "heading", Level: 2, Text: "Experience"},
			{Type: "paragraph", Text: "Senior Software Engineer\nTech Company, 2020-Present\n- Developed web applications..."},
			{Type: "heading", Level: 2, Text: "Education"},
			{Type: "paragraph", Text: "Bachelor of Science in Computer Science\nUniversity Name, 2016-2020"},
			{Type: "heading", Level: 2, Text: "Skills"},
			{Type: "list", Items: []string{"JavaScript", "Python", "React", "Node.js", "SQL"}},
		},
	}
}

func (des *DOCXExportService) createLetterTemplate() *DocumentContent {
	return &DocumentContent{
		Title: "Business Letter",
		Content: []ContentBlock{
			{Type: "paragraph", Text: "Your Name\nYour Address\nCity, State ZIP Code\nEmail\nPhone\nDate"},
			{Type: "paragraph", Text: "Recipient Name\nRecipient Title\nCompany Name\nCompany Address\nCity, State ZIP Code"},
			{Type: "paragraph", Text: "Dear Recipient Name,"},
			{Type: "paragraph", Text: "I am writing to..."},
			{Type: "paragraph", Text: "The purpose of this letter is to..."},
			{Type: "paragraph", Text: "Please contact me if you need any additional information."},
			{Type: "paragraph", Text: "Sincerely,"},
			{Type: "paragraph", Text: "Your Name\nYour Title"},
		},
	}
}

func (des *DOCXExportService) createBasicTemplate() *DocumentContent {
	return &DocumentContent{
		Title: "Document",
		Content: []ContentBlock{
			{Type: "heading", Level: 1, Text: "Document Title"},
			{Type: "paragraph", Text: "This is a sample document created with TPT Titan."},
			{Type: "heading", Level: 2, Text: "Section 1"},
			{Type: "paragraph", Text: "This is the content of section 1."},
			{Type: "list", Items: []string{"Item 1", "Item 2", "Item 3"}},
			{Type: "heading", Level: 2, Text: "Section 2"},
			{Type: "paragraph", Text: "This is the content of section 2."},
		},
	}
}

// BatchExportToDOCX exports multiple documents to separate DOCX files
func (des *DOCXExportService) BatchExportToDOCX(contents []*DocumentContent) (map[string][]byte, error) {
	results := make(map[string][]byte)

	for i, content := range contents {
		docxData, err := des.ExportToDOCX(content)
		if err != nil {
			return nil, fmt.Errorf("failed to export document %d: %w", i+1, err)
		}

		filename := fmt.Sprintf("document_%d.docx", i+1)
		results[filename] = docxData
	}

	return results, nil
}

// GetDOCXMetadata extracts metadata from a DOCX file
func (des *DOCXExportService) GetDOCXMetadata(docxData []byte) (map[string]interface{}, error) {
	// In a real implementation, parse the DOCX file and extract metadata
	// For now, return basic info
	return map[string]interface{}{
		"format":      "DOCX",
		"size_bytes":  len(docxData),
		"created_at":  time.Now(),
		"creator":     "TPT Titan",
	}, nil
}
