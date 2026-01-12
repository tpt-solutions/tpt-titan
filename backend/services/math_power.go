package services

import "math"

// Power and root functions

func (sms *SpreadsheetMathService) sqrt(args []float64) (float64, error) {
	return math.Sqrt(args[0]), nil
}

func (sms *SpreadsheetMathService) power(args []float64) (float64, error) {
	return math.Pow(args[0], args[1]), nil
}

func (sms *SpreadsheetMathService) abs(args []float64) (float64, error) {
	return math.Abs(args[0]), nil
}
