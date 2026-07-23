package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpreadsheetMathService_Power(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"power(2, 3)", []float64{2, 3}, 8, false},
		{"power(5, 2)", []float64{5, 2}, 25, false},
		{"power(10, 0)", []float64{10, 0}, 1, false},
		{"power(2, -1)", []float64{2, -1}, 0.5, false},
		{"power(0, 5)", []float64{0, 5}, 0, false},
		{"power(1, 10)", []float64{1, 10}, 1, false},
		{"too many args", []float64{2, 3, 4}, 8, false}, // Should ignore extra args
		{"too few args", []float64{2}, 0, false},        // Should handle gracefully
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.power(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Sqrt(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"sqrt(4)", []float64{4}, 2, false},
		{"sqrt(9)", []float64{9}, 3, false},
		{"sqrt(2)", []float64{2}, 1.414213562373095, false},
		{"sqrt(0)", []float64{0}, 0, false},
		{"sqrt(1)", []float64{1}, 1, false},
		{"too many args", []float64{4, 2}, 2, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.sqrt(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Abs(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"abs(5)", []float64{5}, 5, false},
		{"abs(-5)", []float64{-5}, 5, false},
		{"abs(0)", []float64{0}, 0, false},
		{"abs(-3.14)", []float64{-3.14}, 3.14, false},
		{"abs(2.5)", []float64{2.5}, 2.5, false},
		{"too many args", []float64{5, 3}, 5, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.abs(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Round(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"round(3.14)", []float64{3.14}, 3, false},
		{"round(3.5)", []float64{3.5}, 4, false},
		{"round(2.7)", []float64{2.7}, 3, false},
		{"round(3.14159, 2)", []float64{3.14159, 2}, 3.14, false},
		{"round(3.14159, 3)", []float64{3.14159, 3}, 3.142, false},
		{"round(1.23456, 0)", []float64{1.23456, 0}, 1, false},
		{"round(-2.7)", []float64{-2.7}, -3, false},
		{"round(-2.3)", []float64{-2.3}, -2, false},
		{"too many args", []float64{3.14, 2, 1}, 3.14, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.round(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Ceiling(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"ceiling(3.14)", []float64{3.14}, 4, false},
		{"ceiling(3.0)", []float64{3.0}, 3, false},
		{"ceiling(-2.7)", []float64{-2.7}, -2, false},
		{"ceiling(-2.1)", []float64{-2.1}, -2, false},
		{"ceiling(5.9)", []float64{5.9}, 6, false},
		{"too many args", []float64{3.14, 2}, 4, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.ceiling(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSpreadsheetMathService_Floor(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"floor(3.14)", []float64{3.14}, 3, false},
		{"floor(3.0)", []float64{3.0}, 3, false},
		{"floor(-2.7)", []float64{-2.7}, -3, false},
		{"floor(-2.1)", []float64{-2.1}, -3, false},
		{"floor(5.9)", []float64{5.9}, 5, false},
		{"too many args", []float64{3.14, 2}, 3, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.floor(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
