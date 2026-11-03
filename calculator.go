package gorpn

import (
	"fmt"
)

// Calculator is the main entry point for expression evaluation
type Calculator struct {
	variables map[string]float64
	functions map[string]Function
}

// New creates a new Calculator instance
func New() *Calculator {
	return &Calculator{
		variables: make(map[string]float64),
		functions: make(map[string]Function),
	}
}

// SetVariable sets a variable value
func (c *Calculator) SetVariable(name string, value float64) {
	c.variables[name] = value
}

// SetVariables sets multiple variables at once
func (c *Calculator) SetVariables(vars map[string]float64) {
	for k, v := range vars {
		c.variables[k] = v
	}
}

// ClearVariables removes all user-defined variables
func (c *Calculator) ClearVariables() {
	c.variables = make(map[string]float64)
}

// RegisterFunction registers a custom function
func (c *Calculator) RegisterFunction(name string, argCount int, fn func(args []float64) (float64, error)) {
	c.functions[name] = Function{
		ArgCount: argCount,
		Fn:       fn,
	}
}

// Evaluate evaluates an expression and returns the result
func (c *Calculator) Evaluate(expression string) (float64, error) {
	return c.EvaluateWithVars(expression, nil)
}

// EvaluateWithVars evaluates an expression with additional variables.
// These variables take precedence over previously set variables, which take
// precedence over built-in constants.
func (c *Calculator) EvaluateWithVars(expression string, vars map[string]float64) (float64, error) {
	expr, err := c.Compile(expression)
	if err != nil {
		return 0, err
	}
	return expr.Eval(vars)
}

// Expr is a parsed expression that can be evaluated many times. Create one with
// Compile; the zero value is not usable. It is never modified after Compile, so
// concurrent Eval calls are safe as long as nothing changes the Calculator it
// came from at the same time.
type Expr struct {
	calc      *Calculator
	rpn       []token
	argCounts map[int]int
}

// Compile parses an expression once, for repeated evaluation against this
// Calculator's variables and functions.
func (c *Calculator) Compile(expression string) (*Expr, error) {
	tokens, err := newTokenizer(expression).Tokenize()
	if err != nil {
		return nil, fmt.Errorf("tokenization error: %w", err)
	}

	rpn, argCounts, err := shuntingYardWithArgCounts(tokens)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return &Expr{calc: c, rpn: rpn, argCounts: argCounts}, nil
}

// Compile parses an expression once, for repeated evaluation with built-in
// functions and constants and the variables passed to Eval.
func Compile(expression string) (*Expr, error) {
	return New().Compile(expression)
}

// Eval evaluates the expression. vars take precedence over the Calculator's
// variables, which take precedence over built-in constants. Each call sees the
// Calculator's variables and functions as they are when it starts.
func (e *Expr) Eval(vars map[string]float64) (float64, error) {
	scope, functions := resolve(e.rpn, vars, e.calc.variables, e.calc.functions)
	result, err := evaluateRPNWithArgCounts(e.rpn, e.argCounts, scope, functions)
	if err != nil {
		return 0, fmt.Errorf("evaluation error: %w", err)
	}

	return result, nil
}

// Eval is a convenience function for one-off evaluation without creating a Calculator
func Eval(expression string, vars map[string]float64) (float64, error) {
	calc := New()
	return calc.EvaluateWithVars(expression, vars)
}

// EvalWithFunctions evaluates an expression with variables and custom functions
func EvalWithFunctions(expression string, vars map[string]float64, functions map[string]Function) (float64, error) {
	calc := New()
	// No copy needed: nothing else can reach calc to modify the map, and
	// evaluation only reads it
	calc.functions = functions
	return calc.EvaluateWithVars(expression, vars)
}
