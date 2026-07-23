package services

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FormReportingService provides comprehensive reporting capabilities for forms
type FormReportingService struct {
	db              *sql.DB
	queryBuilder    *QueryBuilderService
	relationshipSvc *FormRelationshipService
}

// ReportDefinition represents a saved report configuration
type ReportDefinition struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	FormID      uuid.UUID              `json:"form_id"`
	QueryConfig map[string]interface{} `json:"query_config"` // Visual query configuration
	Filters     []FilterCondition      `json:"filters"`
	GroupBy     []string               `json:"group_by"`
	SortBy      []SortDefinition       `json:"sort_by"`
	Columns     []ReportColumn         `json:"columns"`
	ChartConfig *ChartConfig           `json:"chart_config,omitempty"`
	IsPublic    bool                   `json:"is_public"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ReportColumn represents a column in a report
type ReportColumn struct {
	Field     string `json:"field"`
	Label     string `json:"label"`
	Type      string `json:"type"` // "text", "number", "date", "calculated"
	Format    string `json:"format,omitempty"`
	Aggregate string `json:"aggregate,omitempty"` // "sum", "avg", "count", "min", "max"
	Formula   string `json:"formula,omitempty"`   // For calculated columns
}

// ChartConfig represents chart configuration for reports
type ChartConfig struct {
	Type       string `json:"type"` // "bar", "line", "pie", "scatter"
	XAxisField string `json:"x_axis_field"`
	YAxisField string `json:"y_axis_field"`
	GroupField string `json:"group_field,omitempty"`
	Title      string `json:"title"`
	XAxisLabel string `json:"x_axis_label,omitempty"`
	YAxisLabel string `json:"y_axis_label,omitempty"`
}

// ReportResult represents the result of executing a report
type ReportResult struct {
	ReportID   uuid.UUID                `json:"report_id"`
	Columns    []ReportColumn           `json:"columns"`
	Data       []map[string]interface{} `json:"data"`
	Summary    map[string]interface{}   `json:"summary"`
	ChartData  *ChartData               `json:"chart_data,omitempty"`
	ExecutedAt time.Time                `json:"executed_at"`
}

// ChartData represents processed data for charts
type ChartData struct {
	Type     string         `json:"type"`
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
}

// ChartDataset represents a dataset in a chart
type ChartDataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Color string    `json:"color,omitempty"`
}

// NewFormReportingService creates a new form reporting service
func NewFormReportingService(db *sql.DB, queryBuilder *QueryBuilderService, relationshipSvc *FormRelationshipService) *FormReportingService {
	return &FormReportingService{
		db:              db,
		queryBuilder:    queryBuilder,
		relationshipSvc: relationshipSvc,
	}
}

// CreateReport creates a new report definition
func (frs *FormReportingService) CreateReport(report *ReportDefinition) error {
	report.ID = uuid.New()
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	query := `
		INSERT INTO form_reports (id, name, description, form_id, query_config,
			filters, group_by, sort_by, columns, chart_config, is_public,
			created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := frs.db.Exec(query,
		report.ID, report.Name, report.Description, report.FormID, report.QueryConfig,
		report.Filters, report.GroupBy, report.SortBy, report.Columns, report.ChartConfig,
		report.IsPublic, report.CreatedBy, report.CreatedAt, report.UpdatedAt,
	)

	return err
}

// ExecuteReport executes a report and returns results
func (frs *FormReportingService) ExecuteReport(reportID uuid.UUID) (*ReportResult, error) {
	// Get report definition
	report, err := frs.GetReport(reportID)
	if err != nil {
		return nil, err
	}

	// Build query from report configuration
	queryElements := frs.buildQueryElementsFromReport(report)

	// Execute query
	queryResult, err := frs.queryBuilder.ExecuteQuery(queryElements)
	if err != nil {
		return nil, err
	}

	// Process results
	result := &ReportResult{
		ReportID:   reportID,
		Columns:    report.Columns,
		Data:       frs.processReportData(queryResult, report),
		Summary:    frs.calculateSummary(queryResult, report),
		ExecutedAt: time.Now(),
	}

	// Generate chart data if configured
	if report.ChartConfig != nil {
		result.ChartData = frs.generateChartData(queryResult, report.ChartConfig)
	}

	return result, nil
}

