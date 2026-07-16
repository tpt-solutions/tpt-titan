package services

import "math"

// Trigonometric functions

func (sms *SpreadsheetMathService) sin(args []float64) (float64, error) {
	if len(args) < 1 {
		return 0, nil
	}
	return math.Sin(args[0]), nil
}

func (sms *SpreadsheetMathService) cos(args []float64) (float64, error) {
	if len(args) < 1 {
		return 0, nil
	}
	return math.Cos(args[0]), nil
}

func (sms *SpreadsheetMathService) tan(args []float64) (float64, error) {
	if len(args) < 1 {
		return 0, nil
	}
	return math.Tan(args[0]), nil
}

func (sms *SpreadsheetMathService) asin(args []float64) (float64, error) {
	if len(args) < 1 {
		return 0, nil
	}
	return math.Asin(args[0]), nil
}

func (sms *SpreadsheetMathService) acos(args []float64) (float64, error) {
	if len(args) < 1 {
		return 0, nil
	}
	return math.Acos(args[0]), nil
}

func (sms *SpreadsheetMathService) atan(args []float64) (float64, error) {
	if len(args) < 1 {
		return 0, nil
	}
	return math.Atan(args[0]), nil
}
