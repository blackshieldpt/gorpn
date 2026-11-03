package gorpn

import (
	"fmt"
)

type associativity int

const (
	leftAssociative associativity = iota
	rightAssociative
)

type operatorInfo struct {
	Precedence    int
	Associativity associativity
}

var operators = map[string]operatorInfo{
	"+": {Precedence: 2, Associativity: leftAssociative},
	"-": {Precedence: 2, Associativity: leftAssociative},
	"*": {Precedence: 3, Associativity: leftAssociative},
	"/": {Precedence: 3, Associativity: leftAssociative},
	"%": {Precedence: 3, Associativity: leftAssociative},
	"^": {Precedence: 4, Associativity: rightAssociative},
}

// shuntingYardWithArgCounts converts infix tokens to RPN and tracks function argument counts
func shuntingYardWithArgCounts(tokens []token) ([]token, map[int]int, error) {
	var output []token
	var operatorStack []token
	argCounts := make(map[int]int) // Maps output index to arg count for functions
	var argCountStack []int        // Stack to track argument counts for nested functions
	var wasValueStack []bool       // Stack to track if we just saw a value (for counting args)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case tokenNumber, tokenVariable:
			output = append(output, token)
			// Mark that we saw a value for argument counting
			if len(wasValueStack) > 0 {
				wasValueStack[len(wasValueStack)-1] = true
			}

		case tokenFunction:
			operatorStack = append(operatorStack, token)
			argCountStack = append(argCountStack, 0)
			wasValueStack = append(wasValueStack, false)

		case tokenComma:
			// Pop operators until we find a left parenthesis
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type != tokenLeftParen {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, nil, fmt.Errorf("mismatched parentheses or misplaced comma")
			}
			// Increment argument count if we saw a value
			if len(argCountStack) > 0 && len(wasValueStack) > 0 && wasValueStack[len(wasValueStack)-1] {
				argCountStack[len(argCountStack)-1]++
				wasValueStack[len(wasValueStack)-1] = false
			}

		case tokenOperator:
			op1 := operators[token.Value]
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top.Type != tokenOperator {
					break
				}
				op2 := operators[top.Value]
				if (op1.Associativity == leftAssociative && op1.Precedence <= op2.Precedence) ||
					(op1.Associativity == rightAssociative && op1.Precedence < op2.Precedence) {
					output = append(output, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)

		case tokenLeftParen:
			operatorStack = append(operatorStack, token)
			// If this isn't a function call, push dummy values to stacks
			if len(operatorStack) < 2 || operatorStack[len(operatorStack)-2].Type != tokenFunction {
				// Not a function call, just grouping parentheses
			}

		case tokenRightParen:
			// Pop operators until we find the matching left parenthesis
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type != tokenLeftParen {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, nil, fmt.Errorf("mismatched parentheses: missing '('")
			}
			// Pop the left parenthesis
			operatorStack = operatorStack[:len(operatorStack)-1]

			// If there's a function on top of the stack, pop it to output
			if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type == tokenFunction {
				// Count the last argument if there was a value
				if len(argCountStack) > 0 && len(wasValueStack) > 0 && wasValueStack[len(wasValueStack)-1] {
					argCountStack[len(argCountStack)-1]++
				}

				fn := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]

				outputIdx := len(output)
				output = append(output, fn)

				// Store argument count
				if len(argCountStack) > 0 {
					argCounts[outputIdx] = argCountStack[len(argCountStack)-1]
					argCountStack = argCountStack[:len(argCountStack)-1]
				}
				if len(wasValueStack) > 0 {
					wasValueStack = wasValueStack[:len(wasValueStack)-1]
				}

				// Mark that we saw a value for outer function
				if len(wasValueStack) > 0 {
					wasValueStack[len(wasValueStack)-1] = true
				}
			}
		}
	}

	// Pop remaining operators
	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]
		if top.Type == tokenLeftParen {
			return nil, nil, fmt.Errorf("mismatched parentheses: extra '('")
		}
		output = append(output, top)
	}

	return output, argCounts, nil
}
