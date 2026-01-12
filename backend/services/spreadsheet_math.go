package services

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// SpreadsheetMathService provides mathematical computation and formula evaluation
type SpreadsheetMathService struct {
	functions map[string]MathFunction
}

// MathFunction represents a mathematical function
type MathFunction struct {
	Name        string
	MinArgs     int
	MaxArgs     int
	Description string
	Function    func([]float64) (float64, error)
}

// CellReference represents a cell reference (e.g., A1, B2, Sheet1!A1)
type CellReference struct {
	Sheet string
	Col   int // 0-based column index
	Row   int // 0-based row index
}

// FormulaResult represents the result of evaluating a formula
type FormulaResult struct {
	Value     interface{} `json:"value"`
	DataType  string      `json:"data_type"` // "number", "string", "boolean", "error"
	Error     string      `json:"error,omitempty"`
	DependsOn []string    `json:"depends_on,omitempty"` // Cell references this formula depends on
}

// NewSpreadsheetMathService creates a new spreadsheet math service
func NewSpreadsheetMathService() *SpreadsheetMathService {
	sms := &SpreadsheetMathService{
		functions: make(map[string]MathFunction),
	}

	sms.registerFunctions()
	return sms
}

// registerFunctions registers all built-in mathematical functions
func (sms *SpreadsheetMathService) registerFunctions() {
	functions := []MathFunction{
		// Basic math functions
		{"SUM", 1, -1, "Returns the sum of a series of numbers", sms.sum},
		{"AVERAGE", 1, -1, "Returns the average of a series of numbers", sms.average},
		{"MIN", 1, -1, "Returns the minimum value from a series of numbers", sms.min},
		{"MAX", 1, -1, "Returns the maximum value from a series of numbers", sms.max},
		{"COUNT", 1, -1, "Counts the number of cells that contain numbers", sms.count},
		{"PRODUCT", 1, -1, "Returns the product of a series of numbers", sms.product},

		// Statistical functions
		{"STDEV", 1, -1, "Estimates standard deviation based on a sample", sms.stdev},
		{"MEDIAN", 1, -1, "Returns the median of a series of numbers", sms.median},
		{"MODE", 1, -1, "Returns the most frequently occurring number", sms.mode},

		// Trigonometric functions
		{"SIN", 1, 1, "Returns the sine of an angle", sms.sin},
		{"COS", 1, 1, "Returns the cosine of an angle", sms.cos},
		{"TAN", 1, 1, "Returns the tangent of an angle", sms.tan},
		{"ASIN", 1, 1, "Returns the arcsine of a number", sms.asin},
		{"ACOS", 1, 1, "Returns the arccosine of a number", sms.acos},
		{"ATAN", 1, 1, "Returns the arctangent of a number", sms.atan},

		// Exponential and logarithmic functions
		{"EXP", 1, 1, "Returns e raised to the power of a number", sms.exp},
		{"LN", 1, 1, "Returns the natural logarithm of a number", sms.ln},
		{"LOG", 1, 2, "Returns the logarithm of a number to a specified base", sms.log},
		{"LOG10", 1, 1, "Returns the base-10 logarithm of a number", sms.log10},

		// Power and root functions
		{"SQRT", 1, 1, "Returns the square root of a number", sms.sqrt},
		{"POWER", 2, 2, "Returns a number raised to a power", sms.power},
		{"ABS", 1, 1, "Returns the absolute value of a number", sms.abs},

		// Rounding functions
		{"ROUND", 1, 2, "Rounds a number to a specified number of digits", sms.round},
		{"CEILING", 1, 2, "Rounds a number up to the nearest integer or multiple", sms.ceiling},
		{"FLOOR", 1, 2, "Rounds a number down to the nearest integer or multiple", sms.floor},

		// Conditional functions
		{"IF", 2, 3, "Returns one value if condition is true, another if false", sms.ifFunction},
		{"AND", 1, -1, "Returns TRUE if all arguments are TRUE", sms.andFunction},
		{"OR", 1, -1, "Returns TRUE if any argument is TRUE", sms.orFunction},
		{"NOT", 1, 1, "Reverses the logic of its argument", sms.notFunction},

		// Date and time functions
		{"NOW", 0, 0, "Returns the current date and time", sms.now},
		{"TODAY", 0, 0, "Returns the current date", sms.today},
		{"DATE", 3, 3, "Creates a date from year, month, and day", sms.date},
		{"DATEDIF", 3, 3, "Calculates the difference between two dates", sms.datedif},

		// Text functions
		{"LEN", 1, 1, "Returns the length of a text string", sms.len},
		{"LEFT", 2, 2, "Returns the leftmost characters from a text value", sms.left},
		{"RIGHT", 2, 2, "Returns the rightmost characters from a text value", sms.right},
		{"MID", 3, 3, "Returns a specific number of characters from a text string", sms.mid},
		{"CONCATENATE", 1, -1, "Joins several text strings into one string", sms.concatenate},
		{"UPPER", 1, 1, "Converts text to uppercase", sms.upper},
		{"LOWER", 1, 1, "Converts text to lowercase", sms.lower},

		// Lookup functions
		{"VLOOKUP", 3, 4, "Looks for a value in the leftmost column of a table", sms.vlookup},
		{"HLOOKUP", 3, 4, "Looks for a value in the top row of a table", sms.hlookup},
		{"INDEX", 2, 3, "Returns a value from a table at the intersection of a row and column", sms.index},
		{"MATCH", 2, 3, "Returns the relative position of an item in an array", sms.match},

		// Financial functions
		{"PV", 3, 5, "Returns the present value of an investment", sms.pv},
		{"FV", 3, 5, "Returns the future value of an investment", sms.fv},
		{"PMT", 3, 5, "Returns the periodic payment for an annuity", sms.pmt},
		{"NPV", 2, -1, "Returns the net present value of an investment", sms.npv},
		{"IRR", 1, 2, "Returns the internal rate of return for a series of cash flows", sms.irr},
	}

	for _, fn := range functions {
		sms.functions[strings.ToUpper(fn.Name)] = fn
	}
}

