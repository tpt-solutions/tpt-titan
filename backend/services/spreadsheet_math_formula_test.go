package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpreadsheetMathService_EvaluateFormula(t *testing.T) {
	sms := NewSpreadsheetMathService()

	// Mock cell resolver function
	cellResolver := func(ref string) (interface{}, error) {
		switch ref {
		case "A1":
			return 10.0, nil
		case "A2":
			return 20.0, nil
		case "A3":
			return 30.0, nil
		case "B1":
			return 5.0, nil
		case "B2":
			return 15.0, nil
		default:
			return 0.0, nil
		}
	}

	tests := []struct {
		name          string
		formula       string
		expectedValue interface{}
		expectedType  string
		expectedError bool
	}{
		{
			name:          "Simple number",
			formula:       "42",
			expectedValue: "42",
			expectedType:  "string",
		},
		{
			name:          "Basic sum formula",
			formula:       "=SUM(1,2,3)",
			expectedValue: 6.0,
			expectedType:  "number",
		},
		{
			name:          "Sum with cell references",
			formula:       "=SUM(A1,A2,A3)",
			expectedValue: 60.0,
			expectedType:  "number",
		},
		{
			name:          "Average formula",
			formula:       "=AVERAGE(10,20,30)",
			expectedValue: 20.0,
			expectedType:  "number",
		},
		{
			name:          "Arithmetic expression",
			formula:       "=A1+A2*2",
			expectedValue: 50.0,
			expectedType:  "number",
		},
		{
			name:          "Nested functions",
			formula:       "=SUM(A1,SUM(B1,B2))",
			expectedValue: 30.0,
			expectedType:  "number",
		},
		{
			name:          "Trigonometric function",
			formula:       "=SIN(0)",
			expectedValue: 0.0,
			expectedType:  "number",
		},
		{
			name:          "Power function",
			formula:       "=POWER(2,3)",
			expectedValue: 8.0,
			expectedType:  "number",
		},
		{
			name:          "Rounding function",
			formula:       "=ROUND(3.14159,2)",
			expectedValue: 3.14,
			expectedType:  "number",
		},
		{
			name:          "Cell reference",
			formula:       "=A1",
			expectedValue: 10.0,
			expectedType:  "number",
		},
		{
			name:          "Complex expression",
			formula:       "=SUM(A1:A3)*2+10",
			expectedValue: 130.0,
			expectedType:  "number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.EvaluateFormula(tt.formula, cellResolver)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, result.Value)
				assert.Equal(t, tt.expectedType, result.DataType)
			}
		})
	}
}

func TestSpreadsheetMathService_ParseFunctionCall(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name          string
		expression    string
		expectedFunc  string
		expectedArgs  []string
		expectedError bool
	}{
		{
			name:         "Simple function",
			expression:   "SUM(1,2,3)",
			expectedFunc: "SUM",
			expectedArgs: []string{"1", "2", "3"},
		},
		{
			name:         "Function with cell references",
			expression:   "AVERAGE(A1,B2,C3)",
			expectedFunc: "AVERAGE",
			expectedArgs: []string{"A1", "B2", "C3"},
		},
		{
			name:         "Nested functions",
			expression:   "SUM(A1,SUM(B1,B2))",
			expectedFunc: "SUM",
			expectedArgs: []string{"A1", "SUM(B1,B2)"},
		},
		{
			name:         "Function with spaces",
			expression:   "SUM( 1 , 2 , 3 )",
			expectedFunc: "SUM",
			expectedArgs: []string{"1", "2", "3"},
		},
		{
			name:          "Missing closing parenthesis",
			expression:    "SUM(1,2,3",
			expectedError: true,
		},
		{
			name:          "Missing opening parenthesis",
			expression:    "SUM1,2,3)",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcName, args, err := sms.parseFunctionCall(tt.expression)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFunc, funcName)
				assert.Equal(t, tt.expectedArgs, args)
			}
		})
	}
}

