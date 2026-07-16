package services

import (
	"database/sql"
	"fmt"
	"strings"
)

// QueryBuilderService provides visual query building and execution
type QueryBuilderService struct {
	db *sql.DB
}

// QueryElement represents an element in the visual query builder
type QueryElement struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`      // "table", "field", "filter", "join", "sort", "group"
	Table    string                 `json:"table,omitempty"`
	Field    string                 `json:"field,omitempty"`
	Alias    string                 `json:"alias,omitempty"`
	Operator string                 `json:"operator,omitempty"` // "=", "!=", ">", "<", "LIKE", "IN", etc.
	Value    interface{}            `json:"value,omitempty"`
	Position Position               `json:"position"`
	Children []QueryElement         `json:"children,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// Position represents element position in visual builder
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// QueryResult represents the result of executing a visual query
type QueryResult struct {
	Columns []QueryColumn     `json:"columns"`
	Rows    [][]interface{}   `json:"rows"`
	SQL     string            `json:"sql"`
	Count   int               `json:"count"`
}

// QueryColumn represents a column in query results
type QueryColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Table    string `json:"table"`
	Alias    string `json:"alias,omitempty"`
}

// JoinDefinition represents a join between tables
type JoinDefinition struct {
	LeftTable   string `json:"left_table"`
	RightTable  string `json:"right_table"`
	LeftField   string `json:"left_field"`
	RightField  string `json:"right_field"`
	JoinType    string `json:"join_type"` // "INNER", "LEFT", "RIGHT", "FULL OUTER"
}

// FilterCondition represents a filter condition
type FilterCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Logic    string      `json:"logic,omitempty"` // "AND", "OR"
}

// SortDefinition represents a sort operation
type SortDefinition struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "ASC", "DESC"
}

// GroupDefinition represents a group operation
type GroupDefinition struct {
	Field string `json:"field"`
}

// NewQueryBuilderService creates a new query builder service
func NewQueryBuilderService(db *sql.DB) *QueryBuilderService {
	return &QueryBuilderService{db: db}
}

// BuildSQLFromElements converts visual query elements to SQL
func (qbs *QueryBuilderService) BuildSQLFromElements(elements []QueryElement) (string, error) {
	query := &QueryContext{}

	// Process elements to build query context
	for _, element := range elements {
		switch element.Type {
		case "table":
			query.Tables = append(query.Tables, element.Table)
			if element.Alias != "" {
				query.TableAliases[element.Table] = element.Alias
			}
		case "field":
			field := QueryField{
				Table: element.Table,
				Field: element.Field,
				Alias: element.Alias,
			}
			query.Fields = append(query.Fields, field)
		case "join":
			if joinData, ok := element.Data["join"].(map[string]interface{}); ok {
				join := JoinDefinition{
					LeftTable:  joinData["left_table"].(string),
					RightTable: joinData["right_table"].(string),
					LeftField:  joinData["left_field"].(string),
					RightField: joinData["right_field"].(string),
					JoinType:   joinData["join_type"].(string),
				}
				query.Joins = append(query.Joins, join)
			}
		case "filter":
			filter := FilterCondition{
				Field:    element.Field,
				Operator: element.Operator,
				Value:    element.Value,
				Logic:    "AND", // Default
			}
			if logic, ok := element.Data["logic"].(string); ok {
				filter.Logic = logic
			}
			query.Filters = append(query.Filters, filter)
		case "sort":
			sort := SortDefinition{
				Field:     element.Field,
				Direction: "ASC", // Default
			}
			if direction, ok := element.Data["direction"].(string); ok {
				sort.Direction = direction
			}
			query.Sorts = append(query.Sorts, sort)
		case "group":
			group := GroupDefinition{
				Field: element.Field,
			}
			query.Groups = append(query.Groups, group)
		}
	}

	// Generate SQL
	return qbs.buildSQL(query)
}

// QueryContext holds the complete query information
type QueryContext struct {
	Tables        []string
	TableAliases  map[string]string
	Fields        []QueryField
	Joins         []JoinDefinition
	Filters       []FilterCondition
	Sorts         []SortDefinition
	Groups        []GroupDefinition
	Limit         int
	Offset        int
}

// QueryField represents a field in a query
type QueryField struct {
	Table string
	Field string
	Alias string
}