// EvaluateFormula evaluates a spreadsheet formula
func (sms *SpreadsheetMathService) EvaluateFormula(formula string, cellResolver func(string) (interface{}, error)) (*FormulaResult, error) {
	if !strings.HasPrefix(strings.ToUpper(formula), "=") {
		// Not a formula, return as string
		return &FormulaResult{
			Value:    formula,
			DataType: "string",
		}, nil
	}

	// Remove the = prefix
	expression := formula[1:]

	// Parse and evaluate the expression
	result, err := sms.evaluateExpression(expression, cellResolver)
	if err != nil {
		return &FormulaResult{
			Value:    "#ERROR!",
			DataType: "error",
			Error:    err.Error(),
		}, nil
	}

	dataType := "number"
	switch v := result.(type) {
	case string:
		dataType = "string"
	case bool:
		dataType = "boolean"
	case int, int32, int64, float32, float64:
		dataType = "number"
	}

	return &FormulaResult{
		Value:    result,
		DataType: dataType,
	}, nil
}

// evaluateExpression evaluates a mathematical expression
func (sms *SpreadsheetMathService) evaluateExpression(expr string, cellResolver func(string) (interface{}, error)) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	// Handle function calls
	if strings.Contains(expr, "(") {
		return sms.evaluateFunctionCall(expr, cellResolver)
	}

	// Handle cell references
	if sms.isCellReference(expr) {
		return cellResolver(expr)
	}

	// Handle ranges
	if strings.Contains(expr, ":") {
		return sms.evaluateRange(expr, cellResolver)
	}

	// Handle arithmetic expressions
	return sms.evaluateArithmetic(expr, cellResolver)
}

