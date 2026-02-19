package services

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// ExcelService provides Excel file import/export functionality
type ExcelService struct {
	mathService *SpreadsheetMathService
}

// ExcelImportResult represents the result of importing an Excel file
type ExcelImportResult struct {
	SpreadsheetID uuid.UUID              `json:"spreadsheet_id"`
	Name          string                 `json:"name"`
	Sheets        []ExcelSheet           `json:"sheets"`
	Metadata      map[string]interface{} `json:"metadata"`
	Errors        []string               `json:"errors,omitempty"`
}

// ExcelSheet represents a worksheet in an Excel file
type ExcelSheet struct {
	Name    string                 `json:"name"`
	Data    map[string]interface{} `json:"data"`
	Formulas map[string]string     `json:"formulas,omitempty"`
	Merged  []string               `json:"merged,omitempty"`  // Merged cell ranges
	Styles  map[string]ExcelStyle  `json:"styles,omitempty"`
}

// ExcelStyle represents cell styling information
type ExcelStyle struct {
	Bold         bool    `json:"bold,omitempty"`
	Italic       bool    `json:"italic,omitempty"`
	FontSize     float64 `json:"font_size,omitempty"`
	FontColor    string  `json:"font_color,omitempty"`
	Background   string  `json:"background,omitempty"`
	Border       string  `json:"border,omitempty"`
	TextAlign    string  `json:"text_align,omitempty"`
	VerticalAlign string  `json:"vertical_align,omitempty"`
	Format       string  `json:"format,omitempty"` // Number format
}

// ExcelExportOptions represents options for exporting to Excel
type ExcelExportOptions struct {
	IncludeFormulas bool                   `json:"include_formulas,omitempty"`
	IncludeStyles   bool                   `json:"include_styles,omitempty"`
	SheetName       string                 `json:"sheet_name,omitempty"`
	DataRange       string                 `json:"data_range,omitempty"`
	CustomStyles    map[string]ExcelStyle  `json:"custom_styles,omitempty"`
}

// NewExcelService creates a new Excel service
func NewExcelService(mathService *SpreadsheetMathService) *ExcelService {
	return &ExcelService{
		mathService: mathService,
	}
}

// ImportExcel imports data from an Excel file
func (es *ExcelService) ImportExcel(fileData io.Reader, filename string) (*ExcelImportResult, error) {
	f, err := excelize.OpenReader(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	result := &ExcelImportResult{
		SpreadsheetID: uuid.New(),
		Name:          strings.TrimSuffix(filename, ".xlsx"),
		Sheets:        []ExcelSheet{},
		Metadata:      make(map[string]interface{}),
		Errors:        []string{},
	}

	// Get all sheet names
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return nil, fmt.Errorf("no worksheets found in Excel file")
	}

	// Process each sheet
	for _, sheetName := range sheetNames {
		sheet, err := es.importSheet(f, sheetName)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Sheet %s: %v", sheetName, err))
			continue
		}
		result.Sheets = append(result.Sheets, *sheet)
	}

	// Extract metadata
	if docProps, err := f.GetDocProps(); err == nil {
		result.Metadata["creator"] = docProps.Creator
		result.Metadata["created"] = docProps.Created
		result.Metadata["modified"] = docProps.Modified
	}
	result.Metadata["sheet_count"] = len(sheetNames)

	return result, nil
}

