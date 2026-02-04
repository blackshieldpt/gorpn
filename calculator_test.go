package gorpn

import (
	"math"
	"testing"
)

func TestBasicArithmetic(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"2 + 3", 5},
		{"10 - 4", 6},
		{"3 * 4", 12},
		{"15 / 3", 5},
		{"2 ^ 3", 8},
		{"10 % 3", 1},
		{"2 + 3 * 4", 14},     // Operator precedence
		{"(2 + 3) * 4", 20},   // Parentheses
		{"2 * 3 + 4 * 5", 26}, // Multiple operations
		{"2 ^ 3 ^ 2", 512},    // Right associative power (2^(3^2) = 2^9 = 512)
		{"(2 ^ 3) ^ 2", 64},   // Explicit left grouping
		{"10 / 2 / 5", 1},     // Left associative division
		{"-5 + 3", -2},        // Unary minus
		{"+5 + 3", 8},         // Unary plus
		{"-(2 + 3)", -5},      // Unary minus with parentheses
		{"2 * -3", -6},        // Negative in expression
		{"2.5 + 1.5", 4},      // Decimals
		{".5 + .5", 1},        // Decimals starting with dot
		{"1e3 + 1", 1001},     // Scientific notation
		{"1.5e2", 150},        // Scientific with decimal
		{"1e-3", 0.001},       // Negative exponent
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestVariables(t *testing.T) {
	tests := []struct {
		expr     string
		vars     map[string]float64
		expected float64
	}{
		{"x + y", map[string]float64{"x": 2, "y": 3}, 5},
		{"x * y + z", map[string]float64{"x": 2, "y": 3, "z": 4}, 10},
		{"x ^ 2", map[string]float64{"x": 5}, 25},
		{"2 * x + 3 * y", map[string]float64{"x": 4, "y": 2}, 14},
		{"(x + y) * (x - y)", map[string]float64{"x": 5, "y": 3}, 16},
		{"-x", map[string]float64{"x": 5}, -5},
		{"foo_bar + baz123", map[string]float64{"foo_bar": 1, "baz123": 2}, 3},
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.EvaluateWithVars(tt.expr, tt.vars)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBuiltinConstants(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"pi", math.Pi},
		{"PI", math.Pi},
		{"e", math.E},
		{"E", math.E},
		{"2 * pi", 2 * math.Pi},
		{"e ^ 2", math.E * math.E},
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"sin(0)", 0},
		{"cos(0)", 1},
		{"tan(0)", 0},
		{"sqrt(16)", 4},
		{"sqrt(2)", math.Sqrt(2)},
		{"abs(-5)", 5},
		{"abs(5)", 5},
		{"ceil(4.2)", 5},
		{"floor(4.8)", 4},
		{"round(4.5)", 5},
		{"round(4.4)", 4},
		{"exp(0)", 1},
		{"exp(1)", math.E},
		{"log(1)", 0},
		{"log(e)", 1},
		{"log10(100)", 2},
		{"log2(8)", 3},
		{"pow(2, 10)", 1024},
		{"hypot(3, 4)", 5},
		{"mod(10, 3)", 1},
		{"cbrt(27)", 3},
		{"sin(pi / 2)", 1},
		{"cos(pi)", -1},
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestVariadicFunctions(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"min(3, 1, 4, 1, 5)", 1},
		{"max(3, 1, 4, 1, 5)", 5},
		{"min(5)", 5},
		{"max(5)", 5},
		{"avg(2, 4, 6)", 4},
		{"sum(1, 2, 3, 4, 5)", 15},
		{"sum(10)", 10},
		{"min(1+1, 2+2, 3+3)", 2},
		{"max(sqrt(4), sqrt(9), sqrt(16))", 4},
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNestedFunctions(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"sqrt(pow(3, 2) + pow(4, 2))", 5},
		{"sin(asin(0.5))", 0.5},
		{"log(exp(5))", 5},
		{"abs(min(-5, -10, -3))", 10},
		{"max(abs(-5), abs(-3), abs(-10))", 10},
		{"pow(sqrt(2), 2)", 2},
		{"round(sin(pi / 6) * 2)", 1}, // sin(30°) = 0.5
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCustomFunctions(t *testing.T) {
	calc := New()

	// Register a custom function: double(x)
	calc.RegisterFunction("double", 1, func(args []float64) (float64, error) {
		return args[0] * 2, nil
	})

	// Register a custom function: add3(a, b, c)
	calc.RegisterFunction("add3", 3, func(args []float64) (float64, error) {
		return args[0] + args[1] + args[2], nil
	})

	// Register a custom variadic function: product(...)
	calc.RegisterFunction("product", -1, func(args []float64) (float64, error) {
		result := 1.0
		for _, v := range args {
			result *= v
		}
		return result, nil
	})

	tests := []struct {
		expr     string
		expected float64
	}{
		{"double(5)", 10},
		{"double(2 + 3)", 10},
		{"double(double(3))", 12},
		{"add3(1, 2, 3)", 6},
		{"add3(1+1, 2+2, 3+3)", 12},
		{"product(2, 3, 4)", 24},
		{"product(5)", 5},
		{"double(x) + 1", 11}, // with variable x=5
	}

	calc.SetVariable("x", 5)
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFunctionsWithExpressionArguments(t *testing.T) {
	tests := []struct {
		expr     string
		vars     map[string]float64
		expected float64
	}{
		{"sqrt(x^2 + y^2)", map[string]float64{"x": 3, "y": 4}, 5},
		{"pow(x + 1, y - 1)", map[string]float64{"x": 2, "y": 3}, 9},
		{"max(a * 2, b * 3, c * 4)", map[string]float64{"a": 5, "b": 4, "c": 2}, 12},
		{"sin(x * pi / 180)", map[string]float64{"x": 30}, 0.5}, // 30 degrees
		{"atan2(y, x)", map[string]float64{"x": 1, "y": 1}, math.Pi / 4},
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.EvaluateWithVars(tt.expr, tt.vars)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		expr     string
		vars     map[string]float64
		expected float64
	}{
		{"((2 + 3) * (4 - 1)) / 3", nil, 5},
		{"2 ^ 2 ^ 3", nil, 256},            // Right associative: 2^(2^3) = 2^8 = 256
		{"1 + 2 * 3 ^ 2 - 4 / 2", nil, 17}, // 1 + (2 * 9) - 2 = 17
		{"(a + b) * (c + d) / (e - f)", map[string]float64{"a": 1, "b": 2, "c": 3, "d": 4, "e": 10, "f": 3}, 3},
		{"sqrt((x2 - x1)^2 + (y2 - y1)^2)", map[string]float64{"x1": 0, "y1": 0, "x2": 3, "y2": 4}, 5},
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.EvaluateWithVars(tt.expr, tt.vars)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	calc := New()

	tests := []struct {
		name string
		expr string
	}{
		{"undefined variable", "x + 1"},
		{"undefined function", "foo(5)"},
		{"mismatched parens 1", "(2 + 3"},
		{"mismatched parens 2", "2 + 3)"},
		{"division by zero", "1 / 0"},
		{"sqrt negative", "sqrt(-1)"},
		{"log non-positive", "log(0)"},
		{"log negative", "log(-1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Evaluate(tt.expr)
			if err == nil {
				t.Errorf("expected error for expression: %s", tt.expr)
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test Eval function
	result, err := Eval("2 + 3 * x", map[string]float64{"x": 4})
	if err != nil {
		t.Errorf("Eval error: %v", err)
	}
	if result != 14 {
		t.Errorf("Eval: got %v, expected 14", result)
	}

	// Test EvalWithFunctions
	customFuncs := map[string]Function{
		"triple": {
			ArgCount: 1,
			Fn:       func(args []float64) (float64, error) { return args[0] * 3, nil },
		},
	}
	result, err = EvalWithFunctions("triple(5)", nil, customFuncs)
	if err != nil {
		t.Errorf("EvalWithFunctions error: %v", err)
	}
	if result != 15 {
		t.Errorf("EvalWithFunctions: got %v, expected 15", result)
	}
}

func TestCalculatorState(t *testing.T) {
	calc := New()

	// Set variables
	calc.SetVariable("x", 10)
	calc.SetVariable("y", 20)

	result, _ := calc.Evaluate("x + y")
	if result != 30 {
		t.Errorf("expected 30, got %v", result)
	}

	// SetVariables overwrites
	calc.SetVariables(map[string]float64{"x": 5, "z": 15})
	result, _ = calc.Evaluate("x + z")
	if result != 20 {
		t.Errorf("expected 20, got %v", result)
	}

	// y should still be there
	result, _ = calc.Evaluate("y")
	if result != 20 {
		t.Errorf("expected 20, got %v", result)
	}

	// Clear variables
	calc.ClearVariables()
	_, err := calc.Evaluate("x")
	if err == nil {
		t.Error("expected error after clearing variables")
	}
}

func BenchmarkSimpleExpression(b *testing.B) {
	calc := New()
	for i := 0; i < b.N; i++ {
		calc.Evaluate("2 + 3 * 4")
	}
}

func BenchmarkComplexExpression(b *testing.B) {
	calc := New()
	vars := map[string]float64{"x": 3, "y": 4}
	for i := 0; i < b.N; i++ {
		calc.EvaluateWithVars("sqrt(x^2 + y^2) + sin(pi/4) * cos(pi/4)", vars)
	}
}

func BenchmarkVariadicFunction(b *testing.B) {
	calc := New()
	for i := 0; i < b.N; i++ {
		calc.Evaluate("max(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)")
	}
}
