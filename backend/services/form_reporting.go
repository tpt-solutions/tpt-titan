package services

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FormReportingService provides comprehensive reporting capabilities for forms
type FormReportingService struct {
	db               *sql.DB
	queryBuilder     *QueryBuilderService
	relationshipSvc  *FormRelationshipService
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
	Type      string `json:"type"`      // "text", "number", "date", "calculated"
	Format    string `json:"format,omitempty"`
	Aggregate string `json:"aggregate,omitempty"` // "sum", "avg", "count", "min", "max"
	Formula   string `json:"formula,omitempty"`   // For calculated columns
}

// ChartConfig represents chart configuration for reports
type ChartConfig struct {
	Type       string   `json:"type"`        // "bar", "line", "pie", "scatter"
	XAxisField string   `json:"x_axis_field"`
	YAxisField string   `json:"y_axis_field"`
	GroupField string   `json:"group_field,omitempty"`
	Title      string   `json:"title"`
	XAxisLabel string   `json:"x_axis_label,omitempty"`
	YAxisLabel string   `json:"y_axis_label,omitempty"`
}

// ReportResult represents the result of executing a report
type ReportResult struct {
	ReportID   uuid.UUID              `json:"report_id"`
	Columns    []ReportColumn         `json:"columns"`
	Data       []map[string]interface{} `json:"data"`
	Summary    map[string]interface{} `json:"summary"`
	ChartData  *ChartData             `json:"chart_data,omitempty"`
	ExecutedAt time.Time              `json:"executed_at"`
}

// ChartData represents processed data for charts
type ChartData struct {
	Type   string                   `json:"type"`
	Labels []string                 `json:"labels"`
	Datasets []ChartDataset         `json:"datasets"`
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
		ID:   "table_1",
		Type: "table",
		Table: fmt.Sprintf("form_responses_%s", report.FormID.String()[:8]),
		Position: Position{X: 100, Y: 100},
	})

	// Add fields
	for i, col := range report.Columns {
		element := QueryElement{
			ID:     fmt.Sprintf("field_%d", i+1),
			Type:   "field",
			Table:  fmt.Sprintf("form_responses_%s", report.FormID.String()[:8]),
			Field:  col.Field,
			Alias:  col.Label,
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
			ID:   fmt.Sprintf("group_%d", i+1),
			Type: "group",
			Field: group,
			Position: Position{X: 100, Y: 300 + i*50},
		}
		elements = append(elements, element)
	}

	// Add sort
	for i, sort := range report.SortBy {
		element := QueryElement{
			ID:   fmt.Sprintf("sort_%d", i+1),
			Type: "sort",
			Field: sort.Field,
			Position: Position{X: 100, Y: 400 + i*50},
			Data: map[string]interface{}{"direction": sort.Direction},
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
	// Would use the Excel service to export
	// For now, return CSV as Excel
	return frs.exportToCSV(result)
}

func (frs *FormReportingService) exportToPDF(result *ReportResult) ([]byte, error) {
	// Would generate PDF report
	// For now, return placeholder
	return []byte("PDF export not implemented"), nil
}

func (frs *FormReportingService) exportToJSON(result *ReportResult) ([]byte, error) {
	// Would serialize to JSON
	// For now, return placeholder
	return []byte("JSON export not implemented"), nil
}