// importSheet imports a single worksheet
func (es *ExcelService) importSheet(f *excelize.File, sheetName string) (*ExcelSheet, error) {
	sheet := &ExcelSheet{
		Name:     sheetName,
		Data:     make(map[string]interface{}),
		Formulas: make(map[string]string),
		Merged:   []string{},
		Styles:   make(map[string]ExcelStyle),
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	// Get merged cells
	mergeCells, err := f.GetMergeCells(sheetName)
	if err == nil {
		for _, mergeCell := range mergeCells {
			sheet.Merged = append(sheet.Merged, mergeCell.GetStartAxis()+":"+mergeCell.GetEndAxis())
		}
	}

	// Process each cell
	for rowIdx, row := range rows {
		for colIdx, cellValue := range row {
			if cellValue == "" {
				continue
			}

			// Convert column index to letter (0 = A, 1 = B, etc.)
			colLetter := es.columnIndexToLetter(colIdx)
			cellRef := fmt.Sprintf("%s%d", colLetter, rowIdx+1)

			// Get cell formula if it exists
			formula, err := f.GetCellFormula(sheetName, cellRef)
			if err == nil && formula != "" {
				// Convert Excel formula to our format
				ourFormula := es.convertExcelFormula(formula)
				sheet.Formulas[cellRef] = ourFormula
			}

			// Try to parse as number
			if num, err := strconv.ParseFloat(cellValue, 64); err == nil {
				sheet.Data[cellRef] = num
			} else if date, err := es.parseExcelDate(cellValue); err == nil {
				sheet.Data[cellRef] = date.Format("2006-01-02")
			} else {
				sheet.Data[cellRef] = cellValue
			}

			// Get cell style
			styleID, err := f.GetCellStyle(sheetName, cellRef)
			if err == nil {
				style := es.convertExcelStyle(f, styleID)
				if style != nil {
					sheet.Styles[cellRef] = *style
				}
			}
		}
	}

	return sheet, nil
}

// ExportExcel exports spreadsheet data to Excel format
func (es *ExcelService) ExportExcel(spreadsheetID uuid.UUID, sheets []ExcelSheet, options ExcelExportOptions) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Remove default sheet if we have custom sheets
	if len(sheets) > 0 {
		f.DeleteSheet("Sheet1")
	}

	// Process each sheet
	for sheetIdx, sheet := range sheets {
		var sheetName string
		if options.SheetName != "" && sheetIdx == 0 {
			sheetName = options.SheetName
		} else {
			sheetName = sheet.Name
			if sheetName == "" {
				sheetName = fmt.Sprintf("Sheet%d", sheetIdx+1)
			}
		}

		// Create new sheet
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return nil, fmt.Errorf("failed to create sheet %s: %w", sheetName, err)
		}

		// Set as active sheet for first sheet
		if sheetIdx == 0 {
			f.SetActiveSheet(index)
		}

		// Write data
		maxRow, maxCol := 0, 0
		for cellRef, value := range sheet.Data {
			col, row, err := es.parseCellReference(cellRef)
			if err != nil {
				continue
			}

			if row > maxRow {
				maxRow = row
			}
			if col > maxCol {
				maxCol = col
			}

			axis, _ := excelize.CoordinatesToCellName(col+1, row+1)

			// Handle different data types
			switch v := value.(type) {
			case float64:
				f.SetCellValue(sheetName, axis, v)
			case int:
				f.SetCellValue(sheetName, axis, v)
			case string:
				f.SetCellValue(sheetName, axis, v)
			case bool:
				f.SetCellValue(sheetName, axis, v)
			default:
				f.SetCellValue(sheetName, axis, fmt.Sprintf("%v", v))
			}

			// Apply styles if requested
			if options.IncludeStyles && sheet.Styles != nil {
				if style, exists := sheet.Styles[cellRef]; exists {
					styleID := es.applyExcelStyle(f, style)
					f.SetCellStyle(sheetName, axis, axis, styleID)
				}
			}
		}

		// Write formulas if requested
		if options.IncludeFormulas && sheet.Formulas != nil {
			for cellRef, formula := range sheet.Formulas {
				col, row, err := es.parseCellReference(cellRef)
				if err != nil {
					continue
				}

				axis, _ := excelize.CoordinatesToCellName(col+1, row+1)
				excelFormula := es.convertToExcelFormula(formula)
				f.SetCellFormula(sheetName, axis, excelFormula)
			}
		}

		// Apply merged cells
		if sheet.Merged != nil {
			for _, mergeRange := range sheet.Merged {
				f.MergeCell(sheetName, mergeRange[:strings.Index(mergeRange, ":")], mergeRange[strings.Index(mergeRange, ":")+1:])
			}
		}
	}

	// Set document properties
	f.SetDocProps(&excelize.DocProperties{
		Title:    fmt.Sprintf("TPT Titan Spreadsheet - %s", spreadsheetID.String()[:8]),
		Creator:  "TPT Titan",
		Created:  time.Now().Format(time.RFC3339),
		Modified: time.Now().Format(time.RFC3339),
	})

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write Excel file: %w", err)
	}

	return buffer.Bytes(), nil
}

