import { c as create_ssr_component, f as createEventDispatcher, a as add_attribute, b as escape, d as each, u as null_to_empty, v as validate_component } from "../../../chunks/calendar.js";
import "jspdf";
const TextEditorToolbar = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { documentTitle: documentTitle2 = "Untitled Document" } = $$props;
  let { isSaving = false } = $$props;
  let { saveStatus = "" } = $$props;
  let { hasUnsavedChanges = false } = $$props;
  let { currentDocument = null } = $$props;
  let { availableVoices = [] } = $$props;
  let { isReadingAloud = false } = $$props;
  let { isGeneratingAI = false } = $$props;
  let { isGeneratingSummary = false } = $$props;
  let { canUndo = false } = $$props;
  let { canRedo = false } = $$props;
  createEventDispatcher();
  if ($$props.documentTitle === void 0 && $$bindings.documentTitle && documentTitle2 !== void 0) $$bindings.documentTitle(documentTitle2);
  if ($$props.isSaving === void 0 && $$bindings.isSaving && isSaving !== void 0) $$bindings.isSaving(isSaving);
  if ($$props.saveStatus === void 0 && $$bindings.saveStatus && saveStatus !== void 0) $$bindings.saveStatus(saveStatus);
  if ($$props.hasUnsavedChanges === void 0 && $$bindings.hasUnsavedChanges && hasUnsavedChanges !== void 0) $$bindings.hasUnsavedChanges(hasUnsavedChanges);
  if ($$props.currentDocument === void 0 && $$bindings.currentDocument && currentDocument !== void 0) $$bindings.currentDocument(currentDocument);
  if ($$props.availableVoices === void 0 && $$bindings.availableVoices && availableVoices !== void 0) $$bindings.availableVoices(availableVoices);
  if ($$props.isReadingAloud === void 0 && $$bindings.isReadingAloud && isReadingAloud !== void 0) $$bindings.isReadingAloud(isReadingAloud);
  if ($$props.isGeneratingAI === void 0 && $$bindings.isGeneratingAI && isGeneratingAI !== void 0) $$bindings.isGeneratingAI(isGeneratingAI);
  if ($$props.isGeneratingSummary === void 0 && $$bindings.isGeneratingSummary && isGeneratingSummary !== void 0) $$bindings.isGeneratingSummary(isGeneratingSummary);
  if ($$props.canUndo === void 0 && $$bindings.canUndo && canUndo !== void 0) $$bindings.canUndo(canUndo);
  if ($$props.canRedo === void 0 && $$bindings.canRedo && canRedo !== void 0) $$bindings.canRedo(canRedo);
  return `  <div class="bg-white border-b border-gray-200 px-8 py-4"><div class="max-w-4xl mx-auto flex items-center justify-between"><div class="flex items-center space-x-4"><input placeholder="Document title..." class="text-xl font-semibold bg-transparent border-none outline-none focus:ring-2 focus:ring-blue-500 rounded px-2 py-1"${add_attribute("value", documentTitle2, 0)}> ${hasUnsavedChanges ? `<span class="text-sm text-orange-600" data-svelte-h="svelte-j3vvr2">•</span>` : ``} <span class="text-sm text-gray-500">${escape(saveStatus)}</span></div> <div class="flex items-center space-x-2"> <button class="px-2 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed" ${!canUndo ? "disabled" : ""} title="Undo (Ctrl+Z)">↩️ Undo</button> <button class="px-2 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed" ${!canRedo ? "disabled" : ""} title="Redo (Ctrl+Y or Ctrl+Shift+Z)">↪️ Redo</button> <div class="w-px h-6 bg-gray-300 mx-2"></div> <button class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50" ${isSaving ? "disabled" : ""}>${escape(isSaving ? "Saving..." : "Save")}</button> <div class="relative"><button class="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 flex items-center" data-svelte-h="svelte-f7vzdf">Export ▼</button> ${``}</div> <button class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-wq4x4n">New</button> <button class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700" data-svelte-h="svelte-yag7e5">Open</button> ${currentDocument ? `<button class="px-3 py-1 text-sm bg-purple-600 text-white rounded hover:bg-purple-700" data-svelte-h="svelte-5f90qh">History</button>` : ``} <button class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700" title="Math Help &amp; Examples" data-svelte-h="svelte-z3nccx">Math Help</button> ${availableVoices.length > 0 ? `<button class="px-3 py-1 text-sm bg-purple-600 text-white rounded hover:bg-purple-700 disabled:opacity-50" ${isReadingAloud ? "disabled" : ""} title="Read document aloud">${escape(isReadingAloud ? "🔊 Reading..." : "📖 Read Aloud")}</button>` : ``}  <div class="border-l border-gray-300 pl-4 ml-4 flex items-center space-x-2"><span class="text-xs text-gray-500 font-medium" data-svelte-h="svelte-xyd2b0">AI</span> <div class="relative"><button class="px-3 py-1 text-sm bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-50 relative group" ${isGeneratingAI ? "disabled" : ""} title="Get AI writing suggestions for grammar, style, and clarity improvements" aria-label="Get AI writing suggestions" type="button">${escape(isGeneratingAI ? "💭" : "✨ Suggest")} <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none" data-svelte-h="svelte-bm979e">✨ AI Writing Assistant
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div></div></button> <button class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 relative group ml-1" ${isGeneratingAI ? "disabled" : ""} title="AI content continuation - select text to continue from that point" aria-label="Continue writing with AI" type="button">📝 Continue
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none" data-svelte-h="svelte-1648urj">📝 Continue Writing
							<div class="text-xs text-gray-300 mt-1">Select text to continue from that point</div> <div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div></div></button> <button class="px-3 py-1 text-sm bg-yellow-600 text-white rounded hover:bg-yellow-700 disabled:opacity-50 relative group ml-1" ${isGeneratingSummary ? "disabled" : ""} title="Generate intelligent document summaries - works best with longer documents" aria-label="Generate document summary" type="button">${escape(isGeneratingSummary ? "📋" : "📄 Summary")} <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none" data-svelte-h="svelte-1t087c3">📄 Document Summary
							<div class="text-xs text-gray-300 mt-1">Best with 500+ words</div> <div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div></div></button> <button class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700 relative group ml-1" title="Analyze document for readability, sentiment, and key insights" aria-label="Analyze document" type="button" data-svelte-h="svelte-6nyg3m">📊 Analyze
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">📊 Document Analysis
							<div class="text-xs text-gray-300 mt-1">Readability, sentiment &amp; key phrases</div> <div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div></div></button></div></div></div></div></div>`;
});
class MathService {
  constructor() {
    this.baseURL = "/api/v1";
  }
  // Get auth token from localStorage
  getAuthToken() {
    if (typeof window !== "undefined") {
      return localStorage.getItem("token") || "";
    }
    return "";
  }
  // Get headers for API requests
  getHeaders() {
    const token = this.getAuthToken();
    return {
      "Content-Type": "application/json",
      "Authorization": token ? `Bearer ${token}` : ""
    };
  }
  /**
   * Convert natural language math expression to LaTeX/MathML
   * @param {string} expression - Natural language expression (e.g., "integral of x squared dx")
   * @param {string} fromFormat - Input format: "text", "latex", "mathml"
   * @param {string} toFormat - Output format: "latex", "mathml", "svg"
   * @returns {Promise<Object>} Converted expression with all formats
   */
  async convertExpression(expression, fromFormat = "text", toFormat = "latex") {
    try {
      const response = await fetch(`${this.baseURL}/math/convert`, {
        method: "POST",
        headers: this.getHeaders(),
        body: JSON.stringify({
          expression,
          from_format: fromFormat,
          to_format: toFormat
        })
      });
      if (!response.ok) {
        console.warn("Math API unavailable, using local conversion");
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
      console.error("Math conversion error:", error);
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
        method: "POST",
        headers: this.getHeaders(),
        body: JSON.stringify({ expression })
      });
      if (!response.ok) {
        return this.localValidate(expression);
      }
      const result = await response.json();
      return {
        valid: result.valid,
        message: result.message
      };
    } catch (error) {
      console.error("Validation error:", error);
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
        method: "POST",
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
      console.error("Optimization error:", error);
      return { original: expression, optimized: expression };
    }
  }
  /**
   * Get available mathematical functions
   * @param {string} category - Function category (optional)
   * @returns {Promise<Array>} List of functions
   */
  async getFunctions(category = "") {
    try {
      const url = category ? `${this.baseURL}/math/functions?category=${category}` : `${this.baseURL}/math/functions`;
      const response = await fetch(url, {
        headers: this.getHeaders()
      });
      if (!response.ok) {
        return this.getDefaultFunctions();
      }
      const result = await response.json();
      return result.functions || [];
    } catch (error) {
      console.error("Get functions error:", error);
      return this.getDefaultFunctions();
    }
  }
  /**
   * Get mathematical symbols
   * @param {string} category - Symbol category (optional)
   * @returns {Promise<Array>} List of symbols
   */
  async getSymbols(category = "") {
    try {
      const url = category ? `${this.baseURL}/math/symbols?category=${category}` : `${this.baseURL}/math/symbols`;
      const response = await fetch(url, {
        headers: this.getHeaders()
      });
      if (!response.ok) {
        return this.getDefaultSymbols();
      }
      const result = await response.json();
      return result.symbols || [];
    } catch (error) {
      console.error("Get symbols error:", error);
      return this.getDefaultSymbols();
    }
  }
  /**
   * Get equation templates
   * @param {string} category - Template category (optional)
   * @returns {Promise<Array>} List of templates
   */
  async getTemplates(category = "") {
    try {
      const url = category ? `${this.baseURL}/math/templates?category=${category}` : `${this.baseURL}/math/templates`;
      const response = await fetch(url, {
        headers: this.getHeaders()
      });
      if (!response.ok) {
        return this.getDefaultTemplates();
      }
      const result = await response.json();
      return result.templates || [];
    } catch (error) {
      console.error("Get templates error:", error);
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
      console.error("Search equations error:", error);
      return [];
    }
  }
  /**
   * Export equation to different format
   * @param {Object} expression - Math expression object
   * @param {string} format - Export format: "latex", "mathml", "svg", "png", "pdf"
   * @returns {Promise<Blob>} Exported data as blob
   */
  async exportEquation(expression, format = "latex") {
    try {
      const response = await fetch(`${this.baseURL}/math/export`, {
        method: "POST",
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
      console.error("Export error:", error);
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
      "integral of (.+?) d(x|y|z|t|u|v|w)": (match, expr, var_) => `\\int ${expr} \\, d${var_}`,
      "integral from (.+?) to (.+?) of (.+?) d(x|y|z|t|u|v|w)": (match, from, to, expr, var_) => `\\int_{${from}}^{${to}} ${expr} \\, d${var_}`,
      // Fractions
      "fraction (.+?) over (.+)": (match, num, den) => `\\frac{${num}}{${den}}`,
      "(.+?) divided by (.+)": (match, num, den) => `\\frac{${num}}{${den}}`,
      // Roots
      "square root of (.+)": (match, expr) => `\\sqrt{${expr}}`,
      "cube root of (.+)": (match, expr) => `\\sqrt[3]{${expr}}`,
      "nth root of (.+?) with n=(.+)": (match, expr, n) => `\\sqrt[${n}]{${expr}}`,
      // Powers
      "(.+?) squared": (match, base) => `${base}^{2}`,
      "(.+?) cubed": (match, base) => `${base}^{3}`,
      "(.+?) to the power of (.+)": (match, base, exp) => `${base}^{${exp}}`,
      "e to the power of (.+)": (match, exp) => `e^{${exp}}`,
      // Sums and products
      "sum from (.+?)=(.+?) to (.+?) of (.+)": (match, var_, from, to, expr) => `\\sum_{${var_}=${from}}^{${to}} ${expr}`,
      "sum from (.+?) to (.+?) of (.+)": (match, from, to, expr) => `\\sum_{${from}}^{${to}} ${expr}`,
      "product from (.+?)=(.+?) to (.+?) of (.+)": (match, var_, from, to, expr) => `\\prod_{${var_}=${from}}^{${to}} ${expr}`,
      // Limits
      "limit as (.+?) approaches (.+?) of (.+)": (match, var_, val, expr) => `\\lim_{${var_} \\to ${val}} ${expr}`,
      // Greek letters
      "\\balpha\\b": "\\alpha",
      "\\bbeta\\b": "\\beta",
      "\\bgamma\\b": "\\gamma",
      "\\bdelta\\b": "\\delta",
      "\\bepsilon\\b": "\\epsilon",
      "\\btheta\\b": "\\theta",
      "\\blambda\\b": "\\lambda",
      "\\bmu\\b": "\\mu",
      "\\bpi\\b": "\\pi",
      "\\bsigma\\b": "\\sigma",
      "\\bomega\\b": "\\omega",
      // Operators
      "\\bplus or minus\\b": "\\pm",
      "\\btimes\\b": "\\times",
      "\\bdivided by\\b": "\\div",
      "\\binfinity\\b": "\\infty",
      "\\binfty\\b": "\\infty",
      "\\btherefore\\b": "\\therefore",
      "\\bbecause\\b": "\\because",
      // Functions
      "\\bsin\\b": "\\sin",
      "\\bcos\\b": "\\cos",
      "\\btan\\b": "\\tan",
      "\\blog\\b": "\\log",
      "\\bln\\b": "\\ln",
      "\\bexp\\b": "\\exp",
      "\\bsqrt\\b": "\\sqrt"
    };
    let latex = text;
    let converted = false;
    for (const [pattern, replacement] of Object.entries(expressions)) {
      const regex = new RegExp(pattern, "gi");
      if (regex.test(latex)) {
        latex = latex.replace(regex, replacement);
        converted = true;
      }
    }
    const mathml = this.latexToMathML(latex);
    return {
      original: text,
      converted: latex,
      format: "latex",
      latex,
      mathml,
      success: converted,
      source: "local"
    };
  }
  /**
   * Simple LaTeX to MathML conversion
   * @param {string} latex - LaTeX expression
   * @returns {string} MathML
   */
  latexToMathML(latex) {
    let mathml = latex.replace(/\\int/g, "<mo>&#x222B;</mo>").replace(/\\sum/g, "<mo>&#x2211;</mo>").replace(/\\prod/g, "<mo>&#x220F;</mo>").replace(/\\frac\{([^}]+)\}\{([^}]+)\}/g, "<mfrac><mrow>$1</mrow><mrow>$2</mrow></mfrac>").replace(/\\sqrt\{([^}]+)\}/g, "<msqrt><mrow>$1</mrow></msqrt>").replace(/\\sqrt\[(\d+)\]\{([^}]+)\}/g, "<mroot><mrow>$2</mrow><mn>$1</mn></mroot>").replace(/\^(\{[^}]+\}|\d+)/g, "<msup><mi></mi><mn>$1</mn></msup>").replace(/\\alpha/g, "<mi>&#x03B1;</mi>").replace(/\\beta/g, "<mi>&#x03B2;</mi>").replace(/\\gamma/g, "<mi>&#x03B3;</mi>").replace(/\\delta/g, "<mi>&#x03B4;</mi>").replace(/\\pi/g, "<mi>&#x03C0;</mi>").replace(/\\sigma/g, "<mi>&#x03C3;</mi>").replace(/\\omega/g, "<mi>&#x03C9;</mi>").replace(/\\pm/g, "<mo>&#x00B1;</mo>").replace(/\\times/g, "<mo>&#x00D7;</mo>").replace(/\\div/g, "<mo>&#x00F7;</mo>").replace(/\\infty/g, "<mi>&#x221E;</mi>");
    return `<math xmlns="http://www.w3.org/1998/Math/MathML"><mrow>${mathml}</mrow></math>`;
  }
  /**
   * Local validation fallback
   * @param {string} expression - Expression to validate
   * @returns {Object} Validation result
   */
  localValidate(expression) {
    if (!expression || expression.trim() === "") {
      return { valid: false, message: "Expression is empty" };
    }
    const openCount = (expression.match(/\(/g) || []).length;
    const closeCount = (expression.match(/\)/g) || []).length;
    if (openCount !== closeCount) {
      return { valid: false, message: "Unbalanced parentheses" };
    }
    const openBrace = (expression.match(/\{/g) || []).length;
    const closeBrace = (expression.match(/\}/g) || []).length;
    if (openBrace !== closeBrace) {
      return { valid: false, message: "Unbalanced braces" };
    }
    return { valid: true, message: "Expression appears valid" };
  }
  /**
   * Default functions when API is unavailable
   * @returns {Array} Default functions
   */
  getDefaultFunctions() {
    return [
      { name: "sum", category: "aggregation", description: "Sum of numbers", syntax: "SUM(a, b, ...)" },
      { name: "average", category: "aggregation", description: "Average of numbers", syntax: "AVERAGE(a, b, ...)" },
      { name: "sin", category: "trigonometric", description: "Sine function", syntax: "SIN(x)" },
      { name: "cos", category: "trigonometric", description: "Cosine function", syntax: "COS(x)" },
      { name: "tan", category: "trigonometric", description: "Tangent function", syntax: "TAN(x)" },
      { name: "sqrt", category: "basic", description: "Square root", syntax: "SQRT(x)" },
      { name: "power", category: "basic", description: "Power function", syntax: "POWER(base, exp)" },
      { name: "log", category: "logarithmic", description: "Logarithm", syntax: "LOG(x, [base])" },
      { name: "ln", category: "logarithmic", description: "Natural logarithm", syntax: "LN(x)" },
      { name: "exp", category: "exponential", description: "Exponential function", syntax: "EXP(x)" },
      { name: "integral", category: "calculus", description: "Integration", syntax: "\\int f(x) dx" },
      { name: "derivative", category: "calculus", description: "Differentiation", syntax: "\\frac{d}{dx} f(x)" }
    ];
  }
  /**
   * Default symbols when API is unavailable
   * @returns {Array} Default symbols
   */
  getDefaultSymbols() {
    return [
      { symbol: "\\pi", name: "Pi", category: "greek", unicode: "π" },
      { symbol: "\\alpha", name: "Alpha", category: "greek", unicode: "α" },
      { symbol: "\\beta", name: "Beta", category: "greek", unicode: "β" },
      { symbol: "\\gamma", name: "Gamma", category: "greek", unicode: "γ" },
      { symbol: "\\delta", name: "Delta", category: "greek", unicode: "δ" },
      { symbol: "\\theta", name: "Theta", category: "greek", unicode: "θ" },
      { symbol: "\\sigma", name: "Sigma", category: "greek", unicode: "σ" },
      { symbol: "\\omega", name: "Omega", category: "greek", unicode: "ω" },
      { symbol: "\\infty", name: "Infinity", category: "operators", unicode: "∞" },
      { symbol: "\\sum", name: "Summation", category: "operators", unicode: "∑" },
      { symbol: "\\int", name: "Integral", category: "operators", unicode: "∫" },
      { symbol: "\\pm", name: "Plus-Minus", category: "operators", unicode: "±" },
      { symbol: "\\times", name: "Times", category: "operators", unicode: "×" },
      { symbol: "\\div", name: "Divided by", category: "operators", unicode: "÷" }
    ];
  }
  /**
   * Default templates when API is unavailable
   * @returns {Array} Default templates
   */
  getDefaultTemplates() {
    return [
      {
        id: "pythagorean",
        name: "Pythagorean Theorem",
        category: "geometry",
        latex: "a^{2} + b^{2} = c^{2}",
        description: "The square of the hypotenuse equals the sum of the squares of the other two sides"
      },
      {
        id: "quadratic",
        name: "Quadratic Formula",
        category: "algebra",
        latex: "x = \\frac{-b \\pm \\sqrt{b^{2} - 4ac}}{2a}",
        description: "Solutions to the quadratic equation ax² + bx + c = 0"
      },
      {
        id: "euler",
        name: "Euler's Identity",
        category: "complex_analysis",
        latex: "e^{i\\pi} + 1 = 0",
        description: "Euler's famous identity linking e, i, π, 0, and 1"
      },
      {
        id: "einstein",
        name: "Mass-Energy Equivalence",
        category: "physics",
        latex: "E = mc^{2}",
        description: "Einstein's famous equation relating mass and energy"
      },
      {
        id: "integral",
        name: "Gaussian Integral",
        category: "calculus",
        latex: "\\int_{-\\infty}^{\\infty} e^{-x^{2}} dx = \\sqrt{\\pi}",
        description: "The integral of the Gaussian function over the entire real line"
      }
    ];
  }
  /**
   * Parse natural language math and return all formats
   * @param {string} text - Natural language input
   * @returns {Promise<Object>} Parsed math with all representations
   */
  async parseNaturalMath(text) {
    try {
      const result = await this.convertExpression(text, "text", "latex");
      if (result.success) {
        return {
          text,
          latex: result.latex || result.converted,
          mathml: result.mathml,
          format: "latex",
          source: result.source || "api"
        };
      }
    } catch (error) {
      console.warn("Backend parse failed, using local:", error);
    }
    const local = this.localConvert(text);
    return {
      text,
      latex: local.latex,
      mathml: local.mathml,
      format: "latex",
      source: "local"
    };
  }
  /**
   * Render LaTeX to HTML using KaTeX-style rendering
   * @param {string} latex - LaTeX expression
   * @returns {string} HTML representation
   */
  renderToHTML(latex) {
    return latex.replace(/\\int/g, "∫").replace(/\\sum/g, "∑").replace(/\\prod/g, "∏").replace(/\\frac\{([^}]+)\}\{([^}]+)\}/g, '<span class="fraction"><span class="num">$1</span><span class="den">$2</span></span>').replace(/\\sqrt\{([^}]+)\}/g, "√($1)").replace(/\^(\{[^}]+\}|\d+)/g, "<sup>$1</sup>").replace(/_(\{[^}]+\}|\d+)/g, "<sub>$1</sub>").replace(/\\alpha/g, "α").replace(/\\beta/g, "β").replace(/\\gamma/g, "γ").replace(/\\delta/g, "δ").replace(/\\pi/g, "π").replace(/\\sigma/g, "σ").replace(/\\omega/g, "ω").replace(/\\pm/g, "±").replace(/\\times/g, "×").replace(/\\div/g, "÷").replace(/\\infty/g, "∞");
  }
}
const mathService = new MathService();
const css$2 = {
  code: "textarea.svelte-167pkpq{min-height:1.5em}.math-display.svelte-167pkpq{min-height:1.5em;padding:0.25rem 0}.math-rendered.svelte-167pkpq .fraction{display:inline-flex;flex-direction:column;vertical-align:middle;text-align:center;margin:0 0.2em}.math-rendered.svelte-167pkpq .fraction .num{border-bottom:1px solid currentColor;padding:0 0.2em}.math-rendered.svelte-167pkpq .fraction .den{padding:0 0.2em}.math-rendered.svelte-167pkpq sup{font-size:0.75em;vertical-align:super;line-height:0}.math-rendered.svelte-167pkpq sub{font-size:0.75em;vertical-align:sub;line-height:0}",
  map: `{"version":3,"file":"TextEditorBlockEditor.svelte","sources":["TextEditorBlockEditor.svelte"],"sourcesContent":["<!-- frontend/src/lib/components/TextEditorBlockEditor.svelte -->\\r\\n<script>\\r\\n\\timport { createEventDispatcher, onMount } from 'svelte';\\r\\n\\timport mathService from '$lib/services/math.js';\\r\\n\\r\\n\\texport let blocks = [];\\r\\n\\texport let selectedBlockIndex = 0;\\r\\n\\r\\n\\tconst dispatch = createEventDispatcher();\\r\\n\\r\\n\\t// Track math conversion state\\r\\n\\tlet mathConversions = {};\\r\\n\\tlet convertingMath = {};\\r\\n\\r\\n\\r\\n\\tconst blockTypes = {\\r\\n\\t\\ttext: { icon: '📝', placeholder: 'Type something...' },\\r\\n\\t\\theading1: { icon: 'H1', placeholder: 'Heading 1' },\\r\\n\\t\\theading2: { icon: 'H2', placeholder: 'Heading 2' },\\r\\n\\t\\theading3: { icon: 'H3', placeholder: 'Heading 3' },\\r\\n\\t\\tlist: { icon: '•', placeholder: 'List item' },\\r\\n\\t\\tquote: { icon: '\\"', placeholder: 'Quote' },\\r\\n\\t\\tcode: { icon: '</>', placeholder: 'Code block' },\\r\\n\\t\\tmath: { icon: '∫', placeholder: 'Math expression (e.g., integral of x squared dx)' },\\r\\n\\t\\ttable: { icon: '📊', placeholder: 'Create a table' },\\r\\n\\t\\timage: { icon: '🖼️', placeholder: 'Add an image' }\\r\\n\\t};\\r\\n\\r\\n\\tfunction getBlockStyles(blockType) {\\r\\n\\t\\tconst styles = {\\r\\n\\t\\t\\theading1: 'text-3xl font-bold mb-4',\\r\\n\\t\\t\\theading2: 'text-2xl font-semibold mb-3',\\r\\n\\t\\t\\theading3: 'text-xl font-medium mb-2',\\r\\n\\t\\t\\ttext: 'text-base mb-2',\\r\\n\\t\\t\\tlist: 'text-base mb-1 ml-4',\\r\\n\\t\\t\\tquote: 'text-base mb-2 pl-4 border-l-4 border-gray-300 italic',\\r\\n\\t\\t\\tcode: 'text-sm font-mono bg-gray-100 p-3 rounded mb-2',\\r\\n\\t\\t\\tmath: 'text-base mb-2 font-serif bg-purple-50 p-2 rounded border border-purple-200',\\r\\n\\t\\t\\ttable: 'text-base mb-2',\\r\\n\\t\\t\\timage: 'text-base mb-2'\\r\\n\\t\\t};\\r\\n\\t\\treturn styles[blockType] || styles.text;\\r\\n\\t}\\r\\n\\r\\n\\tasync function parseMathExpression(text) {\\r\\n\\t\\tif (!text || text.trim() === '') return '';\\r\\n\\t\\t\\r\\n\\t\\t// Use the math service for superior natural language processing\\r\\n\\t\\ttry {\\r\\n\\t\\t\\tconst result = await mathService.parseNaturalMath(text);\\r\\n\\t\\t\\treturn result.latex || text;\\r\\n\\t\\t} catch (error) {\\r\\n\\t\\t\\tconsole.error('Math parsing error:', error);\\r\\n\\t\\t\\treturn text;\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tasync function convertMathBlock(blockIndex) {\\r\\n\\t\\tconst block = blocks[blockIndex];\\r\\n\\t\\tif (!block || block.type !== 'math' || !block.content) return;\\r\\n\\t\\t\\r\\n\\t\\tconst content = block.content.trim();\\r\\n\\t\\tif (!content) return;\\r\\n\\t\\t\\r\\n\\t\\t// Check cache first\\r\\n\\t\\tif (mathConversions[content]) {\\r\\n\\t\\t\\tblocks[blockIndex].mathData = mathConversions[content];\\r\\n\\t\\t\\tblocks = [...blocks];\\r\\n\\t\\t\\treturn;\\r\\n\\t\\t}\\r\\n\\t\\t\\r\\n\\t\\tconvertingMath[blockIndex] = true;\\r\\n\\t\\t\\r\\n\\t\\ttry {\\r\\n\\t\\t\\tconst result = await mathService.parseNaturalMath(content);\\r\\n\\t\\t\\tmathConversions[content] = result;\\r\\n\\t\\t\\tblocks[blockIndex].mathData = result;\\r\\n\\t\\t\\tblocks = [...blocks];\\r\\n\\t\\t} catch (error) {\\r\\n\\t\\t\\tconsole.error('Math conversion error:', error);\\r\\n\\t\\t} finally {\\r\\n\\t\\t\\tconvertingMath[blockIndex] = false;\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction renderBlockContent(block) {\\r\\n\\t\\tif (block.type === 'math' && block.mathData) {\\r\\n\\t\\t\\treturn block.mathData.latex || block.content;\\r\\n\\t\\t}\\r\\n\\t\\treturn block.content;\\r\\n\\t}\\r\\n\\r\\n\\tfunction renderMathHTML(block) {\\r\\n\\t\\tif (block.type === 'math' && block.mathData) {\\r\\n\\t\\t\\treturn mathService.renderToHTML(block.mathData.latex);\\r\\n\\t\\t}\\r\\n\\t\\treturn block.content;\\r\\n\\t}\\r\\n\\r\\n\\t// Convert math when block loses focus\\r\\n\\tfunction handleMathBlur(blockIndex) {\\r\\n\\t\\tconvertMathBlock(blockIndex);\\r\\n\\t}\\r\\n\\r\\n\\r\\n\\tfunction handleKeyDown(event, blockIndex) {\\r\\n\\t\\tconst { key } = event;\\r\\n\\r\\n\\t\\tif (key === 'Enter' && !event.shiftKey) {\\r\\n\\t\\t\\tevent.preventDefault();\\r\\n\\t\\t\\tdispatch('addBlock', blockIndex);\\r\\n\\t\\t} else if (key === 'Backspace' && blocks[blockIndex].content === '') {\\r\\n\\t\\t\\tevent.preventDefault();\\r\\n\\t\\t\\tdispatch('deleteBlock', blockIndex);\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction autoResize(node) {\\r\\n\\t\\tconst resize = () => {\\r\\n\\t\\t\\tnode.style.height = 'auto';\\r\\n\\t\\t\\tnode.style.height = node.scrollHeight + 'px';\\r\\n\\t\\t};\\r\\n\\r\\n\\t\\tnode.addEventListener('input', resize);\\r\\n\\t\\tnode.addEventListener('focus', resize);\\r\\n\\t\\tsetTimeout(resize, 0);\\r\\n\\r\\n\\t\\treturn {\\r\\n\\t\\t\\tdestroy() {\\r\\n\\t\\t\\t\\tnode.removeEventListener('input', resize);\\r\\n\\t\\t\\t\\tnode.removeEventListener('focus', resize);\\r\\n\\t\\t\\t}\\r\\n\\t\\t};\\r\\n\\t}\\r\\n<\/script>\\r\\n\\r\\n<div class=\\"space-y-2\\">\\r\\n\\t{#each blocks as block, index}\\r\\n\\t\\t<div\\r\\n\\t\\t\\tclass=\\"group relative {selectedBlockIndex === index ? 'ring-2 ring-blue-500 rounded' : ''}\\"\\r\\n\\t\\t\\ton:click={() => dispatch('selectBlock', index)}\\r\\n\\t\\t>\\r\\n\\t\\t\\t<!-- Block Type Indicator -->\\r\\n\\t\\t\\t<div class=\\"absolute -left-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity\\">\\r\\n\\t\\t\\t\\t<span class=\\"text-xs text-gray-400\\">{blockTypes[block.type]?.icon || '📝'}</span>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<!-- Block Content -->\\r\\n\\t\\t\\t{#if block.type === 'image'}\\r\\n\\t\\t\\t\\t<div class=\\"border-2 border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-500 hover:border-blue-400 transition-colors cursor-pointer\\">\\r\\n\\t\\t\\t\\t\\t<svg class=\\"w-12 h-12 mx-auto mb-3 text-gray-300\\" fill=\\"none\\" stroke=\\"currentColor\\" viewBox=\\"0 0 24 24\\">\\r\\n\\t\\t\\t\\t\\t\\t<path stroke-linecap=\\"round\\" stroke-linejoin=\\"round\\" stroke-width=\\"2\\" d=\\"M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z\\"></path>\\r\\n\\t\\t\\t\\t\\t</svg>\\r\\n\\t\\t\\t\\t\\t<p>Click to add an image</p>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{:else if block.type === 'table'}\\r\\n\\t\\t\\t\\t<div class=\\"border border-gray-300 rounded overflow-hidden\\">\\r\\n\\t\\t\\t\\t\\t<table class=\\"w-full\\">\\r\\n\\t\\t\\t\\t\\t\\t<thead class=\\"bg-gray-50\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<tr>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<th class=\\"border border-gray-300 p-2 text-left\\">Column 1</th>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<th class=\\"border border-gray-300 p-2 text-left\\">Column 2</th>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<th class=\\"border border-gray-300 p-2 text-left\\">Column 3</th>\\r\\n\\t\\t\\t\\t\\t\\t\\t</tr>\\r\\n\\t\\t\\t\\t\\t\\t</thead>\\r\\n\\t\\t\\t\\t\\t\\t<tbody>\\r\\n\\t\\t\\t\\t\\t\\t\\t<tr>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<td class=\\"border border-gray-300 p-2\\">Data 1</td>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<td class=\\"border border-gray-300 p-2\\">Data 2</td>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<td class=\\"border border-gray-300 p-2\\">Data 3</td>\\r\\n\\t\\t\\t\\t\\t\\t\\t</tr>\\r\\n\\t\\t\\t\\t\\t\\t</tbody>\\r\\n\\t\\t\\t\\t\\t</table>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t<div class={getBlockStyles(block.type)}>\\r\\n\\t\\t\\t\\t\\t{#if block.type === 'list'}\\r\\n\\t\\t\\t\\t\\t\\t<span class=\\"mr-2\\">•</span>\\r\\n\\t\\t\\t\\t\\t{:else if block.type === 'quote'}\\r\\n\\t\\t\\t\\t\\t\\t<span class=\\"mr-2 text-gray-400\\">\\"</span>\\r\\n\\t\\t\\t\\t\\t{/if}\\r\\n\\r\\n\\t\\t\\t\\t\\t{#if selectedBlockIndex === index}\\r\\n\\t\\t\\t\\t\\t\\t<textarea\\r\\n\\t\\t\\t\\t\\t\\t\\tbind:value={block.content}\\r\\n\\t\\t\\t\\t\\t\\t\\tplaceholder={blockTypes[block.type]?.placeholder || 'Start typing...'}\\r\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"w-full bg-transparent border-none outline-none resize-none {block.type === 'code' ? 'font-mono' : ''} {block.type === 'math' ? 'font-serif' : ''}\\"\\r\\n\\t\\t\\t\\t\\t\\t\\trows=\\"1\\"\\r\\n\\t\\t\\t\\t\\t\\t\\ton:input={(e) => {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\te.target.style.height = 'auto';\\r\\n\\t\\t\\t\\t\\t\\t\\t\\te.target.style.height = e.target.scrollHeight + 'px';\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tdispatch('contentChange', { index, content: e.target.value });\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t// Clear math data when editing\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tif (block.type === 'math') {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tblock.mathData = null;\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t}\\r\\n\\t\\t\\t\\t\\t\\t\\t}}\\r\\n\\t\\t\\t\\t\\t\\t\\ton:keydown={(e) => handleKeyDown(e, index)}\\r\\n\\t\\t\\t\\t\\t\\t\\ton:focus={() => dispatch('selectBlock', index)}\\r\\n\\t\\t\\t\\t\\t\\t\\ton:blur={() => {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tif (block.type === 'math') {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\thandleMathBlur(index);\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t}\\r\\n\\t\\t\\t\\t\\t\\t\\t}}\\r\\n\\t\\t\\t\\t\\t\\t\\tuse:autoResize\\r\\n\\t\\t\\t\\t\\t\\t/>\\r\\n\\r\\n\\t\\t\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t\\t\\t<div\\r\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"cursor-text\\"\\r\\n\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('selectBlock', index)}\\r\\n\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t{#if block.type === 'math'}\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"math-display font-serif text-lg\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t{#if convertingMath[index]}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-400 italic\\">Converting...</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t{:else if block.mathData}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"math-rendered\\" title=\\"LaTeX: {block.mathData.latex}\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{@html renderMathHTML(block)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">{block.content || blockTypes[block.type]?.placeholder}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t{#if block.mathData}\\r\\n\\t\\t\\t\\t\\t\\t\\t<span class=\\"ml-2 text-xs text-green-600 bg-green-50 px-2 py-1 rounded border border-green-200\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t✓ Math\\r\\n\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t{:else if block.content}\\r\\n\\t\\t\\t\\t\\t\\t\\t<span class=\\"ml-2 text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tPress Enter to convert\\r\\n\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\r\\n\\t\\t\\t\\t\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{block.content || blockTypes[block.type]?.placeholder}\\r\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{/if}\\r\\n\\r\\n\\t\\t\\t<!-- Block Actions (visible on hover) -->\\r\\n\\t\\t\\t<div class=\\"absolute -right-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col space-y-1\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"w-6 h-6 bg-gray-100 hover:bg-gray-200 rounded text-xs\\"\\r\\n\\t\\t\\t\\t\\ttitle=\\"Add block below\\"\\r\\n\\t\\t\\t\\t\\ton:click|stopPropagation={() => dispatch('addBlock', index)}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t+\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t{#if blocks.length > 1}\\r\\n\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\tclass=\\"w-6 h-6 bg-red-100 hover:bg-red-200 rounded text-xs text-red-600\\"\\r\\n\\t\\t\\t\\t\\t\\ttitle=\\"Delete block\\"\\r\\n\\t\\t\\t\\t\\t\\ton:click|stopPropagation={() => dispatch('deleteBlock', index)}\\r\\n\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t{/each}\\r\\n</div>\\r\\n\\r\\n<style>\\r\\n\\ttextarea {\\r\\n\\t\\tmin-height: 1.5em;\\r\\n\\t}\\r\\n\\r\\n\\t.math-display {\\r\\n\\t\\tmin-height: 1.5em;\\r\\n\\t\\tpadding: 0.25rem 0;\\r\\n\\t}\\r\\n\\r\\n\\t.math-rendered :global(.fraction) {\\r\\n\\t\\tdisplay: inline-flex;\\r\\n\\t\\tflex-direction: column;\\r\\n\\t\\tvertical-align: middle;\\r\\n\\t\\ttext-align: center;\\r\\n\\t\\tmargin: 0 0.2em;\\r\\n\\t}\\r\\n\\r\\n\\t.math-rendered :global(.fraction .num) {\\r\\n\\t\\tborder-bottom: 1px solid currentColor;\\r\\n\\t\\tpadding: 0 0.2em;\\r\\n\\t}\\r\\n\\r\\n\\t.math-rendered :global(.fraction .den) {\\r\\n\\t\\tpadding: 0 0.2em;\\r\\n\\t}\\r\\n\\r\\n\\t.math-rendered :global(sup) {\\r\\n\\t\\tfont-size: 0.75em;\\r\\n\\t\\tvertical-align: super;\\r\\n\\t\\tline-height: 0;\\r\\n\\t}\\r\\n\\r\\n\\t.math-rendered :global(sub) {\\r\\n\\t\\tfont-size: 0.75em;\\r\\n\\t\\tvertical-align: sub;\\r\\n\\t\\tline-height: 0;\\r\\n\\t}\\r\\n</style>\\r\\n"],"names":[],"mappings":"AA0QC,uBAAS,CACR,UAAU,CAAE,KACb,CAEA,4BAAc,CACb,UAAU,CAAE,KAAK,CACjB,OAAO,CAAE,OAAO,CAAC,CAClB,CAEA,6BAAc,CAAS,SAAW,CACjC,OAAO,CAAE,WAAW,CACpB,cAAc,CAAE,MAAM,CACtB,cAAc,CAAE,MAAM,CACtB,UAAU,CAAE,MAAM,CAClB,MAAM,CAAE,CAAC,CAAC,KACX,CAEA,6BAAc,CAAS,cAAgB,CACtC,aAAa,CAAE,GAAG,CAAC,KAAK,CAAC,YAAY,CACrC,OAAO,CAAE,CAAC,CAAC,KACZ,CAEA,6BAAc,CAAS,cAAgB,CACtC,OAAO,CAAE,CAAC,CAAC,KACZ,CAEA,6BAAc,CAAS,GAAK,CAC3B,SAAS,CAAE,MAAM,CACjB,cAAc,CAAE,KAAK,CACrB,WAAW,CAAE,CACd,CAEA,6BAAc,CAAS,GAAK,CAC3B,SAAS,CAAE,MAAM,CACjB,cAAc,CAAE,GAAG,CACnB,WAAW,CAAE,CACd"}`
};
function getBlockStyles(blockType) {
  const styles = {
    heading1: "text-3xl font-bold mb-4",
    heading2: "text-2xl font-semibold mb-3",
    heading3: "text-xl font-medium mb-2",
    text: "text-base mb-2",
    list: "text-base mb-1 ml-4",
    quote: "text-base mb-2 pl-4 border-l-4 border-gray-300 italic",
    code: "text-sm font-mono bg-gray-100 p-3 rounded mb-2",
    math: "text-base mb-2 font-serif bg-purple-50 p-2 rounded border border-purple-200",
    table: "text-base mb-2",
    image: "text-base mb-2"
  };
  return styles[blockType] || styles.text;
}
const TextEditorBlockEditor = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { blocks = [] } = $$props;
  let { selectedBlockIndex = 0 } = $$props;
  createEventDispatcher();
  let convertingMath = {};
  const blockTypes = {
    text: {
      icon: "📝",
      placeholder: "Type something..."
    },
    heading1: { icon: "H1", placeholder: "Heading 1" },
    heading2: { icon: "H2", placeholder: "Heading 2" },
    heading3: { icon: "H3", placeholder: "Heading 3" },
    list: { icon: "•", placeholder: "List item" },
    quote: { icon: '"', placeholder: "Quote" },
    code: { icon: "</>", placeholder: "Code block" },
    math: {
      icon: "∫",
      placeholder: "Math expression (e.g., integral of x squared dx)"
    },
    table: {
      icon: "📊",
      placeholder: "Create a table"
    },
    image: { icon: "🖼️", placeholder: "Add an image" }
  };
  function renderMathHTML(block) {
    if (block.type === "math" && block.mathData) {
      return mathService.renderToHTML(block.mathData.latex);
    }
    return block.content;
  }
  if ($$props.blocks === void 0 && $$bindings.blocks && blocks !== void 0) $$bindings.blocks(blocks);
  if ($$props.selectedBlockIndex === void 0 && $$bindings.selectedBlockIndex && selectedBlockIndex !== void 0) $$bindings.selectedBlockIndex(selectedBlockIndex);
  $$result.css.add(css$2);
  return `  <div class="space-y-2">${each(blocks, (block, index) => {
    return `<div class="${"group relative " + escape(
      selectedBlockIndex === index ? "ring-2 ring-blue-500 rounded" : "",
      true
    )}"> <div class="absolute -left-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity"><span class="text-xs text-gray-400">${escape(blockTypes[block.type]?.icon || "📝")}</span></div>  ${block.type === "image" ? `<div class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-500 hover:border-blue-400 transition-colors cursor-pointer" data-svelte-h="svelte-1w4p42z"><svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg> <p>Click to add an image</p> </div>` : `${block.type === "table" ? `<div class="border border-gray-300 rounded overflow-hidden" data-svelte-h="svelte-14cxxhb"><table class="w-full"><thead class="bg-gray-50"><tr><th class="border border-gray-300 p-2 text-left">Column 1</th> <th class="border border-gray-300 p-2 text-left">Column 2</th> <th class="border border-gray-300 p-2 text-left">Column 3</th> </tr></thead> <tbody><tr><td class="border border-gray-300 p-2">Data 1</td> <td class="border border-gray-300 p-2">Data 2</td> <td class="border border-gray-300 p-2">Data 3</td></tr> </tbody></table> </div>` : `<div class="${escape(null_to_empty(getBlockStyles(block.type)), true) + " svelte-167pkpq"}">${block.type === "list" ? `<span class="mr-2" data-svelte-h="svelte-1mk8fyn">•</span>` : `${block.type === "quote" ? `<span class="mr-2 text-gray-400" data-svelte-h="svelte-qvh2wb">&quot;</span>` : ``}`} ${selectedBlockIndex === index ? `<textarea${add_attribute("placeholder", blockTypes[block.type]?.placeholder || "Start typing...", 0)} class="${"w-full bg-transparent border-none outline-none resize-none " + escape(block.type === "code" ? "font-mono" : "", true) + " " + escape(block.type === "math" ? "font-serif" : "", true) + " svelte-167pkpq"}" rows="1">${escape(block.content || "")}</textarea>` : `<div class="cursor-text">${block.type === "math" ? `<div class="math-display font-serif text-lg svelte-167pkpq">${convertingMath[index] ? `<span class="text-gray-400 italic" data-svelte-h="svelte-ikoeps">Converting...</span>` : `${block.mathData ? `<span class="math-rendered svelte-167pkpq" title="${"LaTeX: " + escape(block.mathData.latex, true)}"><!-- HTML_TAG_START -->${renderMathHTML(block)}<!-- HTML_TAG_END --> </span>` : `<span class="text-gray-600">${escape(block.content || blockTypes[block.type]?.placeholder)}</span>`}`}</div> ${block.mathData ? `<span class="ml-2 text-xs text-green-600 bg-green-50 px-2 py-1 rounded border border-green-200" data-svelte-h="svelte-f9ed8x">✓ Math
							</span>` : `${block.content ? `<span class="ml-2 text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded" data-svelte-h="svelte-76ydxc">Press Enter to convert
							</span>` : ``}`}` : `${escape(block.content || blockTypes[block.type]?.placeholder)}`} </div>`} </div>`}`}  <div class="absolute -right-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col space-y-1"><button class="w-6 h-6 bg-gray-100 hover:bg-gray-200 rounded text-xs" title="Add block below" data-svelte-h="svelte-200pkv">+</button> ${blocks.length > 1 ? `<button class="w-6 h-6 bg-red-100 hover:bg-red-200 rounded text-xs text-red-600" title="Delete block" data-svelte-h="svelte-1d8x68u">×
					</button>` : ``}</div> </div>`;
  })} </div>`;
});
const css$1 = {
  code: ".prose.svelte-1u0uidk{color:#374151}",
  map: `{"version":3,"file":"TextEditorMarkdownEditor.svelte","sources":["TextEditorMarkdownEditor.svelte"],"sourcesContent":["<!-- frontend/src/lib/components/TextEditorMarkdownEditor.svelte -->\\r\\n<script>\\r\\n\\timport { createEventDispatcher } from 'svelte';\\r\\n\\r\\n\\texport let markdownContent = '';\\r\\n\\r\\n\\tconst dispatch = createEventDispatcher();\\r\\n\\r\\n\\tconst defaultPreview = \`# Start writing in Markdown...\\r\\n\\r\\n## Features\\r\\n- **Bold text**\\r\\n- *Italic text*\\r\\n- \\\\\`Code snippets\\\\\`\\r\\n- [Links](https://example.com)\\r\\n\\r\\n## Math Expressions\\r\\nYou can write natural math expressions:\\r\\n- integral of x squared dx\\r\\n- fraction a over b\\r\\n- square root of 25\`;\\r\\n<\/script>\\r\\n\\r\\n<div class=\\"grid grid-cols-2 gap-8 h-full\\">\\r\\n\\t<!-- Editor Panel -->\\r\\n\\t<div class=\\"border-r border-gray-200 pr-4\\">\\r\\n\\t\\t<h3 class=\\"text-sm font-medium text-gray-700 mb-2\\">Markdown Source</h3>\\r\\n\\t\\t<textarea\\r\\n\\t\\t\\tbind:value={markdownContent}\\r\\n\\t\\t\\tplaceholder={defaultPreview}\\r\\n\\t\\t\\tclass=\\"w-full h-96 p-3 border border-gray-300 rounded-md font-mono text-sm resize-none focus:ring-blue-500 focus:border-blue-500\\"\\r\\n\\t\\t\\ton:input={() => dispatch('change', markdownContent)}\\r\\n\\t\\t></textarea>\\r\\n\\t</div>\\r\\n\\r\\n\\t<!-- Preview Panel -->\\r\\n\\t<div class=\\"pl-4\\">\\r\\n\\t\\t<h3 class=\\"text-sm font-medium text-gray-700 mb-2\\">Live Preview</h3>\\r\\n\\t\\t<div class=\\"w-full h-96 p-3 border border-gray-300 rounded-md bg-gray-50 overflow-y-auto prose prose-sm max-w-none\\">\\r\\n\\t\\t\\t<div class=\\"whitespace-pre-wrap font-mono text-sm\\">\\r\\n\\t\\t\\t\\t{markdownContent || defaultPreview}\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n</div>\\r\\n\\r\\n<style>\\r\\n\\t.prose {\\r\\n\\t\\tcolor: #374151;\\r\\n\\t}\\r\\n</style>\\r\\n"],"names":[],"mappings":"AA+CC,qBAAO,CACN,KAAK,CAAE,OACR"}`
};
const TextEditorMarkdownEditor = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { markdownContent = "" } = $$props;
  createEventDispatcher();
  const defaultPreview = `# Start writing in Markdown...

