package services

import "math"

// Statistical functions

func (sms *SpreadsheetMathService) stdev(args []float64) (float64, error) {
	if len(args) < 2 {
		return 0, fmt.Errorf("stdev requires at least two arguments")
	}

	mean, _ := sms.average(args)
	sumSquares := 0.0
	for _, arg := range args {
		diff := arg - mean
		sumSquares += diff * diff
	}

	return math.Sqrt(sumSquares / float64(len(args)-1)), nil
}

func (sms *SpreadsheetMathService) median(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("median requires at least one argument")
	}

	// Simple implementation - sort and find middle
	sorted := make([]float64, len(args))
	copy(sorted, args)

	// Bubble sort for simplicity
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2, nil
	}
	return sorted[mid], nil
}

func (sms *SpreadsheetMathService) mode(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("mode requires at least one argument")
	}

	counts := make(map[float64]int)
	for _, arg := range args {
		counts[arg]++
	}

	var mode float64
	maxCount := 0
	for value, count := range counts {
		if count > maxCount {
			maxCount = count
			mode = value
		}
	}

	return mode, nil
}
