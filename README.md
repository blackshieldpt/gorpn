# gorpn

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Tests](https://github.com/blackshieldpt/gorpn/actions/workflows/test.yml/badge.svg)](https://github.com/blackshieldpt/gorpn/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/blackshieldpt/gorpn)](https://goreportcard.com/report/github.com/blackshieldpt/gorpn)
[![GoDoc](https://pkg.go.dev/badge/github.com/blackshieldpt/gorpn)](https://pkg.go.dev/github.com/blackshieldpt/gorpn)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Go package for evaluating mathematical expressions using the Shunting Yard algorithm and Reverse Polish Notation (RPN).

## Features

- **Operator precedence**: `^` > unary `-` > `*` `/` `%` > `+` `-`
- **Associativity**: Right-associative `^` and unary `-`, left-associative others
- **Parentheses**: Nested expressions supported
- **Variables**: Pass values via `map[string]float64`
- **Built-in constants**: `pi`, `e`, `phi`
- **Built-in functions**: 29 mathematical functions
- **Custom functions**: Register your own functions (fixed or variadic arguments)
- **Scientific notation**: `1.5e-3`, `2E10`
- **Compiled expressions**: parse once, evaluate many times

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
fmt.Println(result, err) // 20 <nil>

// With variables
result, err = gorpn.Eval("x * y + z", map[string]float64{
    "x": 10,
    "y": 20,
    "z": 5,
})
fmt.Println(result, err) // 205 <nil>
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
fmt.Println(force, bmi) // 686.7 22.857142857142858

// Clear variables
calc.ClearVariables()
```

### Name Resolution

A name in an expression resolves in this order: variables passed to the call (`EvaluateWithVars`, `Eval`, `EvalWithFunctions`, `Expr.Eval`), then the Calculator's variables, then the built-in constants. So a variable named `pi`, `e` or `phi` overrides the constant. Custom functions override built-in functions of the same name.

Each evaluation looks names up once, when it starts. A custom function that changes the Calculator (`SetVariable`, `ClearVariables`, `RegisterFunction`) or the `vars` map mid-evaluation doesn't affect the rest of that evaluation. The next evaluation sees the change.

A Calculator is not safe for concurrent use while it is being modified. Concurrent evaluations are fine as long as nothing changes the Calculator, or the `vars` and `functions` maps they use, at the same time.

Names are ASCII letters, digits and underscores, and can't start with a digit. Positions in error messages are 0-based byte offsets.

### Compiled Expressions

Parse an expression once and evaluate it many times, e.g. the same formula over many records:

```go
calc := gorpn.New()
calc.SetVariable("rate", 0.23)

expr, err := calc.Compile("price * (1 + rate)")
if err != nil {
    // syntax error: "tokenization error: ..." or "parsing error: ..."
}
for _, price := range []float64{10, 20, 30} {
    total, err := expr.Eval(map[string]float64{"price": price})
    if err != nil {
        // "evaluation error: ...", including invalid numbers such as 1e400
        // and undefined variables or functions
    }
    fmt.Println(total)
}
```

`c.Compile(expr)` followed by `Eval(vars)` gives the same result as `c.EvaluateWithVars(expr, vars)`. Each `Eval` sees the Calculator's variables and functions as they are when it starts. The package-level `gorpn.Compile` uses built-in functions and constants only. An `*Expr` never changes after `Compile`, so goroutines can share it as long as nothing modifies the Calculator at the same time.

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

result, err := gorpn.EvalWithFunctions("fahrenheit(100)", nil, funcs)
fmt.Println(result, err) // 212 <nil>
```

`RegisterFunction` and `EvalWithFunctions` take an exact argument count, or `-1` for any number. Any other negative count makes every call fail. A custom function overrides a built-in function of the same name.

## Operators

| Operator | Description | Precedence | Associativity |
|----------|-------------|------------|---------------|
| `^` | Power | 5 | Right |
| `-` (unary) | Negation | 4 | Right |
| `*` | Multiply | 3 | Left |
| `/` | Divide | 3 | Left |
| `%` | Modulo | 3 | Left |
| `+` | Add | 2 | Left |
| `-` | Subtract | 2 | Left |

Unary `-` and `+` work anywhere an operand can start: `-5 + 3`, `-(2 + 3)`, `2 * -x`, `2 ^ -n`. Unary minus binds tighter than `*` `/` `%` but looser than `^`, so `-2 ^ 2` is `-4`.

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
`ceil`, `floor`, `round`, `trunc`

### Other
`abs`, `hypot`, `mod`

### Variadic
`min`, `max`, `avg`, `sum`. With no arguments, `sum()` is `0`; `min`, `max` and `avg` return an error.

## Examples

```
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

An expression goes through four steps. `Compile` runs the first two once; each evaluation runs the last two.

### 1. Tokenizing

Splits the input into numbers, names, operators, parentheses and commas. It checks that operands and operators alternate, so input like `2 3 +` or `max(1,,2)` is a tokenization error. A `-` or `+` where an operand is expected is unary: `-` becomes a negation operator, `+` is dropped.

### 2. Shunting Yard Algorithm

Converts infix notation to Reverse Polish Notation (postfix), recording how many arguments each function call has:

```
Input:  (2 + 3) * 4
Output: 2 3 + 4 *
```

### 3. Name Resolution

Looks up each variable and function the expression uses, in the order under [Name Resolution](#name-resolution).

### 4. RPN Stack Evaluation

Evaluates postfix using a stack:

```
Token   Stack
2       [2]
3       [2, 3]
+       [5]       <- pop 2,3; push 5
4       [5, 4]
*       [20]      <- pop 5,4; push 20
```

If the final result is NaN or ±Inf, evaluation returns an error instead (see [Error Handling](#error-handling)).

## Error Handling

Every error starts with the stage that failed: `tokenization error:`, `parsing error:` or `evaluation error:`.

```go
_, err := gorpn.Eval("2 +", nil)
if err != nil {
    // "tokenization error: unexpected end of expression"
}

_, err = gorpn.Eval("1 / 0", nil)
if err != nil {
    // "evaluation error: division by zero"
}

_, err = gorpn.Eval("sqrt(-1)", nil)
if err != nil {
    // "evaluation error: error in function sqrt: sqrt of negative number"
}

_, err = gorpn.Eval("5 % 0", nil)
if err != nil {
    // "evaluation error: result is not a finite number: NaN"
}

_, err = gorpn.Eval("x + 1", nil)
if err != nil {
    // "evaluation error: undefined variable: x"
}

_, err = gorpn.Eval("(2 + 3", nil)
if err != nil {
    // "parsing error: mismatched parentheses: extra '('"
}
```

If an expression's result is NaN or ±Inf, you get an error instead (`5 % 0`, `asin(2)`, `10 ^ 400`). Only the final result is checked, so `1 / (1 + exp(1000))` still returns 0. If the expression uses a variable that is NaN or ±Inf, or a custom function that returns one from finite arguments, the result is passed through unchecked: check it yourself if you need a finite number.

## API Reference

### Functions

```go
// One-off evaluation
func Eval(expression string, vars map[string]float64) (float64, error)

// One-off with custom functions
func EvalWithFunctions(expression string, vars map[string]float64, functions map[string]Function) (float64, error)

// Parse once, evaluate many times (built-ins only)
func Compile(expression string) (*Expr, error)
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

func (c *Calculator) Compile(expression string) (*Expr, error)
```

### Expr Type

```go
func (e *Expr) Eval(vars map[string]float64) (float64, error)
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