// convertExcelFormula converts Excel formula syntax to our internal format
func (es *ExcelService) convertExcelFormula(excelFormula string) string {
	// Remove leading = if present
	if strings.HasPrefix(excelFormula, "=") {
		excelFormula = excelFormula[1:]
	}

	// Basic conversions - this would need to be much more comprehensive
	// for full Excel compatibility

	// Convert function names to uppercase
	excelFormula = strings.ToUpper(excelFormula)

	// Convert Excel range notation (commas) to our format (colons for ranges)
	// This is a simplified conversion - full implementation would be complex

	return "=" + excelFormula
}

// convertToExcelFormula converts our formula format to Excel format
func (es *ExcelService) convertToExcelFormula(ourFormula string) string {
	if !strings.HasPrefix(ourFormula, "=") {
		return ourFormula
	}

	formula := ourFormula[1:]

	// Basic conversion - would need full implementation
	return formula
}

// parseExcelDate attempts to parse Excel date formats
func (es *ExcelService) parseExcelDate(cellValue string) (time.Time, error) {
	// Try common date formats
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		time.RFC3339,
	}

	for _, format := range formats {
		if date, err := time.Parse(format, cellValue); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("not a date")
}

// convertExcelStyle converts Excel style to our style format
func (es *ExcelService) convertExcelStyle(f *excelize.File, styleID int) *ExcelStyle {
	style := &ExcelStyle{}

	// Get style from Excel file
	excelStyle, err := f.GetStyle(styleID)
	if err != nil {
		return nil
	}

	// Convert font properties
	if excelStyle.Font != nil {
		style.Bold = excelStyle.Font.Bold
		style.Italic = excelStyle.Font.Italic
		style.FontSize = excelStyle.Font.Size
		if excelStyle.Font.Color != "" {
			style.FontColor = excelStyle.Font.Color
		}
	}

	// Convert fill/background
	if excelStyle.Fill.Type != "" && excelStyle.Fill.Pattern == 1 && len(excelStyle.Fill.Color) > 0 {
		style.Background = excelStyle.Fill.Color[0]
	}

	// Convert alignment
	if excelStyle.Alignment != nil {
		switch excelStyle.Alignment.Horizontal {
		case "left":
			style.TextAlign = "left"
		case "center":
			style.TextAlign = "center"
		case "right":
			style.TextAlign = "right"
		}

		switch excelStyle.Alignment.Vertical {
		case "top":
			style.VerticalAlign = "top"
		case "center":
			style.VerticalAlign = "middle"
		case "bottom":
			style.VerticalAlign = "bottom"
		}
	}

	// Convert number format (NumFmt is int in excelize v2)
	if excelStyle.NumFmt != 0 {
		style.Format = strconv.Itoa(excelStyle.NumFmt)
	}

	return style
}

