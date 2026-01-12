package services

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpreadsheetMathService_Sin(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"sin(0)", []float64{0}, 0, false},
		{"sin(π/2)", []float64{math.Pi / 2}, 1, false},
		{"sin(π)", []float64{math.Pi}, 0, false},
		{"sin(π/4)", []float64{math.Pi / 4}, 0.7071067811865476, false},
		{"too many args", []float64{1, 2}, 0, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.sin(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Cos(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"cos(0)", []float64{0}, 1, false},
		{"cos(π/2)", []float64{math.Pi / 2}, 0, false},
		{"cos(π)", []float64{math.Pi}, -1, false},
		{"cos(π/4)", []float64{math.Pi / 4}, 0.7071067811865476, false},
		{"too many args", []float64{1, 2}, 0, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.cos(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Tan(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"tan(0)", []float64{0}, 0, false},
		{"tan(π/4)", []float64{math.Pi / 4}, 1, false},
		{"tan(π/3)", []float64{math.Pi / 3}, 1.732050807568877, false},
		{"too many args", []float64{1, 2}, 0, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.tan(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Asin(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"asin(0)", []float64{0}, 0, false},
		{"asin(0.5)", []float64{0.5}, math.Pi/6, false},
		{"asin(1)", []float64{1}, math.Pi/2, false},
		{"asin(-1)", []float64{-1}, -math.Pi/2, false},
		{"too many args", []float64{1, 2}, 0, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.asin(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Acos(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"acos(1)", []float64{1}, 0, false},
		{"acos(0.5)", []float64{0.5}, math.Pi/3, false},
		{"acos(0)", []float64{0}, math.Pi/2, false},
		{"acos(-1)", []float64{-1}, math.Pi, false},
		{"too many args", []float64{1, 2}, 0, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.acos(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}

func TestSpreadsheetMathService_Atan(t *testing.T) {
	sms := NewSpreadsheetMathService()

	tests := []struct {
		name     string
		args     []float64
		expected float64
		hasError bool
	}{
		{"atan(0)", []float64{0}, 0, false},
		{"atan(1)", []float64{1}, math.Pi/4, false},
		{"atan(-1)", []float64{-1}, -math.Pi/4, false},
		{"atan(∞)", []float64{math.Inf(1)}, math.Pi/2, false},
		{"too many args", []float64{1, 2}, 0, false}, // Should ignore extra args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sms.atan(tt.args)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-10)
			}
		})
	}
}
