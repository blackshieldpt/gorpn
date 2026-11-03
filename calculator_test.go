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
		name    string
		expr    string
		wantErr string
	}{
		{"undefined variable", "x + 1", "evaluation error: undefined variable: x"},
		{"undefined function", "foo(5)", "evaluation error: undefined function: foo"},
		{"mismatched parens 1", "(2 + 3", "parsing error: mismatched parentheses: extra '('"},
		{"mismatched parens 2", "2 + 3)", "parsing error: mismatched parentheses: missing '('"},
		{"division by zero", "1 / 0", "evaluation error: division by zero"},
		{"sqrt negative", "sqrt(-1)", "evaluation error: error in function sqrt: sqrt of negative number"},
		{"log non-positive", "log(0)", "evaluation error: error in function log: log of non-positive number"},
		{"log negative", "log(-1)", "evaluation error: error in function log: log of non-positive number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Evaluate(tt.expr)
			if err == nil {
				t.Fatalf("expected error for expression: %s", tt.expr)
			}
			// Compare the whole message: several checks can reject the same
			// input, and only the message shows which one did
			if err.Error() != tt.wantErr {
				t.Errorf("got error %q, expected %q", err.Error(), tt.wantErr)
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

func TestUnaryMinus(t *testing.T) {
	vars := map[string]float64{"x": 3, "y": 2}
	tests := []struct {
		expr     string
		expected float64
	}{
		{"2 * -x", -6},            // after * before a variable
		{"10 / -y", -5},           // after / before a variable (was a division-by-zero error)
		{"2 ^ -y", 0.25},          // after ^ before a variable
		{"3 - -x", 6},             // after binary minus
		{"x % -y", 1},             // after %
		{"2 ^ -(1 + 1)", 0.25},    // after ^ before a parenthesis
		{"2 / -(1 + 1)", -1},      // after / before a parenthesis
		{"2 * -sqrt(4)", -4},      // before a function call
		{"-2 ^ 2", -4},            // power binds tighter than unary minus
		{"-y ^ 2", -4},            // same result as the literal form
		{"2 ^ -2", 0.25},          // negated literal exponent
		{"2 ^ -x ^ 2", 1.0 / 512}, // 2^(-(x^2))
		{"-x * 2", -6},            // unary minus binds tighter than *
		{"- -3", 3},               // double negation
		{"max(-x, -1)", -1},       // inside an argument list
		{"-2 ^ 0.5", -math.Sqrt2}, // -(2^0.5), not NaN
	}

	calc := New()
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.EvaluateWithVars(tt.expr, vars)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.IsNaN(result) || math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestUnaryMinusOfZeroIsPositiveZero(t *testing.T) {
	result, err := Eval("-x", map[string]float64{"x": 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 0 || math.Signbit(result) {
		t.Errorf("got %v (signbit %v), expected +0", result, math.Signbit(result))
	}
}

func TestUnaryMinusReachesFunctionDomainCheck(t *testing.T) {
	tests := []struct {
		expr    string
		wantErr string
	}{
		{"sqrt(-1)", "evaluation error: error in function sqrt: sqrt of negative number"},
		{"log(-1)", "evaluation error: error in function log: log of non-positive number"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			_, err := Eval(tt.expr, nil)
			if err == nil || err.Error() != tt.wantErr {
				t.Errorf("got %v, expected %q", err, tt.wantErr)
			}
		})
	}
}

func TestMalformedExpressions(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{"empty expression", "", "tokenization error: unexpected end of expression"},
		{"sign with nothing after it", "-", "tokenization error: unexpected end of expression"},
		{"trailing operator", "2 +", "tokenization error: unexpected end of expression"},
		{"sign before closing paren", "max(1, -)", "tokenization error: unexpected ')' at position 8"},
		{"sign before empty parens", "5 + sum(-())", "tokenization error: unexpected ')' at position 10"},
		{"sign before comma-only parens", "5 + sum(-(,))", "tokenization error: unexpected ',' at position 10"},
		{"sign before operator-only parens", "sum(-(*))", "tokenization error: unexpected operator '*' at position 6"},
		{"leading binary operator", "* 2 3", "tokenization error: unexpected operator '*' at position 0"},
		{"juxtaposed numbers", "2 3 +", "tokenization error: unexpected number at position 2"},
		{"juxtaposed identifiers", "x y", "tokenization error: unexpected identifier at position 2"},
		{"juxtaposed parens", "2 (3) *", "tokenization error: unexpected '(' at position 2"},
		{"empty parens", "()", "tokenization error: unexpected ')' at position 1"},
		{"empty parens inside a call", "sum(())", "tokenization error: unexpected ')' at position 5"},
		{"unary plus alone in a call", "sum(+)", "tokenization error: unexpected ')' at position 5"},
		{"empty middle argument", "max(1,,2)", "tokenization error: unexpected ',' at position 6"},
		{"trailing comma", "max(1,)", "tokenization error: unexpected ')' at position 6"},
		{"leading comma", "max(,1)", "tokenization error: unexpected ',' at position 4"},
		{"empty argument to fixed-arity function", "pow(2,,3)", "tokenization error: unexpected ',' at position 6"},
		{"comma in grouping parens", "max((1,2))", "tokenization error: misplaced comma at position 6: not inside a function call"},
		{"comma at top level", "1,2", "tokenization error: misplaced comma at position 1: not inside a function call"},
		{"unmatched open paren", "(2 + 3", "parsing error: mismatched parentheses: extra '('"},
		{"unmatched close paren", "2 + 3)", "parsing error: mismatched parentheses: missing '('"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Eval(tt.expr, nil)
			if err == nil {
				t.Fatalf("expected error for expression: %q", tt.expr)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("got error %q, expected %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestWellFormedExpressions(t *testing.T) {
	calc := New()
	calc.RegisterFunction("answer", 0, func(args []float64) (float64, error) {
		return 42, nil
	})

	tests := []struct {
		expr     string
		expected float64
	}{
		{"answer()", 42},         // zero-argument call
		{"sum()", 0},             // variadic with no arguments
		{"sum( )", 0},            // whitespace inside an empty call
		{"max((1 + 2), 3)", 3},   // grouping parens as an argument
		{"max(1, (2), 3)", 3},    // grouping parens in the middle
		{"max(-1, -2)", -1},      // unary minus arguments
		{"max(1, min(2, 3))", 2}, // nested call
		{"pow((1 + 1), (3))", 8}, // fixed arity with grouping parens
		{"-(-(2))", 2},           // nested negated groups
		{"-+3", -3},              // unary plus after unary minus
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if math.IsNaN(result) || math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNonFiniteResultsAreErrors(t *testing.T) {
	funcs := map[string]Function{
		"double": {ArgCount: 1, Fn: func(args []float64) (float64, error) { return args[0] * 2, nil }},
	}
	tests := []struct {
		expr    string
		wantErr string
	}{
		{"5 % 0", "evaluation error: result is not a finite number: NaN"},
		{"10 ^ 400", "evaluation error: result is not a finite number: +Inf"},
		{"(-8) ^ (1 / 3)", "evaluation error: result is not a finite number: NaN"},
		{"1e308 * 10", "evaluation error: result is not a finite number: +Inf"},
		{"-(1e308 * 10)", "evaluation error: result is not a finite number: -Inf"},
		{"mod(5, 0)", "evaluation error: result is not a finite number: NaN"},
		{"asin(2)", "evaluation error: result is not a finite number: NaN"},
		{"exp(1000)", "evaluation error: result is not a finite number: +Inf"},
		// A custom function that only passes an overflow through doesn't hide it
		{"double(exp(1000))", "evaluation error: result is not a finite number: +Inf"},
		// A custom function with a finite result doesn't switch the check off
		{"double(1) * 10 ^ 400", "evaluation error: result is not a finite number: +Inf"},
		// A finite variable or constant doesn't switch the check off either
		{"e ^ 1000", "evaluation error: result is not a finite number: +Inf"},
		// Unchanged: these have their own checks, which run first
		{"1 / 0", "evaluation error: division by zero"},
		{"sqrt(-1)", "evaluation error: error in function sqrt: sqrt of negative number"},
		{"log(0)", "evaluation error: error in function log: log of non-positive number"},
		{"log10(0)", "evaluation error: error in function log10: log10 of non-positive number"},
		{"log2(0)", "evaluation error: error in function log2: log2 of non-positive number"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			_, err := EvalWithFunctions(tt.expr, nil, funcs)
			if err == nil {
				t.Fatalf("expected error for expression: %s", tt.expr)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("got error %q, expected %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestNonFiniteIntermediateResultsAreFine(t *testing.T) {
	// Only the final result is checked, so an overflow that the rest of the
	// expression brings back to a finite value is not an error
	tests := []struct {
		expr     string
		expected float64
	}{
		{"1 / (1 + exp(1000))", 0},
		{"atan(10 ^ 400)", math.Pi / 2},
		{"min(10 ^ 400, 3)", 3},
		{"tanh(exp(1000))", 1},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := Eval(tt.expr, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.IsNaN(result) || math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNonFiniteInputsPropagate(t *testing.T) {
	calc := New()
	calc.RegisterFunction("nan", 0, func(args []float64) (float64, error) {
		return math.NaN(), nil
	})
	calc.SetVariables(map[string]float64{"inf": math.Inf(1), "notnum": math.NaN()})

	tests := []struct {
		expr     string
		expected float64
	}{
		{"inf * 2", math.Inf(1)},
		{"2 * inf", math.Inf(1)},
		{"inf - inf", math.NaN()},
		{"notnum + 1", math.NaN()},
		{"abs(inf)", math.Inf(1)},
		{"max(1, inf)", math.Inf(1)},
		{"nan()", math.NaN()},
		{"nan() + 1", math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := calc.Evaluate(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			same := result == tt.expected || (math.IsNaN(result) && math.IsNaN(tt.expected))
			if !same {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}