// applyExcelStyle applies our style format to Excel
func (es *ExcelService) applyExcelStyle(f *excelize.File, style ExcelStyle) int {
	excelStyle := &excelize.Style{}

	// Font
	if style.Bold || style.Italic || style.FontSize > 0 || style.FontColor != "" {
		excelStyle.Font = &excelize.Font{}
		if style.Bold {
			excelStyle.Font.Bold = true
		}
		if style.Italic {
			excelStyle.Font.Italic = true
		}
		if style.FontSize > 0 {
			excelStyle.Font.Size = style.FontSize
		}
		if style.FontColor != "" {
			excelStyle.Font.Color = style.FontColor
		}
	}

	// Fill
	if style.Background != "" {
		excelStyle.Fill = excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{style.Background},
		}
	}

	// Alignment
	if style.TextAlign != "" || style.VerticalAlign != "" {
		excelStyle.Alignment = &excelize.Alignment{}
		if style.TextAlign != "" {
			excelStyle.Alignment.Horizontal = style.TextAlign
		}
		if style.VerticalAlign != "" {
			excelStyle.Alignment.Vertical = style.VerticalAlign
		}
	}

	styleID, _ := f.NewStyle(excelStyle)
	return styleID
}

// Helper functions

// columnIndexToLetter converts column index to Excel column letter (0 = A, 1 = B, etc.)
func (es *ExcelService) columnIndexToLetter(col int) string {
	result := ""
	for col >= 0 {
		result = string(rune('A'+col%26)) + result
		col = col/26 - 1
	}
	return result
}

// parseCellReference parses A1 style cell reference to column/row indices
func (es *ExcelService) parseCellReference(ref string) (col, row int, err error) {
	ref = strings.ToUpper(strings.TrimSpace(ref))

	// Find where numbers start
	i := 0
	for i < len(ref) && (ref[i] < '0' || ref[i] > '9') {
		i++
	}

	if i == 0 || i == len(ref) {
		return 0, 0, fmt.Errorf("invalid cell reference: %s", ref)
	}

	colStr := ref[:i]
	rowStr := ref[i:]

	// Convert column letters to number
	col = 0
	for _, r := range colStr {
		col = col*26 + int(r-'A'+1)
	}
	col-- // 0-based

	// Parse row number
	row, err = strconv.Atoi(rowStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row number in %s", ref)
	}
	row-- // 0-based

	return col, row, nil
}

// GetSupportedFormats returns supported Excel import/export formats
func (es *ExcelService) GetSupportedFormats() map[string]string {
	return map[string]string{
		".xlsx": "Excel 2007+ Workbook",
		".xls":  "Excel 97-2003 Workbook (limited support)",
	}
}

// ValidateExcelFile validates if a file is a valid Excel file
func (es *ExcelService) ValidateExcelFile(fileData io.Reader) error {
	f, err := excelize.OpenReader(fileData)
	if err != nil {
		return fmt.Errorf("invalid Excel file: %w", err)
	}
	defer f.Close()

	// Check if it has at least one worksheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("Excel file contains no worksheets")
	}

	return nil
}

// GetExcelInfo returns basic information about an Excel file
func (es *ExcelService) GetExcelInfo(fileData io.Reader) (map[string]interface{}, error) {
	f, err := excelize.OpenReader(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	info := make(map[string]interface{})

	// Get sheet list
	sheets := f.GetSheetList()
	info["sheet_count"] = len(sheets)
	info["sheet_names"] = sheets

	// Get document properties
	if props, err := f.GetDocProps(); err == nil {
		info["creator"] = props.Creator
		info["created"] = props.Created
		info["modified"] = props.Modified
	}

	// Get basic sheet info
	sheetInfo := make(map[string]interface{})
	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		colCount := 0
		if len(rows) > 0 {
			colCount = len(rows[0])
		}

		sheetInfo[sheetName] = map[string]interface{}{
			"row_count": len(rows),
			"col_count": colCount,
		}
	}
	info["sheets"] = sheetInfo

	return info, nil
}

