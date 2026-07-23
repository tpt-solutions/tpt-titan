// @ts-nocheck
// Math Service for Natural Math Input
// Connects to backend handwriting recognition and math processing services

class MathService {
    constructor() {
        this.baseURL = '/api/v1';
    }

    // Get auth token from localStorage
    getAuthToken() {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('token') || '';
        }
        return '';
    }

    // Get headers for API requests
    getHeaders() {
        const token = this.getAuthToken();
        return {
            'Content-Type': 'application/json',
            'Authorization': token ? `Bearer ${token}` : ''
        };
    }

    /**
     * Convert natural language math expression to LaTeX/MathML
     * @param {string} expression - Natural language expression (e.g., "integral of x squared dx")
     * @param {string} fromFormat - Input format: "text", "latex", "mathml"
     * @param {string} toFormat - Output format: "latex", "mathml", "svg"
     * @returns {Promise<Object>} Converted expression with all formats
     */
    async convertExpression(expression, fromFormat = 'text', toFormat = 'latex') {
        try {
            const response = await fetch(`${this.baseURL}/math/convert`, {
                method: 'POST',
                headers: this.getHeaders(),
                body: JSON.stringify({
                    expression,
                    from_format: fromFormat,
                    to_format: toFormat
                })
            });

            if (!response.ok) {
                // If API fails, fall back to local conversion
                console.warn('Math API unavailable, using local conversion');
                return this.localConvert(expression);
            }

            const result = await response.json();
            return {
                original: result.original,
                converted: result.converted,
                format: result.format,
                success: true
            };
        } catch (error) {
            console.error('Math conversion error:', error);
            // Fallback to local conversion
            return this.localConvert(expression);
        }
    }

    /**
     * Validate a mathematical expression
     * @param {string} expression - Math expression to validate
     * @returns {Promise<Object>} Validation result
     */
    async validateExpression(expression) {
        try {
            const response = await fetch(`${this.baseURL}/math/validate`, {
                method: 'POST',
                headers: this.getHeaders(),
                body: JSON.stringify({ expression })
            });

            if (!response.ok) {
                // Basic local validation
                return this.localValidate(expression);
            }

            const result = await response.json();
            return {
                valid: result.valid,
                message: result.message
            };
        } catch (error) {
            console.error('Validation error:', error);
            return this.localValidate(expression);
        }
    }

    /**
     * Optimize a mathematical expression
     * @param {string} expression - Expression to optimize
     * @returns {Promise<Object>} Optimized expression
     */
    async optimizeExpression(expression) {
        try {
            const response = await fetch(`${this.baseURL}/math/optimize`, {
                method: 'POST',
                headers: this.getHeaders(),
                body: JSON.stringify({ expression })
            });

            if (!response.ok) {
                return { original: expression, optimized: expression };
            }

            const result = await response.json();
            return {
                original: result.original,
                optimized: result.optimized
            };
        } catch (error) {
            console.error('Optimization error:', error);
            return { original: expression, optimized: expression };
        }
    }

    /**
     * Get available mathematical functions
     * @param {string} category - Function category (optional)
     * @returns {Promise<Array>} List of functions
     */
    async getFunctions(category = '') {
        try {
            const url = category 
                ? `${this.baseURL}/math/functions?category=${category}`
                : `${this.baseURL}/math/functions`;
            
            const response = await fetch(url, {
                headers: this.getHeaders()
            });

            if (!response.ok) {
                return this.getDefaultFunctions();
            }

            const result = await response.json();
            return result.functions || [];
        } catch (error) {
            console.error('Get functions error:', error);
            return this.getDefaultFunctions();
        }
    }

    /**
     * Get mathematical symbols
     * @param {string} category - Symbol category (optional)
     * @returns {Promise<Array>} List of symbols
     */
    async getSymbols(category = '') {
        try {
            const url = category 
                ? `${this.baseURL}/math/symbols?category=${category}`
                : `${this.baseURL}/math/symbols`;
            
            const response = await fetch(url, {
                headers: this.getHeaders()
            });

            if (!response.ok) {
                return this.getDefaultSymbols();
            }

            const result = await response.json();
            return result.symbols || [];
        } catch (error) {
            console.error('Get symbols error:', error);
            return this.getDefaultSymbols();
        }
    }

    /**
     * Get equation templates
     * @param {string} category - Template category (optional)
     * @returns {Promise<Array>} List of templates
     */
    async getTemplates(category = '') {
        try {
            const url = category 
                ? `${this.baseURL}/math/templates?category=${category}`
                : `${this.baseURL}/math/templates`;
            
            const response = await fetch(url, {
                headers: this.getHeaders()
            });

            if (!response.ok) {
                return this.getDefaultTemplates();
            }

            const result = await response.json();
            return result.templates || [];
        } catch (error) {
            console.error('Get templates error:', error);
            return this.getDefaultTemplates();
        }
    }

    /**
     * Search equations
     * @param {string} query - Search query
     * @returns {Promise<Array>} Matching equations
     */
    async searchEquations(query) {
        try {
            const response = await fetch(`${this.baseURL}/math/templates/search?q=${encodeURIComponent(query)}`, {
                headers: this.getHeaders()
            });

            if (!response.ok) {
                return [];
            }

            const result = await response.json();
            return result.templates || [];
        } catch (error) {
            console.error('Search equations error:', error);
            return [];
        }
    }

    /**
     * Export equation to different format
     * @param {Object} expression - Math expression object
     * @param {string} format - Export format: "latex", "mathml", "svg", "png", "pdf"
     * @returns {Promise<Blob>} Exported data as blob
     */
    async exportEquation(expression, format = 'latex') {
        try {
            const response = await fetch(`${this.baseURL}/math/export`, {
                method: 'POST',
                headers: this.getHeaders(),
                body: JSON.stringify({
                    expression,
                    format
                })
            });

            if (!response.ok) {
                throw new Error(`Export failed: ${response.statusText}`);
            }

            return await response.blob();
        } catch (error) {
            console.error('Export error:', error);
            throw error;
        }
    }

    /**
     * Local conversion fallback when API is unavailable
     * @param {string} text - Natural language text
     * @returns {Object} Converted expression
     */
    localConvert(text) {
        const expressions = {
            // Integrals
            'integral of (.+?) d(x|y|z|t|u|v|w)': (match, expr, var_) => `\\int ${expr} \\, d${var_}`,
            'integral from (.+?) to (.+?) of (.+?) d(x|y|z|t|u|v|w)': (match, from, to, expr, var_) => 
                `\\int_{${from}}^{${to}} ${expr} \\, d${var_}`,
            
            // Fractions
            'fraction (.+?) over (.+)': (match, num, den) => `\\frac{${num}}{${den}}`,
            '(.+?) divided by (.+)': (match, num, den) => `\\frac{${num}}{${den}}`,
            
            // Roots
            'square root of (.+)': (match, expr) => `\\sqrt{${expr}}`,
            'cube root of (.+)': (match, expr) => `\\sqrt[3]{${expr}}`,
            'nth root of (.+?) with n=(.+)': (match, expr, n) => `\\sqrt[${n}]{${expr}}`,
            
            // Powers
            '(.+?) squared': (match, base) => `${base}^{2}`,
            '(.+?) cubed': (match, base) => `${base}^{3}`,
            '(.+?) to the power of (.+)': (match, base, exp) => `${base}^{${exp}}`,
            'e to the power of (.+)': (match, exp) => `e^{${exp}}`,
            
            // Sums and products
            'sum from (.+?)=(.+?) to (.+?) of (.+)': (match, var_, from, to, expr) => 
                `\\sum_{${var_}=${from}}^{${to}} ${expr}`,
            'sum from (.+?) to (.+?) of (.+)': (match, from, to, expr) => 
                `\\sum_{${from}}^{${to}} ${expr}`,
            'product from (.+?)=(.+?) to (.+?) of (.+)': (match, var_, from, to, expr) => 
                `\\prod_{${var_}=${from}}^{${to}} ${expr}`,
            
            // Limits
            'limit as (.+?) approaches (.+?) of (.+)': (match, var_, val, expr) => 
                `\\lim_{${var_} \\to ${val}} ${expr}`,
            
            // Greek letters
            '\\balpha\\b': '\\alpha',
            '\\bbeta\\b': '\\beta',
            '\\bgamma\\b': '\\gamma',
            '\\bdelta\\b': '\\delta',
            '\\bepsilon\\b': '\\epsilon',
            '\\btheta\\b': '\\theta',
            '\\blambda\\b': '\\lambda',
            '\\bmu\\b': '\\mu',
            '\\bpi\\b': '\\pi',
            '\\bsigma\\b': '\\sigma',
            '\\bomega\\b': '\\omega',
            
            // Operators
            '\\bplus or minus\\b': '\\pm',
            '\\btimes\\b': '\\times',
            '\\bdivided by\\b': '\\div',
            '\\binfinity\\b': '\\infty',
            '\\binfty\\b': '\\infty',
            '\\btherefore\\b': '\\therefore',
            '\\bbecause\\b': '\\because',
            
            // Functions
            '\\bsin\\b': '\\sin',
            '\\bcos\\b': '\\cos',
            '\\btan\\b': '\\tan',
            '\\blog\\b': '\\log',
            '\\bln\\b': '\\ln',
            '\\bexp\\b': '\\exp',
            '\\bsqrt\\b': '\\sqrt'
        };

        let latex = text;
        let converted = false;

        for (const [pattern, replacement] of Object.entries(expressions)) {
            const regex = new RegExp(pattern, 'gi');
            if (regex.test(latex)) {
                latex = latex.replace(regex, replacement);
                converted = true;
            }
        }

        // Generate MathML from LaTeX (simplified)
        const mathml = this.latexToMathML(latex);

        return {
            original: text,
            converted: latex,
            format: 'latex',
            latex: latex,
            mathml: mathml,
            success: converted,
            source: 'local'
        };
    }

    /**
     * Simple LaTeX to MathML conversion
     * @param {string} latex - LaTeX expression
     * @returns {string} MathML
     */
    latexToMathML(latex) {
        // Basic conversion - in production, use a proper library
        let mathml = latex
            .replace(/\\int/g, '<mo>&#x222B;</mo>')
            .replace(/\\sum/g, '<mo>&#x2211;</mo>')
            .replace(/\\prod/g, '<mo>&#x220F;</mo>')
            .replace(/\\frac\{([^}]+)\}\{([^}]+)\}/g, '<mfrac><mrow>$1</mrow><mrow>$2</mrow></mfrac>')
            .replace(/\\sqrt\{([^}]+)\}/g, '<msqrt><mrow>$1</mrow></msqrt>')
            .replace(/\\sqrt\[(\d+)\]\{([^}]+)\}/g, '<mroot><mrow>$2</mrow><mn>$1</mn></mroot>')
            .replace(/\^(\{[^}]+\}|\d+)/g, '<msup><mi></mi><mn>$1</mn></msup>')
            .replace(/\\alpha/g, '<mi>&#x03B1;</mi>')
            .replace(/\\beta/g, '<mi>&#x03B2;</mi>')
            .replace(/\\gamma/g, '<mi>&#x03B3;</mi>')
            .replace(/\\delta/g, '<mi>&#x03B4;</mi>')
            .replace(/\\pi/g, '<mi>&#x03C0;</mi>')
            .replace(/\\sigma/g, '<mi>&#x03C3;</mi>')
            .replace(/\\omega/g, '<mi>&#x03C9;</mi>')
            .replace(/\\pm/g, '<mo>&#x00B1;</mo>')
            .replace(/\\times/g, '<mo>&#x00D7;</mo>')
            .replace(/\\div/g, '<mo>&#x00F7;</mo>')
            .replace(/\\infty/g, '<mi>&#x221E;</mi>');

        return `<math xmlns="http://www.w3.org/1998/Math/MathML"><mrow>${mathml}</mrow></math>`;
    }

    /**
     * Local validation fallback
     * @param {string} expression - Expression to validate
     * @returns {Object} Validation result
     */
    localValidate(expression) {
        if (!expression || expression.trim() === '') {
            return { valid: false, message: 'Expression is empty' };
        }

        // Check for balanced parentheses
        const openCount = (expression.match(/\(/g) || []).length;
        const closeCount = (expression.match(/\)/g) || []).length;
        
        if (openCount !== closeCount) {
            return { valid: false, message: 'Unbalanced parentheses' };
        }

        // Check for balanced braces
        const openBrace = (expression.match(/\{/g) || []).length;
        const closeBrace = (expression.match(/\}/g) || []).length;
        
        if (openBrace !== closeBrace) {
            return { valid: false, message: 'Unbalanced braces' };
        }

        return { valid: true, message: 'Expression appears valid' };
    }

    /**
     * Default functions when API is unavailable
     * @returns {Array} Default functions
     */
    getDefaultFunctions() {
        return [
            { name: 'sum', category: 'aggregation', description: 'Sum of numbers', syntax: 'SUM(a, b, ...)' },
            { name: 'average', category: 'aggregation', description: 'Average of numbers', syntax: 'AVERAGE(a, b, ...)' },
            { name: 'sin', category: 'trigonometric', description: 'Sine function', syntax: 'SIN(x)' },
            { name: 'cos', category: 'trigonometric', description: 'Cosine function', syntax: 'COS(x)' },
            { name: 'tan', category: 'trigonometric', description: 'Tangent function', syntax: 'TAN(x)' },
            { name: 'sqrt', category: 'basic', description: 'Square root', syntax: 'SQRT(x)' },
            { name: 'power', category: 'basic', description: 'Power function', syntax: 'POWER(base, exp)' },
            { name: 'log', category: 'logarithmic', description: 'Logarithm', syntax: 'LOG(x, [base])' },
            { name: 'ln', category: 'logarithmic', description: 'Natural logarithm', syntax: 'LN(x)' },
            { name: 'exp', category: 'exponential', description: 'Exponential function', syntax: 'EXP(x)' },
            { name: 'integral', category: 'calculus', description: 'Integration', syntax: '\\int f(x) dx' },
            { name: 'derivative', category: 'calculus', description: 'Differentiation', syntax: '\\frac{d}{dx} f(x)' }
        ];
    }

    /**
     * Default symbols when API is unavailable
     * @returns {Array} Default symbols
     */
    getDefaultSymbols() {
        return [
            { symbol: '\\pi', name: 'Pi', category: 'greek', unicode: 'π' },
            { symbol: '\\alpha', name: 'Alpha', category: 'greek', unicode: 'α' },
            { symbol: '\\beta', name: 'Beta', category: 'greek', unicode: 'β' },
            { symbol: '\\gamma', name: 'Gamma', category: 'greek', unicode: 'γ' },
            { symbol: '\\delta', name: 'Delta', category: 'greek', unicode: 'δ' },
            { symbol: '\\theta', name: 'Theta', category: 'greek', unicode: 'θ' },
            { symbol: '\\sigma', name: 'Sigma', category: 'greek', unicode: 'σ' },
            { symbol: '\\omega', name: 'Omega', category: 'greek', unicode: 'ω' },
            { symbol: '\\infty', name: 'Infinity', category: 'operators', unicode: '∞' },
            { symbol: '\\sum', name: 'Summation', category: 'operators', unicode: '∑' },
            { symbol: '\\int', name: 'Integral', category: 'operators', unicode: '∫' },
            { symbol: '\\pm', name: 'Plus-Minus', category: 'operators', unicode: '±' },
            { symbol: '\\times', name: 'Times', category: 'operators', unicode: '×' },
            { symbol: '\\div', name: 'Divided by', category: 'operators', unicode: '÷' }
        ];
    }

    /**
     * Default templates when API is unavailable
     * @returns {Array} Default templates
     */
    getDefaultTemplates() {
        return [
            {
                id: 'pythagorean',
                name: 'Pythagorean Theorem',
                category: 'geometry',
                latex: 'a^{2} + b^{2} = c^{2}',
                description: 'The square of the hypotenuse equals the sum of the squares of the other two sides'
            },
            {
                id: 'quadratic',
                name: 'Quadratic Formula',
                category: 'algebra',
                latex: 'x = \\frac{-b \\pm \\sqrt{b^{2} - 4ac}}{2a}',
                description: 'Solutions to the quadratic equation ax² + bx + c = 0'
            },
            {
                id: 'euler',
                name: "Euler's Identity",
                category: 'complex_analysis',
                latex: 'e^{i\\pi} + 1 = 0',
                description: "Euler's famous identity linking e, i, π, 0, and 1"
            },
            {
                id: 'einstein',
                name: 'Mass-Energy Equivalence',
                category: 'physics',
                latex: 'E = mc^{2}',
                description: "Einstein's famous equation relating mass and energy"
            },
            {
                id: 'integral',
                name: 'Gaussian Integral',
                category: 'calculus',
                latex: '\\int_{-\\infty}^{\\infty} e^{-x^{2}} dx = \\sqrt{\\pi}',
                description: 'The integral of the Gaussian function over the entire real line'
            }
        ];
    }

    /**
     * Parse natural language math and return all formats
     * @param {string} text - Natural language input
     * @returns {Promise<Object>} Parsed math with all representations
     */
    async parseNaturalMath(text) {
        // Try backend first
        try {
            const result = await this.convertExpression(text, 'text', 'latex');
            if (result.success) {
                return {
                    text: text,
                    latex: result.latex || result.converted,
                    mathml: result.mathml,
                    format: 'latex',
                    source: result.source || 'api'
                };
            }
        } catch (error) {
            console.warn('Backend parse failed, using local:', error);
        }

        // Fallback to local
        const local = this.localConvert(text);
        return {
            text: text,
            latex: local.latex,
            mathml: local.mathml,
            format: 'latex',
            source: 'local'
        };
    }

    /**
     * Render LaTeX to HTML using KaTeX-style rendering
     * @param {string} latex - LaTeX expression
     * @returns {string} HTML representation
     */
    renderToHTML(latex) {
        // Simple HTML rendering - in production, use KaTeX library
        return latex
            .replace(/\\int/g, '∫')
            .replace(/\\sum/g, '∑')
            .replace(/\\prod/g, '∏')
            .replace(/\\frac\{([^}]+)\}\{([^}]+)\}/g, '<span class="fraction"><span class="num">$1</span><span class="den">$2</span></span>')
            .replace(/\\sqrt\{([^}]+)\}/g, '√($1)')
            .replace(/\^(\{[^}]+\}|\d+)/g, '<sup>$1</sup>')
            .replace(/_(\{[^}]+\}|\d+)/g, '<sub>$1</sub>')
            .replace(/\\alpha/g, 'α')
            .replace(/\\beta/g, 'β')
            .replace(/\\gamma/g, 'γ')
            .replace(/\\delta/g, 'δ')
            .replace(/\\pi/g, 'π')
            .replace(/\\sigma/g, 'σ')
            .replace(/\\omega/g, 'ω')
            .replace(/\\pm/g, '±')
            .replace(/\\times/g, '×')
            .replace(/\\div/g, '÷')
            .replace(/\\infty/g, '∞');
    }
}

export default new MathService();