// buildSQL constructs the SQL query from context
func (qbs *QueryBuilderService) buildSQL(ctx *QueryContext) (string, error) {
	if len(ctx.Tables) == 0 {
		return "", fmt.Errorf("no tables specified")
	}

	var parts []string

	// SELECT clause
	selectClause := "SELECT "
	if len(ctx.Fields) == 0 {
		selectClause += "*"
	} else {
		var fieldStrs []string
		for _, field := range ctx.Fields {
			fieldStr := field.Field
			if field.Table != "" {
				fieldStr = field.Table + "." + fieldStr
			}
			if field.Alias != "" {
				fieldStr += " AS " + field.Alias
			}
			fieldStrs = append(fieldStrs, fieldStr)
		}
		selectClause += strings.Join(fieldStrs, ", ")
	}
	parts = append(parts, selectClause)

	// FROM clause
	fromClause := "FROM " + ctx.Tables[0]
	if alias, exists := ctx.TableAliases[ctx.Tables[0]]; exists {
		fromClause += " " + alias
	}
	parts = append(parts, fromClause)

	// JOIN clauses
	for _, join := range ctx.Joins {
		joinClause := fmt.Sprintf("%s JOIN %s ON %s.%s = %s.%s",
			join.JoinType, join.RightTable,
			join.LeftTable, join.LeftField,
			join.RightTable, join.RightField)
		parts = append(parts, joinClause)
	}

	// WHERE clause
	if len(ctx.Filters) > 0 {
		whereClause := "WHERE "
		var conditions []string
		for i, filter := range ctx.Filters {
			condition := qbs.buildFilterCondition(filter)
			conditions = append(conditions, condition)

			// Add logic operator if not the last condition
			if i < len(ctx.Filters)-1 {
				conditions = append(conditions, filter.Logic)
			}
		}
		whereClause += strings.Join(conditions, " ")
		parts = append(parts, whereClause)
	}

	// GROUP BY clause
	if len(ctx.Groups) > 0 {
		groupClause := "GROUP BY "
		var groupFields []string
		for _, group := range ctx.Groups {
			groupFields = append(groupFields, group.Field)
		}
		groupClause += strings.Join(groupFields, ", ")
		parts = append(parts, groupClause)
	}

	// ORDER BY clause
	if len(ctx.Sorts) > 0 {
		orderClause := "ORDER BY "
		var sortFields []string
		for _, sort := range ctx.Sorts {
			sortFields = append(sortFields, sort.Field+" "+sort.Direction)
		}
		orderClause += strings.Join(sortFields, ", ")
		parts = append(parts, orderClause)
	}

	// LIMIT/OFFSET clauses
	if ctx.Limit > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", ctx.Limit))
		if ctx.Offset > 0 {
			parts = append(parts, fmt.Sprintf("OFFSET %d", ctx.Offset))
		}
	}

	return strings.Join(parts, " "), nil
}

// buildFilterCondition builds a filter condition string
func (qbs *QueryBuilderService) buildFilterCondition(filter FilterCondition) string {
	field := filter.Field

	switch filter.Operator {
	case "equals", "=":
		return fmt.Sprintf("%s = %s", field, qbs.formatValue(filter.Value))
	case "not_equals", "!=":
		return fmt.Sprintf("%s != %s", field, qbs.formatValue(filter.Value))
	case "greater_than", ">":
		return fmt.Sprintf("%s > %s", field, qbs.formatValue(filter.Value))
	case "less_than", "<":
		return fmt.Sprintf("%s < %s", field, qbs.formatValue(filter.Value))
	case "greater_equal", ">=":
		return fmt.Sprintf("%s >= %s", field, qbs.formatValue(filter.Value))
	case "less_equal", "<=":
		return fmt.Sprintf("%s <= %s", field, qbs.formatValue(filter.Value))
	case "like":
		return fmt.Sprintf("%s LIKE %s", field, qbs.formatValue("%"+fmt.Sprintf("%v", filter.Value)+"%"))
	case "starts_with":
		return fmt.Sprintf("%s LIKE %s", field, qbs.formatValue(fmt.Sprintf("%v", filter.Value)+"%"))
	case "ends_with":
		return fmt.Sprintf("%s LIKE %s", field, qbs.formatValue("%"+fmt.Sprintf("%v", filter.Value)))
	case "contains":
		return fmt.Sprintf("%s LIKE %s", field, qbs.formatValue("%"+fmt.Sprintf("%v", filter.Value)+"%"))
	case "in":
		if values, ok := filter.Value.([]interface{}); ok {
			var valueStrs []string
			for _, v := range values {
				valueStrs = append(valueStrs, qbs.formatValue(v))
			}
			return fmt.Sprintf("%s IN (%s)", field, strings.Join(valueStrs, ", "))
		}
	case "is_null":
		return fmt.Sprintf("%s IS NULL", field)
	case "is_not_null":
		return fmt.Sprintf("%s IS NOT NULL", field)
	}

	return fmt.Sprintf("%s = %s", field, qbs.formatValue(filter.Value))
}

// formatValue formats a value for SQL
func (qbs *QueryBuilderService) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// ExecuteQuery executes a visual query and returns results
func (qbs *QueryBuilderService) ExecuteQuery(elements []QueryElement) (*QueryResult, error) {
	sql, err := qbs.BuildSQLFromElements(elements)
	if err != nil {
		return nil, err
	}

	rows, err := qbs.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column information
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Get column types (simplified)
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var queryColumns []QueryColumn
	for i, col := range columns {
		queryCol := QueryColumn{
			Name: col,
			Type: columnTypes[i].DatabaseTypeName(),
		}
		queryColumns = append(queryColumns, queryCol)
	}

	// Read all rows
	var resultRows [][]interface{}
	count := 0
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		resultRows = append(resultRows, values)
		count++
	}

	return &QueryResult{
		Columns: queryColumns,
		Rows:    resultRows,
		SQL:     sql,
		Count:   count,
	}, nil
}

