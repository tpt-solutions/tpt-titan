package services

import "math"

// Exponential and logarithmic functions

func (sms *SpreadsheetMathService) exp(args []float64) (float64, error) {
	return math.Exp(args[0]), nil
}

func (sms *SpreadsheetMathService) ln(args []float64) (float64, error) {
	return math.Log(args[0]), nil
}

func (sms *SpreadsheetMathService) log(args []float64) (float64, error) {
	if len(args) == 2 {
		return math.Log(args[0]) / math.Log(args[1]), nil
	}
	return math.Log10(args[0]), nil
}

func (sms *SpreadsheetMathService) log10(args []float64) (float64, error) {
	return math.Log10(args[0]), nil
}
