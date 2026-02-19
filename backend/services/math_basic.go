package services

import "fmt"

// Basic mathematical functions

func (sms *SpreadsheetMathService) sum(args []float64) (float64, error) {
	result := 0.0
	for _, arg := range args {
		result += arg
	}
	return result, nil
}

func (sms *SpreadsheetMathService) average(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("average requires at least one argument")
	}
	sum := 0.0
	for _, arg := range args {
		sum += arg
	}
	return sum / float64(len(args)), nil
}

func (sms *SpreadsheetMathService) min(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("min requires at least one argument")
	}
	min := args[0]
	for _, arg := range args[1:] {
		if arg < min {
			min = arg
		}
	}
	return min, nil
}

func (sms *SpreadsheetMathService) max(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("max requires at least one argument")
	}
	max := args[0]
	for _, arg := range args[1:] {
		if arg > max {
			max = arg
		}
	}
	return max, nil
}

func (sms *SpreadsheetMathService) count(args []float64) (float64, error) {
	return float64(len(args)), nil
}

func (sms *SpreadsheetMathService) product(args []float64) (float64, error) {
	result := 1.0
	for _, arg := range args {
		result *= arg
	}
	return result, nil
}
