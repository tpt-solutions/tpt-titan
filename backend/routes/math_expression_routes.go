package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tpt-titan/backend/services"
)

// ValidateExpression validates a mathematical expression
func ValidateExpression(c *gin.Context) {
	var req struct {
		Expression string `json:"expression" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hrs := services.NewHandwritingRecognitionService("", "", "")
	valid, message := hrs.ValidateMathematicalExpression(req.Expression)

	c.JSON(http.StatusOK, gin.H{
		"valid":   valid,
		"message": message,
	})
}

// OptimizeExpression optimizes a mathematical expression
func OptimizeExpression(c *gin.Context) {
	var req struct {
		Expression string `json:"expression" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hrs := services.NewHandwritingRecognitionService("", "", "")
	optimized, err := hrs.OptimizeExpression(req.Expression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original":  req.Expression,
		"optimized": optimized,
	})
}

// ConvertExpression converts between different mathematical formats
func ConvertExpression(c *gin.Context) {
	var req struct {
		Expression string `json:"expression" binding:"required"`
		FromFormat string `json:"from_format" binding:"required"` // "text", "latex", "mathml"
		ToFormat   string `json:"to_format" binding:"required"`   // "latex", "mathml", "svg"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hrs := services.NewHandwritingRecognitionService("", "", "")

	// Parse the expression first
	mathExpr, err := hrs.ParseMathematicalExpression(req.Expression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse expression: " + err.Error()})
		return
	}

	// Convert to requested format
	var result string
	switch req.ToFormat {
	case "latex":
		result = mathExpr.LaTeX
	case "mathml":
		result = mathExpr.MathML
	case "svg":
		result = mathExpr.SVG
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported target format"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original":  req.Expression,
		"converted": result,
		"format":    req.ToFormat,
	})
}

// GetMathematicalFunctions returns available mathematical functions
func GetMathematicalFunctions(c *gin.Context) {
	category := c.Query("category") // "basic", "trigonometric", "calculus", etc.

	functions := []gin.H{
		{
			"name":        "sum",
			"category":    "aggregation",
			"description": "Returns the sum of a series of numbers",
			"syntax":      "SUM(number1, [number2], ...)",
			"example":     "SUM(1, 2, 3) = 6",
		},
		{
			"name":        "average",
			"category":    "aggregation",
			"description": "Returns the average of a series of numbers",
			"syntax":      "AVERAGE(number1, [number2], ...)",
			"example":     "AVERAGE(1, 2, 3, 4) = 2.5",
		},
		{
			"name":        "sin",
			"category":    "trigonometric",
			"description": "Returns the sine of an angle",
			"syntax":      "SIN(angle)",
			"example":     "SIN(π/2) = 1",
		},
		{
			"name":        "cos",
			"category":    "trigonometric",
			"description": "Returns the cosine of an angle",
			"syntax":      "COS(angle)",
			"example":     "COS(0) = 1",
		},
		{
			"name":        "sqrt",
			"category":    "basic",
			"description": "Returns the square root of a number",
			"syntax":      "SQRT(number)",
			"example":     "SQRT(9) = 3",
		},
		{
			"name":        "power",
			"category":    "basic",
			"description": "Returns a number raised to a power",
			"syntax":      "POWER(base, exponent)",
			"example":     "POWER(2, 3) = 8",
		},
		{
			"name":        "exp",
			"category":    "exponential",
			"description": "Returns e raised to the power of a number",
			"syntax":      "EXP(number)",
			"example":     "EXP(1) ≈ 2.718",
		},
		{
			"name":        "log",
			"category":    "logarithmic",
			"description": "Returns the logarithm of a number",
			"syntax":      "LOG(number, [base])",
			"example":     "LOG(100) = 2 (base 10)",
		},
		{
			"name":        "integral",
			"category":    "calculus",
			"description": "Calculus integration symbol",
			"syntax":      "\\int expression dx",
			"example":     "\\int x^{2} dx = \\frac{x^{3}}{3} + C",
		},
		{
			"name":        "derivative",
			"category":    "calculus",
			"description": "Calculus derivative symbol",
			"syntax":      "\\frac{d}{dx} expression",
			"example":     "\\frac{d}{dx} x^{2} = 2x",
		},
	}

	// Filter by category if specified
	if category != "" {
		filtered := []gin.H{}
		for _, fn := range functions {
			if cat, ok := fn["category"].(string); ok && cat == category {
				filtered = append(filtered, fn)
			}
		}
		functions = filtered
	}

	c.JSON(http.StatusOK, gin.H{"functions": functions})
}

