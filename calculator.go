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

// EvaluateWithVars evaluates an expression with additional variables
// These variables are merged with any previously set variables
func (c *Calculator) EvaluateWithVars(expression string, vars map[string]float64) (float64, error) {
	// Merge variables
	allVars := make(map[string]float64)
	for k, v := range c.variables {
		allVars[k] = v
	}
	for k, v := range vars {
		allVars[k] = v
	}

	// Tokenize
	tokenizer := NewTokenizer(expression)
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		return 0, fmt.Errorf("tokenization error: %w", err)
	}

	// Convert to RPN with argument count tracking
	rpn, argCounts, err := ShuntingYardWithArgCounts(tokens)
	if err != nil {
		return 0, fmt.Errorf("parsing error: %w", err)
	}

	// Evaluate
	result, err := EvaluateRPNWithArgCounts(rpn, argCounts, allVars, c.functions)
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
	for name, fn := range functions {
		calc.functions[name] = fn
	}
	return calc.EvaluateWithVars(expression, vars)
}
