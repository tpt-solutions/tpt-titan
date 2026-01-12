package services

// Conditional functions

func (sms *SpreadsheetMathService) ifFunction(args []float64) (float64, error) {
	condition := args[0] != 0
	if condition {
		return args[1], nil
	}
	if len(args) > 2 {
		return args[2], nil
	}
	return 0, nil
}

func (sms *SpreadsheetMathService) andFunction(args []float64) (float64, error) {
	for _, arg := range args {
		if arg == 0 {
			return 0, nil
		}
	}
	return 1, nil
}

func (sms *SpreadsheetMathService) orFunction(args []float64) (float64, error) {
	for _, arg := range args {
		if arg != 0 {
			return 1, nil
		}
	}
	return 0, nil
}

func (sms *SpreadsheetMathService) notFunction(args []float64) (float64, error) {
	if args[0] == 0 {
		return 1, nil
	}
	return 0, nil
}