// GetMathematicalSymbols returns available mathematical symbols
func GetMathematicalSymbols(c *gin.Context) {
	category := c.Query("category") // "greek", "operators", "arrows", etc.

	symbols := []gin.H{
		{
			"symbol":      "\\pi",
			"name":        "Pi",
			"category":    "greek",
			"description": "Mathematical constant π ≈ 3.14159",
			"unicode":     "π",
		},
		{
			"symbol":      "\\alpha",
			"name":        "Alpha",
			"category":    "greek",
			"description": "Greek letter α",
			"unicode":     "α",
		},
		{
			"symbol":      "\\beta",
			"name":        "Beta",
			"category":    "greek",
			"description": "Greek letter β",
			"unicode":     "β",
		},
		{
			"symbol":      "\\gamma",
			"name":        "Gamma",
			"category":    "greek",
			"description": "Greek letter γ",
			"unicode":     "γ",
		},
		{
			"symbol":      "\\delta",
			"name":        "Delta",
			"category":    "greek",
			"description": "Greek letter δ",
			"unicode":     "δ",
		},
		{
			"symbol":      "\\theta",
			"name":        "Theta",
			"category":    "greek",
			"description": "Greek letter θ",
			"unicode":     "θ",
		},
		{
			"symbol":      "\\lambda",
			"name":        "Lambda",
			"category":    "greek",
			"description": "Greek letter λ",
			"unicode":     "λ",
		},
		{
			"symbol":      "\\mu",
			"name":        "Mu",
			"category":    "greek",
			"description": "Greek letter μ",
			"unicode":     "μ",
		},
		{
			"symbol":      "\\sigma",
			"name":        "Sigma",
			"category":    "greek",
			"description": "Greek letter σ",
			"unicode":     "σ",
		},
		{
			"symbol":      "\\omega",
			"name":        "Omega",
			"category":    "greek",
			"description": "Greek letter ω",
			"unicode":     "ω",
		},
		{
			"symbol":      "\\infty",
			"name":        "Infinity",
			"category":    "operators",
			"description": "Infinity symbol ∞",
			"unicode":     "∞",
		},
		{
			"symbol":      "\\sum",
			"name":        "Summation",
			"category":    "operators",
			"description": "Summation symbol Σ",
			"unicode":     "Σ",
		},
		{
			"symbol":      "\\prod",
			"name":        "Product",
			"category":    "operators",
			"description": "Product symbol ∏",
			"unicode":     "∏",
		},
		{
			"symbol":      "\\int",
			"name":        "Integral",
			"category":    "operators",
			"description": "Integral symbol ∫",
			"unicode":     "∫",
		},
		{
			"symbol":      "\\partial",
			"name":        "Partial",
			"category":    "operators",
			"description": "Partial derivative symbol ∂",
			"unicode":     "∂",
		},
		{
			"symbol":      "\\nabla",
			"name":        "Nabla",
			"category":    "operators",
			"description": "Gradient operator ∇",
			"unicode":     "∇",
		},
		{
			"symbol":      "\\rightarrow",
			"name":        "Right Arrow",
			"category":    "arrows",
			"description": "Right arrow →",
			"unicode":     "→",
		},
		{
			"symbol":      "\\leftarrow",
			"name":        "Left Arrow",
			"category":    "arrows",
			"description": "Left arrow ←",
			"unicode":     "←",
		},
		{
			"symbol":      "\\leftrightarrow",
			"name":        "Left Right Arrow",
			"category":    "arrows",
			"description": "Left-right arrow ↔",
			"unicode":     "↔",
		},
		{
			"symbol":      "\\Rightarrow",
			"name":        "Right Double Arrow",
			"category":    "arrows",
			"description": "Right double arrow ⇒",
			"unicode":     "⇒",
		},
	}

	// Filter by category if specified
	if category != "" {
		filtered := []gin.H{}
		for _, symbol := range symbols {
			if cat, ok := symbol["category"].(string); ok && cat == category {
				filtered = append(filtered, symbol)
			}
		}
		symbols = filtered
	}

	c.JSON(http.StatusOK, gin.H{"symbols": symbols})
}

