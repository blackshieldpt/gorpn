package gorpn

import (
	"fmt"
	"math"
	"strconv"
)

// Function represents a custom function that can be registered with the calculator
type Function struct {
	// ArgCount is the number of arguments the function takes
	// Use -1 for variadic functions (variable number of arguments)
	ArgCount int
	// Fn is the actual function implementation
	Fn func(args []float64) (float64, error)
}

// builtinFunctions contains the default mathematical functions
var builtinFunctions = map[string]Function{
	"sin": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Sin(args[0]), nil },
	},
	"cos": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Cos(args[0]), nil },
	},
	"tan": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Tan(args[0]), nil },
	},
	"asin": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Asin(args[0]), nil },
	},
	"acos": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Acos(args[0]), nil },
	},
	"atan": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Atan(args[0]), nil },
	},
	"atan2": {
		ArgCount: 2,
		Fn:       func(args []float64) (float64, error) { return math.Atan2(args[0], args[1]), nil },
	},
	"sinh": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Sinh(args[0]), nil },
	},
	"cosh": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Cosh(args[0]), nil },
	},
	"tanh": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Tanh(args[0]), nil },
	},
	"sqrt": {
		ArgCount: 1,
		Fn: func(args []float64) (float64, error) {
			if args[0] < 0 {
				return 0, fmt.Errorf("sqrt of negative number")
			}
			return math.Sqrt(args[0]), nil
		},
	},
	"cbrt": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Cbrt(args[0]), nil },
	},
	"abs": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Abs(args[0]), nil },
	},
	"ceil": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Ceil(args[0]), nil },
	},
	"floor": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Floor(args[0]), nil },
	},
	"round": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Round(args[0]), nil },
	},
	"trunc": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Trunc(args[0]), nil },
	},
	"exp": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Exp(args[0]), nil },
	},
	"exp2": {
		ArgCount: 1,
		Fn:       func(args []float64) (float64, error) { return math.Exp2(args[0]), nil },
	},
	"log": {
		ArgCount: 1,
		Fn: func(args []float64) (float64, error) {
			if args[0] <= 0 {
				return 0, fmt.Errorf("log of non-positive number")
			}
			return math.Log(args[0]), nil
		},
	},
	"log10": {
		ArgCount: 1,
		Fn: func(args []float64) (float64, error) {
			if args[0] <= 0 {
				return 0, fmt.Errorf("log10 of non-positive number")
			}
			return math.Log10(args[0]), nil
		},
	},
	"log2": {
		ArgCount: 1,
		Fn: func(args []float64) (float64, error) {
			if args[0] <= 0 {
				return 0, fmt.Errorf("log2 of non-positive number")
			}
			return math.Log2(args[0]), nil
		},
	},
	"pow": {
		ArgCount: 2,
		Fn:       func(args []float64) (float64, error) { return math.Pow(args[0], args[1]), nil },
	},
	"min": {
		ArgCount: -1, // variadic
		Fn: func(args []float64) (float64, error) {
			if len(args) == 0 {
				return 0, fmt.Errorf("min requires at least one argument")
			}
			result := args[0]
			for _, v := range args[1:] {
				if v < result {
					result = v
				}
			}
			return result, nil
		},
	},
	"max": {
		ArgCount: -1, // variadic
		Fn: func(args []float64) (float64, error) {
			if len(args) == 0 {
				return 0, fmt.Errorf("max requires at least one argument")
			}
			result := args[0]
			for _, v := range args[1:] {
				if v > result {
					result = v
				}
			}
			return result, nil
		},
	},
	"avg": {
		ArgCount: -1, // variadic
		Fn: func(args []float64) (float64, error) {
			if len(args) == 0 {
				return 0, fmt.Errorf("avg requires at least one argument")
			}
			sum := 0.0
			for _, v := range args {
				sum += v
			}
			return sum / float64(len(args)), nil
		},
	},
	"sum": {
		ArgCount: -1, // variadic
		Fn: func(args []float64) (float64, error) {
			sum := 0.0
			for _, v := range args {
				sum += v
			}
			return sum, nil
		},
	},
	"mod": {
		ArgCount: 2,
		Fn:       func(args []float64) (float64, error) { return math.Mod(args[0], args[1]), nil },
	},
	"hypot": {
		ArgCount: 2,
		Fn:       func(args []float64) (float64, error) { return math.Hypot(args[0], args[1]), nil },
	},
}

// builtinConstants contains mathematical constants
var builtinConstants = map[string]float64{
	"pi":  math.Pi,
	"PI":  math.Pi,
	"e":   math.E,
	"E":   math.E,
	"phi": math.Phi,
	"PHI": math.Phi,
}