// evaluateFunctionCall evaluates a function call
func (sms *SpreadsheetMathService) evaluateFunctionCall(expr string, cellResolver func(string) (interface{}, error)) (interface{}, error) {
	// Parse function name and arguments
	funcName, args, err := sms.parseFunctionCall(expr)
	if err != nil {
		return nil, err
	}

	function, exists := sms.functions[strings.ToUpper(funcName)]
	if !exists {
		return nil, fmt.Errorf("unknown function: %s", funcName)
	}

	// Evaluate arguments
	argValues := make([]float64, len(args))
	for i, arg := range args {
		result, err := sms.evaluateExpression(arg, cellResolver)
		if err != nil {
			return nil, err
		}

		// Convert to float64
		switch v := result.(type) {
		case float64:
			argValues[i] = v
		case float32:
			argValues[i] = float64(v)
		case int:
			argValues[i] = float64(v)
		case int32:
			argValues[i] = float64(v)
		case int64:
			argValues[i] = float64(v)
		default:
			return nil, fmt.Errorf("function argument must be numeric")
		}
	}

	// Validate argument count
	if len(argValues) < function.MinArgs || (function.MaxArgs != -1 && len(argValues) > function.MaxArgs) {
		return nil, fmt.Errorf("function %s requires %d-%d arguments, got %d",
			function.Name, function.MinArgs, function.MaxArgs, len(argValues))
	}

	// Call function
	return function.Function(argValues)
}

// parseFunctionCall parses a function call expression
func (sms *SpreadsheetMathService) parseFunctionCall(expr string) (string, []string, error) {
	expr = strings.TrimSpace(expr)

	openParen := strings.Index(expr, "(")
	if openParen == -1 {
		return "", nil, fmt.Errorf("invalid function call: missing opening parenthesis")
	}

	funcName := strings.TrimSpace(expr[:openParen])
	closeParen := strings.LastIndex(expr, ")")
	if closeParen == -1 {
		return "", nil, fmt.Errorf("invalid function call: missing closing parenthesis")
	}

	argsStr := expr[openParen+1 : closeParen]
	args := sms.parseArguments(argsStr)

	return funcName, args, nil
}

