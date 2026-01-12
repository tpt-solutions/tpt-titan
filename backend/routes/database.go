package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DatabaseTableInfo represents metadata about a database table
type DatabaseTableInfo struct {
	Name        string                 `json:"name"`
	Columns     []DatabaseColumnInfo   `json:"columns"`
	PrimaryKey  string                 `json:"primary_key"`
	Constraints []DatabaseConstraint  `json:"constraints"`
	Relationships []DatabaseRelationship `json:"relationships"`
}

// DatabaseColumnInfo represents column metadata
type DatabaseColumnInfo struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Nullable     bool   `json:"nullable"`
	DefaultValue string `json:"default_value,omitempty"`
	MaxLength    int    `json:"max_length,omitempty"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	IsForeignKey bool   `json:"is_foreign_key"`
	ForeignTable string `json:"foreign_table,omitempty"`
	ForeignColumn string `json:"foreign_column,omitempty"`
}

// DatabaseConstraint represents table constraints
type DatabaseConstraint struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"` // 'primary_key', 'foreign_key', 'unique', 'check'
	Columns    []string `json:"columns"`
	References string   `json:"references,omitempty"`
}

// DatabaseRelationship represents foreign key relationships
type DatabaseRelationship struct {
	Name           string `json:"name"`
	SourceTable    string `json:"source_table"`
	SourceColumn   string `json:"source_column"`
	TargetTable    string `json:"target_table"`
	TargetColumn   string `json:"target_column"`
	RelationshipType string `json:"relationship_type"`
}

// DatabaseRecord represents a single record from a table
type DatabaseRecord struct {
	ID     interface{}            `json:"id"`
	Values map[string]interface{} `json:"values"`
}

// GetDatabaseTables returns list of available tables
func GetDatabaseTables(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Get all table names from information_schema
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		AND table_name NOT LIKE 'pg_%'
		AND table_name NOT LIKE 'sql_%'
		ORDER BY table_name
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		tables = append(tables, tableName)
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}

// GetTableInfo returns metadata for a specific table
func GetTableInfo(c *gin.Context) {
	tableName := c.Param("table")
	db := c.MustGet("db").(*sql.DB)

	tableInfo, err := getTableMetadata(db, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get table info: %v", err)})
		return
	}

	c.JSON(http.StatusOK, tableInfo)
}

// GetTableData returns paginated data from a table
func GetTableData(c *gin.Context) {
	tableName := c.Param("table")
	db := c.MustGet("db").(*sql.DB)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get record count"})
		return
	}

	// Get table data
	dataQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY 1 LIMIT $1 OFFSET $2", tableName)
	rows, err := db.Query(dataQuery, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch table data"})
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get column names"})
		return
	}

	var records []DatabaseRecord
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		recordValues := make(map[string]interface{})
		for i, col := range columns {
			recordValues[col] = values[i]
		}

		record := DatabaseRecord{
			Values: recordValues,
		}

		// Try to set ID from first column or primary key
		if len(values) > 0 {
			record.ID = values[0]
		}

		records = append(records, record)
	}

	c.JSON(http.StatusOK, gin.H{
		"records":    records,
		"columns":    columns,
		"total":      totalCount,
		"page":       page,
		"limit":      limit,
		"totalPages": (totalCount + limit - 1) / limit,
	})
}