// ConvertSpreadsheetToExcel converts TPT Titan spreadsheet data to Excel format
func (es *ExcelService) ConvertSpreadsheetToExcel(spreadsheetData map[string]interface{}) ([]byte, error) {
	sheets := []ExcelSheet{}

	// Extract sheets from spreadsheet data
	if sheetsData, ok := spreadsheetData["sheets"].([]interface{}); ok {
		for _, sheetData := range sheetsData {
			if sheetMap, ok := sheetData.(map[string]interface{}); ok {
				sheet := ExcelSheet{
					Name:     sheetMap["name"].(string),
					Data:     sheetMap["data"].(map[string]interface{}),
					Formulas: make(map[string]string),
				}

				if formulas, ok := sheetMap["formulas"].(map[string]interface{}); ok {
					for k, v := range formulas {
						if formulaStr, ok := v.(string); ok {
							sheet.Formulas[k] = formulaStr
						}
					}
				}

				sheets = append(sheets, sheet)
			}
		}
	}

	options := ExcelExportOptions{
		IncludeFormulas: true,
		IncludeStyles:   true,
	}

	return es.ExportExcel(uuid.New(), sheets, options)
}

// BatchImportExcel processes multiple Excel files
func (es *ExcelService) BatchImportExcel(files []io.Reader, filenames []string) ([]*ExcelImportResult, []error) {
	results := []*ExcelImportResult{}
	errors := []error{}

	for i, file := range files {
		filename := ""
		if i < len(filenames) {
			filename = filenames[i]
		}

		result, err := es.ImportExcel(file, filename)
		if err != nil {
			errors = append(errors, fmt.Errorf("file %d: %w", i+1, err))
			continue
		}

		results = append(results, result)
	}

	return results, errors
}

// GetExcelTemplate creates a template Excel file with common formulas and formatting
func (es *ExcelService) GetExcelTemplate(templateType string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	switch templateType {
	case "budget":
		return es.createBudgetTemplate(f)
	case "invoice":
		return es.createInvoiceTemplate(f)
	case "timesheet":
		return es.createTimesheetTemplate(f)
	case "data_analysis":
		return es.createDataAnalysisTemplate(f)
	default:
		return es.createBasicTemplate(f)
	}
}