// parseArguments parses function arguments
func (sms *SpreadsheetMathService) parseArguments(argsStr string) []string {
	var args []string
	var current strings.Builder
	var parenDepth int

	for _, r := range argsStr {
		switch r {
		case ',':
			if parenDepth == 0 {
				args = append(args, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		case '(':
			parenDepth++
			current.WriteRune(r)
		case ')':
			parenDepth--
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}

	return args
}

// isCellReference checks if a string is a cell reference
func (sms *SpreadsheetMathService) isCellReference(ref string) bool {
	// Simple check for A1, B2, etc.
	matched, _ := regexp.MatchString(`^[A-Z]+\d+$`, strings.ToUpper(ref))
	return matched
}

// evaluateRange evaluates a cell range (e.g., A1:B3)
func (sms *SpreadsheetMathService) evaluateRange(rangeExpr string, cellResolver func(string) (interface{}, error)) (interface{}, error) {
	// Parse range (simplified - would need full implementation)
	parts := strings.Split(rangeExpr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", rangeExpr)
	}

	// For now, return array of values
	// In full implementation, this would expand the range
	values := []interface{}{}
	for _, ref := range parts {
		value, err := cellResolver(ref)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

// evaluateArithmetic evaluates arithmetic expressions
func (sms *SpreadsheetMathService) evaluateArithmetic(expr string, cellResolver func(string) (interface{}, error)) (interface{}, error) {
	// Simple arithmetic evaluation (would need full expression parser in production)
	expr = strings.ReplaceAll(expr, " ", "")

	// Handle basic operations
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) == 2 {
			left, err := sms.evaluateSimpleTerm(parts[0], cellResolver)
			if err != nil {
				return nil, err
			}
			right, err := sms.evaluateSimpleTerm(parts[1], cellResolver)
			if err != nil {
				return nil, err
			}
			return toFloat64(left) + toFloat64(right), nil
		}
	}

	if strings.Contains(expr, "-") {
		parts := strings.Split(expr, "-")
		if len(parts) == 2 {
			left, err := sms.evaluateSimpleTerm(parts[0], cellResolver)
			if err != nil {
				return nil, err
			}
			right, err := sms.evaluateSimpleTerm(parts[1], cellResolver)
			if err != nil {
				return nil, err
			}
			return toFloat64(left) - toFloat64(right), nil
		}
	}

	// Single number or cell reference
	return sms.evaluateSimpleTerm(expr, cellResolver)
}

// evaluateSimpleTerm evaluates a simple term
func (sms *SpreadsheetMathService) evaluateSimpleTerm(term string, cellResolver func(string) (interface{}, error)) (float64, error) {
	term = strings.TrimSpace(term)

	// Try to parse as number
	if num, err := strconv.ParseFloat(term, 64); err == nil {
		return num, nil
	}

	// Try as cell reference
	if sms.isCellReference(term) {
		value, err := cellResolver(term)
		if err != nil {
			return 0, err
		}
		return toFloat64(value), nil
	}

	return 0, fmt.Errorf("invalid term: %s", term)
}

// toFloat64 converts various numeric types to float64
func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

// Function implementations are now in separate category files:
// - math_basic.go
// - math_statistical.go
// - math_trigonometric.go
// - math_exponential.go
// - math_power.go
// - math_rounding.go
// - math_conditional.go
// - math_datetime.go
// - math_text.go
// - math_lookup.go
// - math_financial.go

// ParseCellReference parses a cell reference like "A1" into column and row indices
func (sms *SpreadsheetMathService) ParseCellReference(ref string) (*CellReference, error) {
	ref = strings.ToUpper(strings.TrimSpace(ref))

	// Handle sheet references (e.g., "Sheet1!A1")
	var sheet string
	if strings.Contains(ref, "!") {
		parts := strings.Split(ref, "!")
		sheet = parts[0]
		ref = parts[1]
	}

	// Parse column (letters)
	colStr := ""
	rowStr := ""
	for i, r := range ref {
		if unicode.IsLetter(r) {
			colStr += string(r)
		} else {
			rowStr = ref[i:]
			break
		}
	}

	if colStr == "" || rowStr == "" {
		return nil, fmt.Errorf("invalid cell reference: %s", ref)
	}

	// Convert column letters to number (A=0, B=1, ..., Z=25, AA=26, etc.)
	col := 0
	for _, r := range colStr {
		col = col*26 + int(r-'A'+1)
	}
	col-- // Make it 0-based

	// Parse row number
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return nil, fmt.Errorf("invalid row number in cell reference: %s", ref)
	}
	row-- // Make it 0-based

	return &CellReference{
		Sheet: sheet,
		Col:   col,
		Row:   row,
	}, nil
}

// FormatCellReference formats column and row indices into a cell reference
func (sms *SpreadsheetMathService) FormatCellReference(sheet string, col, row int) string {
	// Convert column to letters (0=A, 1=B, ..., 25=Z, 26=AA, etc.)
	colStr := ""
	for col >= 0 {
		colStr = string(rune('A'+col%26)) + colStr
		col = col/26 - 1
	}

	rowStr := strconv.Itoa(row + 1) // Make it 1-based

	if sheet != "" {
		return fmt.Sprintf("%s!%s%s", sheet, colStr, rowStr)
	}
	return fmt.Sprintf("%s%s", colStr, rowStr)
}

// GetAvailableFunctions returns a list of available mathematical functions
func (sms *SpreadsheetMathService) GetAvailableFunctions() map[string]MathFunction {
	return sms.functions
}

// ValidateFormula validates a formula syntax
func (sms *SpreadsheetMathService) ValidateFormula(formula string) error {
	if !strings.HasPrefix(strings.ToUpper(formula), "=") {
		return nil // Not a formula
	}

	expression := formula[1:]

	// Basic validation - check for balanced parentheses
	openCount := strings.Count(expression, "(")
	closeCount := strings.Count(expression, ")")

	if openCount != closeCount {
		return fmt.Errorf("unbalanced parentheses")
	}

	// Check for function names
	// This is a basic check - full validation would be more complex
	for funcName := range sms.functions {
		if strings.Contains(strings.ToUpper(expression), funcName+"(") {
			// Valid function call found
			return nil
		}
	}

	return nil
}