// UpdateTableRecord updates a single record in a table
func UpdateTableRecord(c *gin.Context) {
	tableName := c.Param("table")
	recordID := c.Param("id")
	db := c.MustGet("db").(*sql.DB)

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the update against constraints
	validationErrors := validateRecordUpdate(db, tableName, updateData)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	// Build update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	for column, value := range updateData {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Get primary key info to build WHERE clause
	tableInfo, err := getTableMetadata(db, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get table metadata"})
		return
	}

	whereClause := fmt.Sprintf("%s = $%d", tableInfo.PrimaryKey, argIndex)
	args = append(args, recordID)

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName, strings.Join(setParts, ", "), whereClause)

	_, err = db.Exec(updateQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update record: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

// CreateTableRecord creates a new record in a table
func CreateTableRecord(c *gin.Context) {
	tableName := c.Param("table")
	db := c.MustGet("db").(*sql.DB)

	var recordData map[string]interface{}
	if err := c.ShouldBindJSON(&recordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the record against constraints
	validationErrors := validateRecordUpdate(db, tableName, recordData)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	// Build insert query
	columns := []string{}
	placeholders := []string{}
	args := []interface{}{}
	argIndex := 1

	for column, value := range recordData {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
		args = append(args, value)
		argIndex++
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	_, err := db.Exec(insertQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create record: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Record created successfully"})
}

// DeleteTableRecord deletes a record from a table
func DeleteTableRecord(c *gin.Context) {
	tableName := c.Param("table")
	recordID := c.Param("id")
	db := c.MustGet("db").(*sql.DB)

	// Get primary key info
	tableInfo, err := getTableMetadata(db, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get table metadata"})
		return
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, tableInfo.PrimaryKey)
	_, err = db.Exec(deleteQuery, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete record: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Helper functions

func getTableMetadata(db *sql.DB, tableName string) (*DatabaseTableInfo, error) {
	tableInfo := &DatabaseTableInfo{Name: tableName}

	// Get column information
	columnQuery := `
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable = 'YES' as nullable,
			c.column_default,
			c.character_maximum_length,
			tc.constraint_type,
			kcu.column_name as pk_column
		FROM information_schema.columns c
		LEFT JOIN information_schema.key_column_usage kcu
			ON c.table_name = kcu.table_name
			AND c.column_name = kcu.column_name
			AND kcu.constraint_name LIKE '%_pkey'
		LEFT JOIN information_schema.table_constraints tc
			ON kcu.constraint_name = tc.constraint_name
		WHERE c.table_schema = 'public'
		AND c.table_name = $1
		ORDER BY c.ordinal_position
	`

	rows, err := db.Query(columnQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var col DatabaseColumnInfo
		var constraintType, pkColumn sql.NullString
		var maxLength sql.NullInt64

		err := rows.Scan(
			&col.Name, &col.Type, &col.Nullable, &col.DefaultValue,
			&maxLength, &constraintType, &pkColumn,
		)
		if err != nil {
			continue
		}

		if maxLength.Valid {
			col.MaxLength = int(maxLength.Int64)
		}

		if pkColumn.Valid && pkColumn.String == col.Name {
			col.IsPrimaryKey = true
			tableInfo.PrimaryKey = col.Name
		}

		tableInfo.Columns = append(tableInfo.Columns, col)
	}

	// Get foreign key relationships
	fkQuery := `
		SELECT
			tc.constraint_name,
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		AND tc.table_name = $1
	`

	fkRows, err := db.Query(fkQuery, tableName)
	if err == nil {
		for fkRows.Next() {
			var constraintName, columnName, foreignTable, foreignColumn string
			fkRows.Scan(&constraintName, &columnName, &foreignTable, &foreignColumn)

			// Mark column as foreign key
			for i := range tableInfo.Columns {
				if tableInfo.Columns[i].Name == columnName {
					tableInfo.Columns[i].IsForeignKey = true
					tableInfo.Columns[i].ForeignTable = foreignTable
					tableInfo.Columns[i].ForeignColumn = foreignColumn
					break
				}
			}

			// Add relationship
			relationship := DatabaseRelationship{
				Name:            constraintName,
				SourceTable:     tableName,
				SourceColumn:    columnName,
				TargetTable:     foreignTable,
				TargetColumn:    foreignColumn,
				RelationshipType: "many-to-one",
			}
			tableInfo.Relationships = append(tableInfo.Relationships, relationship)
		}
		fkRows.Close()
	}

	return tableInfo, nil
}

func validateRecordUpdate(db *sql.DB, tableName string, data map[string]interface{}) []string {
	var errors []string

	tableInfo, err := getTableMetadata(db, tableName)
	if err != nil {
		errors = append(errors, "Failed to get table metadata")
		return errors
	}

	for _, column := range tableInfo.Columns {
		if value, exists := data[column.Name]; exists && value != nil {
			// Check data type compatibility
			if err := validateDataType(value, column.Type); err != nil {
				errors = append(errors, fmt.Sprintf("Column %s: %v", column.Name, err))
			}

			// Check foreign key constraints
			if column.IsForeignKey {
				if !foreignKeyExists(db, column.ForeignTable, column.ForeignColumn, value) {
					errors = append(errors, fmt.Sprintf("Foreign key constraint failed for %s", column.Name))
				}
			}

			// Check string length constraints
			if column.MaxLength > 0 {
				if str, ok := value.(string); ok && len(str) > column.MaxLength {
					errors = append(errors, fmt.Sprintf("Column %s exceeds maximum length of %d", column.Name, column.MaxLength))
				}
			}
		} else if !column.Nullable {
			// Check required fields
			errors = append(errors, fmt.Sprintf("Column %s is required", column.Name))
		}
	}

	return errors
}

func validateDataType(value interface{}, dataType string) error {
	switch dataType {
	case "integer", "bigint", "smallint":
		switch value.(type) {
		case int, int32, int64, float64:
			return nil
		default:
			return fmt.Errorf("expected integer, got %T", value)
		}
	case "numeric", "decimal", "real", "double precision":
		switch value.(type) {
		case float64, int, int32, int64:
			return nil
		default:
			return fmt.Errorf("expected number, got %T", value)
		}
	case "boolean":
		switch value.(type) {
		case bool:
			return nil
		default:
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "character varying", "varchar", "text":
		switch value.(type) {
		case string:
			return nil
		default:
			return fmt.Errorf("expected string, got %T", value)
		}
	case "uuid":
		switch value.(type) {
		case string:
			if _, err := uuid.Parse(value.(string)); err != nil {
				return fmt.Errorf("invalid UUID format")
			}
			return nil
		default:
			return fmt.Errorf("expected UUID string, got %T", value)
		}
	default:
		return nil // Allow unknown types
	}
}

func foreignKeyExists(db *sql.DB, table, column string, value interface{}) bool {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", table, column)
	var exists bool
	db.QueryRow(query, value).Scan(&exists)
	return exists
}