func TestSpreadsheetMathService_ParseArguments(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name         string
		argsStr      string
		expectedArgs []string
	}{
		{
			name:         "Simple arguments",
			argsStr:      "1,2,3",
			expectedArgs: []string{"1", "2", "3"},
		},
		{
			name:         "Arguments with spaces",
			argsStr:      " 1 , 2 , 3 ",
			expectedArgs: []string{"1", "2", "3"},
		},
		{
			name:         "Arguments with functions",
			argsStr:      "A1,SUM(B1,B2),C3",
			expectedArgs: []string{"A1", "SUM(B1,B2)", "C3"},
		},
		{
			name:         "Nested function arguments",
			argsStr:      "SUM(A1,A2),AVERAGE(B1,B2,B3)",
			expectedArgs: []string{"SUM(A1,A2)", "AVERAGE(B1,B2,B3)"},
		},
		{
			name:         "Empty string",
			argsStr:      "",
			expectedArgs: []string{},
		},
		{
			name:         "Single argument",
			argsStr:      "A1",
			expectedArgs: []string{"A1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := sms.parseArguments(tt.argsStr)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestSpreadsheetMathService_IsCellReference(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		ref      string
		expected bool
	}{
		{"Valid cell reference", "A1", true},
		{"Valid cell reference uppercase", "B10", true},
		{"Valid cell reference mixed case", "a1", true}, // Case-insensitive per spreadsheet convention
		{"Invalid - no number", "AA", false},
		{"Invalid - no letter", "123", false},
		{"Invalid - special characters", "A1!", false},
		{"Invalid - empty", "", false},
		{"Valid - two letters", "AA1", true},
		{"Valid - large numbers", "A1000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sms.isCellReference(tt.ref)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSpreadsheetMathService_ParseCellReference(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name          string
		ref           string
		expectedCol   int
		expectedRow   int
		expectedError bool
	}{
		{"A1", "A1", 0, 0, false},
		{"B2", "B2", 1, 1, false},
		{"Z1", "Z1", 25, 0, false},
		{"AA1", "AA1", 26, 0, false},
		{"AB2", "AB2", 27, 1, false},
		{"A100", "A100", 0, 99, false},
		{"Invalid reference", "123", 0, 0, true},
		{"Invalid reference 2", "ABC", 0, 0, true},
		{"Empty string", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.ParseCellReference(tt.ref)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCol, result.Col)
				assert.Equal(t, tt.expectedRow, result.Row)
			}
		})
	}
}

func TestSpreadsheetMathService_FormatCellReference(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		sheet    string
		col      int
		row      int
		expected string
	}{
		{"A1", "", 0, 0, "A1"},
		{"B2", "", 1, 1, "B2"},
		{"Z1", "", 25, 0, "Z1"},
		{"AA1", "", 26, 0, "AA1"},
		{"AB2", "", 27, 1, "AB2"},
		{"A100", "", 0, 99, "A100"},
		{"Sheet1!A1", "Sheet1", 0, 0, "Sheet1!A1"},
		{"MySheet!B10", "MySheet", 1, 9, "MySheet!B10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sms.FormatCellReference(tt.sheet, tt.col, tt.row)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSpreadsheetMathService_ValidateFormula(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		formula  string
		expected bool // true if valid, false if invalid
	}{
		{"Valid number", "42", true},
		{"Valid formula", "=SUM(1,2,3)", true},
		{"Valid complex formula", "=SUM(A1:A10)*2+10", true},
		{"Valid nested functions", "=SUM(A1,SUM(B1,B2))", true},
		{"Invalid - unbalanced parentheses", "=SUM(1,2,3", false},
		{"Invalid - missing function", "=UNKNOWN(1,2,3)", true}, // Unknown functions are considered valid syntax
		{"Valid - no equals", "SUM(1,2,3)", true},               // Not a formula, so valid
		{"Valid - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sms.ValidateFormula(tt.formula)

			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