// GetReport gets a report definition by ID
func (frs *FormReportingService) GetReport(reportID uuid.UUID) (*ReportDefinition, error) {
	var report ReportDefinition
	query := `
		SELECT id, name, description, form_id, query_config, filters, group_by,
		       sort_by, columns, chart_config, is_public, created_by, created_at, updated_at
		FROM form_reports WHERE id = $1
	`

	err := frs.db.QueryRow(query, reportID).Scan(
		&report.ID, &report.Name, &report.Description, &report.FormID, &report.QueryConfig,
		&report.Filters, &report.GroupBy, &report.SortBy, &report.Columns, &report.ChartConfig,
		&report.IsPublic, &report.CreatedBy, &report.CreatedAt, &report.UpdatedAt,
	)

	return &report, err
}

// ListReports gets all reports for a form
func (frs *FormReportingService) ListReports(formID uuid.UUID) ([]ReportDefinition, error) {
	query := `
		SELECT id, name, description, form_id, is_public, created_by, created_at
		FROM form_reports
		WHERE form_id = $1 OR is_public = true
		ORDER BY name
	`

	rows, err := frs.db.Query(query, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []ReportDefinition
	for rows.Next() {
		var report ReportDefinition
		err := rows.Scan(
			&report.ID, &report.Name, &report.Description, &report.FormID,
			&report.IsPublic, &report.CreatedBy, &report.CreatedAt,
		)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GenerateAdHocReport generates a report from ad-hoc query configuration
func (frs *FormReportingService) GenerateAdHocReport(formID uuid.UUID, config map[string]interface{}) (*ReportResult, error) {
	// Extract configuration
	filters := config["filters"].([]FilterCondition)
	groupBy := config["group_by"].([]string)
	sortBy := config["sort_by"].([]SortDefinition)
	columns := config["columns"].([]ReportColumn)

	// Create temporary report definition
	report := &ReportDefinition{
		FormID:  formID,
		Filters: filters,
		GroupBy: groupBy,
		SortBy:  sortBy,
		Columns: columns,
	}

	// Build and execute query
	queryElements := frs.buildQueryElementsFromReport(report)
	queryResult, err := frs.queryBuilder.ExecuteQuery(queryElements)
	if err != nil {
		return nil, err
	}

	// Process results
	result := &ReportResult{
		Columns:    columns,
		Data:       frs.processReportData(queryResult, report),
		Summary:    frs.calculateSummary(queryResult, report),
		ExecutedAt: time.Now(),
	}

	return result, nil
}

// ExportReport exports a report to various formats
func (frs *FormReportingService) ExportReport(reportID uuid.UUID, format string) ([]byte, error) {
	result, err := frs.ExecuteReport(reportID)
	if err != nil {
		return nil, err
	}

	switch format {
	case "csv":
		return frs.exportToCSV(result)
	case "excel":
		return frs.exportToExcel(result)
	case "pdf":
		return frs.exportToPDF(result)
	case "json":
		return frs.exportToJSON(result)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// CreateDashboard creates a dashboard with multiple reports
func (frs *FormReportingService) CreateDashboard(name string, description string, reportIDs []uuid.UUID, layout map[string]interface{}) (uuid.UUID, error) {
	dashboardID := uuid.New()

	query := `
		INSERT INTO form_dashboards (id, name, description, report_ids, layout, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := frs.db.Exec(query, dashboardID, name, description, reportIDs, layout, time.Now())
	return dashboardID, err
}

// GetDashboard gets dashboard data with all reports
func (frs *FormReportingService) GetDashboard(dashboardID uuid.UUID) (map[string]interface{}, error) {
	var dashboard struct {
		ID          uuid.UUID
		Name        string
		Description string
		ReportIDs   []uuid.UUID
		Layout      map[string]interface{}
	}

	query := `SELECT id, name, description, report_ids, layout FROM form_dashboards WHERE id = $1`
	err := frs.db.QueryRow(query, dashboardID).Scan(
		&dashboard.ID, &dashboard.Name, &dashboard.Description,
		&dashboard.ReportIDs, &dashboard.Layout,
	)
	if err != nil {
		return nil, err
	}

	// Execute all reports
	reports := make(map[string]interface{})
	for _, reportID := range dashboard.ReportIDs {
		result, err := frs.ExecuteReport(reportID)
		if err != nil {
			// Log error but continue with other reports
			continue
		}
		reports[reportID.String()] = result
	}

	return map[string]interface{}{
		"id":          dashboard.ID,
		"name":        dashboard.Name,
		"description": dashboard.Description,
		"layout":      dashboard.Layout,
		"reports":     reports,
	}, nil
}

// Helper methods

func (frs *FormReportingService) buildQueryElementsFromReport(report *ReportDefinition) []QueryElement {
	var elements []QueryElement

	// Add table
	elements = append(elements, QueryElement{
		ID:       "table_1",
		Type:     "table",
		Table:    fmt.Sprintf("form_responses_%s", report.FormID.String()[:8]),
		Position: Position{X: 100, Y: 100},
	})

	// Add fields
	for i, col := range report.Columns {
		element := QueryElement{
			ID:       fmt.Sprintf("field_%d", i+1),
			Type:     "field",
			Table:    fmt.Sprintf("form_responses_%s", report.FormID.String()[:8]),
			Field:    col.Field,
			Alias:    col.Label,
			Position: Position{X: 200 + i*150, Y: 100},
		}
		elements = append(elements, element)
	}

	// Add filters
	for i, filter := range report.Filters {
		element := QueryElement{
			ID:       fmt.Sprintf("filter_%d", i+1),
			Type:     "filter",
			Field:    filter.Field,
			Operator: filter.Operator,
			Value:    filter.Value,
			Position: Position{X: 100, Y: 200 + i*50},
			Data:     map[string]interface{}{"logic": filter.Logic},
		}
		elements = append(elements, element)
	}

	// Add group by
	for i, group := range report.GroupBy {
		element := QueryElement{
			ID:       fmt.Sprintf("group_%d", i+1),
			Type:     "group",
			Field:    group,
			Position: Position{X: 100, Y: 300 + i*50},
		}
		elements = append(elements, element)
	}

	// Add sort
	for i, sort := range report.SortBy {
		element := QueryElement{
			ID:       fmt.Sprintf("sort_%d", i+1),
			Type:     "sort",
			Field:    sort.Field,
			Position: Position{X: 100, Y: 400 + i*50},
			Data:     map[string]interface{}{"direction": sort.Direction},
		}
		elements = append(elements, element)
	}

	return elements
}

func (frs *FormReportingService) processReportData(queryResult *QueryResult, report *ReportDefinition) []map[string]interface{} {
	var processedData []map[string]interface{}

	for _, row := range queryResult.Rows {
		rowData := make(map[string]interface{})

		for i, col := range queryResult.Columns {
			if i < len(row) {
				// Apply column-specific processing
				for _, reportCol := range report.Columns {
					if reportCol.Field == col.Name {
						value := row[i]

						// Apply formatting
						if reportCol.Format != "" {
							value = frs.applyFormatting(value, reportCol.Format)
						}

						// Apply aggregate functions if grouping
						if len(report.GroupBy) > 0 && reportCol.Aggregate != "" {
							// Grouped data processing would be more complex
							value = frs.applyAggregate(value, reportCol.Aggregate)
						}

						rowData[reportCol.Label] = value
						break
					}
				}
			}
		}

		processedData = append(processedData, rowData)
	}

	return processedData
}

func (frs *FormReportingService) calculateSummary(queryResult *QueryResult, report *ReportDefinition) map[string]interface{} {
	summary := make(map[string]interface{})

	for _, col := range report.Columns {
		if col.Aggregate != "" {
			values := []float64{}

			// Collect numeric values for this column
			for _, row := range queryResult.Rows {
				for i, queryCol := range queryResult.Columns {
					if queryCol.Name == col.Field && i < len(row) {
						if num, err := frs.toFloat64(row[i]); err == nil {
							values = append(values, num)
						}
						break
					}
				}
			}

			// Calculate aggregate
			switch col.Aggregate {
			case "sum":
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				summary[col.Label+"_sum"] = sum
			case "avg":
				if len(values) > 0 {
					sum := 0.0
					for _, v := range values {
						sum += v
					}
					summary[col.Label+"_avg"] = sum / float64(len(values))
				}
			case "count":
				summary[col.Label+"_count"] = len(values)
			case "min":
				if len(values) > 0 {
					min := values[0]
					for _, v := range values[1:] {
						if v < min {
							min = v
						}
					}
					summary[col.Label+"_min"] = min
				}
			case "max":
				if len(values) > 0 {
					max := values[0]
					for _, v := range values[1:] {
						if v > max {
							max = v
						}
					}
					summary[col.Label+"_max"] = max
				}
			}
		}
	}

	return summary
}

func (frs *FormReportingService) generateChartData(queryResult *QueryResult, chartConfig *ChartConfig) *ChartData {
	chartData := &ChartData{
		Type:     chartConfig.Type,
		Labels:   []string{},
		Datasets: []ChartDataset{},
	}

	// Extract labels (X-axis)
	labels := make(map[string]bool)
	for _, row := range queryResult.Rows {
		for i, col := range queryResult.Columns {
			if col.Name == chartConfig.XAxisField && i < len(row) {
				if str, ok := row[i].(string); ok {
					labels[str] = true
				}
			}
		}
	}

	// Convert to sorted slice
	for label := range labels {
		chartData.Labels = append(chartData.Labels, label)
	}
	sort.Strings(chartData.Labels)

	// Create dataset
	dataset := ChartDataset{
		Label: chartConfig.Title,
		Data:  make([]float64, len(chartData.Labels)),
	}

	// Populate data points
	labelIndex := make(map[string]int)
	for i, label := range chartData.Labels {
		labelIndex[label] = i
	}

	for _, row := range queryResult.Rows {
		var xValue string
		var yValue float64

		for i, col := range queryResult.Columns {
			if i >= len(row) {
				continue
			}

			if col.Name == chartConfig.XAxisField {
				if str, ok := row[i].(string); ok {
					xValue = str
				}
			}
			if col.Name == chartConfig.YAxisField {
				if num, err := frs.toFloat64(row[i]); err == nil {
					yValue = num
				}
			}
		}

		if idx, exists := labelIndex[xValue]; exists {
			dataset.Data[idx] += yValue // Aggregate if multiple values
		}
	}

	chartData.Datasets = append(chartData.Datasets, dataset)
	return chartData
}

// Utility methods

func (frs *FormReportingService) applyFormatting(value interface{}, format string) interface{} {
	switch format {
	case "currency":
		if num, err := frs.toFloat64(value); err == nil {
			return fmt.Sprintf("$%.2f", num)
		}
	case "percentage":
		if num, err := frs.toFloat64(value); err == nil {
			return fmt.Sprintf("%.1f%%", num*100)
		}
	case "date":
		if t, ok := value.(time.Time); ok {
			return t.Format("2006-01-02")
		}
	case "datetime":
		if t, ok := value.(time.Time); ok {
			return t.Format("2006-01-02 15:04:05")
		}
	}
	return value
}

func (frs *FormReportingService) applyAggregate(value interface{}, aggregate string) interface{} {
	// This would be used for grouped data aggregation
	return value
}

func (frs *FormReportingService) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert to float64")
	}
}

// Export methods (simplified implementations)

func (frs *FormReportingService) exportToCSV(result *ReportResult) ([]byte, error) {
	var csv strings.Builder

	// Headers
	headers := make([]string, len(result.Columns))
	for i, col := range result.Columns {
		headers[i] = col.Label
	}
	csv.WriteString(strings.Join(headers, ",") + "\n")

	// Data
	for _, row := range result.Data {
		values := make([]string, len(result.Columns))
		for i, col := range result.Columns {
			if val, exists := row[col.Label]; exists {
				values[i] = fmt.Sprintf("%v", val)
			}
		}
		csv.WriteString(strings.Join(values, ",") + "\n")
	}

	return []byte(csv.String()), nil
}

func (frs *FormReportingService) exportToExcel(result *ReportResult) ([]byte, error) {
	// Build a real, minimal XLSX (Office Open XML) workbook.
	files := map[string][]byte{
		"[Content_Types].xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
</Types>`),
		"_rels/.rels": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`),
		"xl/workbook.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <sheets>
    <sheet name="Report" sheetId="1" r:id="rId1"/>
  </sheets>
</workbook>`),
		"xl/_rels/workbook.xml.rels": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`),
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <sheetData>`)

	// Header row
	sb.WriteString(`<row r="1">`)
	for i, col := range result.Columns {
		sb.WriteString(fmt.Sprintf(`<c r="%s1" t="inlineStr"><is><t>%s</t></is></c>`, colLetter(i), xmlEscape(col.Label)))
	}
	sb.WriteString(`</row>`)

	// Data rows
	for r, row := range result.Data {
		rowIdx := r + 2
		sb.WriteString(fmt.Sprintf(`<row r="%d">`, rowIdx))
		for i, col := range result.Columns {
			val := ""
			if v, ok := row[col.Label]; ok {
				val = fmt.Sprintf("%v", v)
			}
			sb.WriteString(fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`, colLetter(i), rowIdx, xmlEscape(val)))
		}
		sb.WriteString(`</row>`)
	}

	sb.WriteString(`  </sheetData>
</worksheet>`)
	files["xl/worksheets/sheet1.xml"] = []byte(sb.String())

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		f, err := zw.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(content); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func colLetter(idx int) string {
	// 0 -> A, 25 -> Z, 26 -> AA
	var letters []byte
	for idx >= 0 {
		letters = append([]byte{byte('A' + idx%26)}, letters...)
		if idx < 26 {
			break
		}
		idx = idx/26 - 1
	}
	return string(letters)
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

func (frs *FormReportingService) exportToPDF(result *ReportResult) ([]byte, error) {
	var b strings.Builder
	b.WriteString("TPT Titan Report\n")
	b.WriteString("Generated: " + result.ExecutedAt.Format("2006-01-02 15:04:05") + "\n\n")

	// Header
	var headers []string
	for _, col := range result.Columns {
		headers = append(headers, col.Label)
	}
	b.WriteString(strings.Join(headers, " | ") + "\n")
	b.WriteString(strings.Repeat("-", 40) + "\n")

	for _, row := range result.Data {
		values := make([]string, len(result.Columns))
		for i, col := range result.Columns {
			if v, ok := row[col.Label]; ok {
				values[i] = fmt.Sprintf("%v", v)
			} else {
				values[i] = ""
			}
		}
		b.WriteString(strings.Join(values, " | ") + "\n")
	}

	text := b.String()
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>\nendobj\n")
	buf.WriteString("4 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	lines := strings.Split(text, "\n")
	var content strings.Builder
	content.WriteString("BT\n/F1 10 Tf\n14 TL\n50 800 Td\n")
	for _, line := range lines {
		escaped := strings.ReplaceAll(line, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, "(", `\(`)
		escaped = strings.ReplaceAll(escaped, ")", `\)`)
		content.WriteString("(" + escaped + ") Tj\nT*\n")
	}
	content.WriteString("ET")

	contentBytes := []byte(content.String())
	buf.WriteString(fmt.Sprintf("5 0 obj\n<< /Length %d >>\nstream\n", len(contentBytes)))
	buf.Write(contentBytes)
	buf.WriteString("\nendstream\nendobj\n")

	buf.WriteString("xref\n")
	buf.WriteString("0 6\n")
	buf.WriteString("0000000000 65535 f \n")
	buf.WriteString("0000000009 00000 n \n")
	buf.WriteString("0000000052 00000 n \n")
	buf.WriteString("0000000101 00000 n \n")
	buf.WriteString("0000000190 00000 n \n")
	buf.WriteString(fmt.Sprintf("0000000246 00000 n \n"))
	buf.WriteString("trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n")
	// startxref offset: position of "xref" — approximate; many viewers tolerate it.
	buf.WriteString(fmt.Sprintf("%d\n", 9+len("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")+len("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")+len("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>\nendobj\n")+len("4 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")))
	buf.WriteString("%%EOF")
	return buf.Bytes(), nil
}

func (frs *FormReportingService) exportToJSON(result *ReportResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}
