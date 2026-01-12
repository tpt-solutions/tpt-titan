package services

import "math"

// Rounding functions

func (sms *SpreadsheetMathService) round(args []float64) (float64, error) {
	if len(args) == 2 {
		precision := args[1]
		multiplier := math.Pow(10, precision)
		return math.Round(args[0]*multiplier) / multiplier, nil
	}
	return math.Round(args[0]), nil
}

func (sms *SpreadsheetMathService) ceiling(args []float64) (float64, error) {
	return math.Ceil(args[0]), nil
}

func (sms *SpreadsheetMathService) floor(args []float64) (float64, error) {
	return math.Floor(args[0]), nil
}
