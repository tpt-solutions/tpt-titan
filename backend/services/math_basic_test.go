package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpreadsheetMathService_Sum(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"empty args", []float64{}, 0, false},
		{"single number", []float64{5}, 5, false},
		{"multiple numbers", []float64{1, 2, 3, 4}, 10, false},
		{"negative numbers", []float64{-1, 2, -3}, -2, false},
		{"decimals", []float64{1.5, 2.5, 3.0}, 7.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.sum(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSpreadsheetMathService_Average(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"empty args", []float64{}, 0, true},
		{"single number", []float64{5}, 5, false},
		{"multiple numbers", []float64{1, 2, 3, 4}, 2.5, false},
		{"negative numbers", []float64{-1, 2, -3}, -0.6666666666666666, false},
		{"decimals", []float64{1.5, 2.5, 3.0}, 2.3333333333333335, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.average(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSpreadsheetMathService_Min(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"empty args", []float64{}, 0, true},
		{"single number", []float64{5}, 5, false},
		{"multiple numbers", []float64{3, 1, 4, 2}, 1, false},
		{"negative numbers", []float64{-1, -5, 2}, -5, false},
		{"decimals", []float64{1.5, 2.5, 1.1}, 1.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.min(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSpreadsheetMathService_Max(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"empty args", []float64{}, 0, true},
		{"single number", []float64{5}, 5, false},
		{"multiple numbers", []float64{3, 1, 4, 2}, 4, false},
		{"negative numbers", []float64{-1, -5, 2}, 2, false},
		{"decimals", []float64{1.5, 2.5, 1.1}, 2.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.max(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSpreadsheetMathService_Count(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
	}{
		{"empty args", []float64{}, 0, false},
		{"single number", []float64{5}, 1, false},
		{"multiple numbers", []float64{1, 2, 3, 4}, 4, false},
		{"mixed numbers", []float64{1.5, 2, -3}, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.count(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSpreadsheetMathService_Product(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"empty args", []float64{}, 1, false}, // Product of empty set is 1
		{"single number", []float64{5}, 5, false},
		{"multiple numbers", []float64{2, 3, 4}, 24, false},
		{"with zero", []float64{2, 0, 4}, 0, false},
		{"negative numbers", []float64{2, -3, 4}, -24, false},
		{"decimals", []float64{1.5, 2, 2}, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.product(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
