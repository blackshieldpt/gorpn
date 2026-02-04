# gorpn

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Tests](https://github.com/blackshieldpt/gorpn/actions/workflows/test.yml/badge.svg)](https://github.com/blackshieldpt/gorpn/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/blackshieldpt/gorpn)](https://goreportcard.com/report/github.com/blackshieldpt/gorpn)
[![GoDoc](https://pkg.go.dev/badge/github.com/blackshieldpt/gorpn)](https://pkg.go.dev/github.com/blackshieldpt/gorpn)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Go package for evaluating mathematical expressions using the Shunting Yard algorithm and Reverse Polish Notation (RPN).

## Features

- **Operator precedence**: `^` > `*` `/` `%` > `+` `-`
- **Associativity**: Right-associative `^`, left-associative others
- **Parentheses**: Nested expressions supported
- **Variables**: Pass values via `map[string]float64`
- **Built-in constants**: `pi`, `e`, `phi`
- **Built-in functions**: 29 mathematical functions
- **Custom functions**: Register your own functions (fixed or variadic arguments)
- **Scientific notation**: `1.5e-3`, `2E10`

## Installation

```bash
go get github.com/blackshieldpt/gorpn
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/blackshieldpt/gorpn"
)

func main() {
    // Simple expression
    result, _ := gorpn.Eval("2 + 3 * 4", nil)
    fmt.Println(result) // 14

    // With variables
    vars := map[string]float64{"x": 3, "y": 4}
    result, _ = gorpn.Eval("sqrt(x^2 + y^2)", vars)
    fmt.Println(result) // 5
}
```

## Usage

### Basic Evaluation

```go
// One-off evaluation
result, err := gorpn.Eval("(2 + 3) * 4", nil)

// With variables
result, err := gorpn.Eval("x * y + z", map[string]float64{
    "x": 10,
    "y": 20,
    "z": 5,
})
```

### Calculator Instance

```go
calc := gorpn.New()

// Set persistent variables
calc.SetVariable("gravity", 9.81)
calc.SetVariables(map[string]float64{"mass": 70, "height": 1.75})

// Evaluate multiple expressions
force, _ := calc.Evaluate("mass * gravity")
bmi, _ := calc.Evaluate("mass / height^2")

// Clear variables
calc.ClearVariables()
```

### Custom Functions

```go
calc := gorpn.New()

// Single argument
calc.RegisterFunction("double", 1, func(args []float64) (float64, error) {
    return args[0] * 2, nil
})

// Multiple arguments
calc.RegisterFunction("clamp", 3, func(args []float64) (float64, error) {
    val, min, max := args[0], args[1], args[2]
    if val < min { return min, nil }
    if val > max { return max, nil }
    return val, nil
})

// Variadic (use -1 for argCount)
calc.RegisterFunction("product", -1, func(args []float64) (float64, error) {
    result := 1.0
    for _, v := range args {
        result *= v
    }
    return result, nil
})

calc.Evaluate("double(5)")           // 10
calc.Evaluate("clamp(15, 0, 10)")    // 10
calc.Evaluate("product(2, 3, 4, 5)") // 120
```

### One-off with Custom Functions

```go
funcs := map[string]gorpn.Function{
    "fahrenheit": {
        ArgCount: 1,
        Fn: func(args []float64) (float64, error) {
            return args[0]*9/5 + 32, nil
        },
    },
}

result, _ := gorpn.EvalWithFunctions("fahrenheit(100)", nil, funcs) // 212
```

## Operators

| Operator | Description | Precedence | Associativity |
|----------|-------------|------------|---------------|
| `^` | Power | 4 | Right |
| `*` | Multiply | 3 | Left |
| `/` | Divide | 3 | Left |
| `%` | Modulo | 3 | Left |
| `+` | Add | 2 | Left |
| `-` | Subtract | 2 | Left |

Unary `-` and `+` are supported: `-5 + 3`, `-(2 + 3)`, `2 * -3`

## Built-in Constants

| Constant | Value |
|----------|-------|
| `pi`, `PI` | 3.14159... |
| `e`, `E` | 2.71828... |
| `phi`, `PHI` | 1.61803... (golden ratio) |

## Built-in Functions

### Trigonometric
`sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`, `sinh`, `cosh`, `tanh`

### Power & Roots
`sqrt`, `cbrt`, `pow`, `exp`, `exp2`

### Logarithmic
`log` (natural), `log10`, `log2`

### Rounding
`abs`, `ceil`, `floor`, `round`, `trunc`

### Other
`hypot`, `mod`

### Variadic
`min`, `max`, `avg`, `sum`

## Examples

```go
// Operator precedence
"2 + 3 * 4"          // 14 (not 20)
"2 ^ 3 ^ 2"          // 512 (2^9, right-associative)

// Nested expressions
"((2 + 3) * (4 - 1)) / 3"  // 5

// Functions with expressions
"sqrt(x^2 + y^2)"          // Pythagorean theorem
"sin(deg * pi / 180)"      // Degrees to radians

// Variadic functions
"max(sqrt(4), sqrt(9), sqrt(16))"  // 4
"avg(10, 20, 30)"                   // 20

// Complex formulas
"P * (1 + r/n)^(n*t)"      // Compound interest
"0.5 * m * v^2"            // Kinetic energy
```

## Algorithm

The package uses a two-phase approach:

### 1. Shunting Yard Algorithm

Converts infix notation to Reverse Polish Notation (postfix):

```
Input:  (2 + 3) * 4
Output: 2 3 + 4 *
```

### 2. RPN Stack Evaluation

Evaluates postfix using a stack:

```
Token   Stack
2       [2]
3       [2, 3]
+       [5]       <- pop 2,3; push 5
4       [5, 4]
*       [20]      <- pop 5,4; push 20
```

## Error Handling

```go
result, err := gorpn.Eval("1 / 0", nil)
if err != nil {
    // "evaluation error: division by zero"
}

result, err = gorpn.Eval("sqrt(-1)", nil)
if err != nil {
    // "evaluation error: error in function sqrt: sqrt of negative number"
}

result, err = gorpn.Eval("x + 1", nil)
if err != nil {
    // "evaluation error: undefined variable: x"
}

result, err = gorpn.Eval("(2 + 3", nil)
if err != nil {
    // "parsing error: mismatched parentheses: extra '('"
}
```

## API Reference

### Functions

```go
// One-off evaluation
func Eval(expression string, vars map[string]float64) (float64, error)

// One-off with custom functions
func EvalWithFunctions(expression string, vars map[string]float64, functions map[string]Function) (float64, error)
```

### Calculator Type

```go
func New() *Calculator

func (c *Calculator) Evaluate(expression string) (float64, error)
func (c *Calculator) EvaluateWithVars(expression string, vars map[string]float64) (float64, error)

func (c *Calculator) SetVariable(name string, value float64)
func (c *Calculator) SetVariables(vars map[string]float64)
func (c *Calculator) ClearVariables()

func (c *Calculator) RegisterFunction(name string, argCount int, fn func(args []float64) (float64, error))
```

### Function Type

```go
type Function struct {
    ArgCount int  // Use -1 for variadic
    Fn       func(args []float64) (float64, error)
}
```

## License

MIT License - Copyright (c) 2026 BlackShield LDA