// GetMathematicalConstants returns common mathematical constants
func GetMathematicalConstants(c *gin.Context) {
	constants := []gin.H{
		{
			"name":        "Pi",
			"symbol":      "\\pi",
			"value":       "3.141592653589793",
			"description": "Ratio of a circle's circumference to its diameter",
			"latex":       "\\pi",
		},
		{
			"name":        "Euler's Number",
			"symbol":      "e",
			"value":       "2.718281828459045",
			"description": "Base of the natural logarithm",
			"latex":       "e",
		},
		{
			"name":        "Golden Ratio",
			"symbol":      "\\phi",
			"value":       "1.618033988749895",
			"description": "Ratio where (a+b)/a = a/b",
			"latex":       "\\phi",
		},
		{
			"name":        "Square Root of 2",
			"symbol":      "\\sqrt{2}",
			"value":       "1.414213562373095",
			"description": "Length of diagonal of unit square",
			"latex":       "\\sqrt{2}",
		},
		{
			"name":        "Natural Logarithm of 2",
			"symbol":      "\\ln{2}",
			"value":       "0.693147180559945",
			"description": "Natural logarithm of 2",
			"latex":       "\\ln{2}",
		},
		{
			"name":        "Euler-Mascheroni Constant",
			"symbol":      "\\gamma",
			"value":       "0.577215664901533",
			"description": "Euler-Mascheroni constant γ",
			"latex":       "\\gamma",
		},
	}

	c.JSON(http.StatusOK, gin.H{"constants": constants})
}

// GetMathematicalTheorems returns famous mathematical theorems and formulas
func GetMathematicalTheorems(c *gin.Context) {
	category := c.Query("category") // "algebra", "geometry", "calculus", etc.

	theorems := []gin.H{
		{
			"name":      "Pythagorean Theorem",
			"category":  "geometry",
			"statement": "In a right-angled triangle, the square of the hypotenuse equals the sum of the squares of the other two sides",
			"formula":   "a² + b² = c²",
			"latex":     "a^{2} + b^{2} = c^{2}",
			"variables": []string{"a", "b", "c"},
		},
		{
			"name":      "Quadratic Formula",
			"category":  "algebra",
			"statement": "Solutions to the quadratic equation ax² + bx + c = 0",
			"formula":   "x = [-b ± √(b² - 4ac)] / 2a",
			"latex":     "x = \\frac{-b \\pm \\sqrt{b^{2} - 4ac}}{2a}",
			"variables": []string{"a", "b", "c", "x"},
		},
		{
			"name":      "Fundamental Theorem of Calculus",
			"category":  "calculus",
			"statement": "The fundamental theorem relates differentiation and integration",
			"formula":   "d/dx ∫f(t)dt = f(x)",
			"latex":     "\\frac{d}{dx} \\int f(t) \\, dt = f(x)",
			"variables": []string{"f", "x", "t"},
		},
		{
			"name":      "Euler's Identity",
			"category":  "complex_analysis",
			"statement": "Euler's famous identity linking e, i, π, 0, and 1",
			"formula":   "e^(iπ) + 1 = 0",
			"latex":     "e^{i\\pi} + 1 = 0",
			"variables": []string{"e", "i", "π"},
		},
		{
			"name":      "Mass-Energy Equivalence",
			"category":  "physics",
			"statement": "Einstein's famous equation relating mass and energy",
			"formula":   "E = mc²",
			"latex":     "E = mc^{2}",
			"variables": []string{"E", "m", "c"},
		},
		{
			"name":      "Binomial Theorem",
			"category":  "algebra",
			"statement": "Expansion of (x + y)^n",
			"formula":   "(x + y)^n = Σ(n choose k) x^(n-k) y^k",
			"latex":     "(x + y)^{n} = \\sum_{k=0}^{n} \\binom{n}{k} x^{n-k} y^{k}",
			"variables": []string{"x", "y", "n", "k"},
		},
		{
			"name":      "Bayes' Theorem",
			"category":  "probability",
			"statement": "Theorem describing probability of an event based on prior knowledge",
			"formula":   "P(A|B) = P(B|A) * P(A) / P(B)",
			"latex":     "P(A|B) = \\frac{P(B|A) \\cdot P(A)}{P(B)}",
			"variables": []string{"A", "B", "P"},
		},
		{
			"name":      "Taylor Series",
			"category":  "calculus",
			"statement": "Representation of a function as an infinite sum of terms",
			"formula":   "f(x) = Σ f^(n)(a) * (x-a)^n / n!",
			"latex":     "f(x) = \\sum_{n=0}^{\\infty} \\frac{f^{(n)}(a)}{n!} (x - a)^{n}",
			"variables": []string{"f", "x", "a", "n"},
		},
	}

	// Filter by category if specified
	if category != "" {
		filtered := []gin.H{}
		for _, theorem := range theorems {
			if cat, ok := theorem["category"].(string); ok && cat == category {
				filtered = append(filtered, theorem)
			}
		}
		theorems = filtered
	}

	c.JSON(http.StatusOK, gin.H{"theorems": theorems})
}