// createBasicTemplate creates a basic Excel template
func (es *ExcelService) createBasicTemplate(f *excelize.File) ([]byte, error) {
	// Set some sample data and formulas
	f.SetCellValue("Sheet1", "A1", "Product")
	f.SetCellValue("Sheet1", "B1", "Price")
	f.SetCellValue("Sheet1", "C1", "Quantity")
	f.SetCellValue("Sheet1", "D1", "Total")

	f.SetCellValue("Sheet1", "A2", "Widget A")
	f.SetCellValue("Sheet1", "B2", 10.99)
	f.SetCellValue("Sheet1", "C2", 5)
	f.SetCellFormula("Sheet1", "D2", "=B2*C2")

	f.SetCellValue("Sheet1", "A3", "Widget B")
	f.SetCellValue("Sheet1", "B3", 15.50)
	f.SetCellValue("Sheet1", "C3", 3)
	f.SetCellFormula("Sheet1", "D3", "=B3*C3")

	f.SetCellValue("Sheet1", "A4", "Total")
	f.SetCellFormula("Sheet1", "D4", "=SUM(D2:D3)")

	// Apply some basic styling
	styleID, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle("Sheet1", "A1", "D1", styleID)

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// createBudgetTemplate creates a budget template
func (es *ExcelService) createBudgetTemplate(f *excelize.File) ([]byte, error) {
	f.SetCellValue("Sheet1", "A1", "Budget Template")
	f.SetCellValue("Sheet1", "A3", "Category")
	f.SetCellValue("Sheet1", "B3", "Budgeted")
	f.SetCellValue("Sheet1", "C3", "Actual")
	f.SetCellValue("Sheet1", "D3", "Difference")

	// Add sample categories
	categories := []string{"Rent", "Groceries", "Utilities", "Entertainment", "Transportation"}
	for i, category := range categories {
		row := i + 4
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), category)
		f.SetCellFormula("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("=B%d-C%d", row, row))
	}

	// Add totals
	f.SetCellValue("Sheet1", "A10", "Total")
	f.SetCellFormula("Sheet1", "B10", "=SUM(B4:B8)")
	f.SetCellFormula("Sheet1", "C10", "=SUM(C4:C8)")
	f.SetCellFormula("Sheet1", "D10", "=SUM(D4:D8)")

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// createInvoiceTemplate creates an invoice template
func (es *ExcelService) createInvoiceTemplate(f *excelize.File) ([]byte, error) {
	f.SetCellValue("Sheet1", "A1", "INVOICE")
	f.SetCellValue("Sheet1", "A3", "Invoice #:")
	f.SetCellValue("Sheet1", "A4", "Date:")
	f.SetCellValue("Sheet1", "A5", "Bill To:")

	f.SetCellValue("Sheet1", "A8", "Description")
	f.SetCellValue("Sheet1", "B8", "Quantity")
	f.SetCellValue("Sheet1", "C8", "Unit Price")
	f.SetCellValue("Sheet1", "D8", "Total")

	// Add sample line items
	for i := 9; i <= 13; i++ {
		f.SetCellFormula("Sheet1", fmt.Sprintf("D%d", i), fmt.Sprintf("=B%d*C%d", i, i))
	}

	f.SetCellValue("Sheet1", "C14", "Subtotal:")
	f.SetCellFormula("Sheet1", "D14", "=SUM(D9:D13)")
	f.SetCellValue("Sheet1", "C15", "Tax (10%):")
	f.SetCellFormula("Sheet1", "D15", "=D14*0.1")
	f.SetCellValue("Sheet1", "C16", "Total:")
	f.SetCellFormula("Sheet1", "D16", "=D14+D15")

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// createTimesheetTemplate creates a timesheet template
func (es *ExcelService) createTimesheetTemplate(f *excelize.File) ([]byte, error) {
	f.SetCellValue("Sheet1", "A1", "Weekly Timesheet")
	f.SetCellValue("Sheet1", "A3", "Employee:")
	f.SetCellValue("Sheet1", "A4", "Week Of:")

	f.SetCellValue("Sheet1", "A6", "Date")
	f.SetCellValue("Sheet1", "B6", "Day")
	f.SetCellValue("Sheet1", "C6", "Hours")
	f.SetCellValue("Sheet1", "D6", "Description")

	// Add days of the week
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i, day := range days {
		row := i + 7
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), day)
	}

	f.SetCellValue("Sheet1", "B14", "Total Hours:")
	f.SetCellFormula("Sheet1", "C14", "=SUM(C7:C13)")

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// createDataAnalysisTemplate creates a data analysis template
func (es *ExcelService) createDataAnalysisTemplate(f *excelize.File) ([]byte, error) {
	f.SetCellValue("Sheet1", "A1", "Data Analysis Template")
	f.SetCellValue("Sheet1", "A3", "Data Points")
	f.SetCellValue("Sheet1", "B3", "Values")

	// Add sample data
	for i := 4; i <= 10; i++ {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), fmt.Sprintf("Point %d", i-3))
		f.SetCellFormula("Sheet1", fmt.Sprintf("B%d", i), fmt.Sprintf("=RAND()*100"))
	}

	f.SetCellValue("Sheet1", "A12", "Average:")
	f.SetCellFormula("Sheet1", "B12", "=AVERAGE(B4:B10)")
	f.SetCellValue("Sheet1", "A13", "Median:")
	f.SetCellFormula("Sheet1", "B13", "=MEDIAN(B4:B10)")
	f.SetCellValue("Sheet1", "A14", "Standard Dev:")
	f.SetCellFormula("Sheet1", "B14", "=STDEV(B4:B10)")
	f.SetCellValue("Sheet1", "A15", "Min:")
	f.SetCellFormula("Sheet1", "B15", "=MIN(B4:B10)")
	f.SetCellValue("Sheet1", "A16", "Max:")
	f.SetCellFormula("Sheet1", "B16", "=MAX(B4:B10)")

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