// EvaluateRPN evaluates an RPN (postfix) expression
func EvaluateRPN(tokens []Token, variables map[string]float64, customFunctions map[string]Function) (float64, error) {
	var stack []float64

	// Merge builtin constants with user variables (user variables take precedence)
	vars := make(map[string]float64)
	for k, v := range builtinConstants {
		vars[k] = v
	}
	for k, v := range variables {
		vars[k] = v
	}

	// Merge builtin functions with custom functions (custom take precedence)
	functions := make(map[string]Function)
	for k, v := range builtinFunctions {
		functions[k] = v
	}
	for k, v := range customFunctions {
		functions[k] = v
	}

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token.Value)
			}
			stack = append(stack, val)

		case TokenVariable:
			val, ok := vars[token.Value]
			if !ok {
				return 0, fmt.Errorf("undefined variable: %s", token.Value)
			}
			stack = append(stack, val)

		case TokenOperator:
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator: %s", token.Value)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token.Value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result = a / b
			case "%":
				result = math.Mod(a, b)
			case "^":
				result = math.Pow(a, b)
			default:
				return 0, fmt.Errorf("unknown operator: %s", token.Value)
			}
			stack = append(stack, result)

		case TokenFunction:
			fn, ok := functions[token.Value]
			if !ok {
				return 0, fmt.Errorf("undefined function: %s", token.Value)
			}

			argCount := fn.ArgCount
			if argCount == -1 {
				// Variadic function - need to determine argument count
				// We use a sentinel approach: count values on stack until we hit something
				// For simplicity, variadic functions consume all available values
				// This is a limitation - in practice, we need arg count markers
				// Let's use a simpler approach: variadic functions are limited to stack size
				argCount = len(stack)
			}

			if len(stack) < argCount {
				return 0, fmt.Errorf("insufficient arguments for function %s: expected %d, got %d",
					token.Value, argCount, len(stack))
			}

			args := make([]float64, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

			result, err := fn.Fn(args)
			if err != nil {
				return 0, fmt.Errorf("error in function %s: %w", token.Value, err)
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: stack has %d values instead of 1", len(stack))
	}

	return stack[0], nil
}

// EvaluateRPNWithArgCounts evaluates an RPN expression with explicit argument counts for functions
func EvaluateRPNWithArgCounts(tokens []Token, argCounts map[int]int, variables map[string]float64, customFunctions map[string]Function) (float64, error) {
	var stack []float64

	// Merge builtin constants with user variables (user variables take precedence)
	vars := make(map[string]float64)
	for k, v := range builtinConstants {
		vars[k] = v
	}
	for k, v := range variables {
		vars[k] = v
	}

	// Merge builtin functions with custom functions (custom take precedence)
	functions := make(map[string]Function)
	for k, v := range builtinFunctions {
		functions[k] = v
	}
	for k, v := range customFunctions {
		functions[k] = v
	}

	for i, token := range tokens {
		switch token.Type {
		case TokenNumber:
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token.Value)
			}
			stack = append(stack, val)

		case TokenVariable:
			val, ok := vars[token.Value]
			if !ok {
				return 0, fmt.Errorf("undefined variable: %s", token.Value)
			}
			stack = append(stack, val)

		case TokenOperator:
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator: %s", token.Value)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token.Value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result = a / b
			case "%":
				result = math.Mod(a, b)
			case "^":
				result = math.Pow(a, b)
			default:
				return 0, fmt.Errorf("unknown operator: %s", token.Value)
			}
			stack = append(stack, result)

		case TokenFunction:
			fn, ok := functions[token.Value]
			if !ok {
				return 0, fmt.Errorf("undefined function: %s", token.Value)
			}

			// Get argument count from the map
			argCount, hasArgCount := argCounts[i]
			if !hasArgCount {
				// Fallback to function's declared arg count
				argCount = fn.ArgCount
				if argCount == -1 {
					argCount = len(stack) // Use all stack for variadic without explicit count
				}
			}

			// For non-variadic functions, verify argument count
			if fn.ArgCount != -1 && argCount != fn.ArgCount {
				return 0, fmt.Errorf("wrong number of arguments for function %s: expected %d, got %d",
					token.Value, fn.ArgCount, argCount)
			}

			if len(stack) < argCount {
				return 0, fmt.Errorf("insufficient arguments for function %s: expected %d, got %d",
					token.Value, argCount, len(stack))
			}

			args := make([]float64, argCount)
			for j := argCount - 1; j >= 0; j-- {
				args[j] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

			result, err := fn.Fn(args)
			if err != nil {
				return 0, fmt.Errorf("error in function %s: %w", token.Value, err)
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: stack has %d values instead of 1", len(stack))
	}

	return stack[0], nil
}