// GetAvailableTables returns list of available tables for query building
func (qbs *QueryBuilderService) GetAvailableTables() ([]map[string]interface{}, error) {
	// Get all user tables (simplified - in real implementation, filter by user permissions)
	rows, err := qbs.db.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name NOT LIKE 'pg_%'
		AND table_name NOT LIKE 'sql_%'
		ORDER BY table_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []map[string]interface{}
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			continue
		}

		// Get table info
		tableInfo := map[string]interface{}{
			"name":    tableName,
			"schema":  "public",
			"type":    "table",
		}

		// Get column information
		columns, err := qbs.getTableColumns(tableName)
		if err == nil {
			tableInfo["columns"] = columns
		}

		tables = append(tables, tableInfo)
	}

	return tables, nil
}

// getTableColumns returns column information for a table
func (qbs *QueryBuilderService) getTableColumns(tableName string) ([]map[string]interface{}, error) {
	rows, err := qbs.db.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = $1 AND table_schema = 'public'
		ORDER BY ordinal_position
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []map[string]interface{}
	for rows.Next() {
		var colName, dataType string
		var isNullable, columnDefault sql.NullString

		err := rows.Scan(&colName, &dataType, &isNullable, &columnDefault)
		if err != nil {
			continue
		}

		column := map[string]interface{}{
			"name":         colName,
			"type":         dataType,
			"nullable":     isNullable.String == "YES",
			"primary_key":  false, // Would need additional query
			"foreign_key":  false, // Would need additional query
		}

		if columnDefault.Valid {
			column["default"] = columnDefault.String
		}

		columns = append(columns, column)
	}

	return columns, nil
}

// ValidateQuery validates a visual query before execution
func (qbs *QueryBuilderService) ValidateQuery(elements []QueryElement) []string {
	var errors []string

	tableCount := 0
	fieldCount := 0

	for _, element := range elements {
		switch element.Type {
		case "table":
			tableCount++
		case "field":
			fieldCount++
		case "join":
			// Validate join has required fields
			if joinData, ok := element.Data["join"].(map[string]interface{}); ok {
				if joinData["left_table"] == nil || joinData["right_table"] == nil ||
				   joinData["left_field"] == nil || joinData["right_field"] == nil {
					errors = append(errors, "Join is missing required fields")
				}
			} else {
				errors = append(errors, "Invalid join configuration")
			}
		case "filter":
			if element.Field == "" || element.Operator == "" {
				errors = append(errors, "Filter is missing required fields")
			}
		}
	}

	if tableCount == 0 {
		errors = append(errors, "Query must include at least one table")
	}

	return errors
}

// GetQuerySuggestions provides suggestions for query building
func (qbs *QueryBuilderService) GetQuerySuggestions(currentElements []QueryElement) []map[string]interface{} {
	var suggestions []map[string]interface{}

	// Analyze current query and suggest improvements
	hasTables := false
	hasFilters := false
	hasSorts := false

	for _, element := range currentElements {
		switch element.Type {
		case "table":
			hasTables = true
		case "filter":
			hasFilters = true
		case "sort":
			hasSorts = true
		}
	}

	if hasTables && !hasFilters {
		suggestions = append(suggestions, map[string]interface{}{
			"type":        "filter",
			"description": "Add filters to narrow down your results",
			"priority":    "high",
		})
	}

	if hasTables && !hasSorts {
		suggestions = append(suggestions, map[string]interface{}{
			"type":        "sort",
			"description": "Add sorting to organize your results",
			"priority":    "medium",
		})
	}

	if len(currentElements) > 5 {
		suggestions = append(suggestions, map[string]interface{}{
			"type":        "group",
			"description": "Consider grouping results for summary data",
			"priority":    "low",
		})
	}

	return suggestions
}

// SaveQuery saves a visual query for later use
func (qbs *QueryBuilderService) SaveQuery(userID string, name string, description string, elements []QueryElement) error {
	sql, err := qbs.BuildSQLFromElements(elements)
	if err != nil {
		return err
	}

	// In a real implementation, save to database
	// For now, just validate
	if sql == "" {
		return fmt.Errorf("invalid query")
	}

	return nil
}

// LoadSavedQueries loads saved queries for a user
func (qbs *QueryBuilderService) LoadSavedQueries(userID string) ([]map[string]interface{}, error) {
	// In a real implementation, load from database
	// For now, return empty list
	return []map[string]interface{}{}, nil
}

// ExportQuery exports a visual query to different formats
func (qbs *QueryBuilderService) ExportQuery(elements []QueryElement, format string) (string, error) {
	sql, err := qbs.BuildSQLFromElements(elements)
	if err != nil {
		return "", err
	}

	switch format {
	case "sql":
		return sql, nil
	case "json":
		// Would serialize elements to JSON
		return "", fmt.Errorf("JSON export not implemented")
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}