## Features
- **Bold text**
- *Italic text*
- \`Code snippets\`
- [Links](https://example.com)

## Math Expressions
You can write natural math expressions:
- integral of x squared dx
- fraction a over b
- square root of 25`;
  if ($$props.markdownContent === void 0 && $$bindings.markdownContent && markdownContent !== void 0) $$bindings.markdownContent(markdownContent);
  $$result.css.add(css$1);
  return `  <div class="grid grid-cols-2 gap-8 h-full"> <div class="border-r border-gray-200 pr-4"><h3 class="text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-1ny1bzh">Markdown Source</h3> <textarea${add_attribute("placeholder", defaultPreview, 0)} class="w-full h-96 p-3 border border-gray-300 rounded-md font-mono text-sm resize-none focus:ring-blue-500 focus:border-blue-500">${escape(markdownContent || "")}</textarea></div>  <div class="pl-4"><h3 class="text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-1f28w2r">Live Preview</h3> <div class="w-full h-96 p-3 border border-gray-300 rounded-md bg-gray-50 overflow-y-auto prose prose-sm max-w-none svelte-1u0uidk"><div class="whitespace-pre-wrap font-mono text-sm">${escape(markdownContent || defaultPreview)}</div></div></div> </div>`;
});
const TextEditorRichText = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { editorElement = null } = $$props;
  createEventDispatcher();
  if ($$props.editorElement === void 0 && $$bindings.editorElement && editorElement !== void 0) $$bindings.editorElement(editorElement);
  return `  <div class="border border-gray-300 rounded-lg overflow-hidden"> <div class="bg-gray-50 px-4 py-2 border-b border-gray-200"><div class="flex items-center space-x-2 text-sm"><select class="border border-gray-300 rounded px-2 py-1"><option value="Arial" data-svelte-h="svelte-1wf7cmk">Arial</option><option value="Times New Roman" data-svelte-h="svelte-ik5oy4">Times New Roman</option><option value="Inter" data-svelte-h="svelte-3k6lry">Inter</option><option value="Georgia" data-svelte-h="svelte-bg3fk6">Georgia</option><option value="Verdana" data-svelte-h="svelte-m38ik">Verdana</option></select> <select class="border border-gray-300 rounded px-2 py-1"><option value="12pt" data-svelte-h="svelte-myodem">12pt</option><option value="14pt" data-svelte-h="svelte-1sf8mkq">14pt</option><option value="16pt" data-svelte-h="svelte-ljywja">16pt</option><option value="18pt" data-svelte-h="svelte-1if3aiy">18pt</option><option value="24pt" data-svelte-h="svelte-4gj6by">24pt</option><option value="32pt" data-svelte-h="svelte-1brm9qm">32pt</option></select> <div class="w-px h-6 bg-gray-300"></div> <button class="${"p-1 hover:bg-white rounded " + escape("", true)}" title="Bold"><strong data-svelte-h="svelte-1dz6n38">B</strong></button> <button class="${"p-1 hover:bg-white rounded " + escape("", true)}" title="Italic"><em data-svelte-h="svelte-vhbnu7">I</em></button> <button class="${"p-1 hover:bg-white rounded " + escape("", true)}" title="Underline"><u data-svelte-h="svelte-1t5mkex">U</u></button> <div class="w-px h-6 bg-gray-300"></div> <button class="${"p-1 hover:bg-white rounded " + escape("bg-blue-100", true)}" title="Align Left">⬅️</button> <button class="${"p-1 hover:bg-white rounded " + escape("", true)}" title="Align Center">⬌</button> <button class="${"p-1 hover:bg-white rounded " + escape("", true)}" title="Align Right">➡️</button> <div class="w-px h-6 bg-gray-300"></div> <button class="p-1 hover:bg-white rounded" title="Insert Link" data-svelte-h="svelte-8btr1g">🔗</button> <button class="p-1 hover:bg-white rounded" title="Find &amp; Replace" data-svelte-h="svelte-g5s8f9">🔍</button></div></div>  <div contenteditable="true" class="min-h-96 p-4 focus:outline-none" placeholder="Start writing..."${add_attribute("this", editorElement, 0)} data-svelte-h="svelte-1kbn0t5"><h1>Welcome to TPT Rich Text Editor</h1> <p>This is a traditional word processor interface with all the formatting options you expect.</p> <p>You can also insert <strong>natural math expressions</strong> like &quot;integral from 0 to π of sin(x) dx&quot; which will be automatically converted to proper mathematical notation.</p></div></div>`;
});
const css = {
  code: ".prose.svelte-1u0uidk{color:#374151}",
  map: `{"version":3,"file":"TextEditorModals.svelte","sources":["TextEditorModals.svelte"],"sourcesContent":["<!-- frontend/src/lib/components/TextEditorModals.svelte -->\\r\\n<script>\\r\\n\\timport { createEventDispatcher } from 'svelte';\\r\\n\\r\\n\\t// Document list modal\\r\\n\\texport let showDocumentList = false;\\r\\n\\texport let documentList = [];\\r\\n\\texport let isLoading = false;\\r\\n\\r\\n\\t// Version history modal\\r\\n\\texport let showVersionHistory = false;\\r\\n\\texport let versionHistory = [];\\r\\n\\r\\n\\t// Math help modal\\r\\n\\texport let showMathHelp = false;\\r\\n\\r\\n\\t// Find & replace modal\\r\\n\\texport let showFindReplaceDialog = false;\\r\\n\\r\\n\\t// AI assistant modal\\r\\n\\texport let showAIAssistant = false;\\r\\n\\texport let aiSuggestions = [];\\r\\n\\r\\n\\t// Document summary modal\\r\\n\\texport let documentSummary = '';\\r\\n\\r\\n\\t// Text analysis modal\\r\\n\\texport let showTextAnalysis = false;\\r\\n\\texport let textAnalysis = null;\\r\\n\\r\\n\\t// Find/replace state\\r\\n\\texport let findText = '';\\r\\n\\texport let replaceText = '';\\r\\n\\texport let findResults = [];\\r\\n\\texport let currentFindIndex = -1;\\r\\n\\texport let isCaseSensitive = false;\\r\\n\\r\\n\\tconst dispatch = createEventDispatcher();\\r\\n\\r\\n<\/script>\\r\\n\\r\\n<!-- ── Document List Modal ─────────────────────────────────────── -->\\r\\n{#if showDocumentList}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-md max-h-96 overflow-y-auto\\">\\r\\n\\t\\t\\t<h3 class=\\"text-lg font-semibold mb-4\\">Open Document</h3>\\r\\n\\t\\t\\t{#if documentList.length === 0}\\r\\n\\t\\t\\t\\t<p class=\\"text-gray-500\\">No documents found.</p>\\r\\n\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t<div class=\\"space-y-2\\">\\r\\n\\t\\t\\t\\t\\t{#each documentList as doc}\\r\\n\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"w-full text-left p-3 border border-gray-200 rounded hover:bg-gray-50\\"\\r\\n\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('loadDocument', doc)}\\r\\n\\t\\t\\t\\t\\t\\t\\tdisabled={isLoading}\\r\\n\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-medium\\">{doc.title}</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-500\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tv{doc.version} • {new Date(doc.updated_at).toLocaleDateString()}\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{/if}\\r\\n\\t\\t\\t<div class=\\"mt-4 flex justify-end\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeDocumentList')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tClose\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n<!-- ── Version History Modal ───────────────────────────────────── -->\\r\\n{#if showVersionHistory}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-lg max-h-96 overflow-y-auto\\">\\r\\n\\t\\t\\t<h3 class=\\"text-lg font-semibold mb-4\\">Version History</h3>\\r\\n\\t\\t\\t<div class=\\"space-y-2\\">\\r\\n\\t\\t\\t\\t{#each versionHistory as version}\\r\\n\\t\\t\\t\\t\\t<div class=\\"flex items-center justify-between p-3 border border-gray-200 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-medium\\">Version {version.version}</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-500\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{new Date(version.created_at).toLocaleString()}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{#if version.is_active}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"ml-2 text-green-600\\">(Current)</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t{#if !version.is_active}\\r\\n\\t\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('restoreVersion', version)}\\r\\n\\t\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tRestore\\r\\n\\t\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t\\t<div class=\\"mt-4 flex justify-end\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeVersionHistory')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tClose\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n<!-- ── Math Help Modal ─────────────────────────────────────────── -->\\r\\n{#if showMathHelp}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto\\">\\r\\n\\t\\t\\t<div class=\\"flex items-center justify-between mb-6\\">\\r\\n\\t\\t\\t\\t<h3 class=\\"text-xl font-semibold text-gray-900\\">Natural Math Input - Better Than LaTeX!</h3>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"text-gray-400 hover:text-gray-600 text-2xl\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeMathHelp')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<div class=\\"space-y-6\\">\\r\\n\\t\\t\\t\\t<div class=\\"bg-blue-50 p-4 rounded-lg\\">\\r\\n\\t\\t\\t\\t\\t<h4 class=\\"text-lg font-medium text-blue-900 mb-2\\">✨ Revolutionary Math Input</h4>\\r\\n\\t\\t\\t\\t\\t<p class=\\"text-blue-800\\">\\r\\n\\t\\t\\t\\t\\t\\tForget complex LaTeX syntax! TPT Titan understands natural language math expressions.\\r\\n\\t\\t\\t\\t\\t\\tSimply type what you would say out loud, and watch it transform into beautiful mathematical notation.\\r\\n\\t\\t\\t\\t\\t</p>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t<!-- Basic Examples -->\\r\\n\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t<h4 class=\\"text-lg font-medium text-gray-900 mb-3\\">📚 Basic Expressions</h4>\\r\\n\\t\\t\\t\\t\\t<div class=\\"grid grid-cols-1 md:grid-cols-2 gap-4\\">\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-gray-50 p-3 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mb-1\\">Type this:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-2 py-1 rounded text-sm\\">fraction a over b</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mt-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-lg\\">a/b</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-gray-50 p-3 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mb-1\\">Type this:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-2 py-1 rounded text-sm\\">square root of x</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mt-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-lg\\">√x</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-gray-50 p-3 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mb-1\\">Type this:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-2 py-1 rounded text-sm\\">pi times r squared</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mt-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-lg\\">π × r²</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-gray-50 p-3 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mb-1\\">Type this:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-2 py-1 rounded text-sm\\">alpha beta gamma</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-600 mt-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-lg\\">αβγ</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t<!-- Advanced Examples -->\\r\\n\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t<h4 class=\\"text-lg font-medium text-gray-900 mb-3\\">🔬 Advanced Mathematics</h4>\\r\\n\\t\\t\\t\\t\\t<div class=\\"space-y-4\\">\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-purple-50 p-4 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-purple-700 mb-2\\">Integrals:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-3 py-2 rounded text-sm block mb-2\\">integral from 0 to infinity of e to the power of negative x squared dx</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-purple-700 mb-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-xl\\">∫₀^∞ e^(-x²) dx</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-green-50 p-4 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-green-700 mb-2\\">Summations:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-3 py-2 rounded text-sm block mb-2\\">sum from i equals 1 to n of x sub i</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-green-700 mb-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-xl\\">∑_&#123;i=1&#125;^n x_i</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"bg-orange-50 p-4 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-orange-700 mb-2\\">Complex fractions:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<code class=\\"bg-white px-3 py-2 rounded text-sm block mb-2\\">fraction x plus y over z minus w</code>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-orange-700 mb-1\\">Becomes:</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"font-serif text-xl\\">x+y/z-w</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t<!-- Quick Reference -->\\r\\n\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t<h4 class=\\"text-lg font-medium text-gray-900 mb-3\\">🚀 Quick Reference</h4>\\r\\n\\t\\t\\t\\t\\t<div class=\\"grid grid-cols-2 md:grid-cols-3 gap-3 text-sm\\">\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>fractions:</strong> \\"fraction a over b\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>roots:</strong> \\"square root of\\", \\"cube root of\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>integrals:</strong> \\"integral of\\", \\"integral from a to b of\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>summations:</strong> \\"sum from i=1 to n of\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>greek:</strong> \\"alpha\\", \\"beta\\", \\"gamma\\", \\"pi\\", \\"sigma\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>operators:</strong> \\"times\\", \\"divided by\\", \\"plus or minus\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>superscripts:</strong> \\"squared\\", \\"cubed\\", \\"to the power of\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>subscripts:</strong> \\"sub\\", \\"to the\\"</div>\\r\\n\\t\\t\\t\\t\\t\\t<div><strong>symbols:</strong> \\"therefore\\", \\"because\\", \\"infinity\\"</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t<!-- Try It Section -->\\r\\n\\t\\t\\t\\t<div class=\\"bg-yellow-50 p-4 rounded\\">\\r\\n\\t\\t\\t\\t\\t<h4 class=\\"text-lg font-medium text-yellow-900 mb-2\\">🎯 Try It Now!</h4>\\r\\n\\t\\t\\t\\t\\t<p class=\\"text-yellow-800 mb-3\\">\\r\\n\\t\\t\\t\\t\\t\\tCreate a new math block and try typing any of these expressions:\\r\\n\\t\\t\\t\\t\\t</p>\\r\\n\\t\\t\\t\\t\\t<div class=\\"grid grid-cols-1 md:grid-cols-2 gap-2 text-sm\\">\\r\\n\\t\\t\\t\\t\\t\\t<button class=\\"text-left p-2 bg-white rounded hover:bg-gray-50\\" on:click={() => dispatch('closeMathHelp')}>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\"fraction x squared plus y squared over z\\"\\r\\n\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t<button class=\\"text-left p-2 bg-white rounded hover:bg-gray-50\\" on:click={() => dispatch('closeMathHelp')}>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\"integral from negative infinity to infinity of e to the power of negative x squared over square root of pi dx\\"\\r\\n\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t<button class=\\"text-left p-2 bg-white rounded hover:bg-gray-50\\" on:click={() => dispatch('closeMathHelp')}>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\"sum from k equals 0 to infinity of fraction 1 over k factorial\\"\\r\\n\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t<button class=\\"text-left p-2 bg-white rounded hover:bg-gray-50\\" on:click={() => dispatch('closeMathHelp')}>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\"limit as x approaches 0 of fraction sin x over x equals 1\\"\\r\\n\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t<!-- Export Options -->\\r\\n\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t<h4 class=\\"text-lg font-medium text-gray-900 mb-3\\">📤 Export Your Math</h4>\\r\\n\\t\\t\\t\\t\\t<p class=\\"text-gray-600 mb-3\\">\\r\\n\\t\\t\\t\\t\\t\\tOnce you've created beautiful math expressions, export them to various formats:\\r\\n\\t\\t\\t\\t\\t</p>\\r\\n\\t\\t\\t\\t\\t<div class=\\"grid grid-cols-2 md:grid-cols-4 gap-3\\">\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"text-center p-3 bg-red-50 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-red-600 font-medium\\">PDF</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-xs text-red-600\\">Vector graphics</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"text-center p-3 bg-blue-50 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-blue-600 font-medium\\">LaTeX</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-xs text-blue-600\\">Source code</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"text-center p-3 bg-green-50 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-green-600 font-medium\\">MathML</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-xs text-green-600\\">Web standard</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"text-center p-3 bg-purple-50 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-purple-600 font-medium\\">SVG</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"text-xs text-purple-600\\">Scalable images</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<div class=\\"mt-6 flex justify-end\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeMathHelp')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tGot it!\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n<!-- ── Find & Replace Modal ────────────────────────────────────── -->\\r\\n{#if showFindReplaceDialog}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-lg\\">\\r\\n\\t\\t\\t<div class=\\"flex items-center justify-between mb-6\\">\\r\\n\\t\\t\\t\\t<h3 class=\\"text-xl font-semibold text-gray-900\\">Find & Replace</h3>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"text-gray-400 hover:text-gray-600 text-2xl\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeFindReplace')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<div class=\\"space-y-4\\">\\r\\n\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t<label class=\\"block text-sm font-medium text-gray-700 mb-2\\">Find</label>\\r\\n\\t\\t\\t\\t\\t<div class=\\"flex space-x-2\\">\\r\\n\\t\\t\\t\\t\\t\\t<input\\r\\n\\t\\t\\t\\t\\t\\t\\ttype=\\"text\\"\\r\\n\\t\\t\\t\\t\\t\\t\\tbind:value={findText}\\r\\n\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"Enter text to find...\\"\\r\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500\\"\\r\\n\\t\\t\\t\\t\\t\\t\\ton:keydown={(e) => {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tif (e.key === 'Enter') dispatch('performFind');\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tif (e.key === 'Escape') dispatch('closeFindReplace');\\r\\n\\t\\t\\t\\t\\t\\t\\t}}\\r\\n\\t\\t\\t\\t\\t\\t/>\\r\\n\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200\\"\\r\\n\\t\\t\\t\\t\\t\\t\\tclass:bg-blue-100={isCaseSensitive}\\r\\n\\t\\t\\t\\t\\t\\t\\tclass:text-blue-700={isCaseSensitive}\\r\\n\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('toggleCaseSensitive')}\\r\\n\\t\\t\\t\\t\\t\\t\\ttitle=\\"Case sensitive\\"\\r\\n\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\tAa\\r\\n\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t<div>\\r\\n\\t\\t\\t\\t\\t<label class=\\"block text-sm font-medium text-gray-700 mb-2\\">Replace with</label>\\r\\n\\t\\t\\t\\t\\t<input\\r\\n\\t\\t\\t\\t\\t\\ttype=\\"text\\"\\r\\n\\t\\t\\t\\t\\t\\tbind:value={replaceText}\\r\\n\\t\\t\\t\\t\\t\\tplaceholder=\\"Enter replacement text...\\"\\r\\n\\t\\t\\t\\t\\t\\tclass=\\"w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500\\"\\r\\n\\t\\t\\t\\t\\t\\ton:keydown={(e) => {\\r\\n\\t\\t\\t\\t\\t\\t\\tif (e.key === 'Enter') dispatch('performReplace');\\r\\n\\t\\t\\t\\t\\t\\t\\tif (e.key === 'Escape') dispatch('closeFindReplace');\\r\\n\\t\\t\\t\\t\\t\\t}}\\r\\n\\t\\t\\t\\t\\t/>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t{#if findResults.length > 0}\\r\\n\\t\\t\\t\\t\\t<div class=\\"flex items-center justify-between bg-blue-50 p-3 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"text-sm text-blue-800\\">\\r\\n\\t\\t\\t\\t\\t\\t\\tFound {findResults.length} occurrence{findResults.length !== 1 ? 's' : ''}\\r\\n\\t\\t\\t\\t\\t\\t\\t{#if currentFindIndex >= 0}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">(currently at {currentFindIndex + 1})</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"flex space-x-2\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"px-2 py-1 text-sm bg-white text-gray-700 rounded hover:bg-gray-100\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('findPrevious')}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ttitle=\\"Previous (Shift+Enter)\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t◀ Previous\\r\\n\\t\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"px-2 py-1 text-sm bg-white text-gray-700 rounded hover:bg-gray-100\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('findNext')}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ttitle=\\"Next (Enter)\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tNext ▶\\r\\n\\t\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t{:else if findText}\\r\\n\\t\\t\\t\\t\\t<div class=\\"text-sm text-gray-500 bg-gray-50 p-3 rounded\\">\\r\\n\\t\\t\\t\\t\\t\\tNo matches found\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<div class=\\"mt-6 flex justify-end space-x-3\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeFindReplace')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tClose\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('performFind')}\\r\\n\\t\\t\\t\\t\\tdisabled={!findText}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tFind\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t{#if findResults.length > 0}\\r\\n\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700\\"\\r\\n\\t\\t\\t\\t\\t\\ton:click={() => dispatch('performReplace')}\\r\\n\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\tReplace\\r\\n\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700\\"\\r\\n\\t\\t\\t\\t\\t\\ton:click={() => dispatch('replaceAll')}\\r\\n\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\tReplace All\\r\\n\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n\\r\\n<!-- ── AI Writing Assistant Modal ─────────────────────────────── -->\\r\\n{#if showAIAssistant}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto\\">\\r\\n\\t\\t\\t<div class=\\"flex items-center justify-between mb-6\\">\\r\\n\\t\\t\\t\\t<h3 class=\\"text-xl font-semibold text-gray-900\\">✨ AI Writing Suggestions</h3>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"text-gray-400 hover:text-gray-600 text-2xl\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeAIAssistant')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t{#if aiSuggestions.length === 0}\\r\\n\\t\\t\\t\\t<div class=\\"text-center py-8\\">\\r\\n\\t\\t\\t\\t\\t<div class=\\"text-4xl mb-4\\">💭</div>\\r\\n\\t\\t\\t\\t\\t<p class=\\"text-gray-600\\">Generating suggestions...</p>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t<div class=\\"space-y-4\\">\\r\\n\\t\\t\\t\\t\\t{#each aiSuggestions as suggestion, index}\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"border border-gray-200 rounded-lg p-4\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex items-start justify-between mb-2\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<h4 class=\\"font-medium text-gray-900\\">Suggestion {index + 1}</h4>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatch('applyAISuggestion', suggestion)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tApply\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<p class=\\"text-gray-700 whitespace-pre-wrap\\">{suggestion}</p>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{/if}\\r\\n\\r\\n\\t\\t\\t<div class=\\"mt-6 flex justify-end\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeAIAssistant')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tClose\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n<!-- ── Document Summary Modal ─────────────────────────────────── -->\\r\\n{#if documentSummary}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-3xl max-h-96 overflow-y-auto\\">\\r\\n\\t\\t\\t<div class=\\"flex items-center justify-between mb-6\\">\\r\\n\\t\\t\\t\\t<h3 class=\\"text-xl font-semibold text-gray-900\\">📄 Document Summary</h3>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"text-gray-400 hover:text-gray-600 text-2xl\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeSummary')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<div class=\\"prose prose-sm max-w-none\\">\\r\\n\\t\\t\\t\\t<div class=\\"whitespace-pre-wrap text-gray-700\\">{documentSummary}</div>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t<div class=\\"mt-6 flex justify-end space-x-3\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => navigator.clipboard.writeText(documentSummary)}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tCopy to Clipboard\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeSummary')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tClose\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n<!-- ── Text Analysis Modal ─────────────────────────────────────── -->\\r\\n{#if showTextAnalysis}\\r\\n\\t<div class=\\"fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50\\">\\r\\n\\t\\t<div class=\\"bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto\\">\\r\\n\\t\\t\\t<div class=\\"flex items-center justify-between mb-6\\">\\r\\n\\t\\t\\t\\t<h3 class=\\"text-xl font-semibold text-gray-900\\">📊 Text Analysis</h3>\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"text-gray-400 hover:text-gray-600 text-2xl\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeTextAnalysis')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t{#if textAnalysis}\\r\\n\\t\\t\\t\\t<div class=\\"grid grid-cols-2 gap-6\\">\\r\\n\\t\\t\\t\\t\\t<div class=\\"space-y-4\\">\\r\\n\\t\\t\\t\\t\\t\\t<h4 class=\\"font-medium text-gray-900\\">Statistics</h4>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"space-y-2 text-sm\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Words:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.word_count || 0}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Characters:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.char_count || 0}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Sentences:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.sentence_count || 0}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Paragraphs:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.paragraph_count || 0}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Reading Time:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.reading_time || 'N/A'}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t\\t<div class=\\"space-y-4\\">\\r\\n\\t\\t\\t\\t\\t\\t<h4 class=\\"font-medium text-gray-900\\">Readability</h4>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"space-y-2 text-sm\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Grade Level:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.grade_level || 'N/A'}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Readability:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.readability_score || 'N/A'}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex justify-between\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-gray-600\\">Complexity:</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"font-medium\\">{textAnalysis.complexity || 'N/A'}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t{#if textAnalysis.key_phrases && textAnalysis.key_phrases.length > 0}\\r\\n\\t\\t\\t\\t\\t<div class=\\"mt-6\\">\\r\\n\\t\\t\\t\\t\\t\\t<h4 class=\\"font-medium text-gray-900 mb-3\\">Key Phrases</h4>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"flex flex-wrap gap-2\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t{#each textAnalysis.key_phrases as phrase}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{phrase}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t{/if}\\r\\n\\r\\n\\t\\t\\t\\t{#if textAnalysis.sentiment}\\r\\n\\t\\t\\t\\t\\t<div class=\\"mt-6\\">\\r\\n\\t\\t\\t\\t\\t\\t<h4 class=\\"font-medium text-gray-900 mb-3\\">Sentiment Analysis</h4>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"flex items-center space-x-4\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex-1 bg-gray-200 rounded-full h-2\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<div\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"bg-green-500 h-2 rounded-full\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tstyle=\\"width: {textAnalysis.sentiment.positive || 0}%\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t></div>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-sm text-gray-600\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tPositive: {textAnalysis.sentiment.positive || 0}%\\r\\n\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"flex items-center space-x-4 mt-2\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"flex-1 bg-gray-200 rounded-full h-2\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<div\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"bg-red-500 h-2 rounded-full\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tstyle=\\"width: {textAnalysis.sentiment.negative || 0}%\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t></div>\\r\\n\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-sm text-gray-600\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tNegative: {textAnalysis.sentiment.negative || 0}%\\r\\n\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t<div class=\\"text-center py-8\\">\\r\\n\\t\\t\\t\\t\\t<p class=\\"text-gray-600\\">No analysis data available</p>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{/if}\\r\\n\\r\\n\\t\\t\\t<div class=\\"mt-6 flex justify-end\\">\\r\\n\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\tclass=\\"px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700\\"\\r\\n\\t\\t\\t\\t\\ton:click={() => dispatch('closeTextAnalysis')}\\r\\n\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\tClose\\r\\n\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t</div>\\r\\n\\t</div>\\r\\n{/if}\\r\\n\\r\\n<style>\\r\\n\\t.prose {\\r\\n\\t\\tcolor: #374151;\\r\\n\\t}\\r\\n</style>\\r\\n"],"names":[],"mappings":"AAulBC,qBAAO,CACN,KAAK,CAAE,OACR"}`
};
const TextEditorModals = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { showDocumentList = false } = $$props;
  let { documentList = [] } = $$props;
  let { isLoading = false } = $$props;
  let { showVersionHistory = false } = $$props;
  let { versionHistory = [] } = $$props;
  let { showMathHelp = false } = $$props;
  let { showFindReplaceDialog = false } = $$props;
  let { showAIAssistant = false } = $$props;
  let { aiSuggestions = [] } = $$props;
  let { documentSummary = "" } = $$props;
  let { showTextAnalysis = false } = $$props;
  let { textAnalysis = null } = $$props;
  let { findText = "" } = $$props;
  let { replaceText = "" } = $$props;
  let { findResults = [] } = $$props;
  let { currentFindIndex = -1 } = $$props;
  let { isCaseSensitive = false } = $$props;
  createEventDispatcher();
  if ($$props.showDocumentList === void 0 && $$bindings.showDocumentList && showDocumentList !== void 0) $$bindings.showDocumentList(showDocumentList);
  if ($$props.documentList === void 0 && $$bindings.documentList && documentList !== void 0) $$bindings.documentList(documentList);
  if ($$props.isLoading === void 0 && $$bindings.isLoading && isLoading !== void 0) $$bindings.isLoading(isLoading);
  if ($$props.showVersionHistory === void 0 && $$bindings.showVersionHistory && showVersionHistory !== void 0) $$bindings.showVersionHistory(showVersionHistory);
  if ($$props.versionHistory === void 0 && $$bindings.versionHistory && versionHistory !== void 0) $$bindings.versionHistory(versionHistory);
  if ($$props.showMathHelp === void 0 && $$bindings.showMathHelp && showMathHelp !== void 0) $$bindings.showMathHelp(showMathHelp);
  if ($$props.showFindReplaceDialog === void 0 && $$bindings.showFindReplaceDialog && showFindReplaceDialog !== void 0) $$bindings.showFindReplaceDialog(showFindReplaceDialog);
  if ($$props.showAIAssistant === void 0 && $$bindings.showAIAssistant && showAIAssistant !== void 0) $$bindings.showAIAssistant(showAIAssistant);
  if ($$props.aiSuggestions === void 0 && $$bindings.aiSuggestions && aiSuggestions !== void 0) $$bindings.aiSuggestions(aiSuggestions);
  if ($$props.documentSummary === void 0 && $$bindings.documentSummary && documentSummary !== void 0) $$bindings.documentSummary(documentSummary);
  if ($$props.showTextAnalysis === void 0 && $$bindings.showTextAnalysis && showTextAnalysis !== void 0) $$bindings.showTextAnalysis(showTextAnalysis);
  if ($$props.textAnalysis === void 0 && $$bindings.textAnalysis && textAnalysis !== void 0) $$bindings.textAnalysis(textAnalysis);
  if ($$props.findText === void 0 && $$bindings.findText && findText !== void 0) $$bindings.findText(findText);
  if ($$props.replaceText === void 0 && $$bindings.replaceText && replaceText !== void 0) $$bindings.replaceText(replaceText);
  if ($$props.findResults === void 0 && $$bindings.findResults && findResults !== void 0) $$bindings.findResults(findResults);
  if ($$props.currentFindIndex === void 0 && $$bindings.currentFindIndex && currentFindIndex !== void 0) $$bindings.currentFindIndex(currentFindIndex);
  if ($$props.isCaseSensitive === void 0 && $$bindings.isCaseSensitive && isCaseSensitive !== void 0) $$bindings.isCaseSensitive(isCaseSensitive);
  $$result.css.add(css);
  return `   ${showDocumentList ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-md max-h-96 overflow-y-auto"><h3 class="text-lg font-semibold mb-4" data-svelte-h="svelte-3av5fm">Open Document</h3> ${documentList.length === 0 ? `<p class="text-gray-500" data-svelte-h="svelte-1l87ux9">No documents found.</p>` : `<div class="space-y-2">${each(documentList, (doc) => {
    return `<button class="w-full text-left p-3 border border-gray-200 rounded hover:bg-gray-50" ${isLoading ? "disabled" : ""}><div class="font-medium">${escape(doc.title)}</div> <div class="text-sm text-gray-500">v${escape(doc.version)} • ${escape(new Date(doc.updated_at).toLocaleDateString())}</div> </button>`;
  })}</div>`} <div class="mt-4 flex justify-end"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-1si84i6">Close</button></div></div></div>` : ``}  ${showVersionHistory ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-lg max-h-96 overflow-y-auto"><h3 class="text-lg font-semibold mb-4" data-svelte-h="svelte-pnms9r">Version History</h3> <div class="space-y-2">${each(versionHistory, (version) => {
    return `<div class="flex items-center justify-between p-3 border border-gray-200 rounded"><div><div class="font-medium">Version ${escape(version.version)}</div> <div class="text-sm text-gray-500">${escape(new Date(version.created_at).toLocaleString())} ${version.is_active ? `<span class="ml-2 text-green-600" data-svelte-h="svelte-1pj9guh">(Current)</span>` : ``} </div></div> ${!version.is_active ? `<button class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-1xe4fqu">Restore
							</button>` : ``} </div>`;
  })}</div> <div class="mt-4 flex justify-end"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-1lu65sr">Close</button></div></div></div>` : ``}  ${showMathHelp ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-1aus7jn">Natural Math Input - Better Than LaTeX!</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-1gq38ux">×</button></div> <div class="space-y-6"><div class="bg-blue-50 p-4 rounded-lg" data-svelte-h="svelte-12x5a4b"><h4 class="text-lg font-medium text-blue-900 mb-2">✨ Revolutionary Math Input</h4> <p class="text-blue-800">Forget complex LaTeX syntax! TPT Titan understands natural language math expressions.
						Simply type what you would say out loud, and watch it transform into beautiful mathematical notation.</p></div>  <div data-svelte-h="svelte-1d09dpr"><h4 class="text-lg font-medium text-gray-900 mb-3">📚 Basic Expressions</h4> <div class="grid grid-cols-1 md:grid-cols-2 gap-4"><div class="bg-gray-50 p-3 rounded"><div class="text-sm text-gray-600 mb-1">Type this:</div> <code class="bg-white px-2 py-1 rounded text-sm">fraction a over b</code> <div class="text-sm text-gray-600 mt-1">Becomes:</div> <div class="font-serif text-lg">a/b</div></div> <div class="bg-gray-50 p-3 rounded"><div class="text-sm text-gray-600 mb-1">Type this:</div> <code class="bg-white px-2 py-1 rounded text-sm">square root of x</code> <div class="text-sm text-gray-600 mt-1">Becomes:</div> <div class="font-serif text-lg">√x</div></div> <div class="bg-gray-50 p-3 rounded"><div class="text-sm text-gray-600 mb-1">Type this:</div> <code class="bg-white px-2 py-1 rounded text-sm">pi times r squared</code> <div class="text-sm text-gray-600 mt-1">Becomes:</div> <div class="font-serif text-lg">π × r²</div></div> <div class="bg-gray-50 p-3 rounded"><div class="text-sm text-gray-600 mb-1">Type this:</div> <code class="bg-white px-2 py-1 rounded text-sm">alpha beta gamma</code> <div class="text-sm text-gray-600 mt-1">Becomes:</div> <div class="font-serif text-lg">αβγ</div></div></div></div>  <div data-svelte-h="svelte-w2e35p"><h4 class="text-lg font-medium text-gray-900 mb-3">🔬 Advanced Mathematics</h4> <div class="space-y-4"><div class="bg-purple-50 p-4 rounded"><div class="text-sm text-purple-700 mb-2">Integrals:</div> <code class="bg-white px-3 py-2 rounded text-sm block mb-2">integral from 0 to infinity of e to the power of negative x squared dx</code> <div class="text-sm text-purple-700 mb-1">Becomes:</div> <div class="font-serif text-xl">∫₀^∞ e^(-x²) dx</div></div> <div class="bg-green-50 p-4 rounded"><div class="text-sm text-green-700 mb-2">Summations:</div> <code class="bg-white px-3 py-2 rounded text-sm block mb-2">sum from i equals 1 to n of x sub i</code> <div class="text-sm text-green-700 mb-1">Becomes:</div> <div class="font-serif text-xl">∑_{i=1}^n x_i</div></div> <div class="bg-orange-50 p-4 rounded"><div class="text-sm text-orange-700 mb-2">Complex fractions:</div> <code class="bg-white px-3 py-2 rounded text-sm block mb-2">fraction x plus y over z minus w</code> <div class="text-sm text-orange-700 mb-1">Becomes:</div> <div class="font-serif text-xl">x+y/z-w</div></div></div></div>  <div data-svelte-h="svelte-weuefl"><h4 class="text-lg font-medium text-gray-900 mb-3">🚀 Quick Reference</h4> <div class="grid grid-cols-2 md:grid-cols-3 gap-3 text-sm"><div><strong>fractions:</strong> &quot;fraction a over b&quot;</div> <div><strong>roots:</strong> &quot;square root of&quot;, &quot;cube root of&quot;</div> <div><strong>integrals:</strong> &quot;integral of&quot;, &quot;integral from a to b of&quot;</div> <div><strong>summations:</strong> &quot;sum from i=1 to n of&quot;</div> <div><strong>greek:</strong> &quot;alpha&quot;, &quot;beta&quot;, &quot;gamma&quot;, &quot;pi&quot;, &quot;sigma&quot;</div> <div><strong>operators:</strong> &quot;times&quot;, &quot;divided by&quot;, &quot;plus or minus&quot;</div> <div><strong>superscripts:</strong> &quot;squared&quot;, &quot;cubed&quot;, &quot;to the power of&quot;</div> <div><strong>subscripts:</strong> &quot;sub&quot;, &quot;to the&quot;</div> <div><strong>symbols:</strong> &quot;therefore&quot;, &quot;because&quot;, &quot;infinity&quot;</div></div></div>  <div class="bg-yellow-50 p-4 rounded"><h4 class="text-lg font-medium text-yellow-900 mb-2" data-svelte-h="svelte-2zu42t">🎯 Try It Now!</h4> <p class="text-yellow-800 mb-3" data-svelte-h="svelte-1012kkc">Create a new math block and try typing any of these expressions:</p> <div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm"><button class="text-left p-2 bg-white rounded hover:bg-gray-50" data-svelte-h="svelte-3wuee0">&quot;fraction x squared plus y squared over z&quot;</button> <button class="text-left p-2 bg-white rounded hover:bg-gray-50" data-svelte-h="svelte-vqhues">&quot;integral from negative infinity to infinity of e to the power of negative x squared over square root of pi dx&quot;</button> <button class="text-left p-2 bg-white rounded hover:bg-gray-50" data-svelte-h="svelte-10lqru1">&quot;sum from k equals 0 to infinity of fraction 1 over k factorial&quot;</button> <button class="text-left p-2 bg-white rounded hover:bg-gray-50" data-svelte-h="svelte-k0g6il">&quot;limit as x approaches 0 of fraction sin x over x equals 1&quot;</button></div></div>  <div data-svelte-h="svelte-pgtslg"><h4 class="text-lg font-medium text-gray-900 mb-3">📤 Export Your Math</h4> <p class="text-gray-600 mb-3">Once you&#39;ve created beautiful math expressions, export them to various formats:</p> <div class="grid grid-cols-2 md:grid-cols-4 gap-3"><div class="text-center p-3 bg-red-50 rounded"><div class="text-red-600 font-medium">PDF</div> <div class="text-xs text-red-600">Vector graphics</div></div> <div class="text-center p-3 bg-blue-50 rounded"><div class="text-blue-600 font-medium">LaTeX</div> <div class="text-xs text-blue-600">Source code</div></div> <div class="text-center p-3 bg-green-50 rounded"><div class="text-green-600 font-medium">MathML</div> <div class="text-xs text-green-600">Web standard</div></div> <div class="text-center p-3 bg-purple-50 rounded"><div class="text-purple-600 font-medium">SVG</div> <div class="text-xs text-purple-600">Scalable images</div></div></div></div></div> <div class="mt-6 flex justify-end"><button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-1i4m1ki">Got it!</button></div></div></div>` : ``}  ${showFindReplaceDialog ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-lg"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-3vl8uq">Find &amp; Replace</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-1hbf56t">×</button></div> <div class="space-y-4"><div><label class="block text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-2few71">Find</label> <div class="flex space-x-2"><input type="text" placeholder="Enter text to find..." class="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"${add_attribute("value", findText, 0)}> <button class="${[
    "px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200",
    (isCaseSensitive ? "bg-blue-100" : "") + " " + (isCaseSensitive ? "text-blue-700" : "")
  ].join(" ").trim()}" title="Case sensitive" data-svelte-h="svelte-x173lu">Aa</button></div></div> <div><label class="block text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-h48976">Replace with</label> <input type="text" placeholder="Enter replacement text..." class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"${add_attribute("value", replaceText, 0)}></div> ${findResults.length > 0 ? `<div class="flex items-center justify-between bg-blue-50 p-3 rounded"><div class="text-sm text-blue-800">Found ${escape(findResults.length)} occurrence${escape(findResults.length !== 1 ? "s" : "")} ${currentFindIndex >= 0 ? `<span class="font-medium">(currently at ${escape(currentFindIndex + 1)})</span>` : ``}</div> <div class="flex space-x-2"><button class="px-2 py-1 text-sm bg-white text-gray-700 rounded hover:bg-gray-100" title="Previous (Shift+Enter)" data-svelte-h="svelte-1qebce9">◀ Previous</button> <button class="px-2 py-1 text-sm bg-white text-gray-700 rounded hover:bg-gray-100" title="Next (Enter)" data-svelte-h="svelte-11ou6ty">Next ▶</button></div></div>` : `${findText ? `<div class="text-sm text-gray-500 bg-gray-50 p-3 rounded" data-svelte-h="svelte-1wwm8wr">No matches found</div>` : ``}`}</div> <div class="mt-6 flex justify-end space-x-3"><button class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50" data-svelte-h="svelte-mb9qxm">Close</button> <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" ${!findText ? "disabled" : ""}>Find</button> ${findResults.length > 0 ? `<button class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700" data-svelte-h="svelte-61yf67">Replace</button> <button class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700" data-svelte-h="svelte-gq9fsg">Replace All</button>` : ``}</div></div></div>` : ``}  ${showAIAssistant ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-k65tzs">✨ AI Writing Suggestions</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-k83yo0">×</button></div> ${aiSuggestions.length === 0 ? `<div class="text-center py-8" data-svelte-h="svelte-6ngbmd"><div class="text-4xl mb-4">💭</div> <p class="text-gray-600">Generating suggestions...</p></div>` : `<div class="space-y-4">${each(aiSuggestions, (suggestion, index) => {
    return `<div class="border border-gray-200 rounded-lg p-4"><div class="flex items-start justify-between mb-2"><h4 class="font-medium text-gray-900">Suggestion ${escape(index + 1)}</h4> <button class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-1ids9iv">Apply
								</button></div> <p class="text-gray-700 whitespace-pre-wrap">${escape(suggestion)}</p> </div>`;
  })}</div>`} <div class="mt-6 flex justify-end"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-fk8vix">Close</button></div></div></div>` : ``}  ${documentSummary ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-3xl max-h-96 overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-z93ldr">📄 Document Summary</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-2am0ki">×</button></div> <div class="prose prose-sm max-w-none svelte-1u0uidk"><div class="whitespace-pre-wrap text-gray-700">${escape(documentSummary)}</div></div> <div class="mt-6 flex justify-end space-x-3"><button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-y9z7mt">Copy to Clipboard</button> <button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-ckn2jr">Close</button></div></div></div>` : ``}  ${showTextAnalysis ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-17e4ruf">📊 Text Analysis</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-a9mdul">×</button></div> ${textAnalysis ? `<div class="grid grid-cols-2 gap-6"><div class="space-y-4"><h4 class="font-medium text-gray-900" data-svelte-h="svelte-iphh8y">Statistics</h4> <div class="space-y-2 text-sm"><div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-1bn4jcy">Words:</span> <span class="font-medium">${escape(textAnalysis.word_count || 0)}</span></div> <div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-dlv14x">Characters:</span> <span class="font-medium">${escape(textAnalysis.char_count || 0)}</span></div> <div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-5u056v">Sentences:</span> <span class="font-medium">${escape(textAnalysis.sentence_count || 0)}</span></div> <div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-ym3c7a">Paragraphs:</span> <span class="font-medium">${escape(textAnalysis.paragraph_count || 0)}</span></div> <div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-134xolg">Reading Time:</span> <span class="font-medium">${escape(textAnalysis.reading_time || "N/A")}</span></div></div></div> <div class="space-y-4"><h4 class="font-medium text-gray-900" data-svelte-h="svelte-16xq337">Readability</h4> <div class="space-y-2 text-sm"><div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-1s1u2vy">Grade Level:</span> <span class="font-medium">${escape(textAnalysis.grade_level || "N/A")}</span></div> <div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-sgfp97">Readability:</span> <span class="font-medium">${escape(textAnalysis.readability_score || "N/A")}</span></div> <div class="flex justify-between"><span class="text-gray-600" data-svelte-h="svelte-17ig99h">Complexity:</span> <span class="font-medium">${escape(textAnalysis.complexity || "N/A")}</span></div></div></div></div> ${textAnalysis.key_phrases && textAnalysis.key_phrases.length > 0 ? `<div class="mt-6"><h4 class="font-medium text-gray-900 mb-3" data-svelte-h="svelte-1g75ha9">Key Phrases</h4> <div class="flex flex-wrap gap-2">${each(textAnalysis.key_phrases, (phrase) => {
    return `<span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">${escape(phrase)} </span>`;
  })}</div></div>` : ``} ${textAnalysis.sentiment ? `<div class="mt-6"><h4 class="font-medium text-gray-900 mb-3" data-svelte-h="svelte-hqvuqb">Sentiment Analysis</h4> <div class="flex items-center space-x-4"><div class="flex-1 bg-gray-200 rounded-full h-2"><div class="bg-green-500 h-2 rounded-full" style="${"width: " + escape(textAnalysis.sentiment.positive || 0, true) + "%"}"></div></div> <span class="text-sm text-gray-600">Positive: ${escape(textAnalysis.sentiment.positive || 0)}%</span></div> <div class="flex items-center space-x-4 mt-2"><div class="flex-1 bg-gray-200 rounded-full h-2"><div class="bg-red-500 h-2 rounded-full" style="${"width: " + escape(textAnalysis.sentiment.negative || 0, true) + "%"}"></div></div> <span class="text-sm text-gray-600">Negative: ${escape(textAnalysis.sentiment.negative || 0)}%</span></div></div>` : ``}` : `<div class="text-center py-8" data-svelte-h="svelte-yl3ap9"><p class="text-gray-600">No analysis data available</p></div>`} <div class="mt-6 flex justify-end"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-120m9fm">Close</button></div></div></div>` : ``}`;
});
const TextEditor = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { editorMode = "blocks" } = $$props;
  let { documentId = null } = $$props;
  let currentDocument = null;
  let documentTitle2 = "Untitled Document";
  let isSaving = false;
  let isLoading = false;
  let saveStatus = "";
  let hasUnsavedChanges = false;
  let blocks = [];
  let markdownContent = "";
  let selectedBlockIndex = 0;
  let richTextEditorElement = null;
  let showDocumentList = false;
  let showVersionHistory = false;
  let showMathHelp = false;
  let showFindReplaceDialog = false;
  let showAIAssistant = false;
  let documentList = [];
  let versionHistory = [];
  let aiSuggestions = [];
  let documentSummary = "";
  let showTextAnalysis = false;
  let textAnalysis = null;
  let findText = "";
  let replaceText = "";
  let findResults = [];
  let currentFindIndex = -1;
  let isCaseSensitive = false;
  let isReadingAloud = false;
  let availableVoices = [];
  let isGeneratingAI = false;
  let isGeneratingSummary = false;
  let canUndo = false;
  let canRedo = false;
  if ($$props.editorMode === void 0 && $$bindings.editorMode && editorMode !== void 0) $$bindings.editorMode(editorMode);
  if ($$props.documentId === void 0 && $$bindings.documentId && documentId !== void 0) $$bindings.documentId(documentId);
  let $$settled;
  let $$rendered;
  let previous_head = $$result.head;
  do {
    $$settled = true;
    $$result.head = previous_head;
    $$rendered = `   ${validate_component(TextEditorToolbar, "TextEditorToolbar").$$render(
      $$result,
      {
        documentTitle: documentTitle2,
        isSaving,
        saveStatus,
        hasUnsavedChanges,
        currentDocument,
        availableVoices,
        isReadingAloud,
        isGeneratingAI,
        isGeneratingSummary,
        canUndo,
        canRedo
      },
      {},
      {}
    )}    <div class="h-full overflow-y-auto p-8 bg-white"><div class="max-w-4xl mx-auto">${editorMode === "blocks" ? `${validate_component(TextEditorBlockEditor, "TextEditorBlockEditor").$$render($$result, { blocks, selectedBlockIndex }, {}, {})}` : `${editorMode === "markdown" ? `${validate_component(TextEditorMarkdownEditor, "TextEditorMarkdownEditor").$$render($$result, { markdownContent }, {}, {})}` : `${editorMode === "richtext" ? `${validate_component(TextEditorRichText, "TextEditorRichText").$$render(
      $$result,
      { editorElement: richTextEditorElement },
      {
        editorElement: ($$value) => {
          richTextEditorElement = $$value;
          $$settled = false;
        }
      },
      {}
    )}` : ``}`}`}</div></div>  ${validate_component(TextEditorModals, "TextEditorModals").$$render(
      $$result,
      {
        showDocumentList,
        documentList,
        isLoading,
        showVersionHistory,
        versionHistory,
        showMathHelp,
        showFindReplaceDialog,
        showAIAssistant,
        aiSuggestions,
        documentSummary,
        showTextAnalysis,
        textAnalysis,
        findResults,
        currentFindIndex,
        isCaseSensitive,
        findText,
        replaceText
      },
      {
        findText: ($$value) => {
          findText = $$value;
          $$settled = false;
        },
        replaceText: ($$value) => {
          replaceText = $$value;
          $$settled = false;
        }
      },
      {}
    )}`;
  } while (!$$settled);
  return $$rendered;
});
let documentTitle = "Untitled Document";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const data = null;
  const form = null;
  const params = null;
  let editorMode = "blocks";
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  let $$settled;
  let $$rendered;
  let previous_head = $$result.head;
  do {
    $$settled = true;
    $$result.head = previous_head;
    $$rendered = `  ${$$result.head += `<!-- HEAD_svelte-1pchwvm_START -->${$$result.title = `<title>${escape(documentTitle)} - TPT Text Editor</title>`, ""}<!-- HEAD_svelte-1pchwvm_END -->`, ""} <div class="h-screen flex flex-col bg-white"> <header class="flex items-center px-6 py-2 border-b border-gray-200 bg-gray-50 shrink-0"><span class="text-sm font-medium text-gray-500 mr-3" data-svelte-h="svelte-mhfb70">Editor Mode:</span> <div class="flex items-center bg-white rounded-lg border border-gray-200 p-1"><button class="${"px-3 py-1 text-sm rounded-md transition-colors " + escape(
      editorMode === "blocks" ? "bg-blue-600 text-white shadow-sm" : "text-gray-600 hover:bg-gray-100",
      true
    )}">Blocks</button> <button class="${"px-3 py-1 text-sm rounded-md transition-colors " + escape(
      editorMode === "markdown" ? "bg-blue-600 text-white shadow-sm" : "text-gray-600 hover:bg-gray-100",
      true
    )}">Markdown</button> <button class="${"px-3 py-1 text-sm rounded-md transition-colors " + escape(
      editorMode === "richtext" ? "bg-blue-600 text-white shadow-sm" : "text-gray-600 hover:bg-gray-100",
      true
    )}">Rich Text</button></div></header>  <div class="flex-1 overflow-hidden flex flex-col">${validate_component(TextEditor, "TextEditor").$$render(
      $$result,
      { editorMode },
      {
        editorMode: ($$value) => {
          editorMode = $$value;
          $$settled = false;
        }
      },
      {}
    )}</div></div>`;
  } while (!$$settled);
  return $$rendered;
});
export {
  Page as default
};
